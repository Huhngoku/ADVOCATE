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

	c := make(chan int)
	d := make(chan int)

	go func() {
		select {
		case c <- 1:
			println("c1")
		case <-d:
			println("d1")
		default:
			println("x1")
		}
	}()

	go func() {
		select {
		case c <- 1:
			println("c2")
		case <-d:
			println("d2")
		default:
			println("x2")
		}
	}()

	go func() {
		select {
		case c <- 1:
			println("c3")
		case <-d:
			println("d3")
		default:
			println("x3")
		}
	}()

	<-c
	d <- 1

	time.Sleep(time.Second)
}
