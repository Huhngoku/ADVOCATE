
# Trace replay

## Problem statement

Given a trace T. Force execution of the program in such a way that we "follow closely" the sequence of events
in trace T.


## Related works

"Fuzzing methods" are somewhat related but the use a randomized method to explore different schedules.
We seek to follow a specific schedule as specified by a program trace.


[Who Goes First? Detecting Go Concurrency Bugs via Message Reordering](https://songlh.github.io/paper/gfuzz.pdf)


[Deadlockfuzzer](https://github.com/ksen007/calfuzzer)


## Replaying a trace

Highlights:

We must enforce the total order of events as found in the trace.

For atomics this is a challenage as there is no easy way to suspend atomic operations (unless we re-implement atomics).

Below is a discussion of the various cases to consider.



#### atomics


Ignored.

But this means that some write-read dependencies might be violated
and the control flow chantges (e.g. entering some other parts of the program code).
This then leads to a "stuck" trace replay.

#### goroutines


The spawn of new go routines should now be executed in the same order as recorded in the trace (should work but more testing needed).


#### Unbuffered channel

~~~~
    T1        T2      T3      T4

1.  snd_A
2.            rcv_A
3.                    snd_B
4.			                  rcv_B
~~~~~~~~~~

Q: Do we need to obey the total order at all?

Claim: Seems we can process the events in any order as we long we check for matching communication partners.


#### Buffered channel

Buffer space 2.

Program run resulting in the following trace.
Point to note, channels do NOT necessarily behave like a queue.
See Go memory model.

~~~~
    T1        T2

1.  snd_A
2.  snd_B
3.           rcv_B
4.			 rcv_A
~~~~~~~~~~

Assume the following "replay run".


~~~~
    T1        T2

1.  snd_A
2.  snd_B             // At this point the buffer is filled as follows [A,B]
                      // where the run-time checks the buffer from left to right.
3.           rcv_B    // At this point, we get stuck, cause we can NOT see the matching partner A.
4.			 rcv_A
~~~~~~~~~~



Issue:

* In the "replay run", buffer elements are kept in a "different" order

* Hence, we may fail to find the matching partner element


Note:

* Buffer space is implemented as a queue

* So, if we execute sends in a squential order, the above situation will not arise.



Variant of the above where the two sends execute in separate threads.
Actual program run.

~~~~
    T1        T2     T3

1.  snd_A
2.          snd_B                [B,A] buffer layout
3.                   rcv_B
4.			         rcv_A
~~~~~~~~~~


In theory, it seems possible that snd_A overtakes snd_B.


Note. Likely, the current instrumentation scheme enforces a total order
and the above will not happen.
Anyway, assume this may happen.


Consider the following replay run.


~~~~
    T1        T2     T3

1.  snd_A
2.          snd_B                [A,B] buffer layout (switched)
3.                   rcv_B       At this point we get stuck cause we can not see the matching partner A.
4.			         rcv_A
~~~~~~~~~~


The total order among sends/receives we record in the trace might be "misleading".


To be safe, we need to "scan" the buffer for matching partners.

Q: Do we still need to process the trace based on the total order?


Example. Buffer size 1.

Actual program run.

~~~~
    T1        T2     T3        T4

1.  snd_A
2.			                   rcv_A
3.          snd_B
4.                   rcv_B
~~~~~~~~~~


Replay run. Say, we pick any order.


~~~~
    T1        T2     T3        T4

3.          snd_B
2.			                   rcv_A   -- suspend
4.                   rcv_B             -- match found
1.  snd_A
                              continue
~~~~~~~~~~



~~~~
    T1        T2     T3        T4

3.          snd_B
1.  snd_A                               // what to do here, suspend
2.			                   rcv_A    // suspend as well
4.                   rcv_B             // match found

    continue with T1 or T4
~~~~~~~~~~


For this example, we can pick any other.


Another example. Buffer size 1.

Actual program run.

~~~~
    T1        T2     T3

1.  snd_A
2.                   rcv_A
3.          snd_B
4.                   rcv_B
~~~~~~~~~~

Consider the following replay run.


~~~~
    T1        T2     T3


3.          snd_B                 // executes
1.  snd_A                         // suspend
2.                   rcv_A        // suspend, we are stuck here !!!
4.                   rcv_B
~~~~~~~~~~


Let's repeat the question.

Q: Do we still need to process the trace based on the total order?

A: Yes, this is necessary for buffered channels, as we might otherwise get stuck!


#### Lock + unlock


Go mutexes behave like buffered channels of size 1.

So, we could "recreate" the above example.


Actual program run.

~~~~
    T1        T2     T3

1.  lock_A
2.                   unlock_A
3.          lock_B
4.                   unlock_B
~~~~~~~~~~


Shows that in general, we need to process lock events based on their total order recorded in the trace.


#### Summary


Total order among trace events is important.

Guarantees that we don't get stuck.

Q: What about communication ids to identify matching partners?
Isn't this already enforced by obeying the total order?

Note. The instrumentation scheme must guarantee that operations take place in the order as recorded!
Also relies on the fact that Go buffered channels are implemented as queues.


## Usage
The trace replay is currently not in the main branch jey, but in its separate
`replay` and `replayDev` branches.

To start the replay, add the following header at the beginning of the
main function:

  ```go
  advocate.EnableReplay(1, true)
  defer advocate.WaitForReplayFinish()
  ```
Replace `1` with the index of the bug you want to replay.\
If a rewritten trace should not return exit codes, but e.g. panic if a 
negative waitGroup counter is detected, of send on a closed channel occurs,
the second argument can be set to `false`.

Also include the following imports:
```go
"advocate"
```

Now the program can be run with the modified go routine, identical to the recording of the trace (remember to export the new gopath).

## Implementation
The following is a description of the current implementation of the trace replay.
It is split into three parts:

- Trace Reading
- Order Enforcement
- State Enforcement

### Trace Reading
First we read in the trace and create a new internal data structure to save
the trace, ordered by `tpre`. For now we ignore atomic events, because they do ~~not~~ change the flow of the program and the implementation of the order
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
func WaitForReplayPath(op Operation, file string, line int) (bool, bool, ReplayElement) {
	if !replayEnabled {
		return false, false, ReplayElement{}
	}

	if AdvocateIgnore(op, file, line) {
		return true, false, ReplayElement{}
	}

	timeoutCounter := 0
	for {
		nextRoutine, next := getNextReplayElement()

		// all elements in the trace have been executed
		if nextRoutine == -1 {
			println("The program tried to execute an operation, although all elements in the trace have already been executed.\nDisable Replay")
			replayEnabled = false
			return false, false, ReplayElement{}
		}

		if next.Time != 0 {
			if (next.Op != op && !correctSelect(next.Op, op)) ||
				next.File != file || next.Line != line {

				timeoutCounter++

				checkForTimeout(timeoutCounter, file, line)
				slowExecution()
				continue
			}
		}

		foundReplayElement(nextRoutine)

		lock(&replayDoneLock)
		replayDone++
		unlock(&replayDoneLock)

		return true, true, next
	}
}
```

This function will check if the calling operation is the next function to be 
executed or not. If it is not, it will hold the execution until this is the case.

First we check if the replay is enabled. If it is not, we do not check if we 
need to hold the exections.

Some operations, like e.g. garbage collection are not predictable and would
lead to stuck executions. For this reason we ignore some of the operations.

If non of these to cases apply, we will check if the current operation is
the next operation, by comparing the code position of the current operation with
the code position of next operation to be executed. If those are not equal,
we will wait for a short moment and check again. Thereby be count how long the
operation is already waiting. If a operation is waiting for a longer time, we
write a message to the term- 0: The replay will ended completely without finding a Replay element
  - 10: Replay Stuck: Long wait time for finishing replay
  - 11: Replay Stuck: Long wait time for running element
  - 12: Replay Stuck: No traced operation has been executed for approx. 20s
  - 13: The program tried to execute an operation, although all elements in the trace have already been executed.
  - 20: Leak: Leaking unbuffered channel or select was unstuck
  - 21: Leak: Leaking buffered channel was unstuck
  - 22: Leak: Leaking Mutex was unstuck
  - 23: Leak: Leaking Cond was unstuck
  - 24: Leak: Leaking WaitGroup was unstuck
  - 30: Send on close
  - 31: Receive on close
  - 32: Negative WaitGroup counter

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


## Exit codes
If a end element is found in the trace and the replay is enabled or the replay is stuck, the 
program will exit with one of the following exit codes:
- 0: The replay will ended completely without finding a Replay element
- 10: Replay Stuck: Long wait time for finishing replay
- 11: Replay Stuck: Long wait time for running element
- 12: Replay Stuck: No traced operation has been executed for approx. 20s
- 13: The program tried to execute an operation, although all elements in the trace have already been executed.
- 20: Leak: Leaking unbuffered channel or select was unstuck
- 21: Leak: Leaking buffered channel was unstuck
- 22: Leak: Leaking Mutex was unstuck
- 23: Leak: Leaking Cond was unstuck
- 24: Leak: Leaking WaitGroup was unstuck
- 30: Send on close
- 31: Receive on close
- 32: Negative WaitGroup counter
