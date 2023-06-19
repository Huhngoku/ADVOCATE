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

at the beginning of the main function. 

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

## Structure of trace T in EBNF

```ebnf
T := L\nT | ""                                                  (trace)
L := "" | {E";"}E                                               (routine local trace)
E := G | M | W | C | S                                          (trace element)
G := "G,"tpre","id                                              (element for creation of new routine)
M := "M,"tpre","tpost","id","rw","opM","exec","suc","pos        (element for operation on sync (rw)mutex)
W := "W,"tpre","tpost",id","opW","exec","delta","val","pos      (element for operation on sync wait group)
C := "C,"tpre","tpost","id","opC","exec","oId","pos             (element for operation on channel)
S := "S,"tpre","tpost","id","cases","exec","chosen","oId","pos  (element for select)
tpre := ℕ                                                       (timer when the operation is started)
tpost := ℕ                                                      (timer when the operation has finished)
id := ℕ                                                         (unique id of the underling object)
rw := "R" | "-"                                                 ("R" if the mutex is an RW mutex, "-" otherwise)
opM := "L" | "LR" | "T" | "TR" | "U" | "UR"                     (operation on the mutex, L: lock, LR: rLock, T: tryLock, TR: tryRLock, U: unlock, UR: rUnlock)
opW := "A" | "W"                                                (operation on the wait group, A: add (delta > 0) or done (delta < 0), W: wait)
opC := "S" | "R" | "C"                                          (operation on the channel, S: send, R: receive, C: close)
exec := "e" | "f"                                               (e: the operation was fully executed, o: the operation was not fully executed, e.g. a mutex was still waiting at a lock operation when the program was terminated or a channel never found an communication partner)
suc := "s" | "f"                                                (the mutex lock was successful ("s") or it failed ("f", only possible for try(r)lock))
pos := file":"line                                              (position in the code, where the operation was executed)
file := 𝕊                                                       (file path of pos)
line := ℕ                                                       (line number of pos)
delta := ℕ                                                      (change of the internal counter of wait group, normally +1 for add, -1 for done)
val := ℕ                                                        (internal counter of the wait group after the operation)
oId := ℕ                                                        (identifier for an communication on the channel, the send and receive (or select) that have communicated share the same oId)
cases := case | {case"."}case                                   (list of cases in select, seperated by .)
case := cId""("r" | "s") | "d"                                  (case in select, consisting of channel id and "r" for receive or "s" for send. "d" shows an existing default case)  
cId := ℕ                                                        (id of channel in select case)
chosen := ℕ0 | "-1"                                             (index of the chosen case in cases, -1 for default case)    
```

Info: 
- \n: newline
- ℕ: natural number not including 0
- ℕ0: natural number including 0
- 𝕊: string containing 1 or more characters
- The tracer contains a global timer for all routines that is incremented every time an timer element (tpre/tpost) is recorded.