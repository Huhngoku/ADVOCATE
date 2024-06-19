package explanation

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func readProgInfo(path string, index int) (map[string]string, error) {
	res := make(map[string]string)

	file, err := os.ReadFile(path + "/advocateCommand.log")
	if err != nil {
		return res, err
	}

	lines := strings.Split(string(file), "\n")

	if len(lines) < 3 {
		return res, errors.New("advocateCommand file is too short")
	}

	res["file"] = lines[0]
	res["name"] = lines[1]

	for i := 2; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}

		if strings.Contains(lines[i], "unitTestOverheadInserter") {
			if strings.Contains(lines[i], "-r true") {
				line := lines[i][:strings.LastIndex(lines[i], " ")]
				res["inserterReplay"] = line + " " + strconv.Itoa(index)
			} else {
				res["inserterRecord"] = lines[i]
			}
		} else if strings.Contains(lines[i], "unitTestOverheadRemover") {
			res["remover"] = lines[i]
		} else if strings.Contains(lines[i], "-run") {
			res["run"] = lines[i]
		}
	}

	return res, nil
}

func getProgInfo(info map[string]string, key string) string {
	if _, ok := info[key]; !ok {
		return "Failed to read command for " + key
	}

	if info[key] == "" {
		return "Failed to read command for " + key
	}

	return info[key]
}
