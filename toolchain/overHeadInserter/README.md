# Explanation
This tool automates adding the Advocate overhead for a given file.
After applying this tool to a file, the preamble will be inserted right after the start of main and also advocate will be added to the imports. It is also possible to add the replay overhead with this tool if provided with specific flags.
The way the program works is essentially just checking for the signature of the main method or imports statements in every line. If the respective line is found either the import or the preamble will be inserted below it.

# Input 
The program takes up 3 parameters as input.

- -f: path to the file that you want to add the overhead to
- -r: (optional) if set to true the replay overhead will be added
- -n: if `r` set to true you can the rewritte_trace number with this parameter

# Output
The output is the adjusted file with the same name the original

# Usage
If a go file contains a main method the tool can be used like so
```sh
go run inserter.go -f filename.go
```
or like so in case you want to add the replay overhead
```sh
go run inserter.go -f filename.go -r true -n 5
```

# Example
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