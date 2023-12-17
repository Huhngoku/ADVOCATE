package main

import (
	"advocate"
	"runtime"
	"time"
)

func main() {
	// runtime.InitAdvocate(0)
	// defer advocate.CreateTrace("trace_rewrite.log")

	trace := advocate.ReadTrace("rewritten_trace.log")
	runtime.EnableReplay(trace)
	defer runtime.WaitForReplayFinish()

	c := make(chan int)
	go func() {
		c <- 1
	}()

	go func() {
		<-c
	}()

	time.Sleep(1 * time.Second)
	close(c)

}
