# Full Workflow On Individual Unit Tests
This script runs the entire workflow on an individual Unit Test.
The full workflow contains:
- preamble handling
- running the unit test with patched runtime
- running all rewritten traces
- evaluating found and reproduced bugs (in progress)
## Usage
It takes two parameters:
- absolut path to advocate root
- the root of the go project
- the package the test lies in
- the file that contains the test
- the name of the test
```sh
./runFullWorkflowOnAllUnitTests -a <path-to-advocate> -f <path-to-folder> -p <package> -tf <path-to-test-file> -t <test-name>
```
## Example
Let's say we want to the kubernetes unit test `TestAdmission` in the package `plugin/pkg/admission/deny`.
The command would be
```sh
./unitTestFullWorkflow.bash -a <path-advocate> -f <path-kubernetes-root> -tf <path-kuberbentes-root>/plugin/pkg/admission/deny/admission_test.go -p plugin/pkg/admission/deny -t TestAdmission     
```
### Output (in progress)
- A list of bugs and diagnostics, the analyzer found
- A csv containing bugs we were actually able to reproduce & resolve