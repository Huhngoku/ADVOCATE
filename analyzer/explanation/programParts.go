package explanation

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

/*
 * Get the positions of the bug elements in the program
 * Args:
 *   traceElem1 (map[int]string): The trace elements of the bug
 * Returns:
 *   map[int][]string: Dict for the code snippets
 */
func getBugPositions(traceElems map[int][]string) (map[int][]string, error) {
	res := make(map[int][]string)

	for i, elem := range traceElems {
		for _, e := range elem {
			pos := strings.Split(e, ":")
			line, _ := strconv.Atoi(pos[1])
			code, err := getProgramCode(pos[0], line, true)
			if err != nil {
				res[i] = append(res[i], "")
			} else {
				res[i] = append(res[i], code)
			}
		}
	}

	return res, nil
}

func getProgramCode(file string, line int, numbers bool) (string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(content), "\n")
	if line < 0 || line >= len(lines) {
		return "", errors.New("line number out of range")
	}

	res := "```go\n"

	start := line - 10
	if start < 0 {
		start = 0
	} else {
		res += "...\n\n"
	}
	end := line + 10
	isEnd := false
	if end >= len(lines) {
		end = len(lines)
		isEnd = true
	}

	res += strings.Join(lines[start:end], "\n")

	if !isEnd {
		res += "\n\n..."
	}
	res += "\n```"

	if !numbers {
		return res, nil
	}

	// add line numbers
	resWithLines := ""
	for i, l := range strings.Split(res, "\n") {
		if i == 0 || i == len(strings.Split(res, "\n"))-1 {
			resWithLines += l + "\n"
			continue
		}
		resWithLines += strconv.Itoa(i+start-2) + " " + l
		if i+start-2 == line {
			resWithLines += "           // <-------\n"
		} else {
			resWithLines += "\n"
		}
	}

	return resWithLines, nil
}
