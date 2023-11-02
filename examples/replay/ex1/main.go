package main

import (
	"cobufi"
	"runtime"
	"sync"
	"time"
)

func main() {
	replay := true

	if !replay {
		// init tracing
		runtime.InitCobufi(0)
		defer cobufi.CreateTrace("trace.log")
	} else {
		// init replay
		trace := cobufi.ReadTrace("trace.log")
		runtime.EnableReplay(trace)
		defer runtime.WaitForReplayFinish()
	}

	m := sync.Mutex{}

	go func() {
		m.Lock()
		println(1)
		m.Unlock()
	}()

	go func() {
		m.Lock()
		println(2)
		m.Unlock()
	}()

	go func() {
		m.Lock()
		println(3)
		m.Unlock()
	}()

	go func() {
		m.Lock()
		println(4)
		m.Unlock()
	}()

	go func() {
		m.Lock()
		println(5)
		m.Unlock()
	}()

	time.Sleep(1 * time.Second)
}
