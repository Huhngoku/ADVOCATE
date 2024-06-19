package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	folderName := flag.String("f", "", "path to the folder")
	flag.Parse()
	if *folderName == "" {
		fmt.Fprintln(os.Stderr, "Usage generateStatistics -f <folder>")
		os.Exit(1)
	}
	codes := []string{
		"A1",
		"A2",
		"A3",
		"A4",
		"A5",
		"P1",
		"P2",
		"P3",
		"L1",
		"L2",
		"L3",
		"L4",
		"L5",
		"L6",
		"L7",
		"L8",
		"L9",
		"L0",
	}
	caseReports := make([]caseReport, 0)
	for _, code := range codes {
		report := getCaseReportForCode(code, *folderName)
		caseReports = append(caseReports, report)
	}
	predictedBugCounts, err := getPredictedBugCounts(*folderName)
	if err != nil {
		fmt.Println(err)
	}
	for i, report := range caseReports {
		caseReports[i].occurenceCount = predictedBugCounts[report.caseCode]
	}
	for _, report := range caseReports {
		report.prettyPrint()
	}
}

type caseReport struct {
	caseCode        string
	occurenceCount  int
	actualExitCodes []string
}

func (c *caseReport) prettyPrint() {
	fmt.Printf("%s:%d:", c.caseCode, c.occurenceCount)
	for _, code := range c.actualExitCodes {
		fmt.Printf("%s,", code)
	}
	fmt.Println("")
}

func getCaseReportForCode(code string, folder string) caseReport {
	toRet := caseReport{
		caseCode:        code,
		occurenceCount:  0,
		actualExitCodes: make([]string, 0),
	}
	files, err := getFiles(folder, "rewrite_info.log")
	if err != nil {
		fmt.Println(err)
	}
	filteredFiles := make([]string, 0)
	for _, file := range files {
		_, caseCode, _, err := parseRewriteInfoFile(file)
		if err != nil {
			fmt.Println(err)
			return toRet
		}
		if caseCode == code {
			filteredFiles = append(filteredFiles, file)
		}
	}
	for _, file := range filteredFiles {
		dir := filepath.Dir(file)
		reorderFiles, err := getFiles(dir, "reorder_output.txt")
		if err != nil {
			fmt.Println(err)
			return toRet
		}
		for _, reorderFile := range reorderFiles {
			file, err := os.Open(reorderFile)
			if err != nil {
				fmt.Println(err)
				return toRet
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			fileContent := ""
			for scanner.Scan() {
				line := scanner.Text()
				fileContent += line
			}
			actualCode, err := extractActualCode(fileContent)
			if err != nil {
				continue
			}
			code := strconv.Itoa(actualCode)
			toRet.actualExitCodes = append(toRet.actualExitCodes, code)
		}
	}
	return toRet
}

func getBugCodes(filePath string) []string {
	bugCodes := make([]string, 0)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		idx := strings.Index(line, ",")
		if idx != -1 {
			bugcode := line[:idx]
			bugCodes = append(bugCodes, bugcode)
		} else {
			fmt.Println("no comma found in line -> invalid format")
		}
	}
	return bugCodes
}

func getFiles(folderPath string, fileName string) ([]string, error) {
	var files []string
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Base(path) == fileName {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func getPredictedBugCounts(folderPath string) (map[string]int, error) {
	codes := []string{
		"A1",
		"A2",
		"A3",
		"A4",
		"A5",
		"P1",
		"P2",
		"P3",
		"L1",
		"L2",
		"L3",
		"L4",
		"L5",
		"L6",
		"L7",
		"L8",
		"L9",
		"L0",
	}
	predictedCodes := make(map[string]int)
	for _, code := range codes {
		predictedCodes[code] = 0
	}

	files, err := getFiles(folderPath, "results_machine.log")
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		bugCodes := getBugCodes(file)
		for _, code := range bugCodes {
			_, ok := predictedCodes[code]
			if ok {
				predictedCodes[code]++
			}
		}
	}

	return predictedCodes, nil
}

func getPredictedExitCodesCounts(folderPath string) (map[string]int, error) {
	files, err := getFiles(folderPath, "rewrite_info.log")
	predictedCodes := make(map[string]int)
	exitCodes := []string{
		"0",
		"10",
		"11",
		"12",
		"13",
		"20",
		"21",
		"22",
		"23",
		"24",
		"30",
		"31",
		"32",
	}
	for _, code := range exitCodes {
		predictedCodes[code] = 0
	}
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		_, _, expectedExitCode, err := parseRewriteInfoFile(file)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		predictedCodes[expectedExitCode]++
	}
	return predictedCodes, nil
}

func parseRewriteInfoFile(filePath string) (string, string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return "", "", "", fmt.Errorf("no data in file")
	}
	line := scanner.Text()
	parts := strings.Split(line, "#")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("expected 3 parts, got %d", len(parts))
	}
	return parts[0], parts[1], parts[2], nil
}

func getActualExitCodes(filePath string) (map[string]int, error) {
	actualCodes := make(map[string]int)
	exitCodes := []string{
		"0",
		"10",
		"11",
		"12",
		"13",
		"20",
		"21",
		"22",
		"23",
		"24",
		"30",
		"31",
		"32",
	}
	for _, code := range exitCodes {
		actualCodes[code] = 0
	}
	files, err := getFiles(filePath, "reorder_output.txt")
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		file, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		fileContent := ""
		for scanner.Scan() {
			line := scanner.Text()
			fileContent += line
		}
		actualCode, err := extractActualCode(fileContent)
		if err != nil {
			continue
		}
		code := strconv.Itoa(actualCode)
		actualCodes[code]++
	}
	return actualCodes, nil
}

func extractActualCode(s string) (int, error) {
	re := regexp.MustCompile(`Exit Replay with code  (\d+)`)
	match := re.FindStringSubmatch(s)
	if match == nil {
		return -1, fmt.Errorf("no exit code found")
	}
	code, err := strconv.Atoi(match[1])
	if err != nil {
		return -1, err
	}
	return code, nil
}
