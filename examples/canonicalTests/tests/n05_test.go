
package main

import (
	"testing"
	"time"
)

func Test05(t *testing.T) {
	n05()
}

// TP send on closed
// TP recv on closed
func n05() {
	c := make(chan int)

	go func() {
		c <- 1
	}()

	go func() {
		<-c
	}()

	time.Sleep(300 * time.Millisecond) // prevent actual send on closed channel
	close(c)
}
