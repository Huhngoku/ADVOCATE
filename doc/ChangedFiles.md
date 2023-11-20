# Changed files

The following is a list of all files in the Go runtime that have been 
added or changed.

Added files:

- src/runtime/cobufi_routine.go
- src/runtime/cobufi_trace.go
- src/runtime/cobufi_util.go
- src/runtime/cobufi_replay.go
- src/runtime/internal/atomic/cobufi_atomic.go
- src/cobufi/cobufi.go
- 

Changed files (marked with COBUFI-CHANGE, außer in .s):

- src/runtime/proc.go
- src/runtime/runtime2.go
- src/runtime/chan.go
- src/runtime/select.go
- src/runtime/internal/atomic/doc.go
- src/runtime/internal/atomic/atomic_amd64.go
- src/runtime/internal/atomic/atomic_amd64.s
- src/runtime/internal/atomic/atomic_arm64.go
- src/runtime/internal/atomic/atomic_arm64.s
- src/sync/mutex.go
- src/sync/rwmutex.go
- src/sync/waitgroup.go
- src/sync/once.go
- cmd/compile/internal/ssagen