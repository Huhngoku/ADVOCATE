# Step by Step

The following is a step py step guid on how to use advocate.

## Preparation

- Download the repository.
- Move into the `ADVOCATE/go-bash/src` directory.
- Run `make.bash` or `make.bat` to build the new runtime. This will create a
`ADVOCATE/go-bash/bin/go` executable.

## Running the program / record the trace
- Add the following header at the beginning of your program:
  ```go
  advocate.InitTracing(0)
  defer advocate.Finish()
  ```
  In most cases the beginning is either the beginning of the main function, or the beginning of a
test function.
- Also import the `advocate` package:
  ```go
  import "advocate"
  ```

- Set the `GOROOT` environment variable to `ADVOCATE/go-patch`.
In Linux this can be done with e.g.:
  ```bash
  export GOROOT=/home/erik/Uni/HiWi/ADVOCATE/go-patch/
  ```
  Replace the path with the path to the `ADVOCATE/go-patch` directory.
- Build the program using the new runtime and run the executable, e.g.
  ```bash
  ~/ADVOCATE/go-bash/bin/go build && ./main
  ```
  or run the test with e.g.
  ```bash
  ~/ADVOCATE/go-bash/bin/go test -run TestAllocate
  ```
  This will create a trace folder `advocateTrace` in the directory of the executable or the trace.

## Analyzing and rewriting the trace
- Move to the `ADVOCATE/analyzer` directory.
- Build the analyzer with `go build`. Make sure to not use the modified runtime.
Use the standard go runtime instead (also make sure, that GOROOT is 
pointing to the standard go runtime).
- Run the analyzer with the path to the trace folder created in the previous step as argument, e.g. 
  ```bash
  ~/ADVOCATE/analyzer/analyzer -t ../prog/advocateTrace
  ```
  This will analyze the trace, find potential bugs and rewrite the trace.\
If you only want to run the analyzer, without rewriting the trace, you can
set the `-x` flags. For all other flags, see the README.md.\
The analyzer will create a `result_machine.log` and a `result_readable.log` file,
containing the results of the analysis.\
If `-x` is not set, it will also create a `rewrite_trace_i` folder, containing
the rewritten trace for each bug found, where a rewrite is possible/implemented.
`i` is always replaced by the index of the bug in the result file.

## Replay the trace

- To replay a rewritten trace, replace the header in the program with:
  ```go
  advocate.EnableReplay(1, true)
  defer advocate.WaitForReplayFinish()
  ```

- Replace `1` with the index of the bug you want to replay.
- If a rewritten trace should not return exit codes, but e.g. panic if a 
negative waitGroup counter is detected, of send on a closed channel occurs,
the second argument can be set to `false`.
- Build the program using the new runtime and run the executable, e.g.
  ```bash
  ~/ADVOCATE/go-bash/bin/go build && ./main
  ```
  or run the test with e.g.
  ```bash
  ~/ADVOCATE/go-bash/bin/go test -run TestAllocate
  ```
  the same way, the program was built and run before. Make sure, that the `GOROOT` environment variable is set to the path of the
modified runtime.