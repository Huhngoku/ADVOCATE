package bugs

import (
	"strconv"
	"strings"
)

type BugType int

const (
	SendOnClosed BugType = iota
	RecvOnClosed
)

type Bug struct {
	Type  BugType
	File1 string
	File2 string
	Line1 int
	Line2 int
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
	return typeStr + "\n\t" + arg1Str + b.File1 + ":" + strconv.Itoa(b.Line1) +
		"\n\t" + arg2Str + b.File2 + ":" + strconv.Itoa(b.Line2)
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
	elems = strings.Split(elems[1], ":")
	bug.File1 = elems[0]
	bug.Line1, _ = strconv.Atoi(elems[1])

	elems = strings.Split(arg2, ": ")
	elems = strings.Split(elems[1], ":")
	bug.File2 = elems[0]
	bug.Line2, _ = strconv.Atoi(elems[1])

	bug.Println()

	return false, bug
}
