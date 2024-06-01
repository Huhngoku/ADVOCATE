package main

import (
	"testing"
	"time"
)

func TestSomething(t *testing.T) {
	n49()
}
func n49() {
	c := make(chan int, 0)
	d := make(chan int, 0)
	e := make(chan int, 0)

	go func() {
		c <- 1
	}()

	go func() {
		e <- 1 // prevents send from d to select
		d <- 1
	}()

	select {}

	<-e
	time.Sleep(100 * time.Millisecond)
}
