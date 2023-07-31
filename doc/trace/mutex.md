# Mutex
The locking and unlocking of sync.(rw)-mutexes is recorded in the trace where it occurred.

## Trace element
The basic form of the trace element is 
```
M,[tpre],[tpost],[id],[rw],[opM],[exec],[suc],[pos]
```
where `M` identifies the element as a mutex element.
The other fields are set as follows:
- [tpre]: This is the value of the global counter when the routine starts 
the execution of the lock or unlock function
- [tpost]: This is the value of the global counter when the routine has finished 
its operation. For lock operations this can be either if the lock was successfully acquired or if the routines continues its execution without 
acquiring the lock in case of a trylock. 
- [id]: This is the unique id identifying this mutex
- [rw]: This field records, whether the mutex is an rw mutex ([rw] = `R`) or not
([rw] = `-`)
- [opM]: This field shows the operation of the element those can be
  - [opM] = `L`: Lock
  - [opM] = `LR`: RLock
  - [opM] = `T`: TryLock
  - [opM] = `TR`: TryRLock
  - [opM] = `U`: Unlock
  - [opM] = `UR`: RUnlock
- [exec]: This field shows, whether the operation was finished ([exec] = `e`) or
not ([exec] = `f`). Failed can e.g. mean, that the routine still tried to
acquire a lock, when the program was terminated. It is also set to `f`, if the 
program panicked during the execution of the operation, e.g. because an 
RUnlock operation was called on a mutex that is currently held by and Lock 
operation.
- [suc]: This field is used to determine if an TryLock was successful ([suc] = `s`)
or not ([suc] = `f`) in acquiring the mutex. For all other operation it is always
set to `s`.
- [pos]: The last field show the position in the code, where the mutex operation 
was executed. It consists of the file and line number separated by a colon (:)

## Example
 The following is an  example containing all the different recorded 
elements:
```go
package main

import (
    "sync"
)

func main() {  // routine 1
    var m sync.Mutex  // id = 1
    var n sync.RWMutex  // id = 2
    m.Lock()
    ...
    m.Unlock()
    go func() {  // routine 2
        suc := m.TryLock()
        if suc {
            ...
            m.Unlock()
        }
    }
    go func() {  // routine 3
        n.Lock()
        ...
        n.Unlock()
    }
    go func() {  // routine 4
        suc := n.TryLock()
        if suc {
            ...
            n.Unlock()
        }
    }
    go func() {  // routine 5
        n.RLock()
        ...
        n.RUnlock()
    }
    go func() {  // routine 6
        suc := n.TryRLock()
        if suc {
            ...
            n.RUnlock()
        }
    }
}
```
The different routines show different operation pairs of locks and unlocks on mutex and rwmutex.
- Routine 1: Lock and Unlock of a mutex
- Routine 2: TryLock and Unlock of a mutex
- Routine 3: Lock and Unlock of a rwmutex
- Routine 4: TryLock and Unlock of a rwmutex
- Routine 5: RLock and RUnlock of a rwmutex
- Routine 6: TryRLock and RUnlock of a rwmutex

If we ignore all the internal operations, assume that all operations are executed
before the program terminates and assume, that all try operations are successful, 
we get the following trace. For simplicity we also assume, that the routine 
are executed consecutively, ignoring that the routines are normally
run concurrent (this only effects the time stamps). We also ignore the elements showing the creation of the new go routines. These elements would all be at the end of the trace of the first routine (first line). This would also 
shift the time steps.
```
M,1,2,1,-,L,e,s,example_file.go:8;M,3,4,1,-,U,e,s,example_file.go:10;
M,5,6,1,-,T,e,s,example_file.go:13;M,7,8,1,-,U,e,s,example_file.go:16;
M,9,10,2,R,L,e,s,example_file.go:20;M,11,12,2,R,U,e,s,example_file.go:22;
M,13,14,2,R,T,e,s,example_file.go:25;M,15,16,2,R,U,e,s,example_file.go:28;
M,17,18,2,R,LR,e,s,example_file.go:32;M,19,20,2,R,UR,e,s,example_file.go:34;
M,21,22,2,R,TR,e,s,example_file.go:37;M,23,24,2,R,UR,e,s,example_file.go:40;
```

## Implementation
The recording of the mutex operations is implemented in the `go-patch/src/sync/mutex.go` and `go-patch/src/sync/rwmutex.go` files in the implementation of the 
Lock, RLock, TryLock, TryRLock, Unlock and RUnlock function.\
It consist of 
two function calls, one at the beginning and one at the end of each function.
The first function call is called before the Operation tries to executed 
and records the id ([id]) and type ([rw]) of the involved mutex, the called operation (opM), the position of the operation in the program ([pos]) and the counter at the beginning of the operation ([tpre]).\
The second function call records the success of the operation. This includes 
the counter at the end of the operation ([tpost]), the information that the 
operation finished ([exec]) and the success of try lock operations ([suc]).\
The implementation of those function calls can be found in 
`go-patch/src/runtime/dedego_trace.go` in the functions `DedegoMutexLockPre`, 
`DedegoMutexLockTry`, `DedegoUnlockPre`, `DedegoPost` and `DedegoPostTry`.