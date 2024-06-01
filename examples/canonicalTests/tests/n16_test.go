
package main

import (
	"testing"
	"time"
	"sync"
)

func Test16(t *testing.T) {
	n16()
}
// TP
func n16() {
	c := make(chan int, 0)
	m := sync.Mutex{}

	go func() {
		t := m.TryLock()
		if t {
			c <- 1
			println("send")
			m.Unlock()
		}
	}()

	go func() {
		<-c
	}()

	time.Sleep(100 * time.Millisecond)
	close(c)
}
