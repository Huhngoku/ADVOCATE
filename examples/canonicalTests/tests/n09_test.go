package main

import (
	"testing"
	"time"
)

func Test09(t *testing.T) {
	n09()
}
// TP send on closed
// TP recv on closed
func n09() {
	c := make(chan struct{}, 1)
	d := make(chan struct{}, 1)

	go func() {
		time.Sleep(300 * time.Millisecond) // prevent actual send on closed channel
		close(c)
		close(d)
	}()

	go func() {
		select {
		case c <- struct{}{}:
		default:
		}

		select {
		case <-d:
		default:
		}
	}()

	d <- struct{}{}
	<-c

	time.Sleep(1 * time.Second) // prevent termination before receive
}
