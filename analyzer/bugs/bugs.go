package bugs

import (
	"analyzer/trace"
	"strconv"
	"strings"
)

type BugType int

const (
	SendOnClosed BugType = iota
	RecvOnClosed
	DoneBeforeAdd
)

type Bug struct {
	Type          BugType
	TraceElement1 *trace.TraceElement
	Pos1          string
	TraceElement2 *trace.TraceElement
	Pos2          string
}

/*
 * Convert the bug to a string
 * Returns:
 *   string: The bug as a string
 */
func (b Bug) ToString() string {
	typeStr := ""
	arg1Str := ""
	arg2Str := ""
	switch b.Type {
	case SendOnClosed:
		typeStr = "Possible Send on closed channel:"
		arg1Str = "close: "
		arg2Str = "send: "
	case RecvOnClosed:
		typeStr = "Possible Receive on closed channel:"
		arg1Str = "close: "
		arg2Str = "recv: "
	default:
		panic("Unknown bug type: " + strconv.Itoa(int(b.Type)))
	}
	return typeStr + "\n\t" + arg1Str + b.Pos1 +
		"\n\t" + arg2Str + b.Pos2
}

/*
 * Print the bug
 */
func (b Bug) Println() {
	println(b.ToString())
}

/*
 * Process the bug that was selected from the analysis results
 * Args:
 *   typeStr (string): The type of the bug
 *   arg1 (string): The first argument of the bug
 *   arg2 (string): The second argument of the bug
 * Returns:
 *   bool: true, if the bug was not a possible, but a actually occuring bug
 */
func ProcessBug(typeStr string, arg1 string, arg2 string) (bool, Bug) {
	bug := Bug{}

	actual := strings.Split(typeStr, " ")[0]
	if actual != "Possible" {
		return true, bug
	}

	switch typeStr {
	case "Possible send on closed channel:":
		bug.Type = SendOnClosed
	case "Possible receive on closed channel:":
		bug.Type = RecvOnClosed
	default:
		panic("Unknown bug type: " + typeStr)
	}

	elems := strings.Split(arg1, ": ")
	bug.Pos1 = elems[1]
	elem, err := trace.GetTraceElementFromPos(bug.Pos1)
	if err != nil {
		panic("Error: " + err.Error())
	}
	bug.TraceElement1 = elem

	elems = strings.Split(arg2, ": ")
	bug.Pos2 = elems[1]
	elem, err = trace.GetTraceElementFromPos(bug.Pos2)
	if err != nil {
		println("Error: " + err.Error())
	}
	bug.TraceElement2 = elem

	bug.Println()

	return false, bug
}
