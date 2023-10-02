package vectorClock

import (
	"analyzer/logging"
	"strconv"
)

// elements for buffered channel internal vector clock
type bufferedVC struct {
	occupied bool
	vc       VectorClock
}

// vector clock for each buffer place in vector clock
// the map key is the channel id. The slice is used for the buffer positions
var bufferedVCs map[int]([]bufferedVC) = make(map[int]([]bufferedVC))

// the current buffer position
var bufferedVCsCount map[int]int = make(map[int]int)

// vc of close on channel
var closeVC map[int]VectorClock = make(map[int]VectorClock)

// last send and receive on channel
var lastSend map[int]VectorClock = make(map[int]VectorClock)
var lastRecv map[int]VectorClock = make(map[int]VectorClock)

// most recent send, used for detection of send on closed
var hasSend map[int]bool = make(map[int]bool)
var mostRecentSend map[int]VectorClock = make(map[int]VectorClock)
var mostRecentSendPosition map[int]string = make(map[int]string)

/*
 * Update and calculate the vector clocks given a send/receive pair on a unbuffered
 * channel.
 * Args:
 * 	routSend (int): the route of the sender
 * 	routRecv (int): the route of the receiver
 * 	id (int): the id of the channel
 * 	pos (string): the position of the send in the program
 * 	vc (map[int]VectorClock): the current vector clocks
 * Returns:
 * 	the vector clock of the send
 */
func Unbuffered(routSend int, routRecv int, id int, pos string, vc map[int]VectorClock) VectorClock {
	vc[routRecv] = vc[routRecv].Sync(vc[routSend])
	vc[routSend] = vc[routRecv].Copy()

	// for detection of send on closed
	hasSend[id] = true
	mostRecentSend[id] = mostRecentSend[id].Sync(vc[routSend]).Copy()
	mostRecentSendPosition[id] = pos

	logging.Log("Set most recent send of "+strconv.Itoa(id)+" to "+mostRecentSend[id].ToString(), logging.DEBUG)

	vc[routSend] = vc[routSend].Inc(routSend)
	vc[routRecv] = vc[routRecv].Inc(routRecv)
	return vc[routRecv].Copy()
}

/*
 * Update and calculate the vector clocks given a send on a buffered channel.
 * Args:
 * 	rout (int): the route of the sender
 * 	id (int): the id of the sender
 * 	size (int): buffer size
 *  pos (string): the position of the send in the program
 * 	vc (map[int]VectorClock): the current vector clocks
 *  fifo (bool): true if the channel buffer is assumed to be fifo
 * Returns:
 * 	the vector clock of the send
 */
func Send(rout int, id int, size int, pos string,
	vc map[int]VectorClock, fifo bool) VectorClock {
	newBufferedVCs(id, size, vc[rout].size)

	count := bufferedVCsCount[id]
	bufferedVCsCount[id]++
	if count > size || bufferedVCs[id][count].occupied {
		logging.Log("Write to occupied buffer position or to big count", logging.ERROR)
	}

	v := bufferedVCs[id][count].vc
	vc[rout] = vc[rout].Sync(v)
	if fifo {
		vc[rout] = vc[rout].Sync(lastSend[id])
		lastSend[id] = vc[rout].Copy()
	}
	bufferedVCs[id][count] = bufferedVC{true, vc[rout].Copy()}

	// for detection of send on closed
	hasSend[id] = true
	mostRecentSend[id] = mostRecentSend[id].Sync(vc[rout])
	mostRecentSendPosition[id] = pos

	vc[rout] = vc[rout].Inc(rout)
	return vc[rout].Copy()
}

/*
 * Update and calculate the vector clocks given a receive on a buffered channel.
 * Args:
 * 	rout (int): the route of the sender
 * 	id (int): the id of the sender
 * 	size (int): buffer size
 * 	vc (map[int]VectorClock): the current vector clocks
 *  fifo (bool): true if the channel buffer is assumed to be fifo
 * Returns:
 * 	the vector clock of the receive
 */
func Recv(rout int, id int, size int, vc map[int]VectorClock, fifo bool) VectorClock {
	newBufferedVCs(id, size, vc[rout].size)
	if bufferedVCsCount[id] == 0 {
		logging.Log("Read operation on empty buffer position", logging.ERROR)
	}
	bufferedVCsCount[id]--

	v := bufferedVCs[id][0].vc
	vc[rout] = vc[rout].Sync(v)
	if fifo {
		vc[rout] = vc[rout].Sync(lastRecv[id])
		lastRecv[id] = vc[rout].Copy()
	}
	bufferedVCs[id] = bufferedVCs[id][1:]
	bufferedVCs[id] = append(bufferedVCs[id], bufferedVC{false, vc[rout].Copy()})

	vc[rout] = vc[rout].Inc(rout)
	return vc[rout].Copy()
}

/*
 * Update and calculate the vector clocks given a close on a channel.
 * Args:
 * 	rout (int): the route of the sender
 * 	id (int): the id of the sender
 * 	pos (string): the position of the close in the program
 * 	vc (map[int]VectorClock): the current vector clocks
 * Returns:
 * 	the vector clock of the close
 */
func Close(rout int, id int, pos string, vc map[int]VectorClock) VectorClock {
	closeVC[id] = vc[rout].Copy()

	// check if there is an earlier send, that could happen concurrently to close
	if hasSend[id] {
		logging.Log("Check for possible send on closed channel "+
			strconv.Itoa(id)+" with "+
			mostRecentSend[id].ToString()+" and "+closeVC[id].ToString(),
			logging.DEBUG)
		happensBefore := GetHappensBefore(closeVC[id], mostRecentSend[id])
		if happensBefore == Concurrent {
			found := "Possible send on closed channel:\n"
			found += "\tclose: " + pos + "\n"
			found += "\tsend: " + mostRecentSendPosition[id]
			logging.Log(found, logging.RESULT)
		}
	}

	vc[rout] = vc[rout].Inc(rout)
	return vc[rout].Copy()
}

/*
 * Update and calculate the vector clocks given a receive on a closed channel.
 * Args:
 * 	rout (int): the route of the sender
 * 	id (int): the id of the sender
 * 	vc (map[int]VectorClock): the current vector clocks
 * Returns:
 * 	the vector clock of the close
 */
func RecvC(rout int, id int, vc map[int]VectorClock) VectorClock {
	vc[rout] = vc[rout].Sync(closeVC[id])
	vc[rout] = vc[rout].Inc(rout)
	return vc[rout].Copy()
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
			bufferedVCs[id][i] = bufferedVC{false, NewVectorClock(numRout)}
		}
	}
}
