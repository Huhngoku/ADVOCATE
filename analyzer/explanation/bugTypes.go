package explanation

import "fmt"

// type (bug / diagnostics)
var bugCrit = map[string]string{
	"A1": "Bug",
	"A2": "Diagnostics",
	"A3": "Bug",
	"A4": "Diagnostics",
	"A5": "Diagnostics",
	"P1": "Bug",
	"P2": "Diagnostic",
	"P3": "Leak",
	"L1": "Leak",
	"L2": "Leak",
	"L3": "Leak",
	"L4": "Leak",
	"L5": "Leak",
	"L6": "Leak",
	"L7": "Leak",
	"L8": "Leak",
	"L9": "Leak",
	"L0": "Leak",
}

var bugNames = map[string]string{
	"A1": "Actual Send on Closed Channel",
	"A2": "Actual Receive on Closed Channel",
	"A3": "Actual Close on Closed Channel",
	"A4": "Concurrent Receive",
	"A5": "Select Case without Partner",

	"P1": "Possible Send on Closed Channel",
	"P2": "Possible Receive on Closed Channel",
	"P3": "Possible Negative WaitGroup cCounter",

	"L1": "Leak of unbuffered Channel with possible partner",
	"L2": "Leak on unbuffered Channel without possible partner",
	"L3": "Leak of buffered Channel with possible partner",
	"L4": "Leak on buffered Channel without possible partner",
	"L5": "Leak on nil channel",
	"L6": "Leak of select with possible partner",
	"L7": "Leak on select without possible partner",
	"L8": "Leak on sync.Mutex",
	"L9": "Leak on sync.WaitGroup",
	"L0": "Leak on sync.Cond",
}

// explanations
// TODO: add missing bug explainations
var bugExplanations = map[string]string{
	"A1": "During the execution of the program, a send on a closed channel occurred.\n" +
		"The occurrence of a send on closed leads to a panic.",
	"A2": "During the execution of the program, a receive on a closed channel occurred.\n",
	"A3": "During the execution of the program, a close on a close channel occurred.\n" +
		"The occurrence of a close on a closed channel leads to a panic.",
	"A4": "During the execution of the program, a channel waited to receive at multiple positions at the same time.\n" +
		"In this case, the actual receiver of a send message is chosen randomly.\n" +
		"This can lead to nondeterministic behavior.",
	"A5": "During the execution of the program, a select was executed, where, based " +
		"on the happens-before relation, at least one case could never be triggered.\n" +
		"This can be a desired behavior, especially considering, that only executed " +
		"operations are considered, but it can also be an hint of an unnecessary select case.",
	"P1": "The analyzer detected a possible send on a closed channel.\n" +
		"Although the send on a closed channel did not occur during the recording, " +
		"it is possible that it will occur, based on the happens before relation.\n" +
		"Such a send on a closed channel leads to a panic.",
	"P2": "The analyzer detected a possible receive on a closed channel.\n" +
		"Although the receive on a closed channel did not occur during the recording, " +
		"it is possible that it will occur, based on the happens before relation." +
		"This is not necessarily a bug, but it can be an indication of a bug.",
	"P3": "The analyzer detected a possible negative WaitGroup counter.\n" +
		"Although the negative counter did not occur during the recording, " +
		"it is possible that it will occur, based on the happens before relation.\n" +
		"A negative counter will lead to a panic.",
	"L1": "The analyzer detected a leak of an unbuffered channel with a possible partner.\n" +
		"A leak of an unbuffered channel is a situation, where a unbuffered channel is " +
		"still blocking at the end of the program.\n" +
		"The partner is a corresponding send or receive operation, which communicated with another operation, " +
		"but could communicated with the stuck operation instead, resolving the deadlock.",
	"L2": "The analyzer detected a leak of an unbuffered channel without a possible partner.\n" +
		"A leak of an unbuffered channel is a situation, where a unbuffered channel is " +
		"still blocking at the end of the program.\n" +
		"The analyzer could not find a partner for the stuck operation, which would resolve the leak.",
	"L3": "The analyzer detected a leak of a buffered channel with a possible partner.\n" +
		"A leak of a buffered channel is a situation, where a buffered channel is " +
		"still blocking at the end of the program.\n" +
		"The partner is a corresponding send or receive operation, which communicated with another operation, " +
		"but could communicated with the stuck operation instead, resolving the leak.",
	"L4": "The analyzer detected a leak of a buffered channel without a possible partner.\n" +
		"A leak of a buffered channel is a situation, where a buffered channel is " +
		"still blocking at the end of the program.\n" +
		"The analyzer could not find a partner for the stuck operation, which would resolve the leak.",
	"L5": "The analyzer detected a leak on a nil channel.\n" +
		"A leak on a nil channel is a situation, where a nil channel is still blocking at the end of the program.\n" +
		"A nil channel is a channel, which was never initialized or set to nil." +
		"An operation on a nil channel will block indefinitely.",
	"L6": "The analyzer detected a leak of a select with a possible partner.\n" +
		"A leak of a select is a situation, where a select is still blocking at the end of the program.\n" +
		"The partner is a corresponding send or receive operation, which communicated with another operation, " +
		"but could communicated with the stuck operation instead, resolving the leak.",
	"L7": "The analyzer detected a leak of a select without a possible partner.\n" +
		"A leak of a select is a situation, where a select is still blocking at the end of the program.\n" +
		"The analyzer could not find a partner for the stuck operation, which would resolve the leak.",
	"L8": "The analyzer detected a leak on a sync.Mutex.\n" +
		"A leak on a sync.Mutex is a situation, where a sync.Mutex lock operations is still blocking at the end of the program.\n" +
		"A sync.Mutex lock operation is a operation, which is blocking, because the lock is already acquired.",
	"L9": "The analyzer detected a leak on a sync.WaitGroup.\n" +
		"A leak on a sync.WaitGroup is a situation, where a sync.WaitGroup is still blocking at the end of the program.\n" +
		"A sync.WaitGroup wait is blocking, because the counter is not zero.",
	"L0": "The analyzer detected a leak on a sync.Cond.\n" +
		"A leak on a sync.Cond is a situation, where a sync.Cond wait is still blocking at the end of the program.\n" +
		"A sync.Cond wait is blocking, because the condition is not met.",
}

// examples
var bugExamples map[string]string = map[string]string{
	"A1": "func main() {\n" +
		"    c := make(chan int)\n" +
		"    close(c)          // <-------\n" +
		"    c <- 1            // <-------\n}",
	"A2": "func main() {\n" +
		"    c := make(chan int)\n" +
		"    close(c)          // <-------\n" +
		"    <-c               // <-------\n}",
	"A3": "func main() {\n" +
		"    c := make(chan int)\n" +
		"    close(c)          // <-------\n" +
		"    close(c)          // <-------\n}",
	"A4": "func main() {\n" +
		"    c := make(chan int, 1)\n\n" +
		"    go func() {\n" +
		"        <-c             // <-------\n" +
		"    }()\n\n" +
		"    go func() {\n" +
		"        <-c             // <-------\n" +
		"    }()\n\n" +
		"    c <- 1\n" +
		"}",
	"A5": "func main() {\n" +
		"    c := make(chan int)\n" +
		"    d := make(chan int)\n" +
		"    go func() {\n" +
		"        <-c\n" +
		"    }()\n\n" +
		"    select{\n" +
		"    case c1 := <- c:\n" +
		"        print(c1)\n" +
		"    case d <- 1:      // <-------\n" +
		"        print(\"d\")\n" +
		"    }\n",
	"P1": "func main() {\n" +
		"    c := make(chan int)\n\n" +
		"    go func() {\n" +
		"        c <- 1          // <-------\n" +
		"    }()\n\n" +
		"    go func() {\n" +
		"        <- c\n" +
		"    }()\n\n" +
		"    close(c)            // <-------\n}",
	"P2": "func main() {\n" +
		"    c := make(chan int)\n\n" +
		"    go func() {\n" +
		"        c <- 1\n" +
		"    }()\n\n" +
		"    go func() {\n" +
		"        <- c            // <-------\n" +
		"    }()\n\n" +
		"    close(c)            // <-------\n}",
	"P3": "func main() {\n" +
		"    var wg sync.WaitGroup\n\n" +
		"    go func() {\n" +
		"        wg.Add(1)       // <-------\n" +
		"    }()\n\n" +
		"    go func() {\n" +
		"        wg.Done()       // <-------\n" +
		"    }()\n\n" +
		"    wg.Wait()\n}",
	"L1": "func main() {\n" +
		"    c := make(chan int)\n\n" +
		"    go func() {\n" +
		"        c <- 1          // <------- Communicates\n" +
		"    }()\n\n" +
		"    go func() {\n" +
		"        <- c            // <------- Communicates, possible partner\n" +
		"    }()\n\n" +
		"    go func() {\n" +
		"        c <- 1          // <------- Leak\n" +
		"    }()\n" +
		"}",
	"L2": "func main() {\n" +
		"    c := make(chan int)\n\n" +
		"    go func() {\n" +
		"        c <- 1          // <------- Leak, no possible partner\n" +
		"    }()\n" +
		"}",
	"L3": "func main() {\n" +
		"    c := make(chan int, 1)\n\n" +
		"    go func() {\n" +
		"        c <- 1          // <------- Communicates\n" +
		"    }()\n\n" +
		"    go func() {\n" +
		"        <- c            // <------- Communicates, possible partner\n" +
		"    }()\n\n" +
		"    go func() {\n" +
		"        c <- 1          // <------- Leak\n" +
		"    }()\n" +
		"}",
	"L4": "func main() {\n" +
		"    c := make(chan int, 1)\n\n" +
		"    go func() {\n" +
		"        c <- 1          // <------- Leak, no possible partner\n" +
		"    }()\n" +
		"}",
	"L5": "func main() {\n" +
		"    var c chan int      // <------- Not initialized -> c = nil\n\n" +
		"    go func() {\n" +
		"        c <- 1          // <------- Leak\n" +
		"    }()\n",
	"L6": "func main() {\n" +
		"    c := make(chan int)\n\n" +
		"    go func() {\n" +
		"        c <- 1          // <------- Communicates\n" +
		"    }()\n\n" +
		"    go func() {\n" +
		"        <- c            // <------- Communicates, possible partner\n" +
		"    }()\n\n" +
		"    go func() {\n" +
		"        select {        // <------- Leak\n" +
		"        case c <- 1:    // <------- Possible partner\n" +
		"        }\n" +
		"    }()\n" +
		"}",
	"L7": "func main() {\n" +
		"    c := make(chan int)\n\n" +
		"    go func() {\n" +
		"        select {        // <------- Leak, no possible partner\n" +
		"        case c <- 1:\n" +
		"        }\n" +
		"    }()\n" +
		"}",
	"L8": "func main() {\n" +
		"    var m sync.Mutex\n\n" +
		"    go func() {\n" +
		"        m.Lock()        // <------- Leak\n" +
		"    }()\n\n" +
		"    m.Lock()            // <------- Lock, no unlock\n" +
		"}",
	"L9": "func main() {\n" +
		"    var wg sync.WaitGroup\n\n" +
		"    wg.Add(1)           // <------- Add, no Done\n" +
		"    wg.Wait()           // <------- Leak\n" +
		"}",
	"L0": "func main() {\n" +
		"    var c sync.Cond\n\n" +
		"    c.Wait()            // <------- Leak, no signal/broadcast\n" +
		"}",
}

var rewriteType = map[string]string{
	"A1": "Actual",
	"A2": "Actual",
	"A3": "Actual",
	"A4": "Actual",
	"A5": "Actual",
	"P1": "Possible",
	"P2": "Possible",
	"P3": "Possible",
	"L1": "LeakPos",
	"L2": "Leak",
	"L3": "LeakPos",
	"L4": "Leak",
	"L5": "Leak",
	"L6": "LeakPos",
	"L7": "Leak",
	"L8": "LeakPos",
	"L9": "LeakPos",
	"L0": "LeakPos",
}

// TODO: describe exit codes
var exitCodeExplanation = map[string]string{
	"0":  "exites without",
	"10": "stuck finish",
	"11": "stuck wait elem",
	"12": "stuck no elem",
	"13": "stuck empty trace",
	"20": "leak unbuf",
	"21": "leak buf",
	"22": "leak mutex",
	"23": "leak cond",
	"24": "leak wg",
	"30": "send close",
	"31": "recv close",
	"32": "negative wg",
	// "41": "cyclic",
}

var objectTypes = map[string]string{
	"CS": "Channel: Send",
	"CR": "Channel: Receive",
	"CC": "Channel: Close",
	"ML": "Mutex: Lock",
	"MR": "Mutex: RLock",
	"MT": "Mutex: TryLock",
	"MY": "Mutex: TryRLock",
	"MU": "Mutex: Unlock",
	"MN": "Mutex: RUnlock",
	"WA": "Waitgroup: Add",
	"WD": "Waitgroup: Done",
	"WW": "Waitgroup: Wait",
	"SS": "Select:",
	"NW": "Conditional Variable: Wait",
	"NB": "Conditional Variable: Broadcast",
	"NS": "Conditional Variable: Signal",
	"OE": "Once: Done Executed",
	"ON": "Once: Done Not Executed (because the once was already executed)",
	"GF": "Routine: Fork",
}

func getBugTypeDescription(bugType string) map[string]string {
	return map[string]string{
		"crit":        bugCrit[bugType],
		"name":        bugNames[bugType],
		"explanation": bugExplanations[bugType],
		"example":     bugExamples[bugType],
	}
}

func printBugTypeDescription(bugType string) {
	fmt.Println(bugCrit[bugType] + ": " + bugNames[bugType] + "\n")
	fmt.Println(bugExplanations[bugType] + "\n")
	fmt.Println(bugExamples[bugType])
}

func getBugElementType(elemType string) string {
	if _, ok := objectTypes[elemType]; !ok {
		return "Unknown element type"
	}
	return objectTypes[elemType]
}
