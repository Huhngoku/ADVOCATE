package main

import (
	"advocate"
	"runtime"
)

func main() {
	traceName := "trace2.log"
	if false {
		// init tracing
		runtime.InitAdvocate(0)
		defer advocate.CreateTrace(traceName)
	} else {
		// init replay
		trace := advocate.ReadTrace(traceName)
		runtime.EnableReplay(trace)
		defer runtime.WaitForReplayFinish()
	}

	c := make(chan int)
	d := make(chan int)

	x := false
	if x {
		close(c)
	}

	close(d)
}
