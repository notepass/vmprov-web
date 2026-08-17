package libvirt

import (
	"context"
	"sync"
)

// FakeClient is a test double for Client.
type FakeClient struct {
	mu     sync.Mutex
	Result *TestResult
	Err    error
	Calls  []Connection
}

// TestConnection records the call and returns the configured result or error.
func (f *FakeClient) TestConnection(_ context.Context, conn Connection) (*TestResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Calls = append(f.Calls, conn)
	return f.Result, f.Err
}
