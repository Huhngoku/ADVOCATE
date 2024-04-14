package analysis

import (
	"analyzer/clock"
	"analyzer/logging"
)

// TODO: make this work with buffered channels
/*
 * Run for channel operation without a post event. Check if the operation has
 * a possible communication partner in mostRecentSend, mostRecentReceive or closeData.
 * If so, add an error or warning to the result.
 * If not, add to leakingChannels, for later check.
 * MARK: Channel Stuck
 * Args:
 *   id (int): The channel id
 *   vc (VectorClock): The vector clock of the operation
 *   tID (string): The trace id
 *   opType (int): An identifier for the type of the operation (send = 0, recv = 1)
 *   buffered (bool): If the channel is buffered
 */
func CheckForLeakChannelStuck(id int, vc clock.VectorClock, tID string, opType int,
	buffered bool) {
	logging.Debug("Checking channel for for leak channel", logging.INFO)

	if !buffered {
		foundPartner := false

		if opType == 0 { // send
			if _, ok := mostRecentReceive[id]; ok {
				if clock.GetHappensBefore(mostRecentReceive[id].vc, vc) == clock.Concurrent {
					found := "Leak on unbuffered channel with possible partner:\n"
					found += "\tchannel: " + tID + "\n"
					found += "\tpartner: " + mostRecentReceive[id].tID
					logging.Result(found, logging.CRITICAL)
					foundPartner = true
				}
			}
		} else if opType == 1 { // recv
			if _, ok := mostRecentSend[id]; ok {
				if clock.GetHappensBefore(mostRecentSend[id].vc, vc) == clock.Concurrent {
					found := "Leak on unbuffered channel with possible partner:\n"
					found += "\tchannel: " + tID + "\n"
					found += "\tpartner: " + mostRecentSend[id].tID
					logging.Result(found, logging.CRITICAL)
					foundPartner = true
				}
			}

			// // This cannot happen:
			// if _, ok := closeData[id]; ok {
			// 	found := "Leak on unbuffered channel with possible partner:\n"
			// 	found += "\tchannel: " + tID + "\n"
			// 	found += "\tpartner: " + closeData[id].tID
			// 	logging.Result(found, logging.CRITICAL)
			// 	foundPartner = true
			// }
		}

		if !foundPartner {
			leakingChannels[id] = append(leakingChannels[id], VectorClockTID2{id, vc, tID, opType, -1})
		}
	} else {
		found := "Leak on buffered channel:\n"
		found += "\tchannel: " + tID + "\n"
		found += "\t"
		logging.Result(found, logging.CRITICAL)
	}
}

/*
 * Run for channel operation with a post event. Check if the operation would be
 * possible communication partner for a stuck operation in leakingChannels.
 * If so, add an error or warning to the result and remove the stuck operation.
 * MARK: ChannelRun
 * Args:
 *   id (int): The channel id
 *   vc (VectorClock): The vector clock of the operation
 *   tID (string): The trace id
 *   opType (int): An identifier for the type of the operation (send = 0, recv = 1, close = 2)
 */
func CheckForLeakChannelRun(id int, vcTID VectorClockTID, opType int) bool {
	logging.Debug("Checking channel for for leak channel", logging.INFO)
	res := false
	if opType == 0 || opType == 2 { // send or close
		for i, vcTID2 := range leakingChannels[id] {
			if vcTID2.val != 1 {
				continue
			}
			if clock.GetHappensBefore(vcTID2.vc, vcTID.vc) == clock.Concurrent {
				found := "Leak on unbuffered channel with possible partner:\n"
				found += "\tchannel: " + vcTID2.tID + "\n"
				found += "\tpartner: " + vcTID.tID
				logging.Result(found, logging.CRITICAL)
				res = true

				// remove the stuck operation from the list. If it is a select, remove all operations with the same val
				if vcTID2.val == -1 {
					leakingChannels[id] = append(leakingChannels[id][:i], leakingChannels[id][i+1:]...)
				} else {
					for j, vcTID3 := range leakingChannels[id] {
						if vcTID3.val == vcTID2.val {
							leakingChannels[id] = append(leakingChannels[id][:j], leakingChannels[id][j+1:]...)
						}
					}
				}
			}
		}
	} else if opType == 1 { // recv
		for i, vcTID2 := range leakingChannels[id] {
			if vcTID2.val != 0 && vcTID2.val != 2 {
				continue
			}
			if clock.GetHappensBefore(vcTID2.vc, vcTID.vc) == clock.Concurrent {
				found := "Leak on unbuffered channel with possible partner:\n"
				found += "\tchannel: " + vcTID2.tID + "\n"
				found += "\tpartner: " + vcTID.tID
				logging.Result(found, logging.CRITICAL)
				res = true

				// remove the stuck operation from the list. If it is a select, remove all operations with the same val
				if vcTID2.val == -1 {
					leakingChannels[id] = append(leakingChannels[id][:i], leakingChannels[id][i+1:]...)
				} else {
					for j, vcTID3 := range leakingChannels[id] {
						if vcTID3.val == vcTID2.val {
							leakingChannels[id] = append(leakingChannels[id][:j], leakingChannels[id][j+1:]...)
						}
					}
				}
			}
		}
	}
	return res
}

/*
 * After all operations have been analyzed, check if there are still leaking
 * operations without a possible partner.
 */
func CheckForLeak() {
	// channel
	for _, vcTIDs := range leakingChannels {
		for _, vcTID := range vcTIDs {
			if vcTID.tID == "" {
				continue
			}

			found := false
			for _, c := range selectCases {
				if c.send {
					println(vcTID.tID, c.vcTID.tID, c.send, vcTID.id, c.id)
				}
				if c.id != vcTID.id {
					continue
				}

				if (c.send && vcTID.typeVal == 0) || (!c.send && vcTID.typeVal == 1) {
					continue
				}

				hb := clock.GetHappensBefore(c.vcTID.vc, vcTID.vc)
				if hb == clock.Concurrent {
					found = true
					break
				}

				if c.buffered {
					{
						if (c.send && hb == clock.Before) || (!c.send && hb == clock.After) {
							found = true
							break
						}
					}
				}
			}

			if !found {
				found := "Leak on channel without possible partner:\n"
				found += "\tchannel: " + vcTID.tID + "\n"
				found += "\tpartner: -"
				logging.Result(found, logging.CRITICAL)
			}
		}
	}
}

/*
 * Run for select operation without a post event. Check if the operation has
 * a possible communication partner in mostRecentSend, mostRecentReceive or closeData.
 * If so, add an error or warning to the result.
 * If not, add all elements to leakingChannels, for later check.
 * MARK: SelectStuck
 * Args:
 *   ids (int): The channel ids
 *   vc (VectorClock): The vector clock of the operation
 *   tID (string): The trace id
 *   opTypes ([]int): An identifier for the type of the operations (send = 0, recv = 1)
 *   idSel (int): The id of the select operation
 *   tPre (int): The tpre of the select operations. Used to connect the operations of the
 *     same select statement in leakingChannels.
 */
func CheckForLeakSelectStuck(ids []int, vc clock.VectorClock, tID string, opTypes []int, tPre int) {
	foundPartner := false
	for i, id := range ids {
		if opTypes[i] == 0 { // send
			if _, ok := mostRecentReceive[id]; ok {
				if clock.GetHappensBefore(vc, mostRecentReceive[id].vc) == clock.Concurrent {
					found := "Leak on unbuffered channel with possible partner:\n"
					found += "\tchannel: " + tID + "\n"
					found += "\tpartner: " + mostRecentReceive[id].tID + "\n"
					logging.Result(found, logging.CRITICAL)
					foundPartner = true
				}
			}
		} else if opTypes[i] == 1 { // recv
			if _, ok := mostRecentSend[id]; ok {
				if clock.GetHappensBefore(vc, mostRecentSend[id].vc) == clock.Concurrent {
					found := "Leak on unbuffered channel with possible partner:\n"
					found += "\tchannel: " + tID + "\n"
					found += "\tpartner: " + mostRecentSend[id].tID
					logging.Result(found, logging.CRITICAL)
					foundPartner = true
				}
			}
			if _, ok := closeData[id]; ok {
				found := "Leak on unbuffered channel with possible partner:\n"
				found += "\tchannel: " + tID + "\n"
				found += "\tpartner: " + closeData[id].tID
				logging.Result(found, logging.CRITICAL)
				foundPartner = true
			}
		}
	}

	if !foundPartner {
		for i, id := range ids {
			// add all select operations to leaking Channels,
			leakingChannels[id] = append(leakingChannels[id], VectorClockTID2{id, vc, tID, opTypes[i], tPre})
		}
	}
}

/*
 * Run for select operation with a post event. Check if the operation would be
 * possible communication partner for a stuck operation in leakingChannels.
 * If so, add an error or warning to the result and remove the stuck operation.
 * MARK: SelectRun
 * Args:
 *   id (int): The channel id
 *   vc (VectorClock): The vector clock of the operation
 *   tID (string): The trace id
 *   opType (int): An identifier for the type of the operation (send = 0, recv = 1, close = 2)
 */
func CheckForLeakSelectRun(ids []int, typeIds []int, vc clock.VectorClock, tID string) {
	for i, id := range ids {
		if CheckForLeakChannelRun(id, VectorClockTID{vc, tID}, typeIds[i]) {
			break
		}
	}
}

/*
 * Run for mutex operation without a post event. Show an error in the results
 * MARK: Mutex
 * Args:
 *   id (int): The mutex id
 *   tID (string): The trace id
 */
func CheckForLeakMutex(id int, tID string) {
	found := "Leak on mutex::\n"
	found += "\tmutex: " + tID + "\n"
	found += "\tlast: " + mostRecentAcquireTotal[id].tID + "\n"
	logging.Result(found, logging.CRITICAL)
}

/*
 * Add the most recent acquire operation for a mutex
 * Args:
 *   id (int): The mutex id
 *   tID (string): The trace id
 *   vc (VectorClock): The vector clock of the operation
 */
func addMostRecentAcquireTotal(id int, tID string, vc clock.VectorClock) {
	mostRecentAcquireTotal[id] = VectorClockTID{vc, tID}
}

/*
 * Run for wait group operation without a post event. Show an error in the results
 * MARK: WaitGroup
 * Args:
 *   tID (string): The trace id
 */
func CheckForLeakWait(tID string) {
	found := "Leak on wait group:\n"
	found += "\twait-group: " + tID + "\n"
	found += "\t"
	logging.Result(found, logging.CRITICAL)
}

/*
 * Run for conditional varable operation without a post event. Show an error in the results
 * MARK: Cond
 * Args:
 *   tID (string): The trace id
 */
func CheckForLeakCond(tID string) {
	found := "Leak on conditional variable:\n"
	found += "\tconditional: " + tID + "\n"
	found += "\t"
	logging.Result(found, logging.CRITICAL)
}
