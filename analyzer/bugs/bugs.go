package bugs

import (
	"analyzer/trace"
	"errors"
	"strconv"
	"strings"
)

type BugType int

const (
	SendOnClosed BugType = iota
	PosRecvOnClosed
	RecvOnClosed // actual send on closed
	CloseOnClosed
	DoneBeforeAdd
	SelectWithoutPartner

	ConcurrentRecv

	MixedDeadlock
	CyclicDeadlock

	RoutineLeakPartner   // chan and select
	RoutineLeakNoPartner // chan and select
	RoutineLeakMutex
	RoutineLeakWaitGroup
	RoutineLeakCond
)

type Bug struct {
	Type          BugType
	TraceElement1 []*trace.TraceElement
	TID1          []string
	TraceElement2 []*trace.TraceElement
	TID2          []string
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
	case PosRecvOnClosed:
		typeStr = "Possible Receive on closed channel:"
		arg1Str = "close: "
		arg2Str = "recv: "
	case RecvOnClosed:
		typeStr = "Found receive on closed channel:"
		arg1Str = "close: "
		arg2Str = "recv: "
	case CloseOnClosed:
		typeStr = "Possible close on closed channel:"
		arg1Str = "close: "
		arg2Str = "close: "
	case DoneBeforeAdd:
		typeStr = "Possible negative waitgroup counter:"
		arg1Str = "add: "
		arg2Str = "done: "
	case ConcurrentRecv:
		typeStr = "Found concurrent Recv on same channel:"
		arg1Str = "recv: "
		arg2Str = "recv: "
	case SelectWithoutPartner:
		typeStr = "Possible select case without partner:"
		arg1Str = "select: "
		arg2Str = ""
	case MixedDeadlock:
		typeStr = "Potential mixed deadlock:"
		arg1Str = "lock: "
		arg2Str = "lock: "
	case CyclicDeadlock:
		typeStr = "Potential cyclic deadlock:"
		arg1Str = "lock: "
		arg2Str = "cycle: "
	case RoutineLeakPartner, RoutineLeakNoPartner:
		typeStr = "Potential routine leak channel:"
		arg1Str = "channel: "
		arg2Str = "partner: "
	case RoutineLeakMutex:
		typeStr = "Potential routine leak mutex:"
		arg1Str = "mutex: "
		arg2Str = ""
	case RoutineLeakWaitGroup:
		typeStr = "Potential routine leak waitgroup:"
		arg1Str = "waitgroup: "
		arg2Str = ""
	case RoutineLeakCond:
		typeStr = "Potential routine leak conditional variable:"
		arg1Str = "conditional: "
		arg2Str = ""

	default:
		panic("Unknown bug type: " + strconv.Itoa(int(b.Type)))
	}
	res := typeStr + "\n\t" + arg1Str
	for i, pos := range b.TID1 {
		if i != 0 {
			res += ";"
		}
		res += pos
	}

	res += "\n\t" + arg2Str

	for i, pos := range b.TID2 {
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
 *   Bug: The bug that was selected
 *   error: An error if the bug could not be processed
 */
func ProcessBug(typeStr string, arg1 string, arg2 string) (bool, Bug, error) {
	bug := Bug{}

	actual := strings.Split(typeStr, " ")[0]
	if actual != "Possible" && actual != "Potential" {
		return true, bug, nil
	}

	// println("Process bug: " + typeStr + " " + arg1 + " " + arg2)

	switch typeStr {
	case "Possible send on closed channel:":
		bug.Type = SendOnClosed
	case "Possible receive on closed channel:":
		bug.Type = PosRecvOnClosed
	case "Found receive on closed channel:":
		bug.Type = RecvOnClosed
	case "Found close on closed channel:":
		bug.Type = CloseOnClosed
	case "Possible negative waitgroup counter:":
		bug.Type = DoneBeforeAdd
	case "Possible select case without partner:":
		bug.Type = SelectWithoutPartner
	case "Found concurrent Recv on same channel:":
		bug.Type = ConcurrentRecv
	case "Potential mixed deadlock:":
		bug.Type = MixedDeadlock
	case "Potential leak with possible partner:":
		bug.Type = RoutineLeakPartner
	case "Potential leak without possible partner:":
		bug.Type = RoutineLeakNoPartner
	case "Potential leak on mutex:":
		bug.Type = RoutineLeakMutex
	case "Potential leak on wait group:":
		bug.Type = RoutineLeakWaitGroup
	case "Potential leak on conditional variable:":
		bug.Type = RoutineLeakCond
	case "Potential cyclic deadlock:":
		bug.Type = CyclicDeadlock
	default:
		return false, bug, errors.New("Unknown bug type: " + typeStr)
	}

	bug.TraceElement2 = make([]*trace.TraceElement, 0)
	bug.TID2 = make([]string, 0)

	elems := strings.Split(arg1, ": ")[1]

	for _, tID := range strings.Split(elems, ";") {
		if strings.TrimSpace(tID) == "" {
			continue
		}

		elem, err := trace.GetTraceElementFromTID(tID)
		if err != nil {
			println("not found: " + tID + " in " + arg1 + " " + arg2)
			return false, bug, err
		}
		bug.TraceElement1 = append(bug.TraceElement1, elem)
		bug.TID1 = append(bug.TID1, tID)
	}

	bug.TraceElement2 = make([]*trace.TraceElement, 0)
	bug.TID2 = make([]string, 0)

	elems = strings.Split(arg2, ": ")[1]

	for _, tID := range strings.Split(elems, ";") {
		if strings.TrimSpace(tID) == "" || strings.TrimSpace(tID) == "-" {
			continue
		}
		elem, err := trace.GetTraceElementFromTID(tID)
		if err != nil {
			return false, bug, err
		}
		bug.TraceElement2 = append(bug.TraceElement2, elem)
		bug.TID2 = append(bug.TID2, tID)
	}

	return false, bug, nil
}
