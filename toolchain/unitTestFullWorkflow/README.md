# Explanation
This script goes through the whole advocate worklfow for a single go test.
It can be used on its own, but is mainly intended to be used within the [runFullWorkflowOnAllUnitTests.bash](../runFullWorkflowOnAllUnitTests/runFullWorkflowOnAllUnitTests.bash) script.
It will automatically manage the usual steps usually associated with the advocate analyze process.
The script does the following
1. Change the directory the the project root and export GOROOT
2. Run the test with the patched go runtime
3. Run the analyzer on the produced trace
4. If rewritten_traces were constructed run those traces and save the result in `reorder_outputs.txt`
# Input
As Input it takes up to 6 parameters. 5 of which are needed at minimum.
- -a: Path to the adovacte repository root
- -f: Path to the project root
- -tf: Path to the test file
- -p: Package path of the test
- -t: Testname you want to execute
- -m: (optional) in case you want to enable module mode. See [Common Problems](#common-problems)
# Output
The script will produce the following
- the advocateTrace folder
- if possible rewritten_trace folders
- `reorder_output`.txt files withing their respective trace folders
# Example
For the case of this example it is assumed that your `ADOCATE` repository is in your home directory.
Let's assume you want to the test `TestPostPath` of `~/prometheus` located in the file
`~/prometheus/notifier/notifier_test.go`.
The way you execute the program is like this
```bash
./unitTestFullWorkflow.bash -a ~/ADVOCATE/ -p ./notifier -f ~/prometheus/ -tf /home/mario/Desktop/prometheus/notifier/notifier_test.go -t TestPostPath
```

# Common Problems
This tool requires a go.mod at the project root otherwise the tests won't run.
This is the case for some repositories (eg Moby).
In this case you need to manually add a go.mod via `go mod init` in the project root and call the program with the flag `-m true` like so
```sh
./unitTestFullWorkflow.bash -a <path-advocate> -f <path-kubernetes-root> -m <true> -tf <path-kuberbentes-root>/plugin/pkg/admission/deny/admission_test.go -p plugin/pkg/admission/deny -t TestAdmission 
```