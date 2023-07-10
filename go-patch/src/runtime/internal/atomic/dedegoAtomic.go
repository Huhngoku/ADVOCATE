package atomic

// DEDEGO-FILE-START

import (
	"sync/atomic"
	"unsafe"
)

var com chan<- AtomicElem
var linked bool
var counter atomic.Int64

type AtomicElem struct {
	Index int64
	Addr  int64
}

func DedegoAtomicLink(c chan<- AtomicElem) {
	com = c
	linked = true
}

func DedegoAtomicUnlink() {
	com = nil
	linked = false
}

//go:nosplit
func DedegoAtomicUInt32(addr *uint32) {
	if linked {
		counter.Add(1)
		// if line number changes, change in runtime/dedegoTrace.go DedegoChanSendPre
		com <- AtomicElem{Index: counter.Load(), Addr: int64(uintptr(unsafe.Pointer(addr)))}
	}
}

// DEDEGO-FILE-END
