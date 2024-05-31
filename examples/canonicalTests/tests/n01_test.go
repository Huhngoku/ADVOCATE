package main

import (
	"testing"
	"time"
)

func TestSomething(t *testing.T) {
	n01()
}

func n01() {
	x := make(chan int)
	ch := make(chan int, 1)

	go func() {
		ch <- 1
		x <- 1
	}()

	<-x
	close(ch)
}