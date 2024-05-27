# Analysis Result

This file explains the machine readable (results_machine.log) and the human
readable (rusult_readable.log) result file of the analysis.

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

`[typeID]` is the type of the found bug.
Thy have the following form:

- P | A | L for possible, actual and leak
- a number for the type of the bug

Therefore the possible types are possible:

- A1: Receive on closed channel
- A2: Close on closed channel
- A3: Concurrent recv
- A4: Select case without partner
- P1: Possible send on closed channel
- P2: Possible receive on closed channel
- P3: Possible negative waitgroup counter
- L1: Leak on unbuffered channel or select with possible partner
- L2: Leak on unbuffered channel or select without possible partner
- L3: Leak on buffered channel with possible partner
- L4: Leak on buffered channel without possible partner
- L5: Leak on mutex
- L6: Leak on waitgroup
- L7: Leak on cond

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
    - MN: RUnlock
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
TODO: Add when implemented

### Close on close
TODO: Add when implemented

### Possible send on closed
A possible send on closed has the following form:
```
P1,T:1:22:12:CC:/home/ex/main.go:123,T:2:22:8:CS:/home/ex/main.go:222
```
It is constructed in the following way:

`P1`: The type of the bug is a possible send on closed channel.

`T:1:22:12:CC:/home/erik/Uni/HiWi/ADVOCATE/examples/ex/ex.go:123`:
- T: The element is a trace element
- 1: The close is in routine 1
- 22: The close is on channel 22
- 12: The close is at time (tpre) 12
- CC: The object is a channel close
- /home/erik/ex/main.go:123: The close is in the file /home/ex/main.go at line 123

`T:2:22:8:CS:/home/ex/main.go:222`:
- T: The element is a trace element
- 2: The send is in routine 2
- 22: The send is on channel 22
- 8: The send is at time (tpre) 8
- CS: The object is a channel send
- /home/ex/main.go:222: The send is in the file /home/ex/main.go:222 at line 222


### Possible recv on closed
TODO: Add when implemented


### Concurrent recv
TODO: Add when implemented

### Leak on unbuffered channel or select
#### With possible partner

TODO: Add when implemented

#### Without possible partner
TODO: Add when implemented

### Leak on buffered channel
#### With possible partner
TODO: Add when implemented

#### Without possible partner
TODO: Add when implemented
### Leak on mutex
TODO: Add when implemented

### Leak on waitgroup
TODO: Add when implemented

### Leak on cond
TODO: Add when implemented

### Select case without partner
TODO: Add when implemented
### Possible negative waitgroup counter
TODO: Add when implemented













## Human readable result file

The result file contains all potential bugs found in the analyzed trace.\
A possible result would be:
```
Possible send on closed channel:
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@47
	send: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@44
Possible receive on closed channel:
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@47
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
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@40
	recv: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@44
```
Close contains the tID of the close operation.
Recv contains the tID of the recv operation.

### Close on close
An actual close on closed has the following form:
```
Found receive on closed channel:
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
```
The first close contains the close that lead to the panic.
The second close contains the first close of the channel.

### Possible send on closed
A possible send on closed has the following form:
```
Found receive on closed channel:
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	send: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
```
The close contains the close on the channel
The send contains the send on the channel, that might be closed.


### Possible recv on closed
A possible send on closed has the following form:
```
Possible receive on closed channel:
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	recv: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
```
The close contains the close on the channel
The recv contains the recv on the channel, that might be closed.


### Concurrent recv
A concurrent recv has the following form:
```
Found concurrent Recv on same channel:
	recv: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	recv: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
```
The first recv contains the recv that is concurrent but in this case after the second recv.
The second recv contains the recv that is concurrent but in this case before the second recv.

### Leak on unbuffered channel or select
#### With possible partner

A leak on an unbuffered channel with a possible partner has the following form:
```
Leak on unbuffered channel or select with possible partner:
	channel: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	partner: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
```
The channel contains the tID of the channel or select that is leaking.\
The partner contains the tID of a possible partner of the channel or select.

#### Without possible partner
A leak on an unbuffered channel without a possible partner has the following form:
```
Leak on unbuffered channel or select with possible partner:
	channel: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	partner: -
```
The channel contains the tID of the channel or select that is leaking.\
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
	partner: -
```
The channel contains the tID of the channel or select that is leaking.\
The partner contains a "-" because there is no possible partner.

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
	wait-group: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44

```
wait-group contains the tID of the waitgroup that is leaking.\
The second line is empty.

### Leak on cond
A leak on a cond has the following form:
```
Leak on conditional variable:
	conditional: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44

```
conditional contains the tID of the conditional variable that is leaking.\
The second line is empty.

### Select case without partner
A select case without partner has the following form:
```
Possible select case without partner or nil case:
	select: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	case: 18,R
```
Select contains the tID of the select that is missing a possible partner.\
Case contains the case that is missing a partner. It consist of the channel number and the direction (S: send, R: recv).

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