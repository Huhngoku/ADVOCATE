# Patch version of the Go Programming language for Dedego

## Changed files
Added files:

- src/runtime/dedego_routine.go
- src/runtime/dedego_trace.go
- src/runtime/dedego_util.go

Changed files:

- src/cmd/cgo/internal/test/testx.go
- src/runtime/proc.go
- src/runtime/runtime2.go
- src/runtime/chan.go
- src/runtime/select.go
- src/sync/mutex.go
- src/sync/rwmutex.go
- src/sync/waitgroup.go

Disabled Tests

- src/cmd/cgo/internal/test/cgo_test.go
- src/cmd/dist/test.go
- src/cmd/go/stript_test.go
- src/cmd/compile/internal/types2/sizeof_test.go
- src/context/x_test.go
- src/crypto/internal/nistec/nistec_test.go
- src/crypto/tls/tls_test.go
- src/go/build/deps_test.
- src/go/types/sizeof_test.go
- src/internal/intern/inter_test.go
- src/log/slog/text_handler_test.go
- src/net/netip/netip_test.go
- src/runtime/crash_cgo_test.go
- src/runtime/sizeof_test.go
- src/runtime/align_test.go

## Trace structure

- One line per routine
- Each line has the format n:T, where n is the id of the routine and T the trace
- The trace consists of the trace elements separated by semicolon.
- The trace elements can have the following structure:
  - Spawn new routine: G, 'id'
    - 'id' (number): id of the new routine
  - Mutex: M,'id','rw','op','exec','suc'
    - 'id' (number): id of the mutex
    - 'rw' (R/-): R if it is a rwmutex, otherwise -
    - 'op' (L/LR/T/TR/U/UR): L if it is a lock, LR if it is a rlock, T if it is a trylock, TR if it is a rtrylock, U if it is an unlock, UR if it is an runlock
    - 'exec' (e/o): e if the operation was successfully finished, o otherwise
    - 'suc' (s/f): s if the trylock was successful, f otherwise
  - WaitGroup: W,'id','op','exec','delta','val'
    - 'id' (number): id of the mutex
    - 'op' (A/W): A if it is an add or Done, W if it is a wait
    - 'exec' (e/o): e if the operation was successfully finished, o otherwise
    - 'delta' (number): delta of the waitgroup, positive for add, negative for done, 0 for wait
    - 'val' (number): value of the waitgroup after the operation
  - Channel: C,'id','op','exec','pId'
    - 'id' (number): id of the mutex
    - 'op' (S/R/C): S if it is a send, R if it is a receive, C if it is a close
    - 'exec' (e/o): e if the operation was successfully finished, o otherwise
    - 'pId' (number): id of the channel with wich the communication took place
  - Select: S,'id','cases','exec','chosen','opId'
    - 'id' (number): id of the mutex
    - 'cases' (string): cases of the select, id and r/s, separated by '.', d for default
    - 'exec' (e/o): e if the operation was successfully finished, o otherwise
    - 'chosen' (number): index of the chosen case in cases (0 indexed, -1 for default)
    - 'opId' (number): id of the operation on the channel