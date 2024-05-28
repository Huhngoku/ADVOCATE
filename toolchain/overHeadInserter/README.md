# Preamble Import Inserter
This tool automates adding the Advocate overhead for a given file.
After applying this tool to a file, the preamble will be inserted right after the start of main and also advocate will be added to the imports.
### Usage
If a go file contains a main method the cool can be used like so
```sh
go run inserter.go -f filename.go
```
### Output
Has not output itself, but file will be modified
### Example
Given a file `file.go`
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
After running `go run inserter.go -f file.go` it will look like this
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