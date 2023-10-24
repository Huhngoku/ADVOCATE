package vectorClock

import (
	"analyzer/logging"
	"strconv"
)

// elements for buffered channel internal vector clock
type bufferedVC struct {
	occupied bool
	oId      int
	vc       VectorClock
}

// vector clock for each buffer place in vector clock
// the map key is the channel id. The slice is used for the buffer positions
var bufferedVCs map[int]([]bufferedVC) = make(map[int]([]bufferedVC))

// the current buffer position
var bufferedVCsCount map[int]int = make(map[int]int)

// vc of close on channel
var closeVC map[int]VectorClock = make(map[int]VectorClock)
var closePos map[int]string = make(map[int]string)

// last send and receive on channel
var lastSend map[int]VectorClock = make(map[int]VectorClock)
var lastRecv map[int]VectorClock = make(map[int]VectorClock)

// last receive for each routine and each channel
var lastRecvRoutine map[int]map[int]VectorClock = make(map[int]map[int]VectorClock)
var lastRecvRoutinePos map[int]map[int]string = make(map[int]map[int]string)

// most recent send, used for detection of send on closed
var hasSend map[int]bool = make(map[int]bool)
var mostRecentSend map[int]VectorClock = make(map[int]VectorClock)
var mostRecentSendPosition map[int]string = make(map[int]string)

// most recent send, used for detection of received on closed
var hasReceived map[int]bool = make(map[int]bool)
var mostRecentReceive map[int]VectorClock = make(map[int]VectorClock)
var mostRecentReceivePosition map[int]string = make(map[int]string)

/*
 * Update and calculate the vector clocks given a send/receive pair on a unbuffered
 * channel.
 * Args:
 * 	routSend (int): the route of the sender
 * 	routRecv (int): the route of the receiver
 * 	id (int): the id of the channel
 * 	pos_send (string): the position of the send in the program
 * 	pos_recv (string): the position of the receive in the program
 * 	vc (map[int]VectorClock): the current vector clocks
 */
func Unbuffered(routSend int, routRecv int, id int, pos_send string, pos_recv string, vc map[int]VectorClock) {
	// checkForConcurrentRecv(routRecv, id, pos_recv, vc)

	vc[routRecv] = vc[routRecv].Sync(vc[routSend])
	vc[routSend] = vc[routRecv].Copy()

	// for detection of send on closed
	hasSend[id] = true
	mostRecentSend[id] = mostRecentSend[id].Sync(vc[routSend]).Copy()
	mostRecentSendPosition[id] = pos_send

	// for detection of receive on closed
	hasReceived[id] = true
	mostRecentReceive[id] = mostRecentReceive[id].Sync(vc[routRecv]).Copy()
	mostRecentReceivePosition[id] = pos_recv

	logging.Debug("Set most recent send of "+strconv.Itoa(id)+" to "+mostRecentSend[id].ToString(), logging.DEBUG)

	vc[routSend] = vc[routSend].Inc(routSend)
	vc[routRecv] = vc[routRecv].Inc(routRecv)
}

/*
 * Update and calculate the vector clocks given a send on a buffered channel.
 * Args:
 * 	rout (int): the route of the sender
 * 	id (int): the id of the sender
 * 	oId (int): the id of the communication
 * 	size (int): buffer size
 *  pos (string): the position of the send in the program
 * 	vc (map[int]VectorClock): the current vector clocks
 *  fifo (bool): true if the channel buffer is assumed to be fifo
 */
func Send(rout int, id int, oId int, size int, pos string,
	vc map[int]VectorClock, fifo bool) {
	newBufferedVCs(id, size, vc[rout].size)

	count := bufferedVCsCount[id]
	bufferedVCsCount[id]++
	if count > size || bufferedVCs[id][count].occupied {
		logging.Debug("Write to occupied buffer position or to big count", logging.ERROR)
	}

	v := bufferedVCs[id][count].vc
	vc[rout] = vc[rout].Sync(v)

	if fifo {
		vc[rout] = vc[rout].Sync(lastSend[id])
		lastSend[id] = vc[rout].Copy()
	}

	bufferedVCs[id][count] = bufferedVC{true, oId, vc[rout].Copy()}

	// for detection of send on closed
	hasSend[id] = true
	mostRecentSend[id] = mostRecentSend[id].Sync(vc[rout])
	mostRecentSendPosition[id] = pos

	vc[rout] = vc[rout].Inc(rout)
}

/*
 * Update and calculate the vector clocks given a receive on a buffered channel.
 * Args:
 * 	rout (int): the route of the sender
 * 	id (int): the id of the sender
 * 	oId (int): the id of the communication
 * 	size (int): buffer size
 *  pos (string): the position of the send in the program
 * 	vc (map[int]VectorClock): the current vector clocks
 *  fifo (bool): true if the channel buffer is assumed to be fifo
 */
func Recv(rout int, id int, oId, size int, pos string, vc map[int]VectorClock, fifo bool) {
	newBufferedVCs(id, size, vc[rout].size)

	// checkForConcurrentRecv(rout, id, pos, vc)

	if bufferedVCsCount[id] == 0 {
		logging.Debug("Read operation on empty buffer position", logging.ERROR)
	}
	bufferedVCsCount[id]--

	if bufferedVCs[id][0].oId != oId {
		logging.Debug("Read operation on wrong buffer position", logging.ERROR)
	}
	v := bufferedVCs[id][0].vc

	vc[rout] = vc[rout].Sync(v)
	if fifo {
		vc[rout] = vc[rout].Sync(lastRecv[id])
		lastRecv[id] = vc[rout].Copy()
	}
	bufferedVCs[id] = bufferedVCs[id][1:]
	bufferedVCs[id] = append(bufferedVCs[id], bufferedVC{false, 0, vc[rout].Copy()})

	// for detection of receive on closed
	hasReceived[id] = true
	mostRecentReceive[id] = mostRecentReceive[id].Sync(vc[rout])
	mostRecentReceivePosition[id] = pos

	vc[rout] = vc[rout].Inc(rout)
}

/*
 * Update and calculate the vector clocks given a close on a channel.
 * Args:
 * 	rout (int): the route of the sender
 * 	id (int): the id of the sender
 * 	pos (string): the position of the close in the program
 * 	vc (map[int]VectorClock): the current vector clocks
 */
func Close(rout int, id int, pos string, vc map[int]VectorClock) {
	closeVC[id] = vc[rout].Copy()
	closePos[id] = pos

	CheckForPotentialCommunicationOnClosedChannel(id, pos)

	vc[rout] = vc[rout].Inc(rout)
}

/*
 * Update and calculate the vector clocks given a receive on a closed channel.
 * Args:
 * 	rout (int): the route of the sender
 * 	id (int): the id of the sender
 * 	pos (string): the position of the close in the program
 * 	vc (map[int]VectorClock): the current vector clocks
 */
func RecvC(rout int, id int, pos string, vc map[int]VectorClock) {
	found := "Receive on closed channel:\n"
	found += "\tclose: " + closePos[id] + "\n"
	found += "\trecv : " + pos
	logging.Result(found, logging.WARNING)

	vc[rout] = vc[rout].Sync(closeVC[id])
	vc[rout] = vc[rout].Inc(rout)
}

/*
 * Create a new map of buffered vector clocks for a channel if not already in
 * bufferedVCs.
 * Args:
 * 	id (int): the id of the channel
 * 	size (int): the buffer size of the channel
 * 	numRout (int): the number of routines
 */
func newBufferedVCs(id int, size int, numRout int) {
	if _, ok := bufferedVCs[id]; !ok {
		bufferedVCs[id] = make([]bufferedVC, size)
		for i := 0; i < size; i++ {
			bufferedVCsCount[id] = 0
			bufferedVCs[id][i] = bufferedVC{false, 0, NewVectorClock(numRout)}
		}
	}
}

/*
Check if a send or receive on a closed channel is possible
It it is possible, print a warning or error
Args:

	id (int): the id of the channel
	pos (string): the position of the close in the program
*/
func CheckForPotentialCommunicationOnClosedChannel(id int, pos string) {
	// check if there is an earlier send, that could happen concurrently to close
	if hasSend[id] {
		logging.Debug("Check for possible send on closed channel "+
			strconv.Itoa(id)+" with "+
			mostRecentSend[id].ToString()+" and "+closeVC[id].ToString(),
			logging.DEBUG)
		happensBefore := GetHappensBefore(closeVC[id], mostRecentSend[id])
		if happensBefore == Concurrent {
			found := "Possible send on closed channel:\n"
			found += "\tclose: " + pos + "\n"
			found += "\tsend : " + mostRecentSendPosition[id]
			logging.Result(found, logging.CRITICAL)
		}
	}
	// check if there is an earlier receive, that could happen concurrently to close
	if hasReceived[id] {
		logging.Debug("Check for possible receive on closed channel "+
			strconv.Itoa(id)+" with "+
			mostRecentReceive[id].ToString()+" and "+closeVC[id].ToString(),
			logging.DEBUG)
		happensBefore := GetHappensBefore(closeVC[id], mostRecentReceive[id])
		if happensBefore == Concurrent || happensBefore == Before {
			found := "Possible receive on closed channel:\n"
			found += "\tclose: " + pos + "\n"
			found += "\trecv : " + mostRecentReceivePosition[id]
			logging.Result(found, logging.WARNING)
		}
	}

}

func checkForConcurrentRecv(routine int, id int, pos string, vc map[int]VectorClock) {
	if _, ok := lastRecvRoutine[routine]; !ok {
		lastRecvRoutine[routine] = make(map[int]VectorClock)
		lastRecvRoutinePos[routine] = make(map[int]string)
	}

	lastRecvRoutine[routine][id] = vc[routine].Copy()
	lastRecvRoutinePos[routine][id] = pos

	for r, elem := range lastRecvRoutine {
		// logging.Debug(lastRecvRoutinePos[r][id]+" "+pos, logging.ERROR)
		// logging.Debug(elem[id].ToString()+" "+vc[routine].ToString(), logging.ERROR)
		if r == routine {
			continue
		}

		if elem[id].clock == nil {
			continue
		}

		happensBefore := GetHappensBefore(elem[id], vc[routine])
		if happensBefore == Concurrent {
			found := "Found concurrent Recv on same channel:\n"
			found += "\trecv: " + pos + "\n"
			found += "\trecv : " + lastRecvRoutinePos[r][id]
			logging.Result(found, logging.CRITICAL)
		}
	}
}
