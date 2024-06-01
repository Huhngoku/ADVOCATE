package main

import (
	"testing"
	"time"
)

func TestSomething(t *testing.T) {
	n50()
}

func n50() {
	c := make(chan int, 0)
	d := make(chan int, 1)
	e := make(chan int, 0)

	go func() {
		c <- 1
	}()

	go func() {
		d <- 1
		e <- 1 // prevents d from sending unbuffered
	}()

	<-e

	select {
	case <-c:
	case <-d:
	}

	time.Sleep(100 * time.Millisecond)
}
