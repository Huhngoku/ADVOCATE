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
Lastly you need to set the goroot-environment-variable

### example main method
## Output
- machine_readable.log
- human readable.log
- advocateTrace
- rewritten_Trace
## Tooling
### Preamble and Import Management
### Analyzing of Repository Main Methods
### Analyzing of Repostiory Unit Tests
## How ADVOCATE works