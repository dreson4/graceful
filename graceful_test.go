package graceful

import (
	"sync/atomic"
	"testing"
)

func TestOnShutdownHandlersAreExecuted(t *testing.T) {
	var counter int32 = 0

	// Register three handlers that increment the counter
	OnShutdown(func() {
		atomic.AddInt32(&counter, 1)
	})
	OnShutdown(func() {
		atomic.AddInt32(&counter, 1)
	})
	OnShutdown(func() {
		atomic.AddInt32(&counter, 1)
	})

	// Trigger the shutdown process
	triggerShutdown()

	// Verify that all handlers were executed
	expected := int32(3)
	if counter != expected {
		t.Errorf("Expected %d handlers to be executed, got %d", expected, counter)
	}
}
