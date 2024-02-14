package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const (
	red = "\033[31m"
	org = "\033[33m"
	grn = "\033[32m"
	end = "\033[0m"
)

var (
	home, _        = os.UserHomeDir()
	pathToAdvocate = home + "/Uni/HiWi/ADVOCATE"
	pathToOther    = home + "/Uni/HiWi/Other"
	analyzerPath   = pathToAdvocate + "/analyzer"
	runEval        = pathToAdvocate + "/eval"
	overviewPath   = pathToAdvocate + "/eval/createOverview"

	// name
	// pathToExec
	// pathToProg
	// basename of executable
	// params (if any)
	programs = [][]string{}

	resPath = runEval

	resSucc      = []string{}
	resFailure   = []string{}
	numberErrors = 0
)

func addConstructed() {
	for i := 1; i <= 44; i++ {
		programs = append(programs, []string{
			"constructed " + strconv.Itoa(i),
			pathToAdvocate + "/examples/constructed/",
			pathToAdvocate + "/examples/constructed/",
			"constructed", "-c", strconv.Itoa(i), "-t", "5"})
	}
}

func addGoBench() {
	const n = 18
	names := [n]string{
		"cockroach1055",
		"cockroach1462",
		"etcd6873",
		"etcd7443",
		"etcd7492",
		"etcd7902",
		"grpc1353",
		"grpc1460",
		"grpc1687",
		"istio16224",
		"Kubernetes1321",
		"Kubernetes6632",
		"Kubernetes10182",
		"Kubernetes26980",
		"Moby28462",
		"serving2137",
		"serving3068",
		"serving5865",
	}

	for i := 0; i < n; i++ {
		programs = append(programs,
			[]string{
				"Gobench " + strconv.Itoa(i+1) + ": " + names[i],
				pathToAdvocate + "/examples/GoBench/",
				pathToAdvocate + "/examples/GoBench/" + names[i] + "/",
				"gobench",
				"-c", strconv.Itoa(i + 1), "-t", "20"})
	}
}

func addMediumPrograms() {
	// bbolt
	programs = append(programs,
		[]string{
			"bbolt",
			pathToOther + "/examples/bbolt/cmd/bbolt/",
			pathToOther + "/examples/bbolt/",
			"bbolt",
			"bench"})

	// gocrawl
	programs = append(programs,
		[]string{
			"gocrawl",
			pathToOther + "/examples/gocrawl/",
			pathToOther + "/examples/gocrawl/",
			"gocrawl"})

	// htcat
	programs = append(programs,
		[]string{
			"htcat",
			pathToOther + "/examples/htcat/cmd/htcat/",
			pathToOther + "/examples/htcat/",
			"htcat"})

	// pgzip
	programs = append(programs,
		[]string{
			"pgzip",
			pathToOther + "/examples/pgzip/",
			pathToOther + "/examples/pgzip/",
			"pgzip"})

	// sorty
	programs = append(programs,
		[]string{
			"sorty",
			pathToOther + "/examples/sorty/",
			pathToOther + "/examples/sorty/",
			"sorty"})
}

func runExecs(pathToExec string, execName string, execArgs []string, progName string) {
	runExec(pathToExec, execName, execArgs, progName, "original", 0)
	runExec(pathToExec, execName, execArgs, progName, "advocate", 0)
	// runExec(pathToExec, execName, execArgs, progName, "replay", 0)
}

func runExec(pathToExec string, execName string, execArgs []string, progName string, variant string, repeat int) {
	commandStr := "./" + execName + "_" + variant
	cmd := exec.Command(commandStr, execArgs...)
	cmd.Dir = pathToExec
	log(EXEC, cmd.String()+" in "+pathToExec)

	out, runtimeOriginal, err := runCommand(cmd)

	if err != nil {
		log(ERROR, err.Error()+":\n"+out)
		if repeat < 5 && repeat != -1 {
			log(ERROR, variant+" of "+progName+" resulted in an error. Trying again...")
			runExec(pathToExec, execName, execArgs, progName, variant, repeat+1)
			return
		} else {
			log(ERROR, variant+" of "+progName+" failed 5 times. Skipping...")
		}
	}
	if out == "Timeout" {
		if repeat < 5 && repeat != -1 {
			log(TIMEOUT, variant+" of "+progName+" timed out. Trying again...")
			runExec(pathToExec, execName, execArgs, progName, variant, repeat+1)
			return
		} else {
			log(TIMEOUT, variant+" of "+progName+" failed 5 times. Skipping...")
		}
	}

	err = writeTime(runtimeOriginal, progName)
	if err != nil {
		log(ERROR, err.Error())
	}
}

/*
 * Run the command
 * Args:
 *   cmd (*exec.Cmd): The command to run
 * Returns:
 *   (string): stderr output
 *   (time.Duration): runtime
 *   (error): error
 */
func runCommand(cmd *exec.Cmd) (string, time.Duration, error) {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	start := time.Now()

	err := cmd.Run()

	duration := time.Since(start)

	if err != nil {
		return stderr.String(), duration, err
	}

	return "", duration, nil
}

func runAnalyzer(progName string, pathToTrace string) (string, error) {
	cmdStr := "./analyzer"
	resOut := resPath + "/" + progName + "/"
	cmdArgs := []string{"-t", pathToTrace, "-x", "-p", "-r", resOut}
	cmd := exec.Command(cmdStr, cmdArgs...)
	cmd.Dir = analyzerPath

	log(ANALYZE, cmd.String()+" in "+analyzerPath)

	out, time, err := runCommand(cmd)

	writeTime(time, progName)

	return out, err
}

func writeTime(time time.Duration, programName string) error {
	timeStr := fmt.Sprintf("%f", time.Seconds())
	file, err := os.OpenFile(resPath+"/"+programName+"/times.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString(timeStr + ",")

	return nil
}

func createOverview(progName string, progPath string, execPath string) error {
	cmdStr := "./overview"
	resOut := resPath + "/" + progName
	cmdArgs := []string{
		"-n", "overview",
		"-p", progPath,
		"-s", resOut,
		"-t", execPath + "trace/",
		"-r", resOut + "/results_readable.log",
		"-d", resOut + "/times.log"}
	cmd := exec.Command(cmdStr, cmdArgs...)
	cmd.Dir = overviewPath

	log(OVERVIEW, cmd.String()+" in "+overviewPath)
	err := cmd.Run()
	return err
}

func createResFolder(progName string) error {
	progFolder := resPath + "/" + progName

	log(CREATE, progName+" at "+progFolder)

	err := os.MkdirAll(progFolder, 0755)
	if err != nil {
		return err
	}

	return nil
}

func finish() error {
	log(CLEANUP, "Create overview at "+resPath+"/overview.log")
	file, err := os.OpenFile(resPath+"/overview.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("Success\n\n")
	for _, suc := range resSucc {
		file.WriteString(suc + "\n")
	}

	file.WriteString("Failed\n\n")
	for _, fail := range resFailure {
		_, err = file.WriteString(fail + "\n")
	}
	return err
}

func setup(all, constructed, gobench bool, medium bool) error {
	// create the res folder if needed
	name := resPath + "/results"
	counter := 1

	for {
		nameStr := name + "_" + strconv.Itoa(counter)
		if _, err := os.Stat(nameStr); os.IsNotExist(err) {
			errDir := os.MkdirAll(nameStr, 0755)
			if errDir != nil {
				return errDir
			}
			resPath = nameStr
			break
		}
		counter++
	}

	log(SETUP, "Create res folder "+resPath)

	if all || constructed {
		addConstructed()
	}
	if all || gobench {
		addGoBench()
	}
	if all || medium {
		addMediumPrograms()
	}
	return nil
}

func main() {
	runConstructed := flag.Bool("c", false, "Run constructed programs")
	runGoBench := flag.Bool("g", false, "Run go benchmarks")
	runMedium := flag.Bool("m", false, "Run medium programs")
	flag.Parse()

	runAll := !*runConstructed && !*runGoBench && !*runMedium

	setup(runAll, *runConstructed, *runGoBench, *runMedium)

	for _, program := range programs {
		log(START, program[0])

		name := program[0]
		pathToTrace := program[1] + "trace/"
		execPath := program[1]
		progPath := program[2]
		execName := program[3]
		execArgs := program[4:]

		err := createResFolder(name)
		if err != nil {
			log(ERROR, err.Error())
			log(FAILED, name)
			continue
		}

		runExecs(execPath, execName, execArgs, name)

		out, err := runAnalyzer(name, pathToTrace)
		if err != nil {
			log(ERROR, err.Error()+":\n"+out)
			log(FAILED, name)
			continue
		}

		err = createOverview(name, progPath, execPath)
		if err != nil {
			log(ERROR, err.Error()+":\n"+out)
			log(FAILED, name)
			continue
		}

		log(DONE, name)
	}

	err := finish()
	if err != nil {
		log(ERROR, err.Error())
		log(FAILED, "Failed to create overfiew")
	} else {
		log(DONE, "Created Overview")
	}

	if numberErrors == 0 {
		log(FINISHS, "Finished with 0 errors")
	} else {
		log(FINISHF, "Finished with "+strconv.Itoa(numberErrors)+" errors")
	}
}

type logType int

const (
	START logType = iota
	CREATE
	EXEC
	ANALYZE
	FAILED
	TIMEOUT
	DONE
	ERROR
	SETUP
	FINISHS
	FINISHF
	CLEANUP
	OVERVIEW
)

func log(lt logType, message string) {
	res := ""
	switch lt {
	case START:
		res = "\n[START   ] " + end + message
	case CREATE:
		res = "[CREATE  ] " + end + message
	case EXEC:
		res = "[EXEC    ] " + end + message
	case ANALYZE:
		res = "[ANALYZE ] " + end + message
	case DONE:
		res = grn + "[DONE    ] " + end + message
		resSucc = append(resSucc, message)
	case FAILED:
		res = red + "[FAILED  ] " + end + message
		resFailure = append(resFailure, message)
	case TIMEOUT:
		res = org + "[TIMEOUT ] " + end + message
	case ERROR:
		res = red + "[ERROR   ] " + end + message
		numberErrors++
	case SETUP:
		res = "[SETUP   ] " + message
	case CLEANUP:
		res = "[CLEANUP ] " + message
	case FINISHS:
		res = "\n" + grn + "[FINISH  ] " + end + message
	case FINISHF:
		res = "\n" + red + "[FINISH  ] " + end + message
	case OVERVIEW:
		res = "[OVERVIEW] " + message
	default:
		res = "[UNKOWN  ] "
	}
	fmt.Println(res)
}
