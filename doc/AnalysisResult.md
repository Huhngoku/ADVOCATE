# Analysis Result

This file explains the machine readable (results_machine.log) and the human
readable (result_readable.log) result file of the analysis.

## Machine readable result file

The result file contains all potential bugs found in the analyzed trace.\
A possible result would be:

The file contains one line for each found problem. The line consists of the
following elements:
```
line: [typeID],[args] | [typeId],[args],[args]
```
with
```
[args]: [arg] | [arg];[arg] | [arg];[arg];[arg] | ...
[arg] : T:[routineId]:[objId]:[tpre]:[objType]:[file]:[line] (trace element)
[arg] : S:[objId]:[objType] (select case)
```
The elements have the following meaning:
S:[objId]:[objType] (select case)
- A1: Receive on closed channel
- A2: Send on closed channel
- A3: Close on closed channel
- A4: Concurrent recv
- A5: Select case without partner
- P1: Possible send on closed channel
- P2: Possible receive on closed channel
- P3: Possible negative waitgroup counter
- L1: Leak on unbuffered channel or select with possible partner
- L2: Leak on unbuffered channel or select without possible partner
- L3: Leak on buffered channel with possible partner
- L4: Leak on buffered channel without possible partner
- L5: Leak on nil channel
- L6: Leak on select with only nil channels
- L7: Leak on mutex
- L8: Leak on waitgroup
- L9: Leak on cond

`[args]` shows the elements involved in the problem. There are either
one or two, while the args them self can contain multiple trace elements or select cases.\
The arg in args are separated by a semicolon.\
Each arg contains the following elements separated by a colon (:)
- `[routineId]` is the id of the routine that contains the operation
- `[objId]` is the id of the object that is involved in the operation
- `[tpre]` is the time of the operation
- `[opjType]` is the type of the element
  - Channel:
    - CS: send
    - CR: receive
    - CC: close
  - Mutex:
    - ML: Lock
    - MR: RLock
    - MT: TryLock
    - MY: TryRLock
    - MU: Unlock
    - MN: RUnlockfmt.Sprintf("T%d:%d:%d:%s:%s:%d", t.routineID, t.objID, t.tPre, t.objType, t.file, t.line)
  - Waitgroup:
    - WA: Add
    - WD: Done
    - WW: Wait
  - Select:
    - SS: Select
  - Cond:
    - NW: Wait
    - NB: Broadcast
    - NS: Signal
  - Once:
    - OE: Done Executed
    - ON: Done Not Executed (because the once was already executed)
  - Routine:
    - GF: Fork
- `[file]` is the file of the operation in the program code
- `[line]` is the line of the operation in the program code

The following are examples for the possible types:

### Receive on close
A receive on closed is an actual receive on a closed channel.
The two args of this case are:

- the receive operation
- the close operation

It has the following form:
```
[[Missing]]
```

### Send on close
A send on closed is an actual send on a closed channel (always leads to panic).
The two args of this case are:

- the send operation
- the close operation

It has the following form:
```
[[Missing]]
```

### Close on close
A send on closed is an actual close on a closed channel (always leads to panic).
The two args of this case are:

- the close operation that leads to the panic
- the send operation that is the first close on the channel

It has the following form:
```
[[Missing]]
```

### Concurrent recv
A concurrent recv shows two receive operations on the same channel that are concurrent.:
The two aregs of this case are:

- the first recv operation
- the second recv operation

It has the following form:
```
[[Missing]]
```

### Select case without partner or nil case
A select case without partner shows a select case that is missing a partner or is a nil case.
The two args of this case are:

- the select operation
- the select case without partner

The select case consists of the channel number and the direction (S: send, R: recv). If the channel is nil, the channel number is -1.

It has the following form:
```
[[Missing]]
```

### Possible send on closed
A possible send on closed is a possible but not actual send on a closed channel.
The two args of this case are:

- the send operation
- the close operation


A possible send on closed has the following form:
```
P1,T:2:22:8:CS:/home/ex/main.go:222,T:1:22:12:CC:/home/ex/main.go:123
```
It is constructed in the following way:

`P1`: The type of the bug is a possible send on closed channel.

`T:2:22:8:CS:/home/ex/main.go:222`:
- T: The element is a trace element
- 2: The send is in routine 2
- 22: The send is on channel 22
- 8: The send is at time (tpre) 8
- CS: The object is a channel send
- /home/ex/main.go:222: The send is in the file /home/ex/main.go:222 at line 222

`T:1:22:12:CC:/home/erik/Uni/HiWi/ADVOCATE/examples/ex/ex.go:123`:
- T: The element is a trace element
- 1: The close is in routine 1
- 22: The close is on channel 22
- 12: The close is at time (tpre) 12
- CC: The object is a channel close
- /home/erik/ex/main.go:123: The close is in the file /home/ex/main.go at line 123

### Possible recv on closed
A possible recv on closed is a possible but not actual recv on a closed channel.
The two args of this case are:

- the recv operation
- the close operation

A possible recv on closed has the following form:
```
[[Missing]]
```


### Possible negative waitgroup counter
A possible negative waitgroup counter is a possible but not actual negative waitgroup counter.
The two args of this case are:
- The list of add operations that might make the counter negative (separated by semicolon)
- The list of done operations that might stop the counter from become negative (separated by semicolon)

A possible negative waitgroup counter has the following form:
```
[[Missing]]
```

### Leak on unbuffered channel or select
#### With possible partner
A leak on an unbuffered channel with a possible partner is a unbuffered channel or select that is leaking,
but has a possible partner.
The two arg of this case is:

- the channel or select that is leaking
- the possible partner of the channel or select

A leak on an unbuffered channel or select with a possible partner has the following form:
```
[[Missing]]
```

#### Without possible partner
A leak on an unbuffered channel without a possible partner is a unbuffered channel that is leaking,
but has no possible partner.
The one arg of this case is: mostRecentAcquireTotal[id]

- the channel that is leaking

A leak on an unbuffered channel or select without a possible partner has the following form:
```
[[Missing]]
```

### Leak on buffered channel
#### With possible partner

A leak on an buffered channel with a possible partner is a buffered channel that is leaking,
but has a possible partner.

The two arg of this case is:

- the channel that is leaking
- the possible partner of the channel

A leak on an buffered channel with a possible partner has the following form:
```
[[Missing]]
```

#### Without possible partner

A leak on an buffered channel without a possible partner is a buffered channel that is leaking,
but has no possible partner.

The one arg of this case is:

- the channel that is leaking

A leak on an buffered channel without a possible partner has the following form:
```
[[Missing]]
```

### Leak on nil channel

A leak on a nil channel is a nil channel trying to communicate.

The one arg of this case is:

- the nil channel that is leaking

A leak on a mutex has the following form:
```
[[Missing]]
```

### Leak on select with only nil channels
A leak on a select with only nil channels has the following form:

The one arg of this case is:

- the select that is leaking

A leak on a select with only nil channels has the following form:
```
[[Missing]]
```


### Leak on mutex

A leak on a mutex is a mutex that is leaking.

The onethe mutex that is leakinglose: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@47
	recv: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1103@43
Possible negative waitgroup counter:
	add: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@77;
	done: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@80;
```
Each found problem consist of three lines. The first line explains the
type of the found bug. The other two line contain the information about the
elements responsible for the problem. The elements always have the
form of
```
[type]: [file]:[line]@[tPre]
```
`[file]:[line]@[tPre]` is called tID

The following is a list of possible elements in the result file.
The results for cyclic and mixed deadlocks are currently disabled and therefore
not described.

### Receive on close
An actual recv on closed has the following form:
```
Found receive on closed channel:
	recv: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@44
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@40
```
Close contains the tID of the close operation.
Recv contains the tID of the recv operation.

### Send on close
An actual send on closed has the following form:
```
Found send on closed channel:
	send: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@44
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@40
```
Close contains the tID of the close operation.
Recv contains the tID of the recv operation.

### Close on close
An actual close on closed has the following form:
```
Found close on closed channel:
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
```
The first close contains the close that lead to the panic.
The second close contains the first close of the channel.


### Concurrent recv
A concurrent recv has the following form:
```
Found concurrent Recv on same channel:
	recv: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	recv: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
```
The first recv contains the recv that is concurrent but in this case after the second recv.
The second recv contains the recv that is concurrent but in this case before the second recv.

### Select case without partner
A select case without partner has the following form:
```
Found select case without partner or nil case:
	select: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	case: 18,R
```
Select contains the tID of the select that is missing a possible partner.\
Case contains the case that is missing a partner. It consist of the channel number and the direction (S: send, R: recv). If the channel is nil, the channel number is -1

### Possible send on closed
A possible send on closed has the following form:
```
Possible send on closed channel:
	send: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
```
The close contains the close on the channel
The send contains the send on the channel, that might be closed.


### Possible recv on closed
A possible send on closed has the following form:
```
Possible receive on closed channel:
	recv: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
```
The close contains the close on the channel
The recv contains the recv on the channel, that might be closed.


### Possible negative waitgroup counter
A possible negative waitgroup counter has the following form:
```
Possible negative waitgroup counter:
	add: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44;/home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1234@45
	done: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40;/home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@42
```
Add contains the tIDs of the add operation that might make the counter negative, as well as all add
operations on the same waitgroup, which, if reordered, lead to the negative wait group counter(separated by semicolon).\
Done contains all done
operations on the same waitgroup, which, if reordered, lead to the negative wait group counter(separated by semicolon).

### Leak on unbuffered channel or select
#### With possible partner

A leak on an unbuffered channel with a possible partner has the following form:
```
Leak on unbuffered channel or select with possible partner:
	channel: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	partner: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
```
or
```
Leak on unbuffered channel or select with possible partner:
	select: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	partner: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
```
The channel/select contains the tID of the channel or select that is leaking.\
The partner contains the tID of a possible partner of the channel or select.

#### Without possible partner
A leak on an unbuffered channel without a possible partner has the following form:
```
Leak on unbuffered channel or select with possible partner:
	channel: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
```
or
```
Leak on unbuffered channel or select with possible partner:
	select: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
```
The channel/select contains the tID of the channel or select that is leaking.\
The partner contains a "-" because there is no possible partner.

### Leak on buffered channel
#### With possible partner

A leak on an buffered channel with a possible partner has the following form:
```
Leak on buffered channel with possible partner:
	channel: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	partner: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
```
The channel contains the tID of the channel or select that is leaking.\
The partner contains the tID of a possible partner of the channel or select.

#### Without possible partner
A leak on an buffered channel without a possible partner has the following form:
```
Leak on buffered channel without possible partner:
	channel: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
```
The channel contains the tID of the channel or select that is leaking.\
The partner contains a "-" because there is no possible partner.

### Leak on nil channel
A leak on a nil channel has the following form:
```
Leak on nil channel:
	channel: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
```
Channel contains the tID of the nil channel that is leaking.

### Leak on select with only nil channels
A leak on a select with only nil channels has the following form:
```
Leak on select with only nil channels:
	select: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
```
Select contains the tID of the select that is leaking.

### Leak on mutex
A leak on a mutex has the following form:
```
Leak on mutex:
	mutex: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	last: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
```
Mutex contains the tID of the mutex that is leaking.\
Last contains the tID of the last lock operation on the mutex.

### Leak on waitgroup
A leak on a waitgroup has the following form:
```
Leak on wait group:
	wait: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44

```
wait-group contains the tID of the waitgroup that is leaking.\
The second line is empty.

### Leak on cond
A leak on a cond has the following form:
```
Leak on conditional variable:
	cond: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44

```
cond contains the tID of the conditional variable that is leaking.\
The second line is empty.

