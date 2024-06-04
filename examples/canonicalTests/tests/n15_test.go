package main

import (
	"sync"
	"testing"
	"time"
)

func Test15(t *testing.T) {
	n15()
}

// TN: no send possible
func n15() {
	c := make(chan int, 0)
	m := sync.Mutex{}

	go func() {
		t := m.TryLock()
		if t {
			c <- 1
			m.Unlock()
		}
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		t := m.TryLock()
		if t {
			m.Unlock()
		}
		<-c
	}()

	time.Sleep(1000 * time.Millisecond)
	close(c)
	time.Sleep(300 * time.Millisecond)
}
