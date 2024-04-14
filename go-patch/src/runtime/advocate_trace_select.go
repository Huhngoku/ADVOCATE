package runtime

type advocateSelectElement struct {
	tPre    uint64                   // global timer before the operation
	tPost   uint64                   // global timer after the operation
	id      uint64                   // id of the select
	cases   []advocateChannelElement // cases of the select
	chosen  int                      // index of the chosen case in cases (0 indexed, -1 for default)
	nsend   int                      // number of send cases
	defa    bool                     // set true if a default case exists
	defaSel bool                     // set true if a default case was chosen
	file    string                   // file where the operation was called
	line    int                      // line where the operation was called
}

func (elem advocateSelectElement) isAdvocateTraceElement() {}

/*
 * Get a string representation of the element
 * Return:
 * 	string representation of the element "S,'tPre','tPost','id','cases','opId','file':'line'"
 *    'tPre' (number): global timer before the operation
 *    'tPost' (number): global timer after the operation
 *    'id' (number): id of the mutex
 *	  'cases' (string): cases of the select, d for default
 *    'chosen' (number): index of the chosen case in cases (0 indexed, -1 for default)
 *	  'opId' (number): id of the operation on the channel
 *    'file' (string): file where the operation was called
 *    'line' (number): line where the operation was called
 */
func (elem advocateSelectElement) toString() string {
	res := "S,"
	res += uint64ToString(elem.tPre) + "," + uint64ToString(elem.tPost) + ","
	res += uint64ToString(elem.id) + ","

	notNil := 0
	for _, ca := range elem.cases { // cases
		if ca.tPre != 0 { // ignore nil cases
			if notNil != 0 {
				res += "~"
			}
			res += ca.toStringSep(".", false)
			notNil++
		}
	}

	if elem.defa { // default
		if notNil != 0 {
			res += "~"
		}
		if elem.defaSel {
			res += "D"
		} else {
			res += "d"
		}
	}

	res += "," + intToString(elem.chosen) // case index

	res += "," + elem.file + ":" + intToString(elem.line)
	return res
}

/*
 * Get the operation
 */
func (elem advocateSelectElement) getOperation() Operation {
	return OperationSelect
}

/*
 * Get the file
 */
func (elem advocateSelectElement) getFile() string {
	return elem.file
}

/*
 * Get the line
 */
func (elem advocateSelectElement) getLine() int {
	return elem.line
}

/*
 * AdvocateSelectPre adds a select to the trace
 * Args:
 * 	cases: cases of the select
 * 	nsends: number of send cases
 * 	block: true if the select is blocking (has no default), false otherwise
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateSelectPre(cases *[]scase, nsends int, block bool) int {
	timer := GetAdvocateCounter()
	if cases == nil {
		return -1
	}

	id := GetAdvocateObjectID()
	caseElements := make([]advocateChannelElement, len(*cases))
	_, file, line, _ := Caller(2)

	for i, ca := range *cases {
		if ca.c != nil { // ignore nil cases
			caseElements[i] = advocateChannelElement{id: ca.c.id,
				op:    OperationChannelRecv,
				qSize: uint32(ca.c.dataqsiz), tPre: timer}
		}
	}

	elem := advocateSelectElement{id: id, cases: caseElements, nsend: nsends,
		defa: !block, file: file, line: line, tPre: timer}
	return insertIntoTrace(elem)
}

/*
 * AdvocateSelectPost adds a post event for select in case of an non-default case
 * Args:
 * 	index: index of the operation in the trace
 * 	c: channel of the chosen case
 * 	chosenIndex: index of the chosen case in the select
 * 	lockOrder: order of the locks
 * 	rClosed: true if the channel was closed at another routine
 */
func AdvocateSelectPost(index int, c *hchan, chosenIndex int,
	lockOrder []uint16, rClosed bool) {

	if index == -1 || c == nil {
		return
	}

	elem := currentGoRoutine().getElement(index).(advocateSelectElement)
	timer := GetAdvocateCounter()

	elem.chosen = chosenIndex
	elem.tPost = timer

	for i, op := range lockOrder {
		opChan := OperationChannelRecv
		if op < uint16(elem.nsend) {
			opChan = OperationChannelSend
		}
		elem.cases[i].op = opChan
	}

	if chosenIndex == -1 { // default case
		elem.defaSel = true
	} else {
		elem.cases[chosenIndex].tPost = timer
		elem.cases[chosenIndex].closed = rClosed

		// set oId
		if elem.cases[chosenIndex].op == OperationChannelSend {
			c.numberSend++
			elem.cases[chosenIndex].opID = c.numberSend
		} else {
			c.numberRecv++
			elem.cases[chosenIndex].opID = c.numberRecv
		}

	}

	currentGoRoutine().updateElement(index, elem)
}

/*
* AdvocateSelectPreOneNonDef adds a new select element to the trace if the
* select has exactly one non-default case and a default case
* Args:
* 	c: channel of the non-default case
* 	send: true if the non-default case is a send, false otherwise
* Return:
* 	index of the operation in the trace
 */
func AdvocateSelectPreOneNonDef(c *hchan, send bool) int {
	if c == nil {
		return -1
	}

	id := GetAdvocateObjectID()
	timer := GetAdvocateCounter()

	opChan := OperationChannelRecv
	if send {
		opChan = OperationChannelSend
	}

	caseElements := make([]advocateChannelElement, 1)
	caseElements[0] = advocateChannelElement{id: c.id,
		qSize: uint32(c.dataqsiz), tPre: timer, op: opChan}

	nSend := 0
	if send {
		nSend = 1
	}

	_, file, line, _ := Caller(2)

	elem := advocateSelectElement{id: id, cases: caseElements, nsend: nSend,
		defa: true, file: file, line: line, tPre: timer}
	return insertIntoTrace(elem)
}

/*
 * AdvocateSelectPostOneNonDef adds the selected case for a select with one
 * non-default and one default case
 * Args:
 * 	index: index of the operation in the trace
 * 	res: 0 for the non-default case, -1 for the default case
 */
func AdvocateSelectPostOneNonDef(index int, res bool, c *hchan) {
	if index == -1 {
		return
	}

	timer := GetAdvocateCounter()
	elem := currentGoRoutine().getElement(index).(advocateSelectElement)

	if res {
		elem.chosen = 0
		elem.cases[0].tPost = timer
		if elem.cases[0].op == OperationChannelSend {
			c.numberSend++
			elem.cases[0].opID = c.numberSend
		} else {
			c.numberRecv++
			elem.cases[0].opID = c.numberRecv
		}
	} else {
		elem.chosen = -1
		elem.defaSel = true
	}
	elem.tPost = timer

	currentGoRoutine().updateElement(index, elem)
}
