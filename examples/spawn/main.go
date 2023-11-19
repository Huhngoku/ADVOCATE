package main

import (
	"cobufi"
	"runtime"
	"time"
)

func main() {
	replay := true

	if !replay {
		// init tracing
		runtime.InitCobufi(0)
		defer cobufi.CreateTrace("trace.log")
	} else {
		// init replay
		trace := cobufi.ReadTrace("trace.log")
		runtime.EnableReplay(trace)
		defer runtime.WaitForReplayFinish()
	}

	go func() {
		println("Routine 1")
	}()

	go func() {
		println("Routine 2")
	}()

	time.Sleep(1 * time.Second)

}
