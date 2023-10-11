package trace

import (
	vc "analyzer/vectorClock"
	"errors"
	"math"
	"strconv"
	"strings"
)

/*
 * traceElementSelect is a trace element for a select statement
 * Fields:
 *   routine (int): The routine id
 *   tpre (int): The timestamp at the start of the event
 *   tpost (int): The timestamp at the end of the event
 *   vpost (vectorClock): The vector clock at the end of the event
 *   id (int): The id of the select statement
 *   cases ([]traceElementSelectCase): The cases of the select statement
 *   containsDefault (bool): Whether the select statement contains a default case
 *   chosenCase (traceElementSelectCase): The chosen case, nil if default case chosen
 *   chosenDefault (bool): if the default case was chosen
 *   pos (string): The position of the select statement in the code
 */
type traceElementSelect struct {
	routine int
	tpre    int
	tpost   int
	// vpre            vc.VectorClock
	vpost           vc.VectorClock
	id              int
	cases           []traceElementChannel
	chosenCase      traceElementChannel
	containsDefault bool
	chosenDefault   bool
	pos             string
}

/*
 * Add a new select statement trace element
 * Args:
 *   routine (int): The routine id
 *   numberOfRoutines (int): The number of routines in the trace
 *   tpre (string): The timestamp at the start of the event
 *   tpost (string): The timestamp at the end of the event
 *   id (string): The id of the select statement
 *   cases (string): The cases of the select statement
 *   pos (string): The position of the select statement in the code
 */
func AddTraceElementSelect(routine int, numberOfRoutines int, tpre string,
	tpost string, id string, cases string, pos string) error {
	tpre_int, err := strconv.Atoi(tpre)
	if err != nil {
		return errors.New("tpre is not an integer")
	}

	tpost_int, err := strconv.Atoi(tpost)
	if err != nil {
		return errors.New("tpost is not an integer")
	}

	id_int, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("id is not an integer")
	}

	elem := traceElementSelect{
		routine: routine,
		tpre:    tpre_int,
		tpost:   tpost_int,
		// vpre:    vc.NewVectorClock(numberOfRoutines),
		vpost: vc.NewVectorClock(numberOfRoutines),
		id:    id_int,
		pos:   pos,
	}

	cs := strings.Split(cases, "~")
	cases_list := make([]traceElementChannel, 0)
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
		case_list := strings.Split(c, ".")
		c_tpre, err := strconv.Atoi(case_list[1])
		if err != nil {
			return errors.New("c_tpre is not an integer")
		}
		c_tpost, err := strconv.Atoi(case_list[2])
		if err != nil {
			return errors.New("c_tpost is not an integer")
		}
		c_id, err := strconv.Atoi(case_list[3])
		if err != nil {
			return errors.New("c_id is not an integer")
		}
		var c_opC opChannel = send
		if case_list[4] == "R" {
			c_opC = recv
		} else if case_list[4] == "C" {
			panic("Close in select case list")
		}

		c_cl, err := strconv.ParseBool(case_list[5])
		if err != nil {
			return errors.New("c_cr is not a boolean")
		}

		c_oId, err := strconv.Atoi(case_list[6])
		if err != nil {
			return errors.New("c_oId is not an integer")
		}
		c_oSize, err := strconv.Atoi(case_list[7])
		if err != nil {
			return errors.New("c_oSize is not an integer")
		}

		elem_case := traceElementChannel{
			routine: routine,
			tpre:    c_tpre,
			tpost:   c_tpost,
			id:      c_id,
			opC:     c_opC,
			cl:      c_cl,
			oId:     c_oId,
			qSize:   c_oSize,
			sel:     &elem,
			pos:     pos,
		}

		cases_list = append(cases_list, elem_case)
		if elem_case.tpost != 0 {
			elem.chosenCase = elem_case
		}
	}

	elem.containsDefault = containsDefault
	elem.chosenDefault = chosenDefault
	elem.cases = cases_list

	return addElementToTrace(&elem)
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (se *traceElementSelect) getRoutine() int {
	return se.routine
}

/*
 * Get the timestamp at the start of the event
 * Returns:
 *   int: The timestamp at the start of the event
 */
func (se *traceElementSelect) getTpre() int {
	return se.tpre
}

/*
 * Get the timestamp at the end of the event
 * Returns:
 *   int: The timestamp at the end of the event
 */
func (se *traceElementSelect) getTpost() int {
	return se.tpost
}

/*
 * Get the timer, that is used for the sorting of the trace
 * Returns:
 *   int: The timer of the element
 */
func (se *traceElementSelect) getTsort() int {
	if se.tpost == 0 {
		// add at the end of the trace
		return math.MaxInt
	}
	return se.tpost
}

/*
 * Get the vector clock at the begin of the event
 * Returns:
 *   vectorClock: The vector clock at the begin of the event
 */
// func (se *traceElementSelect) getVpre() *vc.VectorClock {
// 	return &se.vpre
// }

/*
 * Get the vector clock at the end of the event
 * Returns:
 *   vectorClock: The vector clock at the end of the event
 */
func (se *traceElementSelect) getVpost() *vc.VectorClock {
	return &se.vpost
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (se *traceElementSelect) toString() string {
	res := "S" + "," + strconv.Itoa(se.tpre) + "," +
		strconv.Itoa(se.tpost) + "," + strconv.Itoa(se.id) + ","

	for i, c := range se.cases {
		if i != 0 {
			res += "~"
		}
		res += c.toStringSep(".", false)
	}

	if se.containsDefault {
		if se.chosenDefault {
			res += ".D"
		} else {
			res += ".d"
		}
	}
	res += "," + se.pos
	return res
}

/*
 * Update and calculate the vector clock of the select element.
 * For now, we assume the select acted like the chosen channel operation
 * was just a normal channel operation. For the default, we do not update the vc.
 */
func (se *traceElementSelect) updateVectorClock() {
	if se.chosenDefault { // no update for return
		return
	}
	se.chosenCase.updateVectorClock()
}
