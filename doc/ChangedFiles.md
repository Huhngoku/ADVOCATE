# Changed files

The following is a list of all files in the Go runtime that have been
added or changed.

Added files:

- src/runtime/advocate_routine.go
- src/runtime/advocate_trace.go
- src/runtime/advocate_util.go
- src/runtime/advocate_replay.go
- src/runtime/internal/atomic/advocate_atomic.go
- src/advocate/advocate.go

Changed files (marked with ADVOCATE-CHANGE, außer in .s):

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
- src/sync/cond.go
- src/sync/pool.go
- src/internal/poll/fs_poll_runtime.go
- cmd/compile/internal/ssagen/ssa.go


Additionally some test files have been altered.