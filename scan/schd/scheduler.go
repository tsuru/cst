package schd

import (
	"errors"

	"github.com/tsuru/cst/scan"
)

var (
	// ErrImageHasAlreadyBeenSchedule indicates that current image couldn't be
	// scheduled because it's in a queue to be processed yet.
	ErrImageHasAlreadyBeenSchedule = errors.New(`this image has already been scheduled for scanning`)
)

// Scheduler is a basic interface to scheduling scans.
type Scheduler interface {
	Schedule(string) (scan.Scan, error)
}
