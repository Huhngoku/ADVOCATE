# Recording

To analyze a program we have to record the relevant operations. Which operations are recorded is described in `trace.md` and the `traceElements` directory.

## Modified Go

To record this trace, we use a modified version of the Go runtime and standard library, that includes functionality to record these elements. The modified version can be found in the `go-patch` directory. We have to build the new compiler and runtime. This can be done my navigating into `go-patch/src` and running the `all.bash` or `all.bat` file. Potentially failing tests can be ignored. This will create a `go-patch/bin` directory containing a `./go` file. This file can be used like the normal `go` command to build or run programs (e.g. `./go build` or `./go run main.go`). T make it work we also have to change the `GOROOT` environment variable to point to the new runtime by running e.g.
```bash
export GOROOT=$HOME/CoBuFi-Go/go-patch/
```

## Recording
To record the program, we have to add a header at the beginning of the main function:
```go
// initialize the communication for atomics
runtime.InitAtomics(0)

defer func() {
    // disable the trace
	runtime.DisableTrace()

    // write the trace to a file
	file_name := "trace.log"  // name of the trace file
	os.Remove(file_name)
	file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	numRout := runtime.GetNumberOfRoutines()
	for i := 1; i <= numRout; i++ {
		cobufiChan := make(chan string)
		go func() {
			runtime.TraceToStringByIdChannel(i, cobufiChan)
			close(cobufiChan)
		}()
		for trace := range cobufiChan {
			if _, err := file.WriteString(trace); err != nil {
				panic(err)
			}
		}
		if _, err := file.WriteString("\n"); err != nil {
			panic(err)
		}
	}
	
}()
```

It has to be included before any other code. It is also necessary to import the `runtime`,`io/ioutil` and `os` libraries. Warning: Many auto complete tools import the `std/runtime` instead of the `runtime` library. With this, the recording will not work. 

We now run the program like normal (with the created `./go` program in `go-patch/bin`). The trace file will be automatically created as soon as the program execution finishes. It will be created as `trace.log`.

## Known problems

### Holding Locks

In some cases this will result in an 
```
fatal error: schedule: holding locks
``` 
error. The mainly occurs when using the `fmt.Print` command in the 
program. In this case increase the argument in `InitAtomics` until 
the problem disappears.

### Panic

If the program is stopped because of an panic while the main routine is
sleeping, it is possible, that no trace file is created. If the panic
occurs only occasional, it is possible to run the program again, to get
the trace file (assuming the panic does not occur again). If it occurs
in every run, it is necessary to first fix this panic, before a trace can
be recorded.