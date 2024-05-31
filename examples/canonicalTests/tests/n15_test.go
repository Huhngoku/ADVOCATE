
package main

import (
	"testing"
	"time"
	"sync"
)

func TestSomething(t *testing.T) {
	n15()
}

// TN: no send possible
func n15() {
	c := make(chan int, 0)
	m := sync.Mutex{}

	go func() {
		time.Sleep(100 * time.Millisecond)
		t := m.TryLock()
		if t {
			c <- 1
			m.Unlock()
		}
	}()

	go func() {
		t := m.TryLock()
		if t {
			<-c
			m.Unlock()
		}
	}()

	time.Sleep(600 * time.Millisecond)
	close(c)
	time.Sleep(300 * time.Millisecond)
}