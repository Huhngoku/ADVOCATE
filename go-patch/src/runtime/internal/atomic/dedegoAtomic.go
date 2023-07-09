package atomic

// DEDEGO-FILE-START

import "unsafe"

var com chan<- uintptr
var linked bool

func DedegoAtomicLink(c chan<- uintptr) {
	com = c
	linked = true
}

func DedegoAtomicUnlink() {
	com = nil
	linked = false
}

func DedegoAtomicXadd(addr *int32, delta int32) {
	if linked {
		// if line number changes, change in runtime/dedegoTrace.go DedegoChanSendPre
		com <- uintptr(unsafe.Pointer(addr))
	}
}

// DEDEGO-FILE-END
