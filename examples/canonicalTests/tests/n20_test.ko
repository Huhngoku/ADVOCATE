
package main

import (
	"testing"
	"time"
)

func Test20(t *testing.T) {
	n20()
}
// TP
func n20() {
	ch := make(chan int, 2)

	f := func() {
		ch <- 1
	}

	go func() {
		f()
	}()

	time.Sleep(200 * time.Millisecond)
	close(ch)
}
