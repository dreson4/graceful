# Graceful Shutdown Handler

`graceful` is a Go package designed to help register and execute handlers that should be run during the shutdown of a Go application. 
It provides a convenient way to ensure that necessary cleanup operations are performed gracefully when the application is shutting down due to signals like SIGINT or SIGTERM.

## Installation

To use `graceful` in your Go project, simply import it using:

```go
go get github.com/dreson4/graceful
```

## Usage
In your main function at the end you can add
```go
//blocking operation, it will block execution until app terminates
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

	// Wait for graceful shutdown
	graceful.Wait()

	fmt.Println("Application shutdown")
}
```
