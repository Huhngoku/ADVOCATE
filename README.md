# AdvocateGo
## What is AdvocateGo
AdvocateGo is an analysis tool for Go programs.
It detects concurrency bugs and gives  diagnostic insight.
This is achieved through `happens-before-relation` and `vector-clocks`

Furthermore it is also able to produce traces which can be fed back into the program in order to experience the predicted bug.

A more in detail explanation of how it works can be found [here](./doc/Analysis.md).
### AdvocateGo Step by Step
Simplified flowchart of the AdvocateGo Process
![Flowchart of AdvocateGoProcess](doc/img/flow2.png "Title")

For more detail see this [in depth diagram](./doc/img/architecture_without_time.png)
## Running your first analysis
These steps can also be done automatically with scripts. If you want to know more about using them you can skip straight to the [Tooling](#tooling) section. Doing these steps manually at least once is still encouraged to get a feel for how advocateGo works.
### Step 1: Add Overhead
You need to adjust the main method or unit test you want to analyze slightly in order to analyze them.
The code snippet you need is
```go
import "advocate"
...
// ======= Preamble Start =======
    advocate.InitTracing(0)
    defer advocate.Finish()
// ======= Preamble End =======
...
```
Eg. like this 
```go
import "advocate"
func main(){
    // ======= Preamble Start =======
    advocate.InitTracing(0)
    defer advocate.Finish()
    // ======= Preamble End =======
...
}
```
or like this for a unit test
```go
import "advocate"
...
func TestImportantThings(t *testing.T){
    // ======= Preamble Start =======
    advocate.InitTracing(0)
    defer advocate.Finish()
    // ======= Preamble End =======
...
}
```
### Step 2: Build AdvocateGo-Runtime
Before your newly updated main method or test you will need to build the AdvocateGo-Runtime.
This can be done easily by running
#### Unix
```shell
./src/make.bash
```
#### Windows
```shell
./src/make.bat
```
### Step 2.2: Set Goroot Environment Variable
Lastly you need to set the goroot-environment-variable like so
```shell
export GOROOT=$HOME/ADVOCATE/go-patch/
```
### Step 3: Run your go program!
Now you can finally run your go program with the binary that you build in `Step 1`.
It is located under `./go-patch/bin/go`
Eg. like so
```shell
./go-patch/bin/go run main.go
```
or like this for your tests
```shell
./go-patch/bin/go test
```
## Analyzing Traces
After you run your program you will find that it generated the folder `advocateTrace`.
If you are curious about the structure of said trace, you can find an in depth explanation [here](./doc/Trace.md)
It contains a record of what operation ran in what thread during the execution of your program.

This acts as input for the analyzer located under `./analyzer/analyzer`.
It can be run like so
```shell
./analyzer/analyzer -t advocateTrace
```
### Output
Running the analyzer will generate 3 files for you
- machine_readable.log (good for parsing and further analysis)
- human readable.log (more readable representation of bug predictions)
- rewritten_Trace (a trace in which the bug it was rewritten for would occur)

A more detailed explanation of the file contents can be found under [AnalysisResult.md](./doc/AnalysisResult.md)

### What bugs can be found
AdvocateGo currently supports these bugs

## Replay
### How to replay the program and cause the predicted bug
This process is similar to when we first ran the program. Only the Overhead changes slightly.

Instead want to use this overhead

```go
// ======= Preamble Start =======
advocate.EnableReplay(n)
defer advocate.WaitForReplayFinish()
// ======= Preamble End =======
```

where the variable `n` is the rewritten trace you want to use.
Note that the method looks for the `rewritten_trace` folder in the same directory as the file is located
### Which bugs are supported for replay
A more detailed description of how replays work and a list of what bugs are currently supported for replay can be found under `./doc/TraceReplay.md` and `./doc/TraceReconstruction.md`



## Tooling
There are certain scripts that will come in handy when working with AdvocateGo 
### Preamble and Import Management
There are programs that automatically add and remove the overhead described in Step 1
#### For Main Methods
#### For Unit Tests
### Analyzing an Existing local project
#### Main Method

#### Unit Tests
This script takes an entire repository and automatically runs the analysis + replay on all unit tests. After running you will additionally get a csv file that lists all predicted and confirmed bugs. (ongoing)

Its result and additional information (rewritten traces, logs, etc) will be written to. `advocateResult`

The script + a more detailed description can be found under `./toolchain/runFullWorkflowOnAllUnitTests`
### Analyzing of Github Repositories
Will automatically take the unit tests as input and run the same workflow as with the local script.

The script + a more detailed description can be found under `./toolchain/`

## Warning
It is the users responsibility of the user to make sure, that the input to 
the program, including e.g. API calls are equal for the recording and the 
tracing. Otherwise the replay is likely to get stuck.

Do not change the program code between trace recording and replay. The identification of the operations is based on the file names and lines, where the operations occur. If they get changed, the program will most likely block without terminating. If you need to change the program, you must either rerun the trace recording or change the effected trace elements in the recorded trace.
This also includes the adding of the replay header. Make sure, that it is already in the program (but commented out), when you run the recording.
