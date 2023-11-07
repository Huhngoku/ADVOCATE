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

	c := make(chan int)
	d := make(chan int)
	e := make(chan int)

	go func() {
		<-c
		println("a1")
	}()
	go func() {
		<-c
		println("a2")
	}()
	go func() {
		<-c
		println("a3")
	}()
	go func() {
		<-c
		println("a4")
	}()
	go func() {
		<-c
		println("a5")
	}()
	go func() {
		<-c
		println("a6")
	}()

	c <- 1
	c <- 1
	c <- 1

	time.Sleep(1 * time.Second)

	go func() {
		d <- 1
		println("b1")
	}()
	go func() {
		d <- 1
		println("b2")
	}()
	go func() {
		d <- 1
		println("b3")
	}()
	go func() {
		d <- 1
		println("b4")
	}()
	go func() {
		d <- 1
		println("b5")
	}()
	go func() {
		d <- 1
		println("b6")
	}()

	<-d
	<-d
	<-d

	time.Sleep(1 * time.Second)

	go func() {
		select {
		case <-e:
			println("c1")
		case <-e:
			println("c2")
		case <-e:
			println("c3")
		default:
			println("c4")
		}
	}()

	e <- 1
}
