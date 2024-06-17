package explanation

// type (bug / diagnostics)
var bugCrit = map[string]string{
	"A1": "Bug",
	"A2": "Diagnostics",
	"A3": "Bug",
	"A4": "Diagnostics",
	"A5": "Diagnostics",
}

var bugNames = map[string]string{
	"A1": "Actual Send on Closed Channel",
	"A2": "Actual Receive on Closed Channel",
	"A3": "Actual Close on Closed Channel",
	"A4": "Concurrent Receive",
	"A5": "Select Case without Partner",
}

// explanations
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
}

// examples
var bugExamples map[string]string = map[string]string{
	"A1": "func main() {\n" +
		"\tc := make(chan int)\n" +
		"\tclose(c)\n" +
		"\tc <- 1",
	"A2": "func main() {\n" +
		"\tc := make(chan int)\n" +
		"\tclose(c)\n" +
		"\t<-c",
	"A3": "func main() {\n" +
		"\tc := make(chan int)\n" +
		"\tclose(c)\n" +
		"\tclose(c)",
	"A4": "func main() {\n" +
		"\tc := make(chan int, 1)\n\n" +
		"\tgo func() {\n" +
		"\t\t<-c\n" +
		"\t}()\n\n" +
		"\tgo func() {\n" +
		"\t\t<-c\n" +
		"\t}()\n\n" +
		"\tc <- 1\n" +
		"}",
	"A5": "func main() {\n" +
		"\tc := make(chan int)\n" +
		"\td := make(chan int)\n" +
		"\tgo func() {\n" +
		"\t\t<-c\n" +
		"\t}()\n\n" +
		"\tselect{\n" +
		"\tcase c1 := <- c:\n" +
		"\t\tprint(c1)\n" +
		"\tcase d <- 1:\n" +
		"\t\tprint(\"d\")\n" +
		"\t}\n" +
		"}",
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

func getBugElementType(elemType string) string {
	if _, ok := objectTypes[elemType]; !ok {
		return "Unknown element type"
	}
	return objectTypes[elemType]
}
