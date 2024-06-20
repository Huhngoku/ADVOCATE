package explanation

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getRewriteInfo(bugType string, path string, index int) map[string]string {
	res := make(map[string]string)

	rewPos := rewriteType[bugType]

	res["description"] = ""
	res["exitCode"] = ""
	res["exitCodeExplanation"] = ""
	res["replaySuc"] = "was not possible"

	var err error

	if rewPos == "Actual" {
		res["description"] += "The bug is an actual bug. Therefore no rewrite is possibel."
	} else if rewPos == "Possible" {
		res["description"] += "The bug is a potential bug.\n"
		res["description"] += "The analyzer has tries to rewrite the trace in such a way, "
		res["description"] += "that the bug will be triggered when replaying the trace."
		res["exitCode"], res["exitCodeExplanation"], res["replaySuc"], err = getReplayInfo(path, index)
	} else if rewPos == "LeakPos" {
		res["description"] += "The analyzer found a leak in the recorded trace.\n"
		res["description"] += "The analyzer found a way to resolve the leak, meaning the "
		res["description"] += "leak should not reappear in the rewritten trace."
		res["exitCode"], res["exitCodeExplanation"], res["replaySuc"], err = getReplayInfo(path, index)
	} else if rewPos == "Leak" {
		res["description"] += "The analyzer found a leak in the recorded trace.\n"
		res["description"] += "The analyzer could not find a way to resolve the leak."
		res["description"] += "No rewritten trace was created. This does not need to mean, "
		res["description"] += "that the leak can not be resolved, especially because the "
		res["description"] += "analyzer is only aware of executed operations."
	}

	if err != nil {
		fmt.Println("Error getting replay info: ", err)
	}

	return res

}

func getReplayInfo(path string, index int) (string, string, string, error) {
	if _, err := os.Stat(path + "output.txt"); os.IsNotExist(err) {
		res := "No replay info available. Output.txt does not exist."
		return "", res, "false", errors.New(res)
	}

	// read the output file
	content, err := os.ReadFile(path + "output.txt")
	if err != nil {
		res := "No replay info available. Could not read output.txt file"
		return "", res, "false", errors.New(res)
	}

	// find all line, that either start with "Reading trace from "
	// or with "Exit Replay with code"
	traceNumbers := make([]int, 0)
	linesWithCode := make([]string, 0)
	lines := strings.Split(string(content), "\n")

	prefixTrace := "Reading trace from rewritten_trace_"
	prefixCode := "Exit Replay with code"
	prefixPanic := "panic: "

	for _, line := range lines {
		if strings.HasPrefix(line, prefixTrace) {
			line = strings.TrimPrefix(line, prefixTrace)
			line = strings.TrimSpace(line)
			traceNumber, err := strconv.Atoi(line)
			if err != nil {
				res := "Invalid format in output.txt. Could not convert trace number to int"
				return "", res, "failed", errors.New(res)
			}
			traceNumbers = append(traceNumbers, traceNumber)
		} else if strings.HasPrefix(line, prefixCode) {
			line = strings.TrimPrefix(line, prefixCode)
			line = strings.TrimSpace(line)
			line = strings.Split(line, " ")[0]
			line = strings.TrimSpace(line)
			linesWithCode = append(linesWithCode, line)
		} else if strings.HasPrefix(line, prefixPanic) {
			linesWithCode = append(linesWithCode, line)
		}
	}

	if len(traceNumbers) != len(linesWithCode) {
		res := "Invalid format in output.txt. Number of trace numbers does not match number of exit codes."
		return "", res, "failed", errors.New(res)
	}

	// find the line, that corresponds to the index
	foundIndex := -1
	for i, traceNumber := range traceNumbers {
		if traceNumber == index {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		res := "No replay info available. Could not find trace number in output.txt"
		return "", res, "failed", errors.New(res)
	}

	exitCode := linesWithCode[foundIndex]
	replaySuc := "failed"
	if !strings.HasPrefix(exitCode, prefixPanic) {
		exitCodeInt, err := strconv.Atoi(exitCode)
		if err != nil {
			res := "Invalid format in output.txt. Could not convert exit code to int"
			return "", res, "failed", errors.New(res)
		}
		if exitCodeInt >= 30 {
			replaySuc = "was successful"
		}
	} else {
		replaySuc = "panicked"
	}

	return exitCode, exitCodeExplanation[exitCode], replaySuc, nil
}
