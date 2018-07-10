package schd

import "github.com/tsuru/cst/scan"

// DefaultScheduler implements a Scheduler interface.
type DefaultScheduler struct{}

// Schedule registers a new analysis of a given image. It returns the complete
// entry of scan if successful else retuns an error instance to indicate the
// wrong state.
func (ds *DefaultScheduler) Schedule(image string) (scan.Scan, error) {

	return scan.Scan{}, nil
}
