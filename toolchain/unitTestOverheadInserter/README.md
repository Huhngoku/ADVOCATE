# Unit Test Preamble Import Inserter
This tool automates adding the Advocate overhead for a given file.
Applying this tool to a given file and testname will handle the preamble inertion.
Also Additional flags can be provided to also handle the replay overhead insertion
### Usage
If a go file contains a main method the cool can be used like so
```sh
go run inserter.go -f filename.go
go run unitTestOverheadInserter.go -f <file> -t <test-name> 
```
or like so if you want the replay overhead
```sh
go run unitTestOverheadInserter.go -f <file> -t <test-name> -r true -n <trace-number>
```
### Output
Has not output itself, but file will be modified
### Example
#### Trace recording
Given a file `file.go`
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
After running `go run inserter.go -f file.go` it will look like this
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
    advocate.EnableReplay(n, true)
    defer advocate.WaitForReplayFinish()
    // ======= Preamble End =======
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
```