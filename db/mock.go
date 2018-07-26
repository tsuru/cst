package db

import "github.com/tsuru/cst/scan"

// MockStorage implements a Storage interface for testing purposes.
type MockStorage struct {
	MockAppendResultToScanByID  func(string, scan.Result) error
	MockClose                   func()
	MockGetScansByImage         func(string) ([]scan.Scan, error)
	MockHasScheduledScanByImage func(string) bool
	MockSave                    func(scan.Scan) error
	MockUpdateScanStatusByID    func(string, scan.Status) error
	MockPing                    func() bool
}

// AppendResultToScanByID is a mock implementation for testing purposes.
func (ms *MockStorage) AppendResultToScanByID(id string, result scan.Result) error {

	if ms.MockAppendResultToScanByID != nil {
		return ms.MockAppendResultToScanByID(id, result)
	}

	return nil
}

// Close is a mock implementation for testing purposes.
func (ms *MockStorage) Close() {

	if ms.MockClose != nil {
		ms.MockClose()
	}
}

// GetScansByImage is a mock implementation for testing purposes.
func (ms *MockStorage) GetScansByImage(image string) ([]scan.Scan, error) {

	if ms.MockGetScansByImage != nil {
		return ms.MockGetScansByImage(image)
	}

	return []scan.Scan{}, nil
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

// UpdateScanStatusByID is a mock implementation for testing purposes.
func (ms *MockStorage) UpdateScanStatusByID(id string, status scan.Status) error {

	if ms.MockUpdateScanStatusByID != nil {
		return ms.MockUpdateScanStatusByID(id, status)
	}

	return nil
}

// Ping is a mock implementation for testing purposes.
func (ms *MockStorage) Ping() bool {

	if ms.MockPing != nil {
		return ms.MockPing()
	}

	return false
}
