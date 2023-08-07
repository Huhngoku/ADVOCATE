package atomic

// DEDEGO-FILE-START

import (
	"unsafe"
)

var com chan<- AtomicElem
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

func DedegoAtomicLink(c chan<- AtomicElem) {
	com = c
	linked = true
}

func DedegoAtomicUnlink() {
	com = nil
	linked = false
}

//go:nosplit
func DedegoAtomic64Load(addr *uint64) {
	if linked {
		counter += 1
		com <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: LoadOp}
	}
}

//go:nosplit
func DedegoAtomic64Store(addr *uint64) {
	if linked {
		counter += 1
		com <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: StoreOp}
	}
}

//go:nosplit
func DedegoAtomic64Add(addr *uint64) {
	if linked {
		counter += 1
		com <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: AddOp}
	}
}

//go:nosplit
func DedegoAtomic64Swap(addr *uint64) {
	if linked {
		counter += 1
		com <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: SwapOp}
	}
}

//go:nosplit
func DedegoAtomic64CompSwap(addr *uint64) {
	if linked {
		counter += 1
		com <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: CompSwapOp}
	}
}

//go:nosplit
func DedegoAtomic32Load(addr *uint32) {
	if linked {
		counter += 1
		com <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: LoadOp}
	}
}

//go:nosplit
func DedegoAtomic32Store(addr *uint32) {
	if linked {
		counter += 1
		com <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: StoreOp}
	}
}

//go:nosplit
func DedegoAtomic32Add(addr *uint32) {
	if linked {
		counter += 1
		com <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: AddOp}
	}
}

//go:nosplit
func DedegoAtomic32Swap(addr *uint32) {
	if linked {
		counter += 1
		com <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: SwapOp}
	}
}

//go:nosplit
func DedegoAtomic32CompSwap(addr *uint32) {
	if linked {
		counter += 1
		com <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: CompSwapOp}
	}
}

//go:nosplit
func DedegoAtomicUIntPtr(addr *uintptr) {
	if linked {
		counter += 1
		com <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: LoadOp}
	}
}

//go:nosplit
func DedegoAtomicPtr(addr unsafe.Pointer) {
	if linked {
		counter += 1
		com <- AtomicElem{Index: counter, Addr: uint64(uintptr(addr)),
			Operation: LoadOp}
	}
}

// DEDEGO-FILE-END
