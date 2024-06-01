package main

import (
	"testing"
	"time"
)

func TestSomething(t *testing.T) {
	n39()
}

// ============= Leaking ==============

func n39() {
	c := make(chan int, 0)

	go func() {
		c <- 1
	}()

	time.Sleep(100 * time.Millisecond)
}
