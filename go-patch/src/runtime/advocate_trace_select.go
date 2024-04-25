package runtime

/*
 * AdvocateSelectPre adds a select to the trace
 * Args:
 * 	cases: cases of the select
 * 	nsends: number of send cases
 * 	block: true if the select is blocking (has no default), false otherwise
 * 	lockOrder: internal order of the locks
 * Return:
 * 	index of the operation in the trace
 */
func AdvocateSelectPre(cases *[]scase, nsends int, block bool, lockOrder []uint16) int {
	timer := GetNextTimeStep()
	if cases == nil {
		return -1
	}

	id := GetAdvocateObjectID()
	caseElements := ""
	_, file, line, _ := Caller(2)

	for i, ca := range *cases {
		if ca.c == nil { // ignore nil cases
			continue
		}

		if len(caseElements) > 0 {
			caseElements += "~"
		}

		chanOp := "R"
		if lockOrder[i] < uint16(nsends) {
			chanOp = "S"
		}

		caseElements += "C." + uint64ToString(timer) + ".0." +
			uint64ToString(ca.c.id) + "." + chanOp + ".f.0." +
			uint32ToString(uint32(ca.c.dataqsiz))
	}

	if !block {
		caseElements += "~d"
	}

	elem := "S," + uint64ToString(timer) + ",0," + uint64ToString(id) + "," +
		caseElements + "," + file + ":" + intToString(line)

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
func AdvocateSelectPost(index int, c *hchan, chosenIndex int, rClosed bool) {

	if index == -1 || c == nil {
		return
	}

	elem := currentGoRoutine().getElement(index)
	println(elem)
	timer := GetNextTimeStep()

	split := splitStringAtCommas(elem, []int{2, 3, 4, 5})

	split[1] = uint64ToString(timer) // set tpost of select

	println(split[3])
	cases := splitStringAtSeparator(split[3], '~', nil)

	if chosenIndex == -1 { // default case
		if cases[len(cases)-1] != "d" {
			panic("default case on select without default")
		}
		cases[len(cases)-1] = "D"
	} else {
		// set tpost and cl of chosen case

		// split into C,[tpre] - [tPost] - [id] - [opC] - [cl] - [opID] - [qSize]
		chosenCaseSplit := splitStringAtSeparator(cases[chosenIndex], '.', []int{2, 3, 4, 5, 6, 7})
		chosenCaseSplit[1] = uint64ToString(timer)
		if rClosed {
			chosenCaseSplit[4] = "t"
		}

		// set oId
		if chosenCaseSplit[3] == "S" {
			c.numberSend++
			chosenCaseSplit[5] = uint64ToString(c.numberSend)
		} else {
			c.numberRecv++
			chosenCaseSplit[5] = uint64ToString(c.numberRecv)
		}

		cases[chosenIndex] = mergeStringSep(chosenCaseSplit, ".")
	}

	split[3] = mergeStringSep(cases, "~")
	elem = mergeString(split)

	println(elem)

	currentGoRoutine().updateElement(index, elem)
}

// MARK: OneNonDef

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

	opChan := "R"
	if send {
		opChan = "S"
	}

	caseElements := "C," + uint64ToString(timer) + ",0," + uint64ToString(c.id) +
		"," + opChan + ",f,0," + uint32ToString(uint32(c.dataqsiz))

	_, file, line, _ := Caller(2)

	elem := "S," + uint64ToString(timer) + ",0," + uint64ToString(id) + "," +
		caseElements + "~d," + file + ":" + intToString(line)

	return insertIntoTrace(elem)
}

/*
 * AdvocateSelectPostOneNonDef adds the selected case for a select with one
 * non-default and one default case
 * Args:
 * 	index: index of the operation in the trace
 * 	res: true for channel, false for default
 */
func AdvocateSelectPostOneNonDef(index int, res bool, c *hchan) {
	if index == -1 {
		return
	}

	timer := GetNextTimeStep()
	elem := currentGoRoutine().getElement(index)

	split := splitStringAtCommas(elem, []int{2, 3, 4, 5})

	// update tPost
	split[1] = uint64ToString(timer)

	// update cases
	cases := splitStringAtSeparator(split[4], '~', nil)
	if res { // channel case
		// split into C,[tpre] - [tPost] - [id] - [opC] - [cl] - [opID] - [qSize]
		chosenCaseSplit := splitStringAtSeparator(cases[0], '.', []int{2, 3, 4, 5, 6, 7})
		chosenCaseSplit[1] = uint64ToString(timer)

		if chosenCaseSplit[3] == "S" {
			c.numberSend++
			chosenCaseSplit[5] = uint64ToString(c.numberSend)
		} else {
			c.numberRecv++
			chosenCaseSplit[5] = uint64ToString(c.numberRecv)
		}
		cases[0] = mergeStringSep(chosenCaseSplit, ".")
	} else { // default case
		cases[1] = "D"
	}
	split[4] = mergeStringSep(cases, "~")

	elem = mergeString(split)

	currentGoRoutine().updateElement(index, elem)
}
