package trace

import (
	"errors"
	"strconv"
	"strings"
)

/*
 * traceElementSelectCase is a trace element for a select statement
 * Fields:
 *   channel (int): The id of the channel
 *   op (int, enum): The operation on the channel
 */
type traceElementSelectCase struct {
	channel int
	op      opChannel
}

func (elem traceElementSelectCase) toString() string {
	return strconv.Itoa(elem.channel) + "." + strconv.Itoa(int(elem.op))
}

/*
 * traceElementSelect is a trace element for a select statement
 * Fields:
 *   routine (int): The routine id
 *   tpre (int): The timestamp at the start of the event
 *   tpost (int): The timestamp at the end of the event
 *   id (int): The id of the select statement
 *   cases ([]traceElementSelectCase): The cases of the select statement
 *   containsDefault (bool): Whether the select statement contains a default case
 *   exec (int, enum): The execution status of the operation
 *   chosend (traceElementSelectCase): The case that was chosen
 *   oId (int): The id of the communication
 *   pos (string): The position of the select statement in the code
 */
type traceElementSelect struct {
	routine         int
	tpre            int
	tpost           int
	id              int
	cases           []traceElementSelectCase
	containsDefault bool
	exec            bool
	chosend         traceElementSelectCase
	oId             int
	pos             string
}

func AddTraceElementSelect(routine int, tpre string, tpost string, id string,
	cases string, exec string, chosen string, oId string, pos string) error {
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

	cs := strings.Split(cases, ".")
	containsDefault := false
	cases_list := make([]traceElementSelectCase, 0)
	for _, c := range cs {
		if c == "d" {
			containsDefault = true
			break
		}

		channelId, err := strconv.Atoi(c[:len(c)-1])
		if err != nil {
			return err
		}

		op := c[len(c)-1]
		var opC opChannel = 0
		switch op {
		case 'r':
			opC = recv
		case 's':
			opC = send
		default:
			return errors.New("op is not a valid operation")
		}

		elem := traceElementSelectCase{
			channel: channelId,
			op:      opC,
		}

		cases_list = append(cases_list, elem)
	}

	exec_bool, err := strconv.ParseBool(exec)
	if err != nil {
		return errors.New("exec is not a boolean")
	}

	chosen_int, err := strconv.Atoi(chosen)
	if err != nil {
		return errors.New("chosen is not an integer")
	}

	chosenCase := cases_list[chosen_int]

	oId_int, err := strconv.Atoi(oId)
	if err != nil {
		return errors.New("oId is not an integer")
	}

	elem := traceElementSelect{routine, tpre_int, tpost_int, id_int, cases_list,
		containsDefault, exec_bool, chosenCase, oId_int, pos}

	return addElementToTrace(routine, elem)
}

/*
 * Get the routine of the element
 * Returns:
 *   int: The routine of the element
 */
func (elem traceElementSelect) getRoutine() int {
	return elem.routine
}

/*
 * Get the timestamp at the start of the event
 * Returns:
 *   int: The timestamp at the start of the event
 */
func (elem traceElementSelect) getTpre() int {
	return elem.tpre
}

/*
 * Get the timestamp at the end of the event
 * Returns:
 *   int: The timestamp at the end of the event
 */
func (elem traceElementSelect) getTpost() int {
	return elem.tpost
}

/*
 * Get the simple string representation of the element
 * Returns:
 *   string: The simple string representation of the element
 */
func (elem traceElementSelect) getSimpleString() string {
	res := "S" + strconv.Itoa(elem.id) + "," + strconv.Itoa(elem.tpre) + "," +
		strconv.Itoa(elem.tpost) + " "

	for i, c := range elem.cases {
		if i != 0 {
			res += "."
		}
		res += c.toString()
	}

	if elem.containsDefault {
		res += ".d"
	}

	res += "," + strconv.Itoa(elem.oId) + "," + elem.chosend.toString()
	return res
}
