# Explanation
This shell script automatically uses every unit test in a go project to generate an advocate go analysis for each of them. Note this only works if the project root contains a go.mod. If this is not the case you can check the [Common Problem](#common-problems) section for a potential fix.

That aside the script works like this
1. Change the directory to the project
2. Finds all test files and loops through them.
3. For every test in a file it applies the script [unitTestFullWorkflow.bash](../unitTestFullWorkflow/unitTestFullWorkflow.bash)
4. Stores the result within the a new advocateResult in the root
# Input
The script takes up to 3 parameters
- -a: path to the advocate folder's root
- -f: path to the root of the go project you want to test 
- -m: (optional) boolean to enable mod mode. See [Common Problems](#common-problems)
# Output
The output is an advocateResult directory in the root of the project you are analyzing.
It contains a record of all the logs+traces that were produced during the analysis.
This can be fed to the [generateStatisticsFromAdvocateResult](../generateStatisticsFromAdvocateResult/generateStatistics.go) script to get a more summarized overview of the bugs found.
# Usage
## Normal Project
For a given go project contained in the folder `myGoProject` you can simply run 
```bash
./runFullWorkflowOnAllUnitTests.bash -a <path-to-advocate> -f <path-tomyGoProject>
```
## Project that doesn't yet have a go.mod
In the unusual case that no go.mod is provided at the root you can try to create one yourself using `go mod init`.
In this case because the project is still not set up correct you need to add the additional flag `-m true` and run the program. Like so
```bash
./runFullWorkflowOnAllUnitTests.bash -a <path-to-advocate> -f <path-tomyGoProject> -m true
```
# Example
It is assumed that the `ADVOCATE` repository is also in your home directory.
Note that it is important for the script to run correctly that you end absolute folder paths with `/` like `~/ADVOCATE/` and not `~/ADVOCATE`
## Normal Project
A project that is relatively small and can be run on most local machines is `prometheus`.
You can get the project by running `git clone https://github.com/prometheus/prometheus ~`.
After the project is cloned the project into your home directory you can run the script like so.
```bash
./runFullWorkflowOnAllUnitTests.bash -a ~/ADVOCATE/ -f ~/prometheus/
```
The script will go through all the tests and sometimes skip the ones it is unable to run the analysis on.
## Project without go.mod
Moby is one such case where if you were to run the script as is it would result in errors, because the go project doesn't contain a `go.mod` at its root.
To analyze repositories like `moby` regardless use these steps
1. Clone ̀`moby` or a repository of your choice like `git clone https://github.com/moby/moby ~`
2. cd into the project with `cd ~/moby` and run `go mod init` to create a `go.mod`
3. Run the script with the dedicated flags like so
```bash
./runFullWorkflowOnAllUnitTests.bash -a ~/ADVOCATE -f ~/moby -m true
```
# Common Problems
This tool requires a go.mod at the project root otherwise the tests won't run.
This is the case for some repositories (eg Moby).
In this case you need to manually add a go.mod via `go mod init` in the project root and call the program with the flag `-m true` like so
```sh
./runFullWorkflowOnAllUnitTests -a <path-to-advocate> -f <path-to-folder> -m true
```