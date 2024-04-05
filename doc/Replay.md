# Replay
We want to be able to replay a program as it was recorded in a trace file. 

## Usage

To start the replay, add the following header at the beginning of the
main function:

```go
advocate.EnableReplay()
defer advocate.WaitForReplayFinish()
```

Include the following imports:
```go
"advocate"
```

The trace files must be in the `/rewritten_trace` or `/trace` folder. If both folder exist, `/rewritten_trace` is used.

Now the program can be run with the modified go routine, identical to the recording of the trace (remember to export the new gopath).

## Implementation
The following is a description of the current implementation of the trace replay.
It is split into three parts:

- Trace Reading
- Order Enforcement
- State Enforcement

### Trace Reading
First we read in the trace and create a new internal data structure to save
the trace, ordered by `tpre`. For now we ignore atomic events, because they do not change the flow of the program and the implementation of the order
enforcement would be challenging. For each element we store

- the operation
- the tpre
- the file in the program where it occurred
- the line in the program where it occurred
- whether the operation was completely executed (tpost not 0)

For TryLock operations and Once we also store

- whether the operation was successful

For channel and select operations we find the communication partner and store

- the file in the program where the partner operation occurred
- the line in the program where the partner operation occurred


### Order Enforcement
Order enforcement makes sure, that the elements that are recorded in the trace
are run in the correct global order.

For the most operations we use the file and line number to connect an operation
in the trace with an operation in the program code that is to be replayed. The
only operations which cannot use this yet are atomic operations.

If a traced operation in the replaying trace starts, it calls the following function
```go
func WaitForReplay(op ReplayOperation, skip int) (bool, chan ReplayElement) {
	if !replayEnabled {
		return false, nil
	}

	_, file, line, _ := Caller(skip)
	c := make(chan ReplayElement, 1<<16)

	go func() {
		for {
			next := getNextReplayElement()
			if (next.Op != op && !correctSelect(next.Op, op)) ||
				next.File != file || next.Line != line {
				// TODO: sleep here to not waste CPU
				continue
			}
			c <- next
			lock(&replayLock)
			replayIndex++
			unlock(&replayLock)
			return
		}
	}()

	return true, c
}
```

This function will start a new go routine and then return a channel. The channel
is modified in such a way, that it is ignored for replay and tracing purposes (size set to 1<<16).

The only problem with that is, that we would record the creation of new go
routine, if we would record a new trace while we run the replay. I hope, that I
will be able to fix that in the future.

The begin of each operation is modified to call this WaitForReplay and then
try to receive at the returned channel. A simplified version of this looks like
```go
var replayElem ReplayElement
if !c.advocateIgnore {  // c is not an channel from the replay mechanism
    enabled, waitChan := WaitForReplay(AdvocateReplayChannelSend, 3)
    if enabled {
        replayElem = <-waitChan
    }
}
```
WaitForReplay will no continually check, what the next operation to be executed
is and, when the operation, file and line are identical, send on the channel,
to start the execution of the operation.

## State enforcement
This second part makes sure, that the state of the program is equal to the state
in the recorded trace. This includes

- blocking blocked operation
- making sure, that successful operations are successful and unsuccessful once are not
- making sure, that channel partners are correct
- making sure, that select cases are correct.

Many of those should already be enforced automatically because of the order enforcement, bu we implement additional safeguards to make sure, that a shift in
not recorded operations does not allow those operations to change there behavior.

### Blocking blocked operations
Operations that did not execute in the recorded file, e.g. because a mutex was
still blocked at the end or a channel never found a partner, are not supposed to
be executed during replay. A simplified version of this looks as follows:
```go
if enabled {  // replay is running
    ...
    if replayElem.Blocked {
        BlockForever()
    }
}
```
It is included in the `if enabled` section from the
The `BlockForever` function will block the execution of the operation and the
routine where it is contained, until the program terminated. The if block can
also contain additional operations that are necessary, to get the same trace
outcome as in the recorded trace. For channel send this would e.g. look like
```go
if enabled {  // replay is running
    ...
    if replayElem.Blocked {
        lock(&c.numberSendMutex)
        c.numberSend++
        unlock(&c.numberSendMutex)
        _ = AdvocateChanSendPre(c.id, c.numberSend, c.dataqsiz)
        BlockForever()
    }
}
```

<!--
### Making sure, that successful operations are successful and unsuccessful once are not
This is only relevant for Try(R)Lock operations and Once. For operations that
were not successful we will, after the necessary steps to record the operation
are taken, just return. This is equal to an unsuccessful operations. For an once.Do
this would look as follows:
```go
if envable {  // replay is running
    ...
    if !elem.Suc {
        if o.id == 0 {
            o.id = runtime.GetAdvocateObjectId()
        }
        index := runtime.AdvocateOncePre(o.id)
        runtime.AdvocateOncePost(index, false)
        return
    }
}
```
Operations that were successful in the trace can now not be blocked by incorrectly
executed unsuccessful operations and therefor do not need any additional
change.
-->

### Making sure, that channel partners are correct
As described in Trace Reading, we direly store in each trace element of the
line and file of the partner operation. When go tries to send or receive
an element on a channel, it will create a `*sudog` object, that is then passed
to a partner. This element is extended to also store the file and line of the
partner, where this message should arrive. When checking if a communication
partner is available, i.e. if a `*sudog` object is available, we now check if
the position information of a possible communication is identical to the information
of the calling operation. If it is identical, the communication is identical to
the one in the trace, and it can continue as normal. If the information is not correct, the communication is treated, as if no communication partner was available.
This is implemented in the `dequeue` function as
```go
if replayEnabled && !sgp.replayEnabled {  // replay is enabled and the channel is not part of the replay mechanism
    if !(rElem.File == "") && !sgp.c.advocateIgnore {
        if sgp.pFile != rElem.File || sgp.pLine != rElem.Line {
            return nil  // reject communication
        }
    }
}
```

### Making sure, that select cases are correct.
We must make sure, that the correct case in a select is executed.
If the default case is supposed to be executed, we can immediately force the
execution of the default case, before the select checks, if other channels
would be available by adding
```go
if enabled && replayElem.Op == AdvocateReplaySelectDefault {
    selunlock(scases, lockorder)
    casi = -1
    goto retc
}
```
before the check which channel could be executed. This will imminently execute
the default case.

The same check for the channel communication partners as described in `Making sure, that channel partners are correct` will force select cases, to find the actually executed channel pair before being able to execute. This will stop incorrect cases to execute.
If a select contains the same case twice, i.e.
```go
select {
    case <-c:
        ...
    case <-c:
        ...
}
```
this will still select one of this cases by random. To make sure, that those select statements will also be replayed correctly, we use the internal index `casi` for the cases, used in the implementation of the select statement. This case is not identical to the ordering of the select cases but is still deterministic. For this reason it is possible to use this index as an identifier for a specific case. From this, when the select determines, if a select case is usable, we reject every case, for which the index is not correct.

## Disable
If the (rewritten) trace contains a stop signal `X`, the replay is disabled.
This means, the the program is from this moment on allowed, to run freely 
without interference.

## Timeout

For multiple reasons, including 
- nondeterministic execution paths, e.g. if the execution path depends on random numbers
- the program execution path depends on the order of not tracked operations
- the program execution pats depends on outside input, that has not been reproduced exactly
- the program was altered between recording and replay
- the trace/replay mechanism can contain bugs (hopefully not)

the replay can get stuck. In this case, the replay is waiting at an event, 
that is not in trace, or is waiting for a trace element, that is not part of 
the current execution path. 

One example would be the following program
```go
package main

import( 
    "math/rand"
)

func main() {
	c := make(chan int)
	d := make(chan int)

    rand.Seed(time.Now().UnixNano())
    r := rand.Float64()

	if r >= 0.5 {
		close(c)
	}

	close(d)
}

```

### Element is not in trace
If in the recording run `r < 0.5`, meaning `close(c)` was not
executed, but in the replay run `r > 0.5`, meaning the program will try to
execute `close(c)`, the program will get stuck, because `close(c)` is
not in the trace.

To detect this, we record the positions of all possible events in the trace in a
`map[string][]int` (file -> []line) when reading the trace.
If a waiting element has not been executed after a certain time (approx 5s),
it is checked, if the position of the waiting element is somewhere in the
trace using the map. If it is not, the program determines, that a different path than in the recording is being executed and the program panics.

### Additional elements in trace
If we assume the opposite case, that `r > 0.5` during the recording, meaning
`close(c)` was executed and recorded, but in the replay it is not, the
replay will wait at `close(d)` indefinitely, because it assumes, that `close(c)`
has to be executed before. To detect such cases, we approximate how long
a operation is already waiting. If this wait time exceeds approx. 10s, the
program will print a warning to the user. This warning will be repeated every
10s. In this case the user can decide, whether they want to cancel the program
execution, or continue letting it run, hoping the execution will resume.

### Additional time after finishing main routine
In a normal program run, the program is terminated as soon as the main 
routine terminates, even if other routines are still running. To prevent recorded 
trace elements in other routines from not being replayed, the replay stops
the main routine from terminating before all trace elements have been replayed.
If the wait time after the main routine has terminated exceeds approx. 10s, 
the program will also print a warning. In this case it is again the decision of
the user, whether they want to terminate the program or continue its execution, 
hoping it will resume.

### Measuring time
Unfortunately it is not possible to use the time package in the replay (leads
to cyclic import). For this reason, the elapsed time must be approximated
in a different way. A waiting operation uses a loop in a go routine, to check
periodically whether it is the next element to be played. After the check, it
waits for a short time by running an empty counting loop. The waiting time
is approximated by counting the number of checks (approx. 50 checks per second).
This means, that the actual waiting time is dependent on the execution speed
of the computer. Because the length of the waiting time is not really important,
this works, but it is not an accurate measure of time.