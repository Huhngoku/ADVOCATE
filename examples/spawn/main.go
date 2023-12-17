package main

import (
	"advocate"
	"runtime"
	"time"
)

func main() {
	replay := true

	if !replay {
		// init tracing
		runtime.InitAdvocate(0)
		defer advocate.CreateTrace("trace.log")
	} else {
		// init replay
		trace := advocate.ReadTrace("trace.log")
		runtime.EnableReplay(trace)
		defer runtime.WaitForReplayFinish()
	}

	go func() {
		println("Routine 1")
	}()

	go func() {
		println("Routine 2")
	}()

	println("Main routine")

	time.Sleep(1 * time.Second)

	println("Main routine exit")

}
