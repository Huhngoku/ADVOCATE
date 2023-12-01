package main

import (
	"advocate"
	"runtime"
)

func main() {
	runtime.InitAdvocate(0)
	defer advocate.CreateTrace("trace.log")

	trace := advocate.ReadTrace("trace.log")
	runtime.EnableReplay(trace)
	defer runtime.WaitForReplayFinish()

}
