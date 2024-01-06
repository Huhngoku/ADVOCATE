package main

import (
	"advocate"
	"runtime"
	"sync"
)

func main() {
	runtime.InitAdvocate(0)
	defer advocate.CreateTrace("trace_new.log")

	v := sync.Mutex{}
	w := sync.Mutex{}
	x := sync.Mutex{}
	y := sync.Mutex{}
	z := sync.Mutex{}

	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		v.Lock()
		w.Lock()
		w.Unlock()
		v.Unlock()
		y.Lock()
		z.Lock()
		z.Unlock()
		x.Lock()
		x.Unlock()
		y.Unlock()
		wg.Done()
	}()

	go func() {
		w.Lock()
		x.Lock()
		x.Unlock()
		w.Unlock()
		wg.Done()
	}()

	x.Lock()
	v.Lock()
	v.Unlock()
	x.Unlock()

	wg.Wait()
}
