package trace

import (
	"analyzer/analysis"
	"analyzer/clock"
	"errors"
	"math"
	"strconv"
	"strings"
)

/*
 * TraceElementSelect is a trace element for a select statement
 * MARK: Struct
 * Fields:
 *   routine (int): The routine id
 *   tpre (int): The timestamp at the start of the event
 *   tpost (int): The timestamp at the end of the event
 *   id (int): The id of the select statement
 *   cases ([]traceElementSelectCase): The cases of the select statement
 *   chosenIndex (int): The internal index of chosen case
 *   containsDefault (bool): Whether the select statement contains a default case
 *   chosenCase (traceElementSelectCase): The chosen case, nil if default case chosen
 *   chosenDefault (bool): if the default case was chosen
 *   pos (string): The position of the select statement in the code
 *   tID (string): The id of the trace element, contains the position and the tpre
 */
type TraceElementSelect struct {
	routine         int
	tPre            int
	tPost           int
	id              int
	cases           []TraceElementChannel
	chosenCase      TraceElementChannel
	chosenIndex     int
	containsDefault bool
	chosenDefault   bool
	pos             string
	tID             string
	vc              clock.VectorClock
}

/*
 * Add a new select statement trace element
 * MARK: New
 * Args:
 *   routine (int): The routine id
 *   tPre (string): The timestamp at the start of the event
 *   tPost (string): The timestamp at the end of the event
 *   id (string): The id of the select statement
 *   cases (string): The cases of the select statement
 *   chosenIndex (string): The internal index of chosen case
 *   pos (string): The position of the select statement in the code
 */
func AddTraceElementSelect(routine int, tPre string,
	tPost string, id string, cases string, chosenIndex string, pos string) error {

	tPreInt, err := strconv.Atoi(tPre)
	if err != nil {
		return errors.New("tpre is not an integer")
	}

	tPostInt, err := strconv.Atoi(tPost)
	if err != nil {
		return errors.New("tpost is not an integer")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("id is not an integer")
	}

	chosenIndexInt, err := strconv.Atoi(chosenIndex)
	if err != nil {
		return errors.New("chosenIndex is not an integer")
	}

	tID := pos + "@" + tPre

	elem := TraceElementSelect{
		routine:     routine,
		tPre:        tPreInt,
		tPost:       tPostInt,
		id:          idInt,
		chosenIndex: chosenIndexInt,
		pos:         pos,
		tID:         tID,
	}

	cs := strings.Split(cases, "~")
	casesList := make([]TraceElementChannel, 0)
	containsDefault := false
	chosenDefault := false
	for _, c := range cs {
		if c == "" {
			continue
		}

		if c == "d" {
			containsDefault = true
			break
		}
		if c == "D" {
			containsDefault = true
			chosenDefault = true
			break
		}

		// read channel operation
		caseList := strings.Split(c, ".")
		cTPre, err := strconv.Atoi(caseList[1])
		if err != nil {
			return errors.New("c_tpre is not an integer")
		}
		cTPost, err := strconv.Atoi(caseList[2])
		if err != nil {
			return errors.New("c_tpost is not an integer")
		}

		cID := -1
		if caseList[3] != "*" {
			cID, err = strconv.Atoi(caseList[3])
			if err != nil {
				println("Error: " + caseList[3])
				return errors.New("c_id is not an integer")
			}
		}
		var cOpC = Send
		if caseList[4] == "R" {
			cOpC = Recv
		} else if caseList[4] == "C" {
			panic("Close in select case list")
		}

		cCl, err := strconv.ParseBool(caseList[5])
		if err != nil {
			return errors.New("c_cr is not a boolean")
		}

		cOID, err := strconv.Atoi(caseList[6])
		if err != nil {
			return errors.New("c_oId is not an integer")
		}
		cOSize, err := strconv.Atoi(caseList[7])
		if err != nil {
			return errors.New("c_oSize is not an integer")
		}

		tIDStr := pos + "@" + strconv.Itoa(cTPre)

		elemCase := TraceElementChannel{
			routine: routine,
			tPre:    cTPre,
			tPost:   cTPost,
			id:      cID,
			opC:     cOpC,
			cl:      cCl,
			oID:     cOID,
			qSize:   cOSize,
			sel:     &elem,
			pos:     pos,
			tID:     tIDStr,
		}

		casesList = append(casesList, elemCase)
		if elemCase.tPost != 0 {
			elem.chosenCase = elemCase
		}
	}

	elem.containsDefault = containsDefault
	elem.chosenDefault = chosenDefault
	elem.cases = casesList

	// check if partner was already processed, otherwise add to channelWithoutPartner
	if tPostInt != 0 {
		id := elem.chosenCase.id
		oID := elem.chosenCase.oID
		if _, ok := channelWithoutPartner[id][oID]; ok {
			elem.chosenCase.partner = channelWithoutPartner[id][oID]
			channelWithoutPartner[elem.chosenCase.id][oID].partner = &elem.chosenCase
			delete(channelWithoutPartner[id], oID)
		} else {
			if _, ok := channelWithoutPartner[id]; !ok {
				channelWithoutPartner[id] = make(map[int]*TraceElementChannel)
			}

			channelWithoutPartner[id][oID] = &elem.chosenCase
		}
	}

	return AddElementToTrace(&elem)
}

// MARK: Getter

/*
 * Get the id of the element
 * Returns:
 *   int: The id of the element
 */
func (se *TraceElementSelect) GetID() int {
	return se.id
}

/*
 * Get the cases of the select statement
 * Returns:
 *   []traceElementChannel: The cases of the select statement
 */
func (se *TraceElementSelect) GetCases() []TraceElementChannel {
	return se.cases
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (se *TraceElementSelect) GetRoutine() int {
	return se.routine
}

/*
 * Get the timestamp at the start of the event
 * Returns:
 *   int: The timestamp at the start of the event
 */
func (se *TraceElementSelect) GetTPre() int {
	return se.tPre
}

/*
 * Get the timestamp at the end of the event
 * Returns:
 *   int: The timestamp at the end of the event
 */
func (se *TraceElementSelect) getTpost() int {
	return se.tPost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (se *TraceElementSelect) GetTSort() int {
	if se.tPost == 0 {
		// add at the end of the trace
		return math.MaxInt
	}
	return se.tPost
}

/*
 * Get the position of the operation.
 * Returns:
 *   string: The position of the element
 */
func (se *TraceElementSelect) GetPos() string {
	return se.pos
}

/*
 * Get the tID of the element.
 * Returns:
 *   string: The tID of the element
 */
func (se *TraceElementSelect) GetTID() string {
	return se.tID
}

/*
 * Get the vector clock of the element
 * Returns:
 *   VectorClock: The vector clock of the element
 */
func (se *TraceElementSelect) GetVC() clock.VectorClock {
	return se.vc
}

/*
 * Get the communication partner of the select
 * Returns:
 *   *TraceElementChannel: The communication partner of the select or nil
 */
func (se *TraceElementSelect) GetPartner() *TraceElementChannel {
	if se.chosenCase.tPost != 0 {
		return se.chosenCase.partner
	}
	return nil
}

// MARK: Setter

/*
 * Set the tPre and tPost of the element
 * Args:
 *   time (int): The tPre and tPost of the element
 */
func (se *TraceElementSelect) SetT(time int) {
	se.tPre = time
	se.tPost = time

	se.chosenCase.tPost = time

	for _, c := range se.cases {
		c.tPre = time
	}
}

/*
 * Set the tpre of the element.
 * Args:
 *   tPre (int): The tpre of the element
 */
func (se *TraceElementSelect) SetTPre(tPre int) {
	se.tPre = tPre
	if se.tPost != 0 && se.tPost < tPre {
		se.tPost = tPre
	}

	for _, c := range se.cases {
		c.SetTPre2(tPre)
	}
}

/*
 * Set the tpre of the element. Do not update the chosen case
 * Args:
 *   tPre (int): The tpre of the element
 */
func (se *TraceElementSelect) SetTPre2(tPre int) {
	se.tPre = tPre
	if se.tPost != 0 && se.tPost < tPre {
		se.tPost = tPre
	}

	for _, c := range se.cases {
		c.SetTPre2(tPre)
	}
}

/*
 * Set tPost
 * Args:
 *   tSort (int): The timer of the element
 */
func (se *TraceElementSelect) SetTPost(tPost int) {
	se.tPost = tPost
	se.chosenCase.SetTPost2(tPost)
}

/*
 * Set tPost. Do not update the chosen case
 * Args:
 *   tSort (int): The timer of the element
 */
func (se *TraceElementSelect) SetTPost2(tPost int) {
	se.tPost = tPost
}

/*
 * Set the timer, that is used for the sorting of the trace
 * Args:
 *   tSort (int): The timer of the element
 */
func (se *TraceElementSelect) SetTSort(tSort int) {
	se.SetTPre(tSort)
	se.tPost = tSort
}

/*
 * Set the timer, that is used for the sorting of the trace. Do not update the chosen case
 * Args:
 *   tSort (int): The timer of the element
 */
func (se *TraceElementSelect) SetTSort2(tSort int) {
	se.SetTPre2(tSort)
	se.tPost = tSort
}

/*
 * Set the timer, that is used for the sorting of the trace, only if the original
 * value was not 0
 * Args:
 * tSort (int): The timer of the element
 */
func (se *TraceElementSelect) SetTWithoutNotExecuted(tSort int) {
	se.SetTPre(tSort)
	if se.tPost != 0 {
		se.tPost = tSort
	}
	se.chosenCase.SetTWithoutNotExecuted2(tSort)
}

/*
 * Set the timer, that is used for the sorting of the trace, only if the original
 * value was not 0. Do not update the chosen case
 * Args:
 * tSort (int): The timer of the element
 */
func (se *TraceElementSelect) SetTWithoutNotExecuted2(tSort int) {
	se.SetTPre2(tSort)
	if se.tPost != 0 {
		se.tPost = tSort
	}
}

/*
 * Get the simple string representation of the element
 * MARK: ToString
 * Returns:
 *   string: The simple string representation of the element
 */
func (se *TraceElementSelect) ToString() string {
	res := "S" + "," + strconv.Itoa(se.tPre) + "," +
		strconv.Itoa(se.tPost) + "," + strconv.Itoa(se.id) + ","

	notNil := 0
	for _, ca := range se.cases { // cases
		if ca.tPre != 0 { // ignore nil cases
			if notNil != 0 {
				res += "~"
			}
			res += ca.toStringSep(".", false)
			notNil++
		}
	}

	if se.containsDefault {
		if notNil != 0 {
			res += "~"
		}
		if se.chosenDefault {
			res += "D"
		} else {
			res += "d"
		}
	}
	res += "," + strconv.Itoa(se.chosenIndex)
	res += "," + se.pos
	return res
}

/*
 * Update and calculate the vector clock of the select element.
 * MARK: VectorClock
 */
func (se *TraceElementSelect) updateVectorClock() {
	leak := se.chosenDefault || se.tPost == 0

	if leak {
		currentVCHb[se.routine] = currentVCHb[se.routine].Inc(se.routine)
	} else {
		// update the vector clock
		se.chosenCase.updateVectorClock()
	}

	if analysisCases["selectWithoutPartner"] {
		// check for select case without partner
		ids := make([]int, 0)
		buffered := make([]bool, 0)
		sendInfo := make([]bool, 0)
		for _, c := range se.cases {
			ids = append(ids, c.id)
			buffered = append(buffered, c.qSize > 0)
			sendInfo = append(sendInfo, c.opC == Send)
		}

		analysis.CheckForSelectCaseWithoutPartnerSelect(ids, buffered, sendInfo,
			currentVCHb[se.routine], se.tID, se.chosenIndex)
	}

	if leak {
		se.vc = currentVCHb[se.routine].Copy()
	} else {
		se.vc = se.chosenCase.vc.Copy()
	}

	for _, c := range se.cases {
		c.vc = se.vc.Copy()
	}

	if analysisCases["leak"] {
		for _, c := range se.cases {
			analysis.CheckForLeakChannelRun(se.routine, c.id,
				analysis.VectorClockTID{
					Vc:      se.vc.Copy(),
					TID:     se.tID,
					Routine: se.routine},
				int(c.opC), c.IsBuffered())
		}
	}
}

// MARK: Copy

/*
 * Copy the element
 * Returns:
 *   TraceElement: The copy of the element
 */
func (se *TraceElementSelect) Copy() TraceElement {
	cases := make([]TraceElementChannel, 0)
	for _, c := range se.cases {
		cases = append(cases, *c.Copy().(*TraceElementChannel))
	}

	chosenCase := *se.chosenCase.Copy().(*TraceElementChannel)

	return &TraceElementSelect{
		routine:         se.routine,
		tPre:            se.tPre,
		tPost:           se.tPost,
		id:              se.id,
		cases:           cases,
		chosenCase:      chosenCase,
		chosenIndex:     se.chosenIndex,
		containsDefault: se.containsDefault,
		chosenDefault:   se.chosenDefault,
		pos:             se.pos,
		tID:             se.tID,
		vc:              se.vc.Copy(),
	}
}
