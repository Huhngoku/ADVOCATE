# Example
This tool removes the advocate overhead for a testfile.
It does so by looking through every line and if it finds the advocate import or preamble somwhere skips them.

# Input
This script takes 2 arguments.

- -f: Filepath to the testfile
- -t: Testname

# Output
The output is the adjusted file with the same name the original

# Usage
Given a file `some_test.go` you can use the script like so
```bash
go run unitTestOverheadRemover.go -f <path-to-some_test.go> -t <test-name>
```
### Example
Given a file `file.go`
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
After the script ran over the file it will look like:
```go
package main

import (
    "testing"
    "fmt"
    "time"
)

func TestSomething(t *testing.T) {
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