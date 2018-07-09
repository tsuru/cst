package db

import "github.com/tsuru/cst/scan"

// MockStorage implements a Storage interface for testing purposes.
type MockStorage struct {
	MockClose func()
	MockSave  func(scan.Scan) error
}

// Close is a mock implementation for testing purposes.
func (ms *MockStorage) Close() {

	if ms.MockClose != nil {
		ms.MockClose()
	}
}

// Save is a mock implementation for testing purposes.
func (ms *MockStorage) Save(s scan.Scan) error {

	if ms.MockSave != nil {
		return ms.MockSave(s)
	}

	return nil
}
