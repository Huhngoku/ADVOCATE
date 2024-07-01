package main

import (
	"sync"
	"testing"
)

func Test26(t *testing.T) {
	n26()
}

// possible negative wait counter
// TN
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

	wg.Done()
	wg.Done()

}
