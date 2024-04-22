package runtime

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
	timer := GetNextTimeStep()
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
	timer := GetNextTimeStep()

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
	timer := GetNextTimeStep()

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

	timer := GetNextTimeStep()
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
