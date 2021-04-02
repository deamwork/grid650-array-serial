package serial_comm

import (
	"fmt"

	"github.com/XSAM/go-hybrid/errorw"
	"github.com/XSAM/go-hybrid/log"
	"github.com/tarm/serial"
	"go.uber.org/zap"
)

const (
	Time = iota
	Text
	Perf
)

const (
	TimeFormatter = "t%sm"
	TextFormatter = "c%sw"
	PerfFormatter = "h%s,%s,%s,%s,%s,%s,%s,%s,%s,%se"
)

type Serial struct {
	config *serial.Config
	port   *serial.Port
}

func NewSerial(device string, baud int) *Serial {
	return &Serial{
		config: &serial.Config{
			Name: device,
			Baud: baud,
		},
	}
}

func (s *Serial) Connect() error {
	port, err := serial.OpenPort(s.config)
	if err != nil {
		return err
	}

	s.port = port

	return nil
}

func (s Serial) Write(data string, kind int) error {
	var formatter string

	switch kind {
	case Time:
		formatter = TimeFormatter
	case Text:
		formatter = TextFormatter
	case Perf:
		formatter = PerfFormatter
	default:
		return errorw.NewMessage("no valid kind")
	}

	buf := []byte(fmt.Sprintf(formatter, data))
	if _, err := s.port.Write(buf); err != nil {
		return err
	}

	return nil
}

func (s *Serial) GracefulStop() {
	s.Close()
}

func (s *Serial) Close() error {
	if s.port != nil {
		return s.port.Close()
	}

	log.BgLogger().Error("core.emitter", zap.String("port_status", "port is no longer exist"))

	return nil
}
