package atomic

// COBUFI-FILE-START

import (
	"unsafe"
)

var chanRecording chan<- AtomicElem
var linked bool
var counter uint64

const (
	LoadOp = iota
	StoreOp
	AddOp
	SwapOp
	CompSwapOp
)

type AtomicElem struct {
	Index     uint64
	Addr      uint64
	Operation int
}

func CobufiAtomicLink(cRecording chan<- AtomicElem) {
	chanRecording = cRecording
	linked = true
}

func CobufiAtomicUnlink() {
	chanRecording = nil
	linked = false
}

//go:nosplit
func CobufiAtomic64Load(addr *uint64) {
	if linked {
		counter += 1
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: LoadOp}
	}
}

//go:nosplit
func CobufiAtomic64Store(addr *uint64) {
	if linked {
		counter += 1
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: StoreOp}
	}
}

//go:nosplit
func CobufiAtomic64Add(addr *uint64) {
	if linked {
		counter += 1
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: AddOp}
	}
}

//go:nosplit
func CobufiAtomic64Swap(addr *uint64) {
	if linked {
		counter += 1
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: SwapOp}
	}
}

//go:nosplit
func CobufiAtomic64CompSwap(addr *uint64) {
	if linked {
		counter += 1
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: CompSwapOp}
	}
}

//go:nosplit
func CobufiAtomic32Load(addr *uint32) {
	if linked {
		counter += 1
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: LoadOp}
	}
}

//go:nosplit
func CobufiAtomic32Store(addr *uint32) {
	if linked {
		counter += 1
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: StoreOp}
	}
}

//go:nosplit
func CobufiAtomic32Add(addr *uint32) {
	if linked {
		counter += 1
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: AddOp}
	}
}

//go:nosplit
func CobufiAtomic32Swap(addr *uint32) {
	if linked {
		counter += 1
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: SwapOp}
	}
}

//go:nosplit
func CobufiAtomic32CompSwap(addr *uint32) {
	if linked {
		counter += 1
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: CompSwapOp}
	}
}

//go:nosplit
func CobufiAtomicUIntPtr(addr *uintptr) {
	if linked {
		counter += 1
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: LoadOp}
	}
}

//go:nosplit
func CobufiAtomicPtr(addr unsafe.Pointer) {
	if linked {
		counter += 1
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(addr)),
			Operation: LoadOp}
	}
}

// COBUFI-FILE-END
