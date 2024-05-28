# Unit Test Gen Trace
With this script it is possible to automatically run all Go Unit Tests within a folder while collecting traces and storing them in a folder for later analysis.
## Usage
The script expects:
- absolute path to advocate root
- path to go project root
### command
```shell
./unitTestGenTrace.bash -a <path-advocate> -f <path-kubernetes>
```
### Output
Traces are storted in a folder like structure under advocateResult.
The naming conventions is as follows: advocateResult > packageName > fileName > testFunctionName
The file contains the output of the analyzer