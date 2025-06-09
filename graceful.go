package graceful

import (
	"context"
	"os/signal"
	"sync"
	"syscall"
)

var (
	activeJobs       = 0
	handlers         []func()
	mu               sync.Mutex
	shutdownComplete chan struct{}
	jobsCond         = sync.NewCond(&mu)
	isShuttingDown   bool
	initOnce         sync.Once
	shutdownCtx      context.Context
)

func Initialize() {
	initOnce.Do(func() {
		activeJobs = 0
		handlers = make([]func(), 0)
		isShuttingDown = false
		shutdownComplete = make(chan struct{})
		var stop context.CancelFunc
		shutdownCtx, stop = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-shutdownCtx.Done()
			stop()
			triggerShutdown()
		}()
	})
}

func ShutdownContext() context.Context {
	return shutdownCtx
}

// OnShutdown - Register a new handler to be executed on shutdown.
func OnShutdown(handler func()) {
	mu.Lock()
	defer mu.Unlock()
	handlers = append(handlers, handler)
}

// Wait for shutdown to complete before exiting the main function.
func Wait() {
	if shutdownComplete == nil {
		panic("called graceful.Wait() without calling graceful.Start()")
	}
	<-shutdownComplete
}

func runHandlers() {
	mu.Lock()
	handlersCopy := append([]func(){}, handlers...)
	mu.Unlock()

	for _, handler := range handlersCopy {
		handler()
	}
}

func triggerShutdown() {
	mu.Lock()
	if isShuttingDown {
		mu.Unlock()
		return
	}
	isShuttingDown = true
	mu.Unlock()

	func() {
		mu.Lock()
		defer mu.Unlock()
		for activeJobs > 0 {
			jobsCond.Wait()
		}
	}()

	runHandlers()
	close(shutdownComplete)
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
