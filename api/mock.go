package api

// MockWebServer implements a WebServer interface for testing purposes.
type MockWebServer struct {
	MockStart    func() error
	MockShutdown func() error
}

// Start is a mock implementation for testing purposes.
func (mws *MockWebServer) Start() error {

	if mws.MockStart != nil {
		return mws.MockStart()
	}

	return nil
}

// Shutdown is a mock implementation for testing purposes.
func (mws *MockWebServer) Shutdown() error {

	if mws.MockShutdown != nil {
		return mws.MockShutdown()
	}

	return nil
}
