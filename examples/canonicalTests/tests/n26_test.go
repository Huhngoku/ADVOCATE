package main

import (
	"sync"
	"testing"
	"time"
)

func Test26(t *testing.T) {
	n26()
}

// no possible negative wait counter
func n26() {
	var wg sync.WaitGroup
	c := make(chan int, 0)
	d := make(chan int, 0)

	go func() {
		wg.Add(1)
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		wg.Add(1)
		wg.Done()
		d <- 1
	}()

	go func() {
		wg.Add(1)
		<-d
		c <- 1
	}()

	<-c

	time.Sleep(100 * time.Millisecond) // prevent negative wait counter
	wg.Done()
	wg.Done()

	time.Sleep(200 * time.Millisecond)
}
