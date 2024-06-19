# Explanation
This tool automates adding the Advocate overhead for a given file.
Applying this tool to a given file and testname will handle the preamble inertion.
Also Additional flags can be provided to also handle the replay overhead insertion.
It works by checking each line and if it sees the function signature it will insert the overhead below it.

# Input
The program takes up 4 parameters as input.
- -f: path to the file that you want to add the overhead to
- -t: the name of the test
- -r: (optional) if set to true the replay overhead will be added
- -n: (optional) if `r` set to true you can the rewritte_trace number with this parameter
# Output
The output is the adjusted file with the same name the original

# Usage
Given a testfile `some_test.go` you can use the script like so
```sh
go run unitTestOverheadInserter.go -f some_test.go -t <test-name> 
```
or like so if you want the replay overhead
```sh
go run unitTestOverheadInserter.go -f <file> -t <test-name> -r true -n <trace-number>
```
# Example
Given a file `some_test.go`
```go
package main

import (
    "testing"
    "fmt"
    "time"
)

func TestSomething(t *testing.T) {
	fmt.Println("Hello from TestSomething")
	c := make(chan int, 0)

	go func() {
		c <- 1
	}()

	go func() {
		<-c
	}()

	time.Sleep(10 * time.Millisecond)
	close(c)
}
```
If you ran 
```bash
go run unitTestOverheadInserter.go -f some_test.go -t TestSomething
```
It would change the file to look like this
```go
package main

import (
    "testing"
    "advocate"
    "fmt"
    "time"
)

func TestSomething(t *testing.T) {
	// ======= Preamble Start =======
	advocate.InitTracing(0)
	defer advocate.Finish()
	// ======= Preamble End =======
	c := make(chan int, 0)

	go func() {
		c <- 1
	}()

	go func() {
		<-c
	}()

	time.Sleep(10 * time.Millisecond)
	close(c)
}
```
In case of adding the replay overhead instead it will look like this
```shell
package main

import (
    "testing"
    "advocate"
    "fmt"
    "time"
)

func TestSomething(t *testing.T) {
    // ======= Preamble Start =======
    advocate.EnableReplay(1, true)
    defer advocate.WaitForReplayFinish()
    // ======= Preamble End =======
	c := make(chan int, 0)

	go func() {
		c <- 1
	}()

	go func() {
		<-c
	}()

	time.Sleep(10 * time.Millisecond)
	close(c)
}
```