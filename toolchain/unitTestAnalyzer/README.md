# Unit Test Analyzer
With this script it is possible to automatically run all Go Unit Tests within a folder while collecting traces and storing them in a folder for later analysis.

## Usage with Kubernetes
```shell
./unitTestAnalyzer.bash -p [pathToAdvocate]ADVOCATE/go-patch/bin/go -g [pathToAdvocate]/ADVOCATE/go-patch/ -i [pathToAdvocate]/ADVOCATE/toolchain/unitTestOverheadInserter/unitTestOverheadInserter -r /[pathToAdvocate]/ADVOCATE/toolchain/unitTestOverheadRemover/unitTestOverheadRemover -f [pathToKubernetesRepoRoot]
```
## Output
Traces are storted in a folder like structure under advocateResult.
The naming conventions is as follows: advocateResult > packageName > fileName > testFunctionName