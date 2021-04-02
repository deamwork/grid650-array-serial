package serial_comm

import (
	"time"
)

// ClockSync impl emit time to specific device with grid650 spec
// this function will cast input into the proper time format automatically.
func (s *Serial) ClockSync(newTime string) error {
	var targetTime time.Time

	now := time.Now()
	nt, err := time.Parse(time.RFC3339, newTime)

	if len(newTime) < 1 {
		// if not set custom time, use now
		targetTime = now
	} else if err != nil {
		// if set time but cannot be parsed, return error
		return err
	} else {
		// if set a correct time and pass the parser, use it
		targetTime = nt
	}

	// coded & send
	if err := s.Write(targetTime.Format("2006,01,02,15,04,05"), Time); err != nil {
		return err
	}

	return nil
}
