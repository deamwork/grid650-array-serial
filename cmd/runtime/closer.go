package runtime

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/XSAM/go-hybrid/log"
	"go.uber.org/zap"

	"github.com/deamwork/grid650-array-serial/pkg/httpserver"
)

var _ closer = (*httpserver.HTTPServer)(nil)

type closer interface {
	GracefulStop()
}

var closers = make([]closer, 0)

func registerCloseHandler(c closer) {
	closers = append(closers, c)
}

func handleSysSignal() {
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	sig := <-sigs

	log.BgLogger().Warn("signal", zap.String("msg", "got syscall signal, program will be quit"), zap.String("signal", sig.String()))

	for _, closer := range closers {
		closer.GracefulStop()
	}

	os.Exit(0)
}
