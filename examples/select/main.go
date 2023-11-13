package main

import (
	"cobufi"
	"runtime"
	"time"
)

func main() {
	replay := false

	if replay {
		trace := cobufi.ReadTrace("select_trace.log")
		runtime.EnableReplay(trace)
		defer runtime.WaitForReplayFinish()
	} else {
		runtime.InitCobufi(0)
		defer cobufi.CreateTrace("select_trace.log")
	}

	c := make(chan int)
	d := make(chan int)
	e := make(chan int)
	f := make(chan int)

	// recv send default
	go func() {
		select {
		case d <- 1:
			println(11)
		case f <- 1:
			println(12)
		case <-c:
			println(13)
			// do nothing
			// do nothing
		case <-e:
			println(4)
			// do nothing
			// do nothing
		}

		select {
		case d <- 1:
			println(21)
		case f <- 1:
			println(22)
		case <-c:
			println(23)
			// do nothing
			// do nothing
		case <-e:
			println(24)
			// do nothing
			// do nothing
		}

		select {
		case d <- 1:
			println(31)
		case f <- 1:
			println(32)
		case <-c:
			println(33)
			// do nothing
			// do nothing
		case <-e:
			println(34)
			// do nothing
			// do nothing
		}

		select {
		case d <- 1:
			println(41)
		case f <- 1:
			println(42)
		case <-c:
			println(43)
			// do nothing
			// do nothing
		case <-e:
			println(44)
			// do nothing
			// do nothing
		}

		select {
		case d <- 1:
			println(51)
		case f <- 1:
			println(52)
		case <-c:
			println(53)
			// do nothing
			// do nothing
		case <-e:
			println(54)
			// do nothing
			// do nothing
		default:
			println("D")
			// do nothing
		}
	}()

	time.Sleep(100 * time.Millisecond)

	go func() {
		c <- 1
	}()
	time.Sleep(100 * time.Millisecond)
	go func() {
		<-d
	}()
	time.Sleep(100 * time.Millisecond)
	go func() {
		e <- 1
	}()
	time.Sleep(100 * time.Millisecond)
	go func() {
		<-f
	}()

}
