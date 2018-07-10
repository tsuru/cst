package db

import "github.com/tsuru/cst/scan"

// MockStorage implements a Storage interface for testing purposes.
type MockStorage struct {
	MockClose                   func()
	MockHasScheduledScanByImage func(string) bool
	MockSave                    func(scan.Scan) error
}

// Close is a mock implementation for testing purposes.
func (ms *MockStorage) Close() {

	if ms.MockClose != nil {
		ms.MockClose()
	}
}

// HasScheduledScanByImage is a mock implementation for testing purposes.
func (ms *MockStorage) HasScheduledScanByImage(image string) bool {

	if ms.MockHasScheduledScanByImage != nil {
		return ms.MockHasScheduledScanByImage(image)
	}

	return false
}

// Save is a mock implementation for testing purposes.
func (ms *MockStorage) Save(s scan.Scan) error {

	if ms.MockSave != nil {
		return ms.MockSave(s)
	}

	return nil
}
