package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	programName := flag.String("n", "", "Name of the program that was analyzed")
	pathToProgram := flag.String("p", "", "Path to the program that was analyzed")
	pathToStats := flag.String("s", "", "Path to the position where the stats file should be created")
	pathToTrace := flag.String("t", "", "Path to the trace folder")
	pathToResult := flag.String("r", "", "Path to the readable result file")
	pathToTime := flag.String("d", "", "Path to a file with the time durations")
	flag.Parse()

	programNameStr, pathToStatsStr, err := checkFlags(*programName, *pathToStats)
	if err != nil {
		println(err.Error())
	}

	statsPath := pathToStatsStr + "/" + programNameStr + "_stats.md"
	err = createStatsFile(statsPath, programNameStr)
	if err != nil {
		panic(err)
	}

	err = parseProgram(*pathToProgram, statsPath)
	if err != nil {
		println(err.Error())
	}

	err = parseTrace(*pathToTrace, statsPath)
	if err != nil {
		println(err.Error())
	}

	err = writeTimes(*pathToTime, statsPath)
	if err != nil {
		println(err.Error())
	}

	err = copyResults(*pathToResult, statsPath)
	if err != nil {
		println(err.Error())
	}
}

func checkFlags(programName, pathToStats string) (string, string, error) {
	var err error

	if programName == "" {
		programName = time.Now().String()
		err = errors.New("No program name provided. Using current time as program name")
	}

	if pathToStats == "" {
		pathToStats = "~/Desktop/"
		err = errors.New("No path to stats file provided. Using Desktop as path")
	}

	return programName, pathToStats, err

}

func createStatsFile(statsPath string, programName string) error {
	// delete file if it exists
	if _, err := os.Stat(statsPath); err == nil {
		err = os.Remove(statsPath)
		if err != nil {
			return err
		}
	}

	_, err := os.Create(statsPath)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(statsPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString("# " + programName + " Stats\n\n")

	return nil
}

// ========================= Program =========================

func parseProgram(programPath string, statsPath string) error {
	if programPath == "" {
		writeProgramStats(statsPath, map[string]int{"numberFiles": 0, "numberLines": 0, "numberNonEmptyLines": 0})
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

	return writeProgramStats(statsPath, res)
}

func writeProgramStats(statsPath string, stats map[string]int) error {
	file, err := os.OpenFile(statsPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("## Program\n")
	if stats["numberFiles"] == 0 {
		file.WriteString("No program found\n")
		return nil
	}

	file.WriteString("| Info | Value |\n| - | - |\n")
	file.WriteString("| Number of go files | " + strconv.Itoa(stats["numberFiles"]) + " |\n")
	file.WriteString("| Number of lines | " + strconv.Itoa(stats["numberLines"]) + " |\n")
	file.WriteString("| Number of non-empty lines | " + strconv.Itoa(stats["numberNonEmptyLines"]) + " |\n\n\n")
	return nil
}

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

// ========================= Trace =========================

func parseTrace(tracePath string, statsPath string) error {
	if tracePath == "" {
		writeTraceStats(statsPath, map[string]int{"numberFiles": 0, "numberLines": 0, "numberNonEmptyLines": 0})
		return errors.New("Please provide a path to the trace folder. Set with -t [file]")
	}

	res := map[string]int{
		"numberRoutines": 0,
		"numberOfSpawns": 0,

		"numberAtomics":          0,
		"numberAtomicOperations": 0,

		"numberChannels":          0,
		"numberChannelOperations": 0,

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

		return parseTraceFile(path, res, known)
	})

	if err != nil {
		println(err.Error())
	}

	return writeTraceStats(statsPath, res)
}

func parseTraceFile(tracePath string, stats map[string]int, known map[string][]string) error {
	// open the file
	file, err := os.Open(tracePath)
	if err != nil {
		return err
	}

	stats["numberRoutines"]++

	// read the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
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
				if !contains(known["channel"], fields[3]) {
					stats["numberChannels"]++
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

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func writeTraceStats(statsPath string, stats map[string]int) error {
	file, err := os.OpenFile(statsPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("## Trace\n")
	if stats["numberRoutines"] == 0 {
		file.WriteString("No trace found\n")
		return nil
	}

	file.WriteString("| Info | Value |\n| - | - |\n")
	file.WriteString("| Number of routines | " + strconv.Itoa(stats["numberRoutines"]) + " |\n")
	file.WriteString("| Number of spawns | " + strconv.Itoa(stats["numberOfSpawns"]) + " |\n")
	file.WriteString("| Number of atomics | " + strconv.Itoa(stats["numberAtomics"]) + " |\n")
	file.WriteString("| Number of atomic operations | " + strconv.Itoa(stats["numberAtomicOperations"]) + " |\n")
	file.WriteString("| Number of channels | " + strconv.Itoa(stats["numberChannels"]) + " |\n")
	file.WriteString("| Number of channel operations | " + strconv.Itoa(stats["numberChannelOperations"]) + " |\n")
	file.WriteString("| Number of selects | " + strconv.Itoa(stats["numberSelects"]) + " |\n")
	file.WriteString("| Number of select cases | " + strconv.Itoa(stats["numberSelectCases"]) + " |\n")
	file.WriteString("| Number of select channel operations | " + strconv.Itoa(stats["numberSelectChanOps"]) + " |\n")
	file.WriteString("| Number of select default operations | " + strconv.Itoa(stats["numberSelectDefaultOps"]) + " |\n")
	file.WriteString("| Number of mutexes | " + strconv.Itoa(stats["numberMutexes"]) + " |\n")
	file.WriteString("| Number of mutex operations | " + strconv.Itoa(stats["numberMutexOperations"]) + " |\n")
	file.WriteString("| Number of wait groups | " + strconv.Itoa(stats["numberWaitGroups"]) + " |\n")
	file.WriteString("| Number of wait group operations | " + strconv.Itoa(stats["numberWaitGroupOperations"]) + " |\n")
	file.WriteString("| Number of cond vars | " + strconv.Itoa(stats["numberCondVars"]) + " |\n")
	file.WriteString("| Number of cond var operations | " + strconv.Itoa(stats["numberCondVarOperations"]) + " |\n")
	file.WriteString("| Number of once | " + strconv.Itoa(stats["numberOnce"]) + "| \n")
	file.WriteString("| Number of once operations | " + strconv.Itoa(stats["numberOnceOperations"]) + " |\n\n\n")

	return nil
}

// ========================= Times =========================
func writeTimes(pathToTime string, statsPath string) error {
	fileStats, err := os.OpenFile(statsPath, os.O_APPEND|os.O_WRONLY, 0644)
	fileStats.WriteString("## Times\n")
	if err != nil {
		return err
	}

	if pathToTime == "" {
		fileStats.WriteString("No time file provided\n\n\n")
		return errors.New("No time file provided. Set with -d [file]")
	}

	fileTime, err := os.Open(pathToTime)
	if err != nil {
		fileStats.WriteString("Could not read time\n\n\n")
		return err
	}

	scanner := bufio.NewScanner(fileTime)
	for scanner.Scan() {
		line := scanner.Text()
		times := strings.Split(line, ",")
		println(len(times))
		if len(times) != 5 && len(times) != 4 {
			fileStats.WriteString("Invalid time file\n")
			fileStats.WriteString(line)
			return errors.New("Invalid time file")
		}

		timeOriginal, _ := strconv.ParseFloat(times[0], 64)
		timeAdvocate, _ := strconv.ParseFloat(times[1], 64)
		timeReplay, _ := strconv.ParseFloat(times[2], 64)

		overheadAdvocate := max(0, (timeAdvocate-timeOriginal)/timeOriginal*100)
		overheadReplay := max(0, (timeReplay-timeOriginal)/timeOriginal*100)

		overheadAdvocateStr := fmt.Sprintf("%f", overheadAdvocate) + " %"
		overheadReplayStr := fmt.Sprintf("%f", overheadReplay) + " %"

		fileStats.WriteString("| Info | Value |\n| - | - |\n")
		fileStats.WriteString("| Time for run without ADVOCATE | " + times[0] + " s |\n")
		fileStats.WriteString("| Time for run with ADVOCATE | " + times[1] + " s |\n")
		fileStats.WriteString("| Overhead of ADVOCATE | " + overheadAdvocateStr + " |\n")
		if len(times) == 4 {
			fileStats.WriteString("| Analysis | " + times[2] + " s |\n\n\n")
		} else {
			fileStats.WriteString("| Replay without changes | " + times[2] + " s |\n")
			fileStats.WriteString("| Overhead of Replay | " + overheadReplayStr + " s |\n")
			fileStats.WriteString("| Analysis | " + times[3] + " s |\n\n\n")
		}
		return nil
	}

	return nil
}

// ========================= Results =========================

func copyResults(pathToResult string, statsPath string) error {
	fileStats, err := os.OpenFile(statsPath, os.O_APPEND|os.O_WRONLY, 0644)
	fileStats.WriteString("## Results\n")
	if err != nil {
		return err
	}
	defer fileStats.Close()

	if pathToResult == "" {
		fileStats.WriteString("No result file provided\n\n\n")
		return errors.New("No result file provided. Set with -r [file]")
	}

	fileResult, err := os.Open(pathToResult)
	if err != nil {
		return err
	}
	defer fileResult.Close()

	scanner := bufio.NewScanner(fileResult)
	for scanner.Scan() {
		fileStats.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
