package explanation

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// create an overview over an analyzed, and if possible replayed
// bug. It is mostly meant to give an explanation of a found
// bug to people, who are not used to the internal structure an
// representation of the analyzer.

// It creates one file. This file has the following element:
// - The type of bug found
// - maybe an minimal example for the bug type
// - The test/program, where the bug was found
// - if possible, the command to run the program
// - if possible, the command to replay the bug
// - position of the bug elements
// - code of the bug elements in the trace (+- 10 lines)
// - info about replay (was it possible or not)

/*
 * The function CreateOverview creates an overview over a bug found by the analyzer.
 * It reads the results of the analysis, the code of the bug elements and the replay info.
 * It then writes all this information into a file.
 * Args:
 *    path: the path to the folder, where the results of the analysis and the trace are stored
 *    index: the index of the bug in the results
 * Returns:
 *    error: if an error occurred
 */
func CreateOverview(path string, index int) error {
	bugType, bugPos, bugElemType, err := readAnalysisResults(path, index)
	if err != nil {
		return err
	}

	// get the code info (main file, test name, commands)
	progInfo, err := readProgInfo(path, index)
	if err != nil {
		fmt.Println("Error reading prog info: ", err)
	}

	// get the bug type description
	bugTypeDescription := getBugTypeDescription(bugType)

	// get the code of the bug elements
	code, err := getBugPositions(bugPos)
	if err != nil {
		fmt.Println("Error getting bug positions: ", err)
	}

	// get the replay info
	replay := getRewriteInfo(bugType, path, index)

	err = writeFile(path, index, bugTypeDescription, bugPos, bugElemType, code,
		replay, progInfo)

	copyTrace(path, index)
	copyRewrite(path, index)

	return err

}

func readAnalysisResults(path string, index int) (string, map[int][]string, map[int]string, error) {
	file, err := os.ReadFile(path + "results_machine.log")
	if err != nil {
		return "", nil, nil, err
	}

	lines := strings.Split(string(file), "\n")

	index-- // the index is 1-based

	if index >= len(lines) {
		return "", nil, nil, errors.New("index out of range")
	}

	bugStr := string(lines[index])
	bugFields := strings.Split(bugStr, ",")
	bugType := bugFields[0]

	bugPos := make(map[int][]string)
	bugElemType := make(map[int]string)

	for i := 1; i < len(bugFields); i++ {
		bugElems := strings.Split(bugFields[i], ";")
		if len(bugElems) == 0 {
			continue
		}

		bugPos[i] = make([]string, 0)

		for j, elem := range bugElems {
			fields := strings.Split(elem, ":")

			if fields[0] != "T" {
				continue
			}

			if j == 0 {
				bugElemType[i] = getBugElementType(fields[4])
			}

			file := fields[5]
			line := fields[6]
			pos := file + ":" + line
			bugPos[i] = append(bugPos[i], pos)
		}
	}

	return bugType, bugPos, bugElemType, nil

}

func writeFile(path string, index int, description map[string]string,
	positions map[int][]string, bugElemType map[int]string, code map[int][]string,
	replay map[string]string, progInfo map[string]string) error {
	// if in path, the folder "bugs" does not exist, create it
	if _, err := os.Stat(path + "bugs"); os.IsNotExist(err) {
		err := os.Mkdir(path+"bugs", 0755)
		if err != nil {
			return err
		}
	}

	folderName := path + "bugs/bug_" + fmt.Sprint(index)
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err := os.Mkdir(folderName, 0755)
		if err != nil {
			return err
		}
	}

	// create the file
	file, err := os.Create(folderName + "/README.md")
	if err != nil {
		return err
	}

	res := ""

	// write the bug type description
	res += "# " + description["crit"] + ": " + description["name"] + "\n\n"
	res += description["explanation"] + "\n\n"
	res += "## Minimal Example\n"
	res += "The following code is a minimal example for the bug type. It is not the code where the bug was found.\n\n```go\n"
	res += description["example"] + "\n```\n\n"

	// write the positions of the bug
	res += "## Test/Program\n"
	res += "The bug was found in the following test/program:\n\n"
	res += "Test: " + progInfo["name"] + "\n\n"
	res += "File: " + progInfo["file"] + "\n\n"

	res += "## Commands\n"
	res += "The following commands can be used to run and record the program:\n\n"

	res += "```bash\n"
	res += getProgInfo(progInfo, "inserterRecord") + "\n"
	res += getProgInfo(progInfo, "run") + "\n"
	res += getProgInfo(progInfo, "remover") + "\n"
	res += "```\n\n"

	res += "The following command can be used to replay the bug:\n\n"
	res += "Be aware, that the folder rewritten_trace_" + fmt.Sprint(index) + " must exist "
	res += "and contain the rewritten trace. It must be in the same folder as the recorded trace.\n\n"

	res += "```bash\n"
	res += getProgInfo(progInfo, "inserterReplay") + "\n"
	res += getProgInfo(progInfo, "run") + "\n"
	res += getProgInfo(progInfo, "remover") + "\n"
	res += "```\n\n"

	// write the code of the bug elements
	res += "## Bug Elements\n"
	res += "The full trace of the recording can be found in the `advocateTrace` folder.\n\n"
	res += "The bug elements are located at the following positions:\n\n"

	for key, _ := range positions {
		res += "###  "
		res += bugElemType[key] + "\n"

		for j, pos := range positions[key] {
			code := code[key][j]
			res += pos + "\n\n"
			res += code + "\n\n"
		}
	}

	// write the info about the replay, if possible including the command to read the bug
	res += "## Replay the bug\n"
	res += replay["description"] + "\n\n"

	replayPossible := replay["replaySuc"] != "was not possible"

	if replayPossible {
		res += "The rewritten trace can be found in the `rewritten_trace` folder.\n\n"
	}

	res += "Replaying " + replay["replaySuc"] + ".\n"
	if replayPossible {
		res += "It exited with the following code: "
		res += replay["exitCode"] + "\n"
		res += replay["exitCodeExplanation"] + "\n\n"
	}

	_, err = file.WriteString(res)
	return err

}
