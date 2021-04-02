package runtime

import (
	"github.com/XSAM/go-hybrid/cmdutil"
	"github.com/XSAM/go-hybrid/log"
	"github.com/XSAM/go-hybrid/metadata"
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
	//cmd.AddCommand(runCmd())
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
		HTTPListen: "0.0.0.0:80",
		ConfigFile: "../config/config.yaml",
		Device:     "/dev/tty.usbmodem14233301",
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
		Args: cobra.ExactArgs(1),
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

			log.BgLogger().Info("core.emitter", zap.Any("sync_time", args[0]))

			// populate config & connect
			conn := serial.NewSerial(rtConfig.Device.Name, rtConfig.Device.Baud)
			if err := conn.Connect(); err != nil {
				log.BgLogger().Fatal("core.comm", zap.Error(err))
			}
			registerCloseHandler(conn)

			// send message
			if err := conn.ClockSync(args[0]); err != nil {
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

			if len(flag.HTTPListen) > 0 {
				var httpServer *httpserver.HTTPServer

				// setup listener
				httpServer = httpserver.New(flag.HTTPListen)
				log.BgLogger().Info("core.config.rpc", zap.String("msg", "Using insecure tcp"))

				registerCloseHandler(httpServer)

				// starts http server without blocking main thread.
				go httpServer.Serve()
				log.BgLogger().Info("core.config.http", zap.String("msg", "serial-comm server is ready."), zap.String("addr", flag.HTTPListen))
			}

			handleSysSignal()
		},
	}

	return &cmd
}
