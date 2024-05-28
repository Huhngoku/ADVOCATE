# Analyze Traces
This script recursively looks for `advocateTrace` folders an runs the analyzer tool on them.
Initially this tool was written to evaluate a large amount of traces as produced by the `toolchain/unitTestGenTrace`
## Usage
The program expects the absolute path to:
- the analyzer binary (`ADVOCATE/analyzer/analyzer`)
- a folder containing advocate traces
### Command
```shell
./analyzeTraces.bash <path-analyzer> <path-folder>
```
### Output
The output will be a chain of analyzer outputs