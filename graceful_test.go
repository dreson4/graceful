package graceful

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestOnShutdownHandlersAreExecuted(t *testing.T) {
	Initialize()
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
	Wait()

	// Verify that all handlers were executed
	expected := int32(3)
	if counter != expected {
		t.Errorf("Expected %d handlers to be executed, got %d", expected, counter)
	}
}

func TestJobCounterLogic(t *testing.T) {
	Initialize()
	var (
		jobStarted   int32
		jobCompleted int32
	)

	// Define a job that increments jobStarted upon start,
	// waits, then signals completion
	job := func() {
		if IncrementJobCounter() {
			atomic.AddInt32(&jobStarted, 1)
			// Simulate job work...
			DecrementJobCounter()
			atomic.AddInt32(&jobCompleted, 1)
		}
	}

	// Start a few jobs
	go job()
	go job()
	go job()

	// Allow some time for jobs to start
	time.Sleep(100 * time.Millisecond)

	// Initiate shutdown process
	triggerShutdown()

	// Attempt to start another job which should fail as shutdown has started
	if IncrementJobCounter() {
		t.Fatal("Was able to start a job after shutdown process began")
	}

	// Wait for the shutdown to complete
	Wait()

	// Verify that all started jobs have completed
	if jobStarted != jobCompleted {
		t.Errorf("Not all jobs completed: started %d, completed %d", jobStarted, jobCompleted)
	}
}
