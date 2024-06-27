package complete

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getTraceElements(resultFolderPath string) (map[string][]int, error) {
	res := make(map[string][]int)

	// for each subfolder in resultFolderPath, not recursively
	subfolder, err := getSubfolders(resultFolderPath)
	if err != nil {
		println("Error in getting subfolders")
		return nil, err
	}

	for _, folder := range subfolder {
		importLine := -1
		overheadLine := -1
		overheadFile := ""
		resLocal := make(map[string][]int)

		err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				println("Error in walking trace: ", filepath.Clean(path))
				return err
			}

			fileName := filepath.Base(path)

			if info.IsDir() && fileName == "rewritten_trace" {
				return filepath.SkipDir
			}

			// read command line
			if fileName == "advocateCommand.log" {
				overheadFile, importLine, overheadLine, err = readCommandFile(path)
				if err != nil {
					println("Error in reading command: ", filepath.Clean(path))
					return err
				}
				return nil
			}

			// read trace file
			if !strings.HasPrefix(fileName, "trace_") || !strings.HasSuffix(fileName, ".log") {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				println("Error in reading trace: ", filepath.Clean(path))
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

				if _, ok := resLocal[file]; !ok {
					resLocal[file] = make([]int, 0)
				}
				resLocal[file] = append(resLocal[file], line)
			}

			return nil
		})

		if err != nil {
			println("Error in walking trace")
			return nil, err
		}

		// fix lines of trace with overhead
		for i, line := range resLocal[overheadFile] {
			if line >= importLine {
				resLocal[overheadFile][i]--
			}
			if line >= overheadLine {
				resLocal[overheadFile][i] -= 4
			}
		}

		// add resLocal into res
		for file, lines := range resLocal {
			if _, ok := res[file]; !ok {
				res[file] = make([]int, 0)
			}

			for _, line := range lines {
				if !contains(res[file], line) {
					res[file] = append(res[file], line)
				}
			}
		}
	}

	return res, nil
}

func getSubfolders(path string) ([]string, error) {
	var subfolders []string

	// Öffnen des Verzeichnisses
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Lesen des Verzeichnisinhalts
	files, err := dir.Readdir(-1) // -1 bedeutet, alle Einträge lesen
	if err != nil {
		return nil, err
	}

	// Filtern der Unterordner
	for _, file := range files {
		if file.IsDir() {
			subfolderPath := filepath.Join(path, file.Name())
			subfolders = append(subfolders, subfolderPath)
		}
	}

	return subfolders, nil
}

func readCommandFile(path string) (string, int, int, error) {
	importLine := -1
	overheadLine := -1
	overheadFile := ""

	// read the command file
	content, err := os.ReadFile(path)
	if err != nil {
		println("Error in reading command: ", filepath.Clean(path))
		return overheadFile, importLine, overheadLine, err
	}
	// find the line starting with Import added at line:
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if i == 0 {
			overheadFile = line
			continue
		}

		if strings.Contains(line, "Import added at line: ") {
			line := strings.TrimPrefix(line, "Import added at line: ")
			importLine, err = strconv.Atoi(line)
			if err != nil {
				println("Error in converting import line: ", line)
				return overheadFile, importLine, overheadLine, err
			}
		} else if strings.Contains(line, "Overhead added at line: ") {
			line := strings.TrimPrefix(line, "Overhead added at line: ")
			overheadLine, err = strconv.Atoi(line)
			if err != nil {
				println("Error in converting overhead line: ", line)
				return overheadFile, importLine, overheadLine, err
			}
		}
	}

	if importLine == -1 || overheadLine == -1 {
		println("Error in reading import or overhead line")
		return overheadFile, importLine, overheadLine, errors.New("Error in reading import or overhead line")
	}

	return overheadFile, importLine, overheadLine, nil
}
