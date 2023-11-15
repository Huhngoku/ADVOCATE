package main

import (
	"cobufi"
	"runtime"
)

func main() {
	replay := false

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

	}()

	go func() {

	}()

	go func() {

	}()

	go func() {

	}()

	go func() {

	}()
}
