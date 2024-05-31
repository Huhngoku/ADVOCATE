
package main

import (
	"testing"
	"time"
)

func TestSomething(t *testing.T) {
	n10()
}

// FN
func n10() {
	c := make(chan struct{}, 0)

	go func() {
		time.Sleep(200 * time.Millisecond) // prevent actual send on closed channel
		close(c)
	}()

	go func() {
		select {
		case c <- struct{}{}:
		default:
		}
	}()

	time.Sleep(500 * time.Millisecond) // make sure, that the default values are taken
}
