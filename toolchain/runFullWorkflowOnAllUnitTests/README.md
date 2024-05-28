# Run Full Workflow On All Unit Tests
This script wraps the Fullworkflow for an individual test in `unitTestFullWorkflow`.
And applies it to every test it can find in a repository
## Usage
It takes two parameters:
- absolut path to advocate root
- absolut path to go project root
It will then look for all unit test functions and apply the full workflow on them.
The full workflow contains:
- preamble handling
- running the unit test with patched runtime
- running all rewritten traces
- evaluating found and reproduced bugs (in progress)
```sh
./runFullWorkflowOnAllUnitTests -a <path-to-advocate> -f <path-to-folder> 
```
## Example
First we need an repository to test to run the program on.
For this I have chosen the kubernetes repository.
`git clone https://github.com/kubernetes/kubernetes`

After cloning the repository its path can be passed to the program via the -f flag and it will handle the analysis of all the unit tests

### Output (in progress)
The output will be chained unitTestFullWorkflow.bash`
An `advocateResult` folder will be created in the Root, that will store the statistics for each individual test.
The statistics contain:
- A list of bugs and diagnostics, the analyzer found
- A csv containing bugs we were actually able to reproduce & resolve