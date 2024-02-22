package graceful

import (
	"context"
	"os/signal"
	"sync"
	"syscall"
)

var (
	handlers []func()
	mu       sync.Mutex
	done     chan struct{}
)

func init() {
	done = make(chan struct{})
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ctx.Done()
		stop()
		runHandlers()
		close(done)
	}()
}

// OnShutdown - Register a new handler to be executed on shutdown.
func OnShutdown(handler func()) {
	mu.Lock()
	defer mu.Unlock()
	handlers = append(handlers, handler)
}

// Wait for shutdown to complete before exiting the main function.
// This function should be called at the end of the main function.
func Wait() {
	<-done
}

func runHandlers() {
	mu.Lock()
	defer mu.Unlock()

	for _, handler := range handlers {
		handler()
	}
}
