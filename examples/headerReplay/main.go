package main

import (
	"cobufi"
	"runtime"
)

func main() {
	runtime.InitCobufi(0)
	defer cobufi.CreateTrace("trace.log")

	trace := cobufi.ReadTrace("trace.log")
	runtime.EnableReplay(trace)
	defer runtime.WaitForReplayFinish()

}
