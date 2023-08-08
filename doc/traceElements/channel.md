# Channel
The sending and receiving on and the closing of channels is recorded in the 
trace of the routine where they occur.

## Trace element
The basic form of the trace element is 
```
C,[tpre],[tpost],[id],[opC],[exec],[oId],[qSize],[qCountPre],[qCoundPost],[pos] 
```
where `C` identifies the element as a channel element. The other fields are 
set as follows:
- [tpre] $\in \mathbb N$ : This is the value of the global counter when the routine starts 
the execution of the operation
- [tpost]$\in \mathbb N$: This is the value of the global counter when the channel has finished its operation. For close we get [tpost] = [tpre]
- [id]$\in \mathbb N$: This shows the unique id of the channel
- [opC]: This field shows the operation that was executed:
    - [opC] = `S`: send
    - [opC] = `R`: receive
    - [opC] = `C`: close
- [exec]: This field shows, whether the operation was finished ([exec] = `t`) or
not ([exec] = `f`). Failed can e.g. mean, that the routine still tried to send or receive when the program was terminated
- [oId] $\in \mathbb N$: This field shows the communication id. This can be used to connect corresponding communications. If a send and a receive on the same channel (same channel id) have the same [oId], a message was send from the send to the receive. For close this is always `0`
- [qSize] $\in \mathbb N_0$: This is the size of the channel. For unbuffered channels this is `0`.
- [qCountPre] $\in \mathbb N_0$: This is the amount of elements 
in the queue of the channel before the operation was executed.
- [qCountPost] $\in \mathbb N_0$: This is the amount of elements in the queue of the channel after the operation was executed. For close it is always [qCountPre] = [qCountPost]
- [pos]: The last field show the position in the code, where the mutex operation 
was executed. It consists of the file and line number separated by a colon (:)
## Example
The following is an example for a program with different channel operations:
```go
package main
func main() {    // routine 1
    c := make(chan int, 2)  // id = 4
	d := make(chan int, 0)  // id = 5

	go func() {  // routine 2
		c <- 1   // line 7
		c <- 2   // line 8
		d <- 1   // line 9
	}()

	<-d          // line 12
	<-c          // line 13
	<-c          // line 14

	close(d)     // line 16
}
```
If we ignore all internal operations, we would get the following trace:
```
G,1,2;C,8,9,5,R,t,1,0,0,0,example_file.go:12;C,10,11,4,R,t,1,2,2,1,example_file.go:13;C,12,13,4,R,t,2,2,1,0,example_file.go:14;C,14,14,5,C,t,0,0,0,0,example_file.go:16
C,2,3,4,S,t,1,2,0,1,example_file.go:7;C,4,5,4,S,t,2,2,1,2,example_file.go:8;C,6,7,5,S,e,1,0,0,0,example_file.go:9
```

## Implementation
The recording of the channel operations is done in the 
`go-patch/src/runtime/chan.go` file in the `chansend`, `chanrecv` and `closechan` function. Additionally the 
`hchan` struct in the same file is ammended by the following fields:
- `id`: identifier for the channel
- `numberSend`: number of completed send operations
- `numberRecv`: number of completed reveive operations

`numberSend` and `numberRecv` are later set as `oId` in the corresponding trace elements. The send operations are implemented as a FIFO-queue. We can therefore count the number of elements added to the queue and removed from the
queue, to determine, which send and receive operations are
communication partners. Because of mutexes, that are already present in the original channel implementation,
it is not possible to mix up these numbers.\
For the send and receive operations three record functions are added. The first one (`DedegoChanSendPre`/`DedegoChanRecvPre`) at the beginning of the operation, which records [tpre], [id], [opC], [qSize], [qCountPre] and [pos].\
The other two functions are called at the end of the
operation, after the send or receive was fully executed.
These functions record [qCountPost] (`DedegoChanPostQCount`)
as well as [tpost] and [exec] (`DedegoChanPost`).\
A close on a channel cannot block, it only needs one recording function. This function (`DedegoChanClose`) records all needed values. For [tpre] and [tpost] the same 
value is set. The same is true for [qCountPre] and [qCountPost].