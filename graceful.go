package graceful

import (
	"context"
	"os/signal"
	"sync"
	"syscall"
)

var (
	activeJobs     = 0
	handlers       []func()
	mu             sync.Mutex
	done           chan struct{}
	jobsCond       = sync.NewCond(&mu)
	isShuttingDown bool
)

var didInit bool

func Initialize() {
	didInit = true
	activeJobs = 0
	handlers = make([]func(), 0)
	isShuttingDown = false
	done = make(chan struct{})
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ctx.Done()
		stop()
		triggerShutdown()
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
	if !didInit {
		panic("called graceful.Wait() without calling graceful.Initialize()")
	}
	<-done
}

func runHandlers() {
	mu.Lock()
	defer mu.Unlock()

	for _, handler := range handlers {
		handler()
	}
}

func triggerShutdown() {
	mu.Lock()
	if isShuttingDown {
		mu.Unlock()
		return
	}
	isShuttingDown = true // Ensure no new jobs can start.
	for activeJobs > 0 {
		jobsCond.Wait()
	}
	mu.Unlock()
	runHandlers()
	close(done)
}

// IncrementJobCounter - Call this when starting a job to prevent shutdown until job is finished.
func IncrementJobCounter() bool {
	mu.Lock()
	defer mu.Unlock()
	if isShuttingDown {
		return false
	}
	activeJobs++
	return true
}

// DecrementJobCounter - Call this when a job is finished.
func DecrementJobCounter() {
	mu.Lock()
	defer mu.Unlock()
	activeJobs--
	if activeJobs == 0 {
		jobsCond.Broadcast()
	}
}
