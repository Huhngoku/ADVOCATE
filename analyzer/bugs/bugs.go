package bugs

import (
	"analyzer/trace"
	"fmt"
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
	TraceElement2 []*trace.TraceElement
	Pos2          []string
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
	case DoneBeforeAdd:
		typeStr = "Possible negative waitgroup counter:"
		arg1Str = "done: "
		arg2Str = "done: "
	default:
		panic("Unknown bug type: " + strconv.Itoa(int(b.Type)))
	}
	res := typeStr + "\n\t" + arg1Str + b.Pos1 +
		"\n\t" + arg2Str
	for i, pos := range b.Pos2 {
		if i != 0 {
			res += ";"
		}
		res += pos
	}
	return res
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
	case "Possible negative waitgroup counter:":
		bug.Type = DoneBeforeAdd
		panic("Not implemented yet")
	default:
		panic("Unknown bug type: " + typeStr)
	}

	bug.Pos1 = strings.Split(arg1, ": ")[1]
	elem, err := trace.GetTraceElementFromPos(bug.Pos1)
	if err != nil {
		panic("Error: " + err.Error())
	}
	bug.TraceElement1 = elem

	bug.TraceElement2 = make([]*trace.TraceElement, 1)
	bug.Pos2 = make([]string, 1)

	elems := strings.Split(arg2, ": ")[1]
	fmt.Println(strings.Split(elems, ";"))
	for _, pos := range strings.Split(elems, ";") {
		elem, err = trace.GetTraceElementFromPos(pos)
		if err != nil {
			println("Error: " + err.Error())
		}
		bug.TraceElement2 = append(bug.TraceElement2, elem)
		bug.Pos2 = append(bug.Pos2, pos)
	}

	bug.Println()

	return false, bug
}
