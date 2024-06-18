package explanation

import (
	"errors"
	"os"
	"strings"
)

func readCommandInfo(path string) (map[string]string, []string, error) {
	res := make(map[string]string)
	res["FileName"] = ""
	res["TestName"] = ""

	commands := make([]string, 0)

	file, err := os.ReadFile(path + "commands.txt")
	if err != nil {
		return nil, nil, err
	}

	lines := strings.Split(string(file), "\n")
	if len(lines) < 2 {
		return nil, nil, errors.New("Invalid format in commands.txt. Not enough lines")
	}
	res["FileName"] = lines[0]
	res["TestName"] = lines[1]

	for i := 2; i < len(lines); i++ {
		commands = append(commands, lines[i])
	}

	return res, commands, nil
}
