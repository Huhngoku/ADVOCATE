
package main

import (
	"testing"
	"time"
	"sync"
)

func Test17(t *testing.T) {
	n17()
}


// FN
func n17() {
	c := make(chan int, 0)
	m := sync.Mutex{}

	go func() {
		time.Sleep(100 * time.Millisecond)
		t := m.TryLock()
		if t {
			c <- 1
			<-c
			m.Unlock()
		}
	}()

	m.Lock()
	time.Sleep(300 * time.Millisecond)
	close(c)
	m.Unlock()

	time.Sleep(100 * time.Millisecond)
}
