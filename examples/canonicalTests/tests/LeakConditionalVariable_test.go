package main

import (
	"sync"
	"testing"
	"time"
)

func Test56(t *testing.T) {
	n56()
}

// leak because of conditional variable
func n56() {
	c := sync.NewCond(&sync.Mutex{})

	// wait for signal
	go func() {
		c.L.Lock()
		c.Wait()
		c.L.Unlock()
	}()

	time.Sleep(200 * time.Millisecond)
}
