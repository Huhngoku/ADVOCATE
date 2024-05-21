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

	LeakUnbufChanPartner
	LeakUnbufChanNoPartner
	LeakBufChanPartner
	LeakBufChanNoPartner
	LeakSelectPartnerUnbuf
	LeakSelectPartnerBuf
	LeakSelectNoPartner
	LeakMutex
	LeakWaitGroup
	LeakCond
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
		typeStr = "Possible send on closed channel:"
		arg1Str = "close: "
		arg2Str = "send: "
	case PosRecvOnClosed:
		typeStr = "Possible receive on closed channel:"
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
	case SelectWithoutPartner:
		typeStr = "Possible select case without partner:"
		arg1Str = "select: "
		arg2Str = "partner: "
	case ConcurrentRecv:
		typeStr = "Found concurrent Recv on same channel:"
		arg1Str = "recv: "
		arg2Str = "recv: "
	case MixedDeadlock:
		typeStr = "Possible mixed deadlock:"
		arg1Str = "lock: "
		arg2Str = "lock: "
	case CyclicDeadlock:
		typeStr = "Possible cyclic deadlock:"
		arg1Str = "lock: "
		arg2Str = "cycle: "
	case LeakUnbufChanPartner, LeakUnbufChanNoPartner:
		typeStr = "Leak of unbuffered channel:"
		arg1Str = "channel: "
		arg2Str = "partner: "
	case LeakBufChanPartner:
		typeStr = "Leak of buffered channel with partner:"
		arg1Str = "channel: "
		arg2Str = "partner: "
	case LeakBufChanNoPartner:
		typeStr = "Leak of buffered channel without partner:"
		arg1Str = "channel: "
		arg2Str = "partner: "
	case LeakSelectPartnerUnbuf:
		typeStr = "Leak of select with unbuffered partner:"
		arg1Str = "select: "
		arg2Str = "partner: "
	case LeakSelectPartnerBuf:
		typeStr = "Leak of select with buffered partner:"
		arg1Str = "select: "
		arg2Str = "partner: "
	case LeakSelectNoPartner:
		typeStr = "Leak of select without partner:"
		arg1Str = "select: "
		arg2Str = "partner: "
	case LeakMutex:
		typeStr = "Leak of mutex:"
		arg1Str = "mutex: "
		arg2Str = "last: "
	case LeakWaitGroup:
		typeStr = "Leak of waitgroup:"
		arg1Str = "waitgroup: "
		arg2Str = ""
	case LeakCond:
		typeStr = "Leak of conditional variable:"
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

	if len(b.TID2) == 0 {
		res += "-"
	}

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
	if actual != "Possible" && actual != "Leak" {
		return true, bug, nil
	}

	// println("Process bug: " + typeStr + " " + arg1 + " " + arg2)

	containsArg2 := true

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
		containsArg2 = false
	case "Found concurrent Recv on same channel:":
		bug.Type = ConcurrentRecv
	case "Possible mixed deadlock:":
		bug.Type = MixedDeadlock
	case "Leak on unbuffered channel with possible partner:":
		bug.Type = LeakUnbufChanPartner
	case "Leak on unbuffered channel without possible partner:":
		bug.Type = LeakUnbufChanNoPartner
		containsArg2 = false
	case "Leak on buffered channel with possible partner:":
		bug.Type = LeakBufChanPartner
	case "Leak on buffered channel without possible partner:":
		bug.Type = LeakBufChanNoPartner
		containsArg2 = false
	case "Leak on select with possible buffered partner:":
		bug.Type = LeakSelectPartnerBuf
	case "Leak on select with possible unbuffered partner:":
		bug.Type = LeakSelectPartnerUnbuf
	case "Leak on select without possible partner:":
		bug.Type = LeakSelectNoPartner
		containsArg2 = false
	case "Leak on mutex:":
		bug.Type = LeakMutex
	case "Leak on wait group:":
		bug.Type = LeakWaitGroup
		containsArg2 = false
	case "Leak on conditional variable:":
		bug.Type = LeakCond
		containsArg2 = false
	case "Possible cyclic deadlock:":
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
			println("Could not find: " + tID + " in trace")
			return false, bug, err
		}
		bug.TraceElement1 = append(bug.TraceElement1, elem)
		bug.TID1 = append(bug.TID1, tID)
	}

	bug.TraceElement2 = make([]*trace.TraceElement, 0)
	bug.TID2 = make([]string, 0)

	if arg2 == "" || arg2 == "\t" || !containsArg2 {
		return false, bug, nil
	}

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
