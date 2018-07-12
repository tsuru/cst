package scan

// MockScanner is a mock implementation for testing purposes.
type MockScanner struct {
	MockScan func(string) Result
}

// Scan is a mock implementation for testing purposes.
func (ms *MockScanner) Scan(image string) Result {

	if ms.MockScan != nil {
		return ms.MockScan(image)
	}

	return Result{}
}
