# Analysis Result

This file explains the machine readable result file (results_machine.log) of the analysis.
The human readable result file (result_readable.log) has the same general structure,
but contains additional formatting.

The result file contains all potential bugs found in the analyzed trace.
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

## Results
The following is a list of possible elements in the result file.
The results for cyclic and mixed deadlocks are currently disabled and therefore
not described.

## Receive on close
An actual recv on closed has the following form:
```
Found receive on closed channel:
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@40
	recv: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@44
```
Close contains the tID of the close operation.
Recv contains the tID of the recv operation.

## Close on close
An actual close on closed has the following form:
```
Found receive on closed channel:
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
```
The first close contains the close that lead to the panic.
The second close contains the first close of the channel.

## Possible send on closed
A possible send on closed has the following form:
```
Found receive on closed channel:
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	send: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
```
The close contains the close on the channel
The send contains the send on the channel, that might be closed.


## Possible recv on closed
A possible send on closed has the following form:
```
Possible receive on closed channel:
	close: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	recv: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
```
The close contains the close on the channel
The recv contains the recv on the channel, that might be closed.


## Concurrent recv
A concurrent recv has the following form:
```
Found concurrent Recv on same channel:
	recv: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	recv: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
```
The first recv contains the recv that is concurrent but in this case after the second recv.
The second recv contains the recv that is concurrent but in this case before the second recv.

## Leak on unbuffered channel or select
### With possible partner

A leak on an unbuffered channel with a possible partner has the following form:
```
Leak on unbuffered channel or select with possible partner:
	channel: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	partner: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
```
The channel contains the tID of the channel or select that is leaking.\
The partner contains the tID of a possible partner of the channel or select.

### Without possible partner
A leak on an unbuffered channel without a possible partner has the following form:
```
Leak on unbuffered channel or select with possible partner:
	channel: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	partner: -
```
The channel contains the tID of the channel or select that is leaking.\
The partner contains a "-" because there is no possible partner.

## Leak on buffered channel
<!-- TODO: check partner-->
A leak on a buffered channel has the following form:
```
Leak on buffered channel:
	channel: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44

```
The channel contains the tID of the channel that is leaking.\
The second line is empty.

## Leak on mutex
A leak on a mutex has the following form:
```
Leak on mutex:
	mutex: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	last: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1095@40
```
Mutex contains the tID of the mutex that is leaking.\
Last contains the tID of the last lock operation on the mutex.

## Leak on waitgroup
A leak on a waitgroup has the following form:
```
Leak on wait group:
	wait-group: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44

```
wait-group contains the tID of the waitgroup that is leaking.\
The second line is empty.

## Leak on cond
A leak on a cond has the following form:
```
Leak on conditional variable:
	conditional: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44

```
conditional contains the tID of the conditional variable that is leaking.\
The second line is empty.

## Select case without partner
A select case without partner has the following form:
```
Possible select case without partner or nil case:
	select: /home/erik/Uni/HiWi/ADVOCATE/examples/constructed/constructed.go:1100@44
	case: 18,R
```
Select contains the tID of the select that is missing a possible partner.\
Case contains the case that is missing a partner. It consist of the channel number and the direction (S: send, R: recv).

## Possible negative waitgroup counter
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