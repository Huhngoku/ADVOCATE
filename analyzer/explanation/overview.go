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

func CreateOverview(path string, index int) error {
	bugType, bugPos, bugElemType, err := ReadAnalysisResults(path, index)
	if err != nil {
		return err
	}

	// get the bug type description
	bugTypeDescription := getBugTypeDescription(bugType)

	code, err := getBugPositions(bugPos)

	err = writeFile(path, index, bugTypeDescription, bugPos, bugElemType, code)

	return err

}

func ReadAnalysisResults(path string, index int) (string, map[int][]string, map[int]string, error) {
	file, err := os.ReadFile(path + "results_machine.log")
	if err != nil {
		return "", nil, nil, err
	}

	lines := strings.Split(string(file), "\n")

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
	positions map[int][]string, bugElemType map[int]string, code map[int][]string) error {
	// if in path, the folder "bugs" does not exist, create it
	if _, err := os.Stat(path + "bugs"); os.IsNotExist(err) {
		err := os.Mkdir(path+"bugs", 0755)
		if err != nil {
			return err
		}
	}

	// create the file
	file, err := os.Create(path + "bugs/bug_" + fmt.Sprint(index) + ".md")
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
	// TODO: get the test/program name

	// write the code of the bug elements
	res += "## Bug Elements\n"
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

	// write the command to run the program
	res += "## Run the program\n"
	res += "To run the program, use the following command:\n\n"
	// TODO: get the command to run the program

	// write the info about the replay, if possible including the command to read the bug
	res += "## Replay the bug\n"
	// TODO: get the info / command about replay the bug

	_, err = file.WriteString(res)

	return err
}
