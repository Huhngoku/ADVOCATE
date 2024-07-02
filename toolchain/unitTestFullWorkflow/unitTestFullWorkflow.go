package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	advocateRoot := flag.String("a", "", "Path to the advocate root")
	projectRoot := flag.String("f", "", "Path to the project root folder")
	moduleModeString := flag.String("m", "", "Module mode")
	testName := flag.String("t", "", "Name of the test")
	testFile := flag.String("tf", "", "Path to the test file")
	flag.Parse()
	moduleMode := false
	if *moduleModeString == "true" {
		moduleMode = true
	}
	//check if argument is missing else print usage
	if *advocateRoot == "" || *projectRoot == "" || *testName == "" || *testFile == "" {
		fmt.Println("Please provide all the arguments")
		fmt.Println("Usage: go run unitTestFullWorkflow -a <advocate folder> -f <project root folder> -t <test name> -tf <test file>")
		//print current arguments
		fmt.Println("Advocate Folder:", *advocateRoot)
		fmt.Println("Project Root Folder:", *projectRoot)
		fmt.Println("Test Name:", *testName)
		fmt.Println("Test File:", *testFile)
		//exit program with 1
		return
	}
	pathToAnalyzer := filepath.Join(*advocateRoot, "analyzer", "analyzer")
	pathToPatchedGoRuntime := filepath.Join(*advocateRoot, "go-patch", "bin", "go")
	pathToGoRoot := filepath.Join(*advocateRoot, "go-patch")
	pathToOverheadInserter := filepath.Join(*advocateRoot, "toolchain", "unitTestOverheadInserter", "unitTestOverheadInserter")
	pathToOverheadRemover := filepath.Join(*advocateRoot, "toolchain", "unitTestOverheadRemover", "unitTestOverheadRemover")

	//change directory to the project root folder
	err := os.Chdir(*projectRoot)
	if err != nil {
		fmt.Println("Error changing directory:", err)
		return
	}
	fmt.Println("In directory:", *projectRoot)
	//export GOROOT=$advocateRoot/go-patch
	err = os.Setenv("GOROOT", pathToGoRoot)
	if err != nil {
		fmt.Println("Error setting GOROOT:", err)
		return
	}
	fmt.Println("GOROOT set to:", pathToGoRoot)
	//create file called advocateCommand.log
	advocateCommandLog, err := os.Create("advocateCommand.log")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	//append filepath to command log
	_, err = advocateCommandLog.WriteString(*testFile + "\n")
	//append testname to commad log
	_, err = advocateCommandLog.WriteString(*testName + "\n")
	//append module mode to command log
	fmt.Println("Remove Overhead just in case")
	//execute overhead remover: Usage $pathToOverheadRemover -f $file -t $testName
	cmd := exec.Command(pathToOverheadRemover, "-f", *testFile, "-t", *testName)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error removing overhead:", err)
		return
	}
	//Add overhead: Usage $pathToOverheadInserter -f $file -t $testName
	fmt.Println("Add Overhead")
	cmd = exec.Command(pathToOverheadInserter, "-f", *testFile, "-t", *testName)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error adding overhead:", err)
		return
	}
	fmt.Println("Run test")
	//check for module mode
	if moduleMode {
		// echo "$pathToPatchedGoRuntime test -count=1 -run=$testName -mod=mod ./$package" >>advocateCommand.log
		// $pathToPatchedGoRuntime test -count=1 -run=$testName -mod=mod "./$package"
		_, err = advocateCommandLog.WriteString(pathToPatchedGoRuntime + " test -count=1 -run=" + *testName + " -mod=mod ./" + *testFile + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		cmd = exec.Command(pathToPatchedGoRuntime, "test", "-count=1", "-run="+*testName, "-mod=mod", "./"+*testFile)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Remove Overhead")
			cmd = exec.Command(pathToOverheadRemover, "-f", *testFile, "-t", *testName)
			fmt.Println("Error in running therefor overhead removed and full workflow stopped", err)
			return
		}
	} else {
		//do the same without -mod=mod
		_, err = advocateCommandLog.WriteString(pathToPatchedGoRuntime + " test -count=1 -run=" + *testName + " ./" + *testFile + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		cmd = exec.Command(pathToPatchedGoRuntime, "test", "-count=1", "-run="+*testName, "./"+*testFile)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Remove Overhead")
			cmd = exec.Command(pathToOverheadRemover, "-f", *testFile, "-t", *testName)
			fmt.Println("Error in running therefor overhead removed and full workflow stopped", err)
			return
		}
	}
	// echo "Remove Overhead"
	// echo "$pathToOverheadRemover -f $file -t $testName" >>advocateCommand.log
	// $pathToOverheadRemover -f $file -t $testName
	// echo "$pathToAnalyzer -t $dir/$package/advocateTrace" >>advocateCommand.log
	// $pathToAnalyzer -t "$dir/$package/advocateTrace"
	// rewritten_traces=$(find "./$package" -type d -name "rewritten_trace*")
	// for trace in $rewritten_traces; do
	// 	rtracenum=$(echo $trace | grep -o '[0-9]*$')
	// 	echo "Apply reorder overhead"
	// 	echo $pathToOverheadInserter -f $file -t $testName -r true -n "$rtracenum" >>advocateCommand.log
	// 	$pathToOverheadInserter -f $file -t $testName -r true -n "$rtracenum"  >>advocateCommand.log
	// 	if [ "$modulemode" == "true" ]; then
	// 		echo "$pathToPatchedGoRuntime test -count=1 -run=$testName -mod=mod ./$package" >>advocateCommand.log
	// 		$pathToPatchedGoRuntime test -count=1 -run=$testName -mod=mod "./$package" 2>&1 | tee -a "$trace/reorder_output.txt"
	// 	else
	// 		echo "$pathToPatchedGoRuntime test -count=1 -run=$testName ./$package" >>advocateCommand.log
	// 		$pathToPatchedGoRuntime test -count=1 -run=$testName "./$package" 2>&1 | tee -a "$trace/reorder_output.txt"
	// 	fi
	// 	echo "Remove reorder overhead"
	// 	echo "$pathToOverheadRemover -f $file -t $testName" >>advocateCommand.log
	// 	$pathToOverheadRemover -f $file -t $testName
	// done
	fmt.Println("Remove Overhead")
	_, err = advocateCommandLog.WriteString(pathToOverheadRemover + " -f " + *testFile + " -t " + *testName + "\n")
	cmd = exec.Command(pathToOverheadRemover, "-f", *testFile, "-t", *testName)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error removing overhead:", err)
		return
	}
	_, err 
}
