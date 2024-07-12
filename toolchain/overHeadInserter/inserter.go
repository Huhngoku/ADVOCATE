package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	fileName := flag.String("f", "", "path to the file")
	replayOverheadString := flag.String("r", "false", "replay overhead")
	replayNum := flag.String("n", "1", "replay number")
	flag.Parse()
	replayOverhead := false
	if *replayOverheadString == "true" {
		replayOverhead = true
	}
	if *fileName == "" {
		fmt.Fprintln(os.Stderr, "Please provide a file name")
		fmt.Fprintln(os.Stderr, "Usage: preambleInserter -f <file>")
		os.Exit(1)
	}
	if _, err := os.Stat(*fileName); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "File %s does not exist\n", *fileName)
		os.Exit(1)

	}
	exists, err := mainMethodExists(*fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if !exists {
		fmt.Fprintln(os.Stderr, "Main Method not found in file")
		os.Exit(1)
	}

	addOverhead(*fileName, replayOverhead, *replayNum)
}
func mainMethodExists(fileName string) (bool, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return false, err
	}
	defer file.Close()
	regexStr := "func main\\(\\) {"
	regex, err := regexp.Compile(regexStr)
	if err != nil {
		return false, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if regex.MatchString(line) {
			return true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}
func addOverhead(fileName string, replayOverhead bool, replayNumber string) {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	importAdded := false
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)

		if strings.Contains(line, "package main") {
			lines = append(lines, "import \"advocate\"")
			importAdded = true
		} else if strings.Contains(line, "import \"") && !importAdded {
			lines = append(lines, "import \"advocate\"")
			importAdded = true
		} else if strings.Contains(line, "import (") && !importAdded {
			lines = append(lines, "\t\"advocate\"")
			importAdded = true
		}

		if strings.Contains(line, "func main() {") {
			if replayOverhead {
				lines = append(lines, fmt.Sprintf(`	// ======= Preamble Start =======
	advocate.EnableReplay(%s, true)
	defer advocate.WaitForReplayFinish()
	// ======= Preamble End =======`, replayNumber))
			} else {
				lines = append(lines, `	// ======= Preamble Start =======
	advocate.InitTracing(0)
	defer advocate.Finish()
	// ======= Preamble End =======`)
			}
		}
	}

	file.Truncate(0)
	file.Seek(0, 0)
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()
}
