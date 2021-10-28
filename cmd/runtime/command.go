package runtime

import (
	"fmt"
	"time"

	"github.com/XSAM/go-hybrid/cmdutil"
	"github.com/XSAM/go-hybrid/log"
	"github.com/XSAM/go-hybrid/metadata"
	"github.com/deamwork/grid650-array-serial/pkg/playback"
	"github.com/deamwork/grid650-array-serial/pkg/playback/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/deamwork/grid650-array-serial/config"
	"github.com/deamwork/grid650-array-serial/pkg/httpserver"
	serial "github.com/deamwork/grid650-array-serial/pkg/serial-comm"
)

var (
	flag Flag
)

func Start() {
	cmd := rootCmd()
	cmd.AddCommand(cmdutil.VersionCmd())
	cmd.AddCommand(configCmd())
	cmd.AddCommand(sendCmd())
	cmd.AddCommand(timeCmd())
	// TODO: impl this
	cmd.AddCommand(runCmd())
	cmd.Execute()
}

func rootCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:  metadata.AppName(),
		Long: "grid 650 array module serial-comm server",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if flag.Environment.DevelopmentMode {
				cfg := log.DevelopmentAndTextConfig()
				cfg.ZapConfig.EncoderConfig.TimeKey = "ts"
				log.BuildAndSetBgLogger(cfg)
			} else {
				cfg := log.ProductionAndTextConfig()
				cfg.ZapConfig.EncoderConfig.TimeKey = "ts"
				cfg.ZapConfig.DisableCaller = false
				log.BuildAndSetBgLogger(cfg)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	flag = Flag{
		//HTTPListen: "0.0.0.0:80",
		//ConfigFile: "../config/config.yaml",
		//Device:     "/dev/tty.usbmodem14233301",
		Environment: Environment{
			DevelopmentMode: false,
			JSONLogStyle:    false,
		},
	}
	cmdutil.ResolveFlagVariable(&cmd, &flag)

	return &cmd
}

func sendCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "send",
		Short: "Send one custom text to the grid650 array device",
		Long: "Send one custom text to the grid650 array device.\n" +
			"You must use quote if your text contains space(s)\n" +
			"You can adjust text position with space(s).\n" +
			"Limitation of the text length is 250 ASCII characters.",
		Args:    cobra.ExactArgs(1),
		Example: "./grid650-array-serial send \"i am grid 650\"",
		Run: func(cmd *cobra.Command, args []string) {
			rtConfig, err := config.LoadGrid650ArraySerialConfig(flag.ConfigFile)
			if err != nil {
				log.BgLogger().Fatal("core.config.parser", zap.Error(err))
				log.BgLogger().Error("err", zap.NamedError("dump", err))
			}

			// override config if flag is set
			if len(flag.Device) > 0 && flag.Baud > 0 {
				rtConfig.Device.Name = flag.Device
				rtConfig.Device.Baud = flag.Baud
			}

			log.BgLogger().Info("core.emitter", zap.Any("send_text", args[0]))

			// populate config & connect
			conn := serial.NewSerial(rtConfig.Device.Name, rtConfig.Device.Baud)
			if err := conn.Connect(); err != nil {
				log.BgLogger().Fatal("core.comm", zap.Error(err))
			}
			registerCloseHandler(conn)

			// send message
			if err := conn.TransmitData(args[0]); err != nil {
				log.BgLogger().Error("core.emitter", zap.Error(err))
			}
			log.BgLogger().Info("core.emitter", zap.Bool("send_ok", true))
		},
	}
	return &cmd
}

func timeCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "syncTime",
		Short: "Sync current time to the grid650 array device, only support the RFC3339",
		Long: "Sync current time to the grid650 array device, if you would like to set a time but not now,\n" +
			"you can add a time with RFC3339 timestamp supported string at the end of the arguments.\n" +
			"You can read some examples here: https://tools.ietf.org/html/rfc3339#section-5.8",
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			rtConfig, err := config.LoadGrid650ArraySerialConfig(flag.ConfigFile)
			if err != nil {
				log.BgLogger().Fatal("core.config.parser", zap.Error(err))
				log.BgLogger().Error("err", zap.NamedError("dump", err))
			}

			// override config if flag is set
			if len(flag.Device) > 0 && flag.Baud > 0 {
				rtConfig.Device.Name = flag.Device
				rtConfig.Device.Baud = flag.Baud
			}

			var argTime string
			if len(args) < 1 {
				argTime = ""
			} else {
				argTime = args[0]
			}

			log.BgLogger().Info("core.emitter", zap.Any("sync_time", argTime))

			// populate config & connect
			conn := serial.NewSerial(rtConfig.Device.Name, rtConfig.Device.Baud)
			if err := conn.Connect(); err != nil {
				log.BgLogger().Fatal("core.comm", zap.Error(err))
			}
			registerCloseHandler(conn)

			// send message
			if err := conn.ClockSync(argTime); err != nil {
				log.BgLogger().Error("core.emitter", zap.Error(err))
			}
			log.BgLogger().Info("core.emitter", zap.Bool("send_ok", true))
		},
	}
	return &cmd
}

func configCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "config",
		Short: "Print config parse result",
		Run: func(cmd *cobra.Command, args []string) {
			log.BgLogger().Info("core.flag.parser", zap.Any("flag", flag))
			rtConfig, err := config.LoadGrid650ArraySerialConfig(flag.ConfigFile)
			if err != nil {
				log.BgLogger().Fatal("core.config.parser", zap.Error(err))
				log.BgLogger().Error("err", zap.NamedError("dump", err))
			}
			log.BgLogger().Info("core.flag.parser", zap.Any("config", rtConfig))
		},
	}
	return &cmd
}

func runCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "run",
		Short: "Run HTTP server",
		Long:  "run server",
		Run: func(cmd *cobra.Command, args []string) {
			rtConfig, err := config.LoadGrid650ArraySerialConfig(flag.ConfigFile)
			if err != nil {
				log.BgLogger().Fatal("core.config.parser", zap.Error(err))
			}
			log.BgLogger().Info("core.flag.parser", zap.Any("config", rtConfig))

			if len(flag.HTTPListen) == 0 {
				flag.HTTPListen = rtConfig.HTTP.Listen
			}

			var httpServer *httpserver.HTTPServer

			// setup listener
			httpServer = httpserver.New(flag.HTTPListen)
			log.BgLogger().Info("core.config.rpc", zap.String("msg", "Using insecure tcp"))

			trackCh := make(chan utils.TrackInfo, 1)
			s := playback.NewSpotifyClient()
			go s.Start(rtConfig.Integration.Spotify.ClientID, rtConfig.Integration.Spotify.ClientSecret, httpServer, trackCh)

			registerCloseHandler(httpServer)

			// override config if flag is set
			if len(flag.Device) > 0 && flag.Baud > 0 {
				rtConfig.Device.Name = flag.Device
				rtConfig.Device.Baud = flag.Baud
			}

			// populate config & connect
			conn := serial.NewSerial(rtConfig.Device.Name, rtConfig.Device.Baud)
			if err := conn.Connect(); err != nil {
				log.BgLogger().Fatal("core.comm", zap.Error(err))
			}

			registerCloseHandler(conn)

			go func() {
				var lastSend string

				for {
					track := <-trackCh
					var song string
					if track.Name == "" {
						// not playing
						song = fmt.Sprintf("SPOTIFY - STAND BY")
					} else {
						// render text
						song = fmt.Sprintf("%s - %s", track.Artists[0], track.Name)
					}

					// prevent repeatedly send
					if lastSend != song {
						log.BgLogger().Info("spotify.track", zap.Any("track", track))
						log.BgLogger().Debug("core.emitter", zap.Any("send_text", song))
						if err := conn.TransmitData(song); err != nil {
							log.BgLogger().Error("core.emitter", zap.Error(err))
						}
						log.BgLogger().Debug("core.emitter", zap.Bool("send_ok", true))

						// save last sent result
						lastSend = song
					}

					time.Sleep(time.Second)
				}
			}()

			go httpServer.Serve()
			log.BgLogger().Info("core.config.http", zap.String("msg", "serial-comm server is ready."), zap.String("addr", flag.HTTPListen))
			// starts http server without blocking main thread.

			handleSysSignal()
		},
	}

	return &cmd
}
