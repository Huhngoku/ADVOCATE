# Preamble Import Inserter
This tool automates adding the Advocate overhead for a given file.
After applying this tool to a file, the preamble will be inserted right after the start of main and also advocate will be added to the imports.
##Example
### Single file
If a go file contains a main method the cool can be used like so
```sh
go run inserter.go -f filename.go
```
