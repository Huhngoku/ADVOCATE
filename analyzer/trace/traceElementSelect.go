package trace

import (
	"analyzer/analysis"
	"errors"
	"math"
	"strconv"
	"strings"
)

/*
 * TraceElementSelect is a trace element for a select statement
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
}

/*
 * Add a new select statement trace element
 * Args:
 *   routine (int): The routine id
 *   tPre (string): The timestamp at the start of the event
 *   tPost (string): The timestamp at the end of the event
 *   id (string): The id of the select statement
 *   cases (string): The cases of the select statement
 *   index (string): The internal index of chosen case
 *   pos (string): The position of the select statement in the code
 */
func AddTraceElementSelect(routine int, tPre string,
	tPost string, id string, cases string, index string, pos string) error {
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

	tID := pos + "@" + tPre

	elem := TraceElementSelect{
		routine: routine,
		tPre:    tPreInt,
		tPost:   tPostInt,
		id:      idInt,
		pos:     pos,
		tID:     tID,
	}

	cs := strings.Split(cases, "~")
	casesList := make([]TraceElementChannel, 0)
	containsDefault := false
	chosenDefault := false
	for _, c := range cs {
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
		cID, err := strconv.Atoi(caseList[3])
		if err != nil {
			return errors.New("c_id is not an integer")
		}
		var cOpC = send
		if caseList[4] == "R" {
			cOpC = recv
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

	elem.chosenIndex, err = strconv.Atoi(index)
	if err != nil {
		return errors.New("index is not an integer")
	}
	elem.containsDefault = containsDefault
	elem.chosenDefault = chosenDefault
	elem.cases = casesList

	if elem.pos == "/home/erikkassubek/Uni/HiWi/Other/examples/gocrawl/popchannel.go:20" {
		println(elem.ToString())
	}

	return AddElementToTrace(&elem)
}

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
 * Set the tpre of the element.
 * Args:
 *   tPre (int): The tpre of the element
 */
func (se *TraceElementSelect) SetTPre(tPre int) {
	se.tPre = tPre
	if se.tPost != 0 && se.tPost < tPre {
		se.tPost = tPre
	}

	se.chosenCase.SetTPre(tPre)
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
 * Set the timer, that is used for the sorting of the trace
 * Args:
 *   tSort (int): The timer of the element
 */
func (se *TraceElementSelect) SetTSort(tSort int) {
	se.tPost = tSort
	se.chosenCase.SetTSort(tSort)
}

/*
 * Set the timer, that is used for the sorting of the trace, only if the original
 * value was not 0
 * Args:
 *   tSort (int): The timer of the element
 */
func (se *TraceElementSelect) SetTSortWithoutNotExecuted(tSort int) {
	if se.tPost != 0 {
		se.tPost = tSort
	}
	se.chosenCase.SetTSortWithoutNotExecuted(tSort)
}

/*
 * Get the simple string representation of the element
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
 * For now, we assume the select acted like the chosen channel operation
 * was just a normal channel operation. For the default, we do not update the vc.
 */
func (se *TraceElementSelect) updateVectorClock() {
	if !se.chosenDefault { // no update for default
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
			sendInfo = append(sendInfo, c.opC == send)
		}

		analysis.CheckForSelectCaseWithoutPartnerSelect(ids, buffered, sendInfo,
			currentVCHb[se.routine], se.tID, se.chosenIndex)
	}
}
