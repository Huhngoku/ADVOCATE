package main

import (
	"sync"

	"github.com/ErikKassubek/exlib"
	ex "github.com/ErikKassubek/exlib"
)

func main() {
	var m1 sync.Mutex
	m2 := sync.RWMutex{}

	// c1 := make(chan int)
	// c2 := make(chan int, 2)

	// ============ Mutex ============
	ex.MutexInValue(m1)
	exlib.MutexInReference(&m1)

	_ = exlib.MutexInOutValue(m1)
	_ = exlib.MutexInOutReference(&m1)

	// ============ RWMutex ============
	exlib.RWMutexInValue(m2)
	exlib.RWMutexInReference(&m2)

	_ = exlib.RWMutexInOutValue(m2)
	_ = exlib.RWMutexInOutReference(&m2)

	// ============ Channel unbuffered ============
	// exlib.ChannelInValue(c1)
	// exlib.ChannelInReference(&c1)

	// _ = exlib.ChannelInOutValue(c1)
	// _ = exlib.ChannelInOutReference(&c1)

	// ============ Channel buffered ============
	// exlib.ChannelInValue(c2)
	// exlib.ChannelInReference(&c2)

	// _ = exlib.ChannelInOutValue(c2)
	// _ = exlib.ChannelInOutReference(&c2)

}
