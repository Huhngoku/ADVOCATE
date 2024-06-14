package explanation

import "fmt"

// type (bug / diagnostics)
var bugCrit = map[string]string{
	"ASendOnClosed":          "bug",
	"AReceiveOnClosed":       "diagnostics",
	"ACloseOnClosed":         "bug",
	"AConcurrentRecv":        "diagnostics",
	"ASelCaseWithoutPartner": "diagnostics",
}

// explanations
var bugExplanations = map[string]string{
	"ASendOnClosed": "Actual send on closed channel.\n" +
		"During the execution of the program, a send on a closed channel occurred.\n" +
		"The occurrence of a send on closed leads to a panic.",
	"ARecvOnClosed": "Actual receive on closed channel.\n" +
		"During the execution of the program, a receive on a closed channel occurred.\n",
	"ACloseOnClosed": "Actual close on closed channel.\n" +
		"During the execution of the program, a close on a close channel occurred.\n" +
		"The occurrence of a close on a closed channel leads to a panic.",
	"AConcurrentRecv": "Concurrent Receive\n" +
		"During the execution of the program, a channel waited for receive at multiple positions at the same time.\n" +
		"In this case, the actual receiver of a send message is chosen randomly.\n" +
		"This can lead to nondeterministic behavior.",
	"ASelCaseWithoutPartner": "Select Case without Partner\n" +
		"During the execution of the program, a select was executed, where, based " +
		"on the happens-before relation, at least one case could never be triggered.\n" +
		"This can be a desired behavior, especially considering, that only executed " +
		"operations are considered, but it can also be an hint of an unnecessary select case.",
}

const expASendOnClosed string = "Actual send on closed channel.\n" +
	"During the execution of the program, a send on a closed channel occurred."

// examples
var bugExamples map[string]string = map[string]string{
	"ASendOnClosed": "func main() {\n" +
		"\tc := make(chan int)\n" +
		"\tclose(c)\n" +
		"\tc <- 1",
	"ARecvOnClosed": "func main() {\n" +
		"\tc := make(chan int)\n" +
		"\tclose(c)\n" +
		"\t<-c",
	"ACloseOnClosed": "func main() {\n" +
		"\tc := make(chan int)\n" +
		"\tclose(c)\n" +
		"\tclose(c)",
	"AConcurrentRecv": "func main() {\n" +
		"\tc := make(chan int, 1)\n\n" +
		"\tgo func() {\n" +
		"\t\t<-c\n" +
		"\t}()\n\n" +
		"\tgo func() {\n" +
		"\t\t<-c\n" +
		"\t}()\n\n" +
		"\tc <- 1\n" +
		"}",
	"ASelCaseWithoutPartner": "func main() {\n" +
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

func printExplanation(bugType string) {
	fmt.Printf("Type: %s\n\n", bugCrit[bugType])
	fmt.Println(bugExplanations[bugType] + "\n")
	fmt.Println("A minimal example of this type of bug would be the following:\n")
	printCode(bugExamples[bugType])
}

func printCode(code string) {
	fmt.Println("```go")
	fmt.Println(code)
	fmt.Println("```")
}

func Run() {
	printExplanation("ASelCaseWithoutPartner")
}
