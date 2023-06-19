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

Disabled Tests (files contain disabled tests): 

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
- src/runtime/metrics_test.go
- src/net/tcpsock_test.go
- src/reflect/all_test.go

## Save the trace

Add

```go
import (
  "runtime",
  "io/ioutil"
  "os"
)

runtime.DedegoInit("path/to/project/root")
defer func() {
  file_name := "dedego.log"
  output := runtime.AllTracesToString()
  err := ioutil.WriteFile(file_name, []byte(output), os.ModePerm)
  if err != nil {
    panic(err)
  }
}()
```

to the beginning of the main function. 

For programs with many recorded 
operations this can lead to memory problems. In this case use

```go
runtime.DedegoInit("path/to/project/root")
defer func() {
file_name := "dedego.log"
file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err != nil {
  panic(err)
}
runtime.DisableTrace()
numRout := runtime.GetNumberOfRoutines()
for i := 0; i < numRout; i++ {
  c := make(chan string)
  go func() {
    runtime.TraceToStringByIdChannel(i, c)
    close(c)
  }()
  for trace := range c {
    if _, err := file.WriteString(trace); err != nil {
      panic(err)
    }
  }
  if _, err := file.WriteString("\n"); err != nil {
    panic(err)
  }
}
file.Close()
}()

```

instead.

## Trace structure

- One line per routine
- Each line contains the trace of one routine. The line number is equal to the routine id
- The trace consists of the trace elements separated by semicolon.
- The trace elements can have the following structure:
  - Spawn new routine: G,'t','id'
    - 't' (number): global timer when the trace was created
    - 'id' (number): id of the new routine
  - Mutex: M,'tpre','tpost','id','rw','op','exec','suc','file':'line'
    - 'tpre' (number): global timer when the operation starts
    - 'tpost' (number): global timer when the operation ends
    - 'id' (number): id of the mutex
    - 'rw' (R/-): R if it is a rwmutex, otherwise -
    - 'op' (L/LR/T/TR/U/UR): L if it is a lock, LR if it is a rlock, T if it is a trylock, TR if it is a rtrylock, U if it is an unlock, UR if it is an runlock
    - 'exec' (e/o): e if the operation was successfully finished, o otherwise
    - 'suc' (s/f): s if the trylock was successful, f otherwise
    - 'pos' (string): position where the operation was called (file:line)
  - WaitGroup: W,'tpre','tpost','id','op','exec','delta','val','file':'line'
    - 'tpre' (number): global timer when the operation starts
    - 'tpost' (number): global timer when the operation ends
    - 'id' (number): id of the mutex
    - 'op' (A/W): A if it is an add or Done, W if it is a wait
    - 'exec' (e/o): e if the operation was successfully finished, o otherwise
    - 'delta' (number): delta of the waitgroup, positive for add, negative for done, 0 for wait
    - 'val' (number): value of the waitgroup after the operation
    - 'pos' (string): position where the operation was called (file:line)
  - Channel: C,'tpre','tpost','id','op','exec','oId','file':'line'
    - 'tpre' (number): global timer when the operation starts
    - 'tpost' (number): global timer when the operation ends
    - 'id' (number): id of the mutex
    - 'op' (S/R/C): S if it is a send, R if it is a receive, C if it is a close
    - 'exec' (e/o): e if the operation was successfully finished, o otherwise
    - 'oId' (number): id of the operation
    - 'pos' (string): position where the operation was called (file:line)
  - Select: S,'tpre','tpost','id','cases','exec','chosen','opId','file':'line
    - 'tpre' (number): global timer when the operation starts
    - 'tpost' (number): global timer when the operation ends
    - 'id' (number): id of the mutex
    - 'cases' (string): cases of the select, id and r/s, separated by '.', d for default
    - 'exec' (e/o): e if the operation was successfully finished, o otherwise
    - 'chosen' (number): index of the chosen case in cases (0 indexed, -1 for default)
    - 'opId' (number): id of the operation on the channel
    - 'pos' (string): position where the operation was called (file:line)