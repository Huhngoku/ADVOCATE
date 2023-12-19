package atomic

// ADVOCATE-FILE-START

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
	// ChanReturn chan bool
}

func AdvocateAtomicLink(cRecording chan<- AtomicElem) {
	chanRecording = cRecording
	linked = true
}

func AdvocateAtomicUnlink() {
	chanRecording = nil
	linked = false
}

func AdvocateAtomic64Load(addr *uint64) {
	if linked {
		counter++
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: LoadOp}
	}
}

func AdvocateAtomic64Store(addr *uint64) {
	if linked {
		counter++
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: StoreOp}
	}
}

func AdvocateAtomic64Add(addr *uint64) {
	if linked {
		counter++
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: AddOp}
	}
}

func AdvocateAtomic64Swap(addr *uint64) {
	if linked {
		counter++
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: SwapOp}
	}
}

func AdvocateAtomic64CompSwap(addr *uint64) {
	if linked {
		counter++
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: CompSwapOp}
	}
}

func AdvocateAtomic32Load(addr *uint32) {
	if linked {
		counter++
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: LoadOp}
	}
}

func AdvocateAtomic32Store(addr *uint32) {
	if linked {
		counter++
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: StoreOp}
	}
}

func AdvocateAtomic32Add(addr *uint32) {
	if linked {
		counter++
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: AddOp}
	}
}

func AdvocateAtomic32Swap(addr *uint32) {
	if linked {
		counter++
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: SwapOp}
	}
}

func AdvocateAtomic32CompSwap(addr *uint32) {
	if linked {
		counter++
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: CompSwapOp}
	}
}

func AdvocateAtomicUIntPtr(addr *uintptr) {
	if linked {
		counter++
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(unsafe.Pointer(addr))),
			Operation: LoadOp}
	}
}

func AdvocateAtomicPtr(addr unsafe.Pointer) {
	if linked {
		counter++
		chanRecording <- AtomicElem{Index: counter, Addr: uint64(uintptr(addr)),
			Operation: LoadOp}
	}
}

// ADVOCATE-FILE-END
