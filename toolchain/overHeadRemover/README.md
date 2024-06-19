# Explanation
This tool removes the advocate overhead from a a given file. 
Similar to the inserter it checks every line and if it finds the inserted lines they will be removed.

# Input
The program takes 1 parameter as input
- -f: path to the file that you want to remove the overhead from

# Output
The output is the adjusted file with the same name the original
# Usage
If a go file contains a main method the cool can be used like so
```sh
go run remover.go -f filename.go
```
# Example
Given a file `file.go`
```go
package main

import (
	"time"
	"advocate"
)

func main() {
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
After running `go run remover.go -f file.go` it will look like this
```go
package main

import (
	"time"
)

func main() {
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