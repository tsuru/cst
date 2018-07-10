package schd

import "github.com/tsuru/cst/scan"

// MockScheduler implements a Scheduler interface for testing purposes.
type MockScheduler struct {
	MockSchedule func(string) (scan.Scan, error)
}

// Schedule is a mock implementation for testing purposes.
func (ms *MockScheduler) Schedule(image string) (scan.Scan, error) {

	if ms.MockSchedule != nil {
		return ms.MockSchedule(image)
	}

	return scan.Scan{}, nil
}
