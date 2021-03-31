package serial_comm

import (
	"fmt"
	"strings"

	"github.com/XSAM/go-hybrid/log"
	"github.com/tarm/serial"
	"go.uber.org/zap"
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

func (s Serial) Connect() error {
	port, err := serial.OpenPort(s.config)
	if err != nil {
		return err
	}

	s.port = port

	return nil
}

func (s Serial) GracefulStop() {
	s.Close()
}

func (s Serial) Close() error {
	if s.port != nil {
		return s.port.Close()
	}

	log.BgLogger().Error("core.emitter", zap.String("port_status", "port is no longer exist"))

	return nil
}

// TransmitData impl emit chars to specific device with grid650 spec
// this function will cast input into uppercase automatically.
func (s Serial) TransmitData(chars string) error {
	// check chars before send
	if err := s.Validate(chars); err != nil {
		return err
	}

	// coded & send
	buf := []byte(fmt.Sprintf("c%sw", strings.ToUpper(chars)))
	if _, err := s.port.Write(buf); err != nil {
		return err
	}

	return nil
}
