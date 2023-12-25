package main

import (
	"advocate"
	"fmt"
	"runtime"
	"sync"
	"time"
)

var (
	condition bool
	mu        sync.Mutex
	cond      *sync.Cond
)

func waitForCondition() {
	mu.Lock()
	defer mu.Unlock()

	for !condition {
		cond.Wait()
	}

	fmt.Println("Condition is true now!")
}

func signalCondition() {
	mu.Lock()
	defer mu.Unlock()

	condition = true
	cond.Signal()
}

func main() {
	runtime.InitAdvocate(10)
	defer advocate.CreateTrace("trace_name.log")

	cond = sync.NewCond(&mu)

	go waitForCondition()

	time.Sleep(time.Second) // Simulate some work

	signalCondition()

	time.Sleep(time.Second) // Allow time for the goroutine to print the message
}
