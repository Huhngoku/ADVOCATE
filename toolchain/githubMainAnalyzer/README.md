# What does this script do
For a given repository this script looks for the main method and runs the advocate go analysis on it
## Example with docker compose
```shell
./githubMainAnalyzer.bash -p [pathToAdvocate]/ADVOCATE/go-patch/bin/go -g [pathToAdvocate]/ADVOCATE/go-patch/ -i [pathToAdvocate]/ADVOCATE/toolchain/overHeadInserter/inserter -r [pathToAdvocate]/ADVOCATE/toolchain/overHeadRemover/remover -t [pathToAdvocate]/ADVOCATE/toolchain/genTrace/genTrace.bash -a [pathToAdvocate]/ADVOCATE/analyzer/analyzer -u https://github.com/docker/compose

```