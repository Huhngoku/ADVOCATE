# Explanation
This bash script enables you to run and analyze the main function of a program with a single command.
It goes through these steps to achieve this.
1. It changes the directory and sets the GOROOT Environment Variable
2. Adds overhead to the file you provided
3. Run the file with the patched go runtime
4. Removes the overhead
5. Runs the analyzer on the created trace
6. Looks for rewritten traces and runs them
7. Cleanup

# Input
The script takes 2 parameters as input
- -a: The path to the ADVOCATE root
- -f: The path to the file that contains the main method
# Output
The script produces these outputs
- an advocateTrace folder
- rewritten_trace folders if available
- reorder_output.txt within the rewritten_trace folders that will be useful for later analysis
# Usage
Given a fail named main.go you would execute the program like so
```bash
./runFullWorkflowMain.bash -a <path-advocate-root> -f <absolute-path-to-file>
```