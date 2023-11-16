package main

import "time"

func main() {
	// replay := false

	// if !replay {
	// 	// init tracing
	// 	runtime.InitCobufi(0)
	// 	defer cobufi.CreateTrace("trace.log")
	// } else {
	// 	// init replay
	// 	trace := cobufi.ReadTrace("trace.log")
	// 	runtime.EnableReplay(trace)
	// 	defer runtime.WaitForReplayFinish()
	// }

	go func() {
		println("Routine 1")
	}()

	go func() {
		println("Routine 2")
	}()

	time.Sleep(1 * time.Second)

}
