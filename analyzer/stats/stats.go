package stats

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Create(pathToProgram *string, pathToTrace *string) {

	if pathToProgram == nil && pathToTrace == nil {
		panic("Please provide at least one of the following flags: -t [file] or -P [file]")
	}

	pathToStats := ""
	if pathToTrace != nil {
		pathToStats = filepath.Dir(*pathToTrace)
	} else {
		pathToStats = *pathToProgram
	}

	pathToCSV := pathToStats + "/stats.csv"
	err := createFile(pathToCSV)
	if err != nil {
		panic(err)
	}

	if pathToProgram != nil {
		parseProgramToCSV(*pathToProgram, pathToCSV)
	}

	if pathToTrace != nil {
		parseTraceToCSV(*pathToTrace, pathToCSV)
	}
}

func createFile(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	return nil
}

// ========================= Program =========================

func parseProgramFile(filePath string) (map[string]int, error) {
	res := make(map[string]int)
	res["numberLines"] = 0
	res["numberNonEmptyLines"] = 0

	// open the file
	file, err := os.Open(filePath)
	if err != nil {
		return res, err
	}
	defer file.Close()

	// read the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		res["numberLines"]++
		if text != "" && text != "\n" && !strings.HasPrefix(text, "//") {
			res["numberNonEmptyLines"]++
		}
	}

	if err := scanner.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func parseProgramToCSV(programPath string, csvPath string) error {
	if programPath == "" {
		writeProgramStatsToCSV(csvPath, map[string]int{"numberFiles": 0, "numberLines": 0, "numberNonEmptyLines": 0})
		return errors.New("Please provide a path to the program that was analyzed. Set with -p [file]")
	}

	res := make(map[string]int)
	res["numberFiles"] = 0
	res["numberLines"] = 0
	res["numberNonEmptyLines"] = 0

	err := filepath.Walk(programPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".go" {
			resFile, err := parseProgramFile(path)
			if err != nil {
				return err
			}

			res["numberFiles"]++
			res["numberLines"] += resFile["numberLines"]
			res["numberNonEmptyLines"] += resFile["numberNonEmptyLines"]
		}

		return nil
	})
	if err != nil {
		return err
	}

	return writeProgramStatsToCSV(csvPath, res)
}

func writeProgramStatsToCSV(statsPath string, stats map[string]int) error {
	file, err := os.OpenFile(statsPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if stats["numberFiles"] == 0 {
		file.WriteString("No program found\n")
		return nil
	}

	file.WriteString("Go files:" + strconv.Itoa(stats["numberFiles"]) + "\n")
	file.WriteString("Lines:" + strconv.Itoa(stats["numberLines"]) + "\n")
	file.WriteString("Non-empty lines:" + strconv.Itoa(stats["numberNonEmptyLines"]) + "\n")
	return nil
}

// ========================= Trace =========================

func parseTraceToCSV(tracePath string, statsPath string) error {
	if tracePath == "" {
		writeTraceCSV(statsPath, map[string]int{"numberFiles": 0, "numberLines": 0, "numberNonEmptyLines": 0})
		return errors.New("Please provide a path to the trace folder. Set with -t [file]")
	}

	res := map[string]int{
		"numberRoutines":         0,
		"numberNonEmptyRoutines": 0,
		"numberOfSpawns":         0,

		"numberAtomics":          0,
		"numberAtomicOperations": 0,

		"numberChannels":           0,
		"numberBufferedChannels":   0,
		"numberUnbufferedChannels": 0,
		"numberChannelOperations":  0,
		"numberBufferedOps":        0,
		"numberUnbufferedOps":      0,

		"numberSelects":          0,
		"numberSelectCases":      0,
		"numberSelectChanOps":    0, // number of executed channel operations in select
		"numberSelectDefaultOps": 0, // number of executed default operations in select

		"numberMutexes":         0,
		"numberMutexOperations": 0,

		"numberWaitGroups":          0,
		"numberWaitGroupOperations": 0,

		"numberCondVars":          0,
		"numberCondVarOperations": 0,

		"numberOnce":           0,
		"numberOnceOperations": 0,
	}

	known := map[string][]string{
		"atomic":    []string{},
		"channel":   []string{},
		"mutex":     []string{},
		"waitGroup": []string{},
		"condVar":   []string{},
		"once":      []string{},
	}

	err := filepath.Walk(tracePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".log" {
			return nil
		}

		if filepath.Base(path) == "times.log" {
			return nil
		}
		return parseTraceFile(path, res, known)
	})

	if err != nil {
		println(err.Error())
	}

	return writeTraceCSV(statsPath, res)
}

func writeTraceCSV(statsPath string, stats map[string]int) error {
	file, err := os.OpenFile(statsPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if stats["numberRoutines"] == 0 {
		file.WriteString("No trace found\n")
		return nil
	}
	file.WriteString("Routines:" + strconv.Itoa(stats["numberRoutines"]) + "\n")
	file.WriteString("Non-empty routines:" + strconv.Itoa(stats["numberNonEmptyRoutines"]) + "\n")
	file.WriteString("Spawns:" + strconv.Itoa(stats["numberOfSpawns"]) + "\n")
	file.WriteString("Atomics:" + strconv.Itoa(stats["numberAtomics"]) + "\n")
	file.WriteString("Channels:" + strconv.Itoa(stats["numberChannels"]) + "\n")
	file.WriteString("Unbuffered channels:" + strconv.Itoa(stats["numberUnbufferedChannels"]) + "\n")
	file.WriteString("Buffered channels:" + strconv.Itoa(stats["numberBufferedChannels"]) + "\n")
	file.WriteString("Selects:" + strconv.Itoa(stats["numberSelects"]) + "\n")
	file.WriteString("Mutexes:" + strconv.Itoa(stats["numberMutexes"]) + "\n")
	file.WriteString("Wait groups:" + strconv.Itoa(stats["numberWaitGroups"]) + "\n")
	file.WriteString("Cond vars:" + strconv.Itoa(stats["numberCondVars"]) + "\n")
	file.WriteString("Once:" + strconv.Itoa(stats["numberOnce"]) + "\n")

	file.WriteString("Atomic operations:" + strconv.Itoa(stats["numberAtomicOperations"]) + "\n")
	file.WriteString("Channel operations:" + strconv.Itoa(stats["numberChannelOperations"]) + "\n")
	file.WriteString("Channel unbuffered operations:" + strconv.Itoa(stats["numberUnbufferedOps"]) + "\n")
	file.WriteString("Channel buffered operations:" + strconv.Itoa(stats["numberBufferedOps"]) + "\n")
	file.WriteString("Select cases:" + strconv.Itoa(stats["numberSelectCases"]) + "\n")
	file.WriteString("Select channel operations:" + strconv.Itoa(stats["numberSelectChanOps"]) + "\n")
	file.WriteString("Select default operations:" + strconv.Itoa(stats["numberSelectDefaultOps"]) + "\n")
	file.WriteString("Mutex operations:" + strconv.Itoa(stats["numberMutexOperations"]) + "\n")
	file.WriteString("Wait group operations:" + strconv.Itoa(stats["numberWaitGroupOperations"]) + "\n")
	file.WriteString("Cond var operations:" + strconv.Itoa(stats["numberCondVarOperations"]) + "\n")
	file.WriteString("Once operations:" + strconv.Itoa(stats["numberOnceOperations"]) + "\n")

	return nil
}

func parseTraceFile(tracePath string, stats map[string]int, known map[string][]string) error {
	// open the file
	file, err := os.Open(tracePath)
	if err != nil {
		return err
	}

	routine, err := getRoutineFromFileName(filepath.Base(tracePath))
	if err == nil {
		stats["numberRoutines"] = max(stats["numberRoutines"], routine)
	}

	scanner := bufio.NewScanner(file)
	const maxCapacity = 3 * 1024 * 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	// read the file
	for scanner.Scan() {
		line := scanner.Text()

		if line != "" {
			stats["numberNonEmptyRoutines"]++
		}

		for _, elem := range strings.Split(line, ";") {
			fields := strings.Split(elem, ",")
			switch fields[0] {
			case "G":
				stats["numberOfSpawns"]++
			case "A":
				stats["numberAtomicOperations"]++
				if !contains(known["atomic"], fields[2]) {
					stats["numberAtomics"]++
					known["atomic"] = append(known["atomic"], fields[2])
				}
			case "C":
				stats["numberChannelOperations"]++
				if fields[7] == "0" {
					stats["numberUnbufferedOps"]++
				} else {
					stats["numberBufferedOps"]++
				}
				if !contains(known["channel"], fields[3]) {
					stats["numberChannels"]++
					if fields[7] == "0" {
						stats["numberUnbufferedChannels"]++
					} else {
						stats["numberBufferedChannels"]++
					}
					known["channel"] = append(known["channel"], fields[3])
				}
			case "S":
				stats["numberSelects"]++
				cases := strings.Split(fields[4], "~")
				stats["numberSelectCases"] += len(cases)
				if cases[len(cases)-1] == "D" {
					stats["numberSelectDefaultOps"]++
				} else {
					stats["numberSelectChanOps"] += len(cases)
				}
			case "M":
				stats["numberMutexOperations"]++
				if !contains(known["mutex"], fields[3]) {
					stats["numberMutexes"]++
					known["mutex"] = append(known["mutex"], fields[3])
				}
			case "W":
				stats["numberWaitGroupOperations"]++
				if !contains(known["waitGroup"], fields[3]) {
					stats["numberWaitGroups"]++
					known["waitGroup"] = append(known["waitGroup"], fields[3])
				}
			case "O":
				stats["numberOnceOperations"]++
				if !contains(known["once"], fields[3]) {
					stats["numberOnce"]++
					known["once"] = append(known["once"], fields[3])
				}
			case "N":
				stats["numberCondVarOperations"]++
				if !contains(known["condVar"], fields[3]) {
					stats["numberCondVars"]++
					known["condVar"] = append(known["condVar"], fields[3])
				}
			default:
				err = errors.New("Unknown trace element: " + fields[0])
			}
		}
	}
	return err
}

func getRoutineFromFileName(fileName string) (int, error) {
	// the file name is "trace_routineID.log"
	// remove the .log at the end
	fileName1 := strings.TrimSuffix(fileName, ".log")
	if fileName1 == fileName {
		return 0, errors.New("File name does not end with .log")
	}

	fileName2 := strings.TrimPrefix(fileName1, "trace_")
	if fileName2 == fileName1 {
		return 0, errors.New("File name does not start with trace_")
	}

	routine, err := strconv.Atoi(fileName2)
	if err != nil {
		return 0, err
	}

	return routine, nil
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
