# Worklow automator
This script automatically manages preamble insertion/removal and exporting goroot and using the patched go runtime to generate a trace
## Usage
This script expects these parameters:
- -p|--patched-go-runtime
- -g|--go-root
- -i|--overhead-inserter
- -r|--overhead-remover)
- -f|--file

The overheadInserter and overheadRemover can also be found in the toolchain directory

## Example
```
./genTrace.bash -p [pathToAdvocate]/ADVOCATE/go-patch/bin/go -g [pathToAdvocate]ADVOCATE/go-patch/ -i [pathToAdvocate] ADVOCATE/toolchain/overHeadInserter/inserter -r [pathToAdvocate] ADVOCATE/toolchain/overHeadRemover/remover -f your_go_program_containing_main_method.go
```