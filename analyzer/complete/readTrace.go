package complete

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// TODO: remove duplicates
func getTraceElements(resultFolderPath string) (map[string][]int, error) {
	res := make(map[string][]int)

	err := filepath.Walk(resultFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fileName := filepath.Base(path)

		if info.IsDir() && strings.HasPrefix(fileName, "rewritten_trace") {
			return filepath.SkipDir
		}

		if strings.HasPrefix(fileName, "trace_") && strings.HasSuffix(fileName, ".log") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			elems := strings.Split(string(content), ";")

			for _, elem := range elems {
				field := strings.Split(elem, ",")
				if len(field) == 0 {
					continue
				}

				if field[0] == "A" || field[0] == "X" {
					continue
				}

				pos := strings.Split(field[len(field)-1], ":")
				if len(pos) != 2 {
					continue
				}

				file := pos[0]
				line, err := strconv.Atoi(pos[1])
				if err != nil {
					continue
				}

				if _, ok := res[file]; !ok {
					res[file] = make([]int, 0)
				}
				res[file] = append(res[file], line)
			}
		}
		return nil
	})

	println("AAAAAAAAA")

	return res, err
}
