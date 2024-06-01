
package main

import (
	"testing"
)

func Test07(t *testing.T) {
	n07()
}

// TN recv/send on closed
func n07() {
	c := make(chan int)

	go func() {
		c <- 1
	}()

	<-c

	close(c)
}
