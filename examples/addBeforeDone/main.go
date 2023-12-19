package main

import (
	"advocate"
	"runtime"
	"sync"
	"time"
)

func main() {
	// init tracing
	runtime.InitAdvocate(0)
	defer advocate.CreateTrace("trace_constructed.log")

	w := sync.WaitGroup{}

	go func() {
		w.Done()
	}()

	go func() {
		w.Add(1)
		w.Add(1)
		w.Done()
		w.Done()
	}()

	time.Sleep(1 * time.Second)
	w.Add(1)
	w.Done() // 1 < 1 + 1 - 2 = 0

	w.Wait()
}
