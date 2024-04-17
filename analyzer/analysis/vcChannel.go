package analysis

import (
	"analyzer/clock"
	"analyzer/logging"
	"strconv"
)

// elements for buffered channel internal vector clock
type bufferedVC struct {
	occupied    bool
	oID         int
	vc          clock.VectorClock
	routineSend int
	tID         string
}

/*
 * Update and calculate the vector clocks given a send/receive pair on a unbuffered
 * channel.
 * Args:
 * 	routSend (int): the route of the sender
 * 	routRecv (int): the route of the receiver
 * 	id (int): the id of the channel
 * 	tID_send (string): the position of the send in the program
 * 	tID_recv (string): the position of the receive in the program
 * 	vc (map[int]VectorClock): the current vector clocks
 *  tPost (int): the timestamp at the end of the event
 */
func Unbuffered(routSend int, routRecv int, id int, tIDSend string,
	tIDRecv string, vc map[int]clock.VectorClock, tPost int) {
	if tPost != 0 {
		if analysisCases["concurrentRecv"] {
			checkForConcurrentRecv(routRecv, id, tIDRecv, vc)
		}

		vc[routSend] = vc[routSend].Inc(routSend)
		vc[routRecv] = vc[routRecv].Inc(routRecv)
		vc[routRecv] = vc[routRecv].Sync(vc[routSend])
		vc[routSend] = vc[routRecv].Copy()

		// for detection of send on closed
		hasSend[id] = true
		mostRecentSend[id] = VectorClockTID{mostRecentSend[id].vc.Sync(vc[routSend]).Copy(), tIDSend}

		// for detection of receive on closed
		hasReceived[id] = true
		mostRecentReceive[id] = VectorClockTID{mostRecentReceive[id].vc.Sync(vc[routRecv]).Copy(), tIDRecv}

		logging.Debug("Set most recent send of "+strconv.Itoa(id)+" to "+mostRecentSend[id].vc.ToString(), logging.DEBUG)

	} else {
		vc[routSend] = vc[routSend].Inc(routSend)
	}

	if analysisCases["mixedDeadlock"] {
		CheckForSelectCaseWithoutPartnerChannel(id, vc[routSend], tIDSend, true, false)
		CheckForSelectCaseWithoutPartnerChannel(id, vc[routRecv], tIDRecv, false, false)
		checkForMixedDeadlock(routSend, routRecv, tIDSend, tIDRecv)
	}

	if analysisCases["selectWithoutPartner"] {
		CheckForSelectCaseWithoutPartnerChannel(id, vc[routSend], tIDSend, true, false)
		CheckForSelectCaseWithoutPartnerChannel(id, vc[routRecv], tIDRecv, false, false)
	}

	if analysisCases["leak"] {
		CheckForLeakChannelRun(id, VectorClockTID{vc[routSend].Copy(), tIDSend}, 0)
		CheckForLeakChannelRun(id, VectorClockTID{vc[routRecv].Copy(), tIDSend}, 1)
	}

}

/*
 * Update and calculate the vector clocks given a send on a buffered channel.
 * Args:
 * 	rout (int): the route of the sender
 * 	id (int): the id of the sender
 * 	oID (int): the id of the communication
 * 	size (int): buffer size
 *  tId (string): the position of the send in the program
 * 	vc (map[int]VectorClock): the current vector clocks
 *  fifo (bool): true if the channel buffer is assumed to be fifo
 *  tPost (int): the timestamp at the end of the event
 */
func Send(rout int, id int, oID int, size int, tID string,
	vc map[int]clock.VectorClock, fifo bool, tPost int) {

	if tPost == 0 {
		vc[rout] = vc[rout].Inc(rout)
		return
	}

	newBufferedVCs(id, size, vc[rout].GetSize())

	count := bufferedVCsCount[id]

	if len(bufferedVCs[id]) <= count {
		panic("BufferedVCsCount is bigger than the buffer size")
	}

	if count > size || bufferedVCs[id][count].occupied {
		logging.Debug("Write to occupied buffer position or to big count", logging.ERROR)
	}

	v := bufferedVCs[id][count].vc
	vc[rout] = vc[rout].Sync(v)

	if fifo {
		vc[rout] = vc[rout].Sync(mostRecentSend[id].vc)
		mostRecentSend[id] = VectorClockTID{vc[rout].Copy(), mostRecentSend[id].tID}
	}

	bufferedVCs[id][count] = bufferedVC{true, oID, vc[rout].Copy(), rout, tID}

	bufferedVCsCount[id]++

	// for detection of send on closed
	hasSend[id] = true
	mostRecentSend[id] = VectorClockTID{mostRecentSend[id].vc.Sync(vc[rout]), tID}

	vc[rout] = vc[rout].Inc(rout)

	if analysisCases["selectWithoutPartner"] {
		CheckForSelectCaseWithoutPartnerChannel(id, vc[rout], tID, true, true)
	}

	if analysisCases["leak"] {
		CheckForLeakChannelRun(id, VectorClockTID{vc[rout].Copy(), tID}, 0)
	}

}

/*
 * Update and calculate the vector clocks given a receive on a buffered channel.
 * Args:
 * 	rout (int): the route of the sender
 * 	id (int): the id of the sender
 * 	oId (int): the id of the communication
 * 	size (int): buffer size
 *  tID (string): the position of the send in the program
 * 	vc (map[int]VectorClock): the current vector clocks
 *  fifo (bool): true if the channel buffer is assumed to be fifo
 *  tPost (int): the timestamp at the end of the event
 */
func Recv(rout int, id int, oID, size int, tID string, vc map[int]clock.VectorClock,
	fifo bool, tPost int) {
	if tPost == 0 {
		vc[rout] = vc[rout].Inc(rout)
		return
	}

	newBufferedVCs(id, size, vc[rout].GetSize())

	if analysisCases["concurrentRecv"] {
		checkForConcurrentRecv(rout, id, tID, vc)
	}

	if bufferedVCsCount[id] == 0 {
		logging.Debug("Read operation on empty buffer position", logging.ERROR)
	}
	bufferedVCsCount[id]--

	if bufferedVCs[id][0].oID != oID {
		found := false
		for i := 1; i < size; i++ {
			if bufferedVCs[id][i].oID == oID {
				found = true
				bufferedVCs[id][0] = bufferedVCs[id][i]
				bufferedVCs[id][i] = bufferedVC{false, 0, vc[rout].Copy(), 0, ""}
				break
			}
		}
		if !found {
			err := "Read operation on wrong buffer position - ID: " + strconv.Itoa(id) + ", OID: " + strconv.Itoa(oID) + ", SIZE: " + strconv.Itoa(size)
			logging.Debug(err, logging.ERROR)
		}
	}
	v := bufferedVCs[id][0].vc
	routSend := bufferedVCs[id][0].routineSend
	tIDSend := bufferedVCs[id][0].tID

	vc[rout] = vc[rout].Sync(v)

	if fifo {
		vc[rout] = vc[rout].Sync(mostRecentReceive[id].vc)
		mostRecentReceive[id] = VectorClockTID{vc[rout].Copy(), mostRecentReceive[id].tID}
	}

	bufferedVCs[id] = bufferedVCs[id][1:]
	bufferedVCs[id] = append(bufferedVCs[id], bufferedVC{false, 0, vc[rout].Copy(), 0, ""})

	// for detection of receive on closed
	hasReceived[id] = true
	mostRecentReceive[id] = VectorClockTID{mostRecentReceive[id].vc.Sync(vc[rout]), tID}

	vc[rout] = vc[rout].Inc(rout)

	if analysisCases["selectWithoutPartner"] {
		CheckForSelectCaseWithoutPartnerChannel(id, vc[rout], tID, true, true)
	}

	if analysisCases["mixedDeadlock"] {
		checkForMixedDeadlock(routSend, rout, tIDSend, tID)
	}
	if analysisCases["leak"] {
		CheckForLeakChannelRun(id, VectorClockTID{vc[rout].Copy(), tID}, 1)
	}
}

/*
 * Update and calculate the vector clocks given a close on a channel.
 * Args:
 * 	rout (int): the route of the operation
 * 	id (int): the id of the channel
 * 	tID (string): the position of the close in the program
 * 	vc (map[int]VectorClock): the current vector clocks
 *  tPost (int): the timestamp at the end of the event
 */
func Close(rout int, id int, tID string, vc map[int]clock.VectorClock, tPost int) {
	if tPost == 0 {
		return
	}

	if analysisCases["closeOnClosed"] {
		checkForClosedOnClosed(id, tID) // must be called before closePos is updated
	}

	closeData[id] = VectorClockTID{vc[rout].Copy(), tID}
	closeRout[id] = rout

	if analysisCases["sendOnClosed"] || analysisCases["receiveOnClosed"] {
		checkForCommunicationOnClosedChannel(id, tID)
	}

	vc[rout] = vc[rout].Inc(rout)

	if analysisCases["selectWithoutPartner"] {
		CheckForSelectCaseWithoutPartnerClose(id, vc[rout])
	}

	if analysisCases["leak"] {
		CheckForLeakChannelRun(id, VectorClockTID{vc[rout].Copy(), tID}, 2)
	}
}

/*
 * Update and calculate the vector clocks given a receive on a closed channel.
 * Args:
 * 	rout (int): the route of the sender
 * 	id (int): the id of the sender
 * 	tID (string): the position of the close in the program
 * 	vc (map[int]VectorClock): the current vector clocks
 *  tPost (int): the timestamp at the end of the event
 *  buffered (bool): true if the channel is buffered
 */
func RecvC(rout int, id int, tID string, vc map[int]clock.VectorClock, tPost int,
	buffered bool) {
	if tPost == 0 {
		return
	}

	foundReceiveOnClosedChannel(closeData[id].tID, tID)

	vc[rout] = vc[rout].Sync(closeData[id].vc)
	vc[rout] = vc[rout].Inc(rout)

	if analysisCases["selectWithoutPartner"] {
		CheckForSelectCaseWithoutPartnerChannel(id, vc[rout], tID, false, buffered)
	}

	if analysisCases["mixedDeadlock"] {
		checkForMixedDeadlock(closeRout[id], rout, closeData[id].tID, tID)
	}
	if analysisCases["leak"] {
		CheckForLeakChannelRun(id, VectorClockTID{vc[rout].Copy(), tID}, 1)
	}
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
			bufferedVCs[id][i] = bufferedVC{false, 0, clock.NewVectorClock(numRout), 0, ""}
		}
	}
}
