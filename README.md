# Graceful Shutdown Handler

`graceful` is a Go package designed to help register and execute handlers that should be run during the shutdown of a Go application. 
It provides a convenient way to ensure that necessary cleanup operations are performed gracefully when the application is shutting down due to signals like SIGINT or SIGTERM.

## Installation

To use `graceful` in your Go project, simply import it using:

```go
go get github.com/dreson4/graceful
```

## Usage
In your main function 
```go
//initialize the graceful package before using it, add this to your main file init or start of main function
graceful.Initialize()

//blocking operation, it will block execution until app terminates, add this at the end of your main function
graceful.Wait()
```

Anywhere within the app register your handler which will be called when the app shuts down, could be a cleanup operation e.t.c
```go
func cleanUp(){

}

graceful.OnShutdown(cleanUp)
```


## Example
```go
package main

import (
	"fmt"
	"time"

	"github.com/dreson4/graceful"
)

func init(){
	graceful.Initialize()
}

func mySuperImportantFunction(){
	graceful.IncrementJobCounter()
	defer graceful.DecrementJobCounter()
	
	//do some important stuff that shouldn't be interrupted
}

func main() {
	fmt.Println("Starting application")

	// Register a sample shutdown handler
	graceful.OnShutdown(func() {
		fmt.Println("Cleanup operation started")
		time.Sleep(2 * time.Second)
		fmt.Println("Cleanup operation completed")
	})

	// Simulate application logic
	time.Sleep(5 * time.Second)
	
	mySuperImportantFunction()

	// Wait for graceful shutdown
	graceful.Wait()

	fmt.Println("Application shutdown")
}
```
