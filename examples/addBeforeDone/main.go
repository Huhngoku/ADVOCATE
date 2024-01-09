package main

import (
	"advocate"
	"runtime"
	"sync"
	"time"
)

func main() {
	// init tracing
	if false {
		runtime.InitAdvocate(0)
		defer advocate.CreateTrace("trace_new.log")
	} else {
		// init replay
		trace := advocate.ReadTrace("rewritten_trace.log")
		runtime.EnableReplay(trace)
		defer runtime.WaitForReplayFinish()
	}

	w := sync.WaitGroup{}

	go func() {
		time.Sleep(1 * time.Second)
		w.Done()
	}()

	go func() {
		w.Add(1)
		w.Add(1)
		w.Done()
		w.Done()
	}()

	w.Add(1)
	w.Add(1)
	w.Add(1)
	w.Done()
	w.Done()

	w.Wait()
	time.Sleep(2 * time.Second)
}
