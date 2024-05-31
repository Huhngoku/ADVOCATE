
package main

import (
	"testing"
	"sync"
)

func Test02(t *testing.T) {
	n03()
}

// Once
// TN
func n03() {
	var once sync.Once
	ch := make(chan int, 1)
	setup := func() {
		ch <- 1
	}

	once.Do(setup)
	close(ch)

}
