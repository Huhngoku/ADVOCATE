# Github Main Analyzer
This program takes a link to a github repository and runs the first main method it finds with ADVOCATE.
It then uses the analyzer to check for potential bugs in the generated trace.
## Usage
The program expects :
- the absolute path to the ADVOCATE root 
- link to a github repository containing a main method
### Command
```shell
./githubMainAnalyzer.bash -a <path-advocate-root> -u <github-repository-url>
```
### Output
The output will be the output of the program itself + the analyzer output afterwards
### Example
```shell
./githubMainAnalyzer.bash -a /home/user/ADVOCATE -u https://github.com/junegunn/fzf
```