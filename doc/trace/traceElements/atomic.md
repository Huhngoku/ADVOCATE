# Atomic
The recording of atomics records atomic operations, both on normal types and on atomic types. This includes Add, CompareAndSwap, Swap, Load and Store operations.

## Info:
The recording of atomic events is currently only implemented and tested for `amd64`.
For `arm64` an untested implementation exists, but there are no guaranties, that
this implementation is runnable.

## Trace element:
The basic form of the trace element is 
```
A,[tpost],[id],[opA]
```
where `A` identifies the element as an atomic operation.
The other fields are set as follows:
- [tpost]: This field shows the value of the internal counter when the operation is executed.
- [id]: This field shows a number representing the variable. It is not possible to give every variable its own unique, consecutive id. For this reason, this id is equal to the memory position of the variable.
- [opA]: This field shows the type of operation. Those can be 
	- `L`: Load
	- `S`: Store
	- `A`: Add
	- `W`: Swap
	- `C`: CompareAndSwap
	- `U`: unknown (should not appear)

Each atomic event is always followed by a channel send event. (see implementation)

## Example
The following is an example containing atomic operations.
```go
package main

func main() {  // routine 1
    var a atomic.Int32
	var b int32

	a.Add(1)
	atomic.StoreInt32(&b, a.Load())
}
```
For the trace we ignore all internal operations. Because of the implementation of the recording, each atomic operation trace element is always followed by an channel operation. These channel operations and there communication partner have been left in the trace.
```txt
G,1,2;A,2,824634851400;C,3,4,3,S,t,1,0,0,0,.../go-patch/src/runtime/internal/atomic/dedego_atomic.go:40;A,7,824634851400;C,8,9,3,S,t,2,0,0,0,.../go-patch/src/runtime/internal/atomic/dedego_atomic.go:40;A,12,824634851404;C,13,14,3,S,t,3,0,0,0,.../go-patch/src/runtime/internal/atomic/dedego_atomic.go:40
C,5,6,3,R,t,1,0,0,0,.../go-patch/src/runtime/dedego_trace.go:196;C,10,11,3,R,t,2,0,0,0,.../go-patch/src/runtime/dedego_trace.go:196;C,15,16,3,R,t,3,0,0,0,.../go-patch/src/runtime/dedego_trace.go:196;C,17,18,3,R,f,4,0,0,0,.../go-patch/src/runtime/dedego_trace.go:196
```

## Implementation
Most of the atomic operations are directly implemented in assembly. The functions that record the atomic elements are added in `go-patch/src/runtime/internal/atomic/atomic_amd64.go`, `go-patch/src/runtime/internal/atomic/atomic_amd64.s`. The used functions are implemented in `go-patch/src/runtime/internal/atomic/dedego_atomic.go` and `go-patch/src/runtime/chan.go`. It was also necessary to delete alias definitions in `go-patch/src/cmd/compile/internal/ssagen/ssa.go`.

It is not possible to directly get the information about the recording of an atomic event into the trace, because of cycling imports.\
For this reason, a background routine is used. This routine is started by the `runtime.InitAtomics` functions, that must be added in the header. This routine constantly read on a channel. The reading is only necessary to empty the channel and has no effect on the actual recording of atomics.

If an atomic event is recorded, the recording function sends a message on this cannel, containing an index, the memory address of the involved variable and the 
type of operation. The index is later used to connect the trace element with the memory address and type.\
The actual recording of the atomic events is done in the `DedegoChanSendPre` function, that also records the normal pre-send on channels. 
Because this function is called every time a message is send, but before the routine actually tries to send, the information about the atomic event is prevented from being held back by delays in the channel. By using this channel method it is also possible to determine, in which routine the atomic operation took place (in the same routine from which the channel send). If the pre-send function detects a channel operation, that started in the `go-patch/src/runtime/internal/atomic/dedego_atomic.go` file, the info about the atomic operation is added to the trace.\
Because of this method each atomic trace element is always followed by a channel send element. 
