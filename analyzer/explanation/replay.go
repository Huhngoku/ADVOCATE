package explanation

import (
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

	if rewPos == "Actual" {
		res["description"] += "The bug is an actual bug. Therefore to rewrite is possibel."
	} else if rewPos == "Potential" {
		res["description"] += "The bug is a potential bug.\n"
		res["description"] += "The analyzer has tries to rewrite the trace in such a way, "
		res["description"] += "that the bug will be triggered when replaying the trace."
	} else if rewPos == "LeakPos" {
		res["description"] += "The analyzer found a leak in the recorded trace.\n"
		res["description"] += "The analyzer found a way to resolve the leak, meaning the "
		res["description"] += "leak should not reappear in the rewritten trace."
		res["exitCode"], res["exitCodeExplanation"], res["replaySuc"] = getReplayInfo(path, index)
	} else if rewPos == "Leak" {
		res["description"] += "The analyzer found a leak in the recorded trace.\n"
		res["description"] += "The analyzer could not find a way to resolve the leak."
		res["description"] += "No rewritten trace was created. This does not need to mean, "
		res["description"] += "that the leak can not be resolved, especially because the "
		res["description"] += "analyzer is only aware of executed operations."
		res["exitCode"], res["exitCodeExplanation"], res["replaySuc"] = getReplayInfo(path, index)
	}

	return res

}

func getReplayInfo(path string, index int) (string, string, string) {
	if _, err := os.Stat(path + "output.txt"); os.IsNotExist(err) {
		return "", "No replay info available. Output.txt does not exist.", "false"
	}

	// read the output file
	content, err := os.ReadFile(path + "output.txt")
	if err != nil {
		return "", "No replay info available.Could not read output.txt file", "false"
	}

	// find all line, that either start with "Reading trace from "
	// or with "Exit Replay with code"
	traceNumbers := make([]int, 0)
	linesWithCode := make([]string, 0)
	lines := strings.Split(string(content), "\n")

	prefixTrace := "Reading trace from rewritten_trace_"
	prefixCode := "Exit Replay with code"

	for _, line := range lines {
		if strings.HasPrefix(line, prefixTrace) {
			line = strings.TrimPrefix(line, prefixTrace)
			line = strings.TrimSpace(line)
			traceNumber, err := strconv.Atoi(line)
			if err != nil {
				return "", "Invalid format in output.txt. Could not convert trace number to int", "failed"
			}
			traceNumbers = append(traceNumbers, traceNumber)
		}
		if strings.HasPrefix(line, prefixCode) {
			line = strings.TrimPrefix(line, prefixCode)
			line = strings.TrimSpace(line)
			line = strings.Split(line, " ")[0]
			line = strings.TrimSpace(line)
			linesWithCode = append(linesWithCode, line)
		}
	}

	if len(traceNumbers) != len(linesWithCode) {
		return "", "Invalid format in output.txt. Number of trace numbers does not match number of exit codes.", "failed"
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
		return "", "No replay info available. Could not find trace number in output.txt", "failed"
	}

	exitCode := linesWithCode[foundIndex]
	exitCodeInt, err := strconv.Atoi(exitCode)
	if err != nil {
		return "", "Invalid format in output.txt. Could not convert exit code to int", "failed"
	}

	replaySuc := "failed"
	if exitCodeInt >= 30 {
		replaySuc = "was successful"
	}

	return exitCode, exitCodeExplanation[exitCode], replaySuc
}
