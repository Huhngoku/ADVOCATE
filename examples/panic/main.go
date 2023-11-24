package main

import (
	"cobufi"
	"runtime"
	"sync"
)

func main() {
	runtime.InitCobufi(0)
	defer cobufi.CreateTrace("trace_name.log")

	m := sync.Mutex{}
	m.Lock()
	m.Unlock()

	c := make(chan int)

	close(c)

	close(c)
}
