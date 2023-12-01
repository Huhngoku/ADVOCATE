package main

import (
	"advocate"
	"runtime"
	"sync"
)

func main() {
	runtime.InitAdvocate(0)
	defer advocate.CreateTrace("trace_name.log")

	m := sync.Mutex{}
	m.Lock()
	m.Unlock()

	c := make(chan int)

	close(c)

	close(c)
}
