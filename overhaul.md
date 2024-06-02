# AdvocateGo
## What is AdvocateGo
AdvocateGo is an analysis tool for Go programs.
It detects concurrency bugs and gives  diagnostic insight.
This is achieved through...

Furthermore it is also able to produce traces which can be fed back into the program in order to experience the predicted bug.
### AdvocateGo Process
![Flowchart of AdvocateGoProcess](doc/img/flow.png "Title")
## Running your first analysis
These steps can also be done automatically with scripts located in `toolchain` but doing these steps at least once are good to get a feel for how advocateGo works. You can find a more detailed tutorial for automation there.
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
## Analyzing Trace
After you run your program you will find that it generated the folder `advocateTrace`.
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

### What bugs can be found

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
The bugs that are currently supported for the replay feature are
- P1: Possible send on closed channel
- P2: Possible receive on closed channel
- P3: Possible negative waitgroup counter
- L1: Leak on unbuffered channel with possible partner
- L3: Leak on buffered channel with possible partner
- L6: Leak on select with possible partner
- L8: Leak on mutex
- L9: Leak on waitgroup
- L0: Leak on cond
## Tooling
### Preamble and Import Management
### Analyzing of Repository Main Methods
### Analyzing of Repostiory Unit Tests
## How ADVOCATE works