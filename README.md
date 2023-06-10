# Dynamic Analysis of Message Passing Go Programs

Warning: The program is currently being completely revised and is not usable at the moment.


## What
Program to run a dynamic analysis of concurrent Go programs to detect 
possible deadlock situations.

### Mutexes
For mutexes cyclic deadlocks as well as deadlocks by double locking can be detected.

Cyclic Deadlocks are the result of cyclicly blocking routines.
The following program shows an example:
```go
func main() {
	x := sync.Mutex{}
	y := sync.Mutex{}

	go func() {
		x.Lock()  // 1
		y.Lock()  // 2
		y.Unlock()
		x.Unlock()
	}()

	y.Lock()  // 3
	x.Lock()  // 4
	x.Unlock()
	y.Lock()
}
```
If (1) and (3) run simultaneously and before (2) or (4) run, 
both routines block on (2) and (4) causing a cyclic deadlock.

Double locking arise if a mutex is locked multiple times by the same routine without unlocking. The following program shows an example:
```go
func main() {
	x := sync.Mutex{}

	x.Lock() 
	x.Lock()
}
```
In this case the routine blocks it self, which leads to a deadlock.

The program is able to differentiate between mutexes and rw-mutexes. The following example therefore does not lead to any problem because RLock operations do not block each other:
```go
func main() {
	x := sync.RWMutex{}
	y := sync.RWMutex{}

	go func() {
		x.RLock()  // 1
		y.Lock()  // 2
		y.Unlock()
		x.Unlock()
	}()

	y.Lock()  // 3
	x.RLock()  // 4
	x.Unlock()
	y.Lock()
}
```
The program is able to detect problems like these.

### Channels
Channels can also lead to blocking situations. Let's use the 
following program as an example:
```go
func main() {
	x := make(chan int)

	go func() {
		x <- 1  // 1
		<-x     // 2
	}()

	go func() {
		x <- 1  // 3
	}()

	<-x         // 4
	time.Sleep(time.Second)
}
```
If (1) communicates with (4) and (3) with (2) everything is fine. But if (3) communicates with (4) (1) has no valid communication partner and will therefore block the routine forever. The program is able to detect situations like these. 
It can to a certain extend also detect blocking problems 
with buffered channels. 

To detect problems caused or hidden by select statements, the program is analyzed multiple times with different preferred select cases in the different runs. 

