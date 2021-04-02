package serial_comm

import (
	"strings"
)

// *Serial.Validate checks if character is in range of permitted chars.
// reference: doc/grid_array_data_structure_{LANG}.pdf
// TODO: impl this
func (s *Serial) Validate(chars string) error {
	return nil
}

// TransmitData impl emit chars to specific device with grid650 spec
// this function will cast input into uppercase automatically.
func (s *Serial) TransmitData(chars string) error {
	// check chars before send
	if err := s.Validate(chars); err != nil {
		return err
	}

	// coded & send
	if err := s.Write(strings.ToUpper(chars), Text); err != nil {
		return err
	}

	return nil
}
