# Select

Channel operations including select statements are recorded in the trace where the select statement is located.

## Info
Select statements with only one case are not recorded as select statements. If the case is an default case, it is not recorded at all. If the case is a send or receive, the select statement is equivalent to just the send/receive statement and therefor recorded as such. 


## Trace element
The basic form of the trace element is 
```
S,[tpre],[tpost],[id],[cases],[exec],[chosen],[oId],[pos]
```
where `S` identifies the element as a select element.
The other fields are set as follows:
- [tpre]: This is the value of the global counter when the routine gets to the select statement.
- [tpost]: This is the value of the global counter when the select has finished, meaning if either on of the select cases has successfully send or received or the default case has been executed
- [id]: This field contains the unique id of the select statement
- [cases]: This field shows the available cases in the trace. Each case is denoted by the id of the involved channel, followed by a `r` for receive cases or `s` for send cases. A default case is shown by a `d` (without any id number). The cases are separated by points (.). The 
order of the cases must not be equal to the order of the cases in the select statement.
- [exec]: This field shows, whether the operation was finished ([exec] = `t`) or not ([exec] = `f`). Failed can e.g. mean, that in a select with no default case, non of the cases could find a communication partner when the program terminates.
- [chosen]: This field shows the index of the chosen case. This index is equal to the position of the chosen case in [cases] (0 based). Because
of an internal reordering of the cases, this index must not be equal to the index of the order of cases in the program code. If the default case 
was selected, [chosen] is set to `-1`.
- [oId]: This is the operation id of the send or receive statement executed in the select (see channel). The send/receive statement is not 
additionally recorded as a separate channel element in the trace.
- [pos]: The last field show the position in the code, where the select is implemented. It consists of the file and line number separated by a colon (:). It the select only contains one case, the line number is 
equal to the line of this case.


## Example
The following is an example containing two select statements, one with and one without a select:
```go
package main
func main() {  // Routine 1
    c := make(chan int, 0)  // id = 1
	d := make(chan int, 0)  // id = 2
	e := make(chan int, 1)  // id = 3

	go func() {  // routine 2
		select {  // id = 4, line 7
		case <-d:
			println("d1")
		case <-e:
			println("e1")
		default:
			println("default")
		}

		c <- 1  // line 16
	}()

	select {  // id = 5, line 19
	case <-c:
		println("c2")
	case <-d:
		println("d2")
	}
}
``` 
If we ignore all internal operations, we get the following trace:
```txt
G,1,2;S,7,8,5,2r.1r,t,1,1,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:19
S,3,4,4,3r.2r.d,t,-1,0,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:16;C,5,6,1,S,t,1,0,0,0,/home/erikkassubek/Uni/dedego/go-patch/bin/main.go:7
```

## Implementation
The recording of the select statement is done in the `selectgo` function in the `go-patch/src/runtime/select.go` file. It contains two function calls. The first one is called shortly after the beginning of the `selectgo` functions to record the available cases and there internal order. The second function call is done when `selectgo` returns, to record
the successful execution of the statement and the selected case.\
Selects with only one case are automatically turned into a pure 
send or receive by the compiler. For this reason, these cases must 
not be additionally recorded.\
Selects with exactly one non-default and one default case are als 
rewritten by the compiler. In this case it is necessary to record 
these cases, because the rewritten version does not use the `select`
go function. For this reason, two additional recorder functions
have been added to the `selectnbsend` and `selectnbrecv` functions 
in the `go-patch/src/runtime/chan.go` file.