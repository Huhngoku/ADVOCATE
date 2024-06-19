# Explanation
The script `generateStatistics.go` aims to summarize analysis results in a comprehensible and digestible manner.
It manages the data for a specific scenario with the `caseReport` dataType, fills them with the data found within a foulder and then printing them out.
For every scenario it does it generates a new caseReport with the `getCaseReportForCode(code string, folder string)` helper function. The function simply looks for all files name `rewrite_info.log` that located in the directory provided and parse the values pesent. 
The `rewrite_info.log` format looks like this
```
number in analysis result#Scenario code#Expected code for rewrittent race
```
So like this for instance `1#L6#20`. 
It then filters the files so that only files that give information about the code requested remain.
For those files it then looks in the same directory for the corresponding `reorder_output.txt` that contains the log of what happened when we tried to execute an reordered trace.
With a simple regex we can extract the actual exit code that was produced.
Because `caseReport` struct looks like this
```
type caseReport struct {
	caseCode        string
	occurenceCount  int
	actualExitCodes []string
}
```
this means we now gathered  all the information except the `occurenceCount`. This information will be added later on.

Since we iteraively did this for all the scenario codes we now have a list of `caseReport` that contains all corresponding exit codes of the rewrites.
To sum up the scenario occurences we use the helper function `getPredictedBugCounts(folderPath string)` that simply searches for all `results_machine.log`,counts the occurences in a map and updates the `caseReport` we obtained earlier.

The reports are then simply prettyPrinted via.
The result has the form
```
ScenarioCode:Occurences:[list of exit codes that were produced]
```
eg like this
```
A1:0:
A2:737:
A3:2:
A4:29:
A5:7380:
P1:0:
P2:7:31,31,31,31,31,12,31,
P3:0:
L1:0:
L2:58:
L3:41:12,
L4:0:
L5:0:
L6:15:20,20,20,20,12,12,12,12,12,12,12,12,
L7:55:
L8:15:12,12,12,12,12,12,
L9:2:12,
L0:1:
```

Note as not all codes have replay implemented. Some of them will be missing.
It is also possible that scenarios were detected but a rewrite was not possible.

In the example above we see a lot of `Exit code 12`.
Exit code 12 is used when replays get stuck. Meaning in those cases we were not able to produce the actual exit code we expected with the rewrite.
# Input
- -f Path to the advocateResult folder that you want to analyze
# Output
The output prints to std.out and looks like this
```
A1:0:
A2:737:
A3:2:
A4:29:
A5:7380:
P1:0:
P2:7:31,31,31,31,31,12,31,
P3:0:
L1:0:
L2:58:
L3:41:12,
L4:0:
L5:0:
L6:15:20,20,20,20,12,12,12,12,12,12,12,12,
L7:55:
L8:15:12,12,12,12,12,12,
L9:2:12,
L0:1:
```
# Usage
`go run generateStatistics.go -f <path-to-advocateResult-folder>`