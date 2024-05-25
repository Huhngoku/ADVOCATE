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

	if id == -1 {
		found := "Leak on nil channel:\n"
		found += "\tChannel: " + tID + "\n"
		found += "\tPartner: -"
		logging.Result(found, logging.CRITICAL)
		return
	}

	// if !buffered {
	foundPartner := false

	if opType == 0 { // send
		for _, mrr := range mostRecentReceive {
			if _, ok := mrr[id]; ok {
				if clock.GetHappensBefore(mrr[id].Vc, vc) == clock.Concurrent {
					found := ""
					if buffered {
						found = "Leak on buffered channel with possible partner:\n"
					} else {
						found = "Leak on unbuffered channel with possible partner:\n"
					}
					found += "\tchannel: " + tID + "\n"
					found += "\tpartner: " + mrr[id].TID
					logging.Result(found, logging.CRITICAL)
					foundPartner = true
				}
			}
		}
	} else if opType == 1 { // recv
		for _, mrs := range mostRecentSend {
			if _, ok := mrs[id]; ok {
				if clock.GetHappensBefore(mrs[id].Vc, vc) == clock.Concurrent {
					found := ""
					if buffered {
						found = "Leak on buffered channel with possible partner:\n"
					} else {
						found = "Leak on unbuffered channel with possible partner:\n"
					}
					found += "\tchannel: " + tID + "\n"
					found += "\tpartner: " + mrs[id].TID
					logging.Result(found, logging.CRITICAL)
					foundPartner = true
				}
			}
		}

		// // This cannot happen:
		// if _, ok := closeData[id]; ok {
		// 	found := "Leak on unbuffered channel or select with possible partner:\n"
		// 	found += "\tchannel: " + tID + "\n"
		// 	found += "\tpartner: " + closeData[id].tID
		// 	logging.Result(found, logging.CRITICAL)
		// 	foundPartner = true
		// }
	}

	if !foundPartner {
		leakingChannels[id] = append(leakingChannels[id], VectorClockTID2{id, vc, tID, opType, -1, buffered, false})
	}
	// } else {
	// 	// BUG: find possible partners, if there are any
	// 	// find possible partners
	// 	found := "Leak on buffered channel without possible partner:\n"
	// 	found += "\tchannel: " + tID + "\n"
	// 	found += "\tpartner: -"
	// 	logging.Result(found, logging.CRITICAL)
	// }
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
 *   buffered (bool): If the channel is buffered
 */
func CheckForLeakChannelRun(id int, vcTID VectorClockTID, opType int, buffered bool) bool {
	logging.Debug("Checking channel for for leak channels", logging.INFO)
	res := false
	if opType == 0 || opType == 2 { // send or close
		for i, vcTID2 := range leakingChannels[id] {
			if vcTID2.val != 1 {
				continue
			}
			if clock.GetHappensBefore(vcTID2.vc, vcTID.Vc) == clock.Concurrent {
				found := ""
				if buffered {
					found = "Leak on buffered channel with possible partner:\n"
				} else {
					found = "Leak on unbuffered channel with possible partner:\n"
				}
				found += "\tchannel: " + vcTID2.tID + "\n"
				found += "\tpartner: " + vcTID.TID
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
			if clock.GetHappensBefore(vcTID2.vc, vcTID.Vc) == clock.Concurrent {
				found := ""
				if buffered {
					found = "Leak on buffered channel with possible partner:\n"
				} else {
					found = "Leak on unbuffered channel with possible partner:\n"
				}
				found += "\tchannel: " + vcTID2.tID + "\n"
				found += "\tpartner: " + vcTID.TID
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
		buffered := false
		for _, vcTID := range vcTIDs {
			if vcTID.tID == "" {
				continue
			}

			found := false
			var partner VectorClockTID
			for _, c := range selectCases {
				if c.id != vcTID.id {
					continue
				}

				if (c.send && vcTID.typeVal == 0) || (!c.send && vcTID.typeVal == 1) {
					continue
				}

				hb := clock.GetHappensBefore(c.vcTID.Vc, vcTID.vc)
				if hb == clock.Concurrent {
					found = true
					if c.buffered {
						buffered = true
					}
					partner = c.vcTID
					break
				}

				if c.buffered {
					if (c.send && hb == clock.Before) || (!c.send && hb == clock.After) {
						found = true
						buffered = true
						partner = c.vcTID
						break
					}
				}
			}

			foundStr := ""
			if found {
				if vcTID.sel {
					if buffered {
						foundStr = "Leak on select with possible buffered partner:\n"
					} else {
						foundStr = "Leak on select with possible unbuffered partner:\n"
					}
					foundStr += "\tselect: " + vcTID.tID + "\n"
				} else {
					if buffered { // BUG: get unbuffered but should be bufferd
						foundStr = "Leak on buffered channel with possible partner:\n"
					} else {
						foundStr = "Leak on unbuffered channel with possible partner:\n"
					}
					foundStr += "\tchannel: " + vcTID.tID + "\n"
				}
				foundStr += "\tpartner: " + partner.TID
				logging.Result(foundStr, logging.CRITICAL)
			} else {
				if vcTID.sel {
					foundStr = "Leak on select without possible partner:\n"
					foundStr += "\tselect: " + vcTID.tID + "\n"
				} else {
					if buffered {
						foundStr = "Leak on buffered channel without possible partner:\n"
					} else {
						foundStr = "Leak on unbuffered channel without possible partner:\n"
					}
					foundStr += "\tchannel: " + vcTID.tID + "\n"
				}
				foundStr += "\tpartner: -"
				logging.Result(foundStr, logging.CRITICAL)
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
 *   buffered ([]bool): If the channels are buffered
 *   vc (VectorClock): The vector clock of the operation
 *   tID (string): The trace id
 *   opTypes ([]int): An identifier for the type of the operations (send = 0, recv = 1)
 *   idSel (int): The id of the select operation
 *   tPre (int): The tpre of the select operations. Used to connect the operations of the
 *     same select statement in leakingChannels.
 */
func CheckForLeakSelectStuck(ids []int, buffered []bool, vc clock.VectorClock, tID string, opTypes []int, tPre int) {
	foundPartner := false

	if len(ids) == 0 {
		found := "Leak on select with only nil channels:\n"
		found += "\tselect: " + tID + "\n"
		found += "\tpartner: -"
		logging.Result(found, logging.CRITICAL)
		return
	}

	for i, id := range ids {
		if opTypes[i] == 0 { // send
			for _, mrr := range mostRecentReceive {
				if _, ok := mrr[id]; ok {
					if clock.GetHappensBefore(vc, mrr[id].Vc) == clock.Concurrent {
						found := ""
						if buffered[i] {
							found = "Leak on select with possible buffered partner:\n"
						} else {
							found = "Leak on select with possible unbuffered partner:\n"
						}
						found += "\tselect: " + tID + "\n"
						found += "\tpartner: " + mrr[id].TID + "\n"
						logging.Result(found, logging.CRITICAL)
						foundPartner = true
					}
				}
			}
		} else if opTypes[i] == 1 { // recv
			for _, mrs := range mostRecentSend {
				if _, ok := mrs[id]; ok {
					if clock.GetHappensBefore(vc, mrs[id].Vc) == clock.Concurrent {
						found := ""
						if buffered[i] {
							found = "Leak on select with possible buffered partner:\n"
						} else {
							found = "Leak on select with possible unbuffered partner:\n"
						}
						found += "\tselect: " + tID + "\n"
						found += "\tpartner: " + mrs[id].TID
						logging.Result(found, logging.CRITICAL)
						foundPartner = true
					}
				}
			}
			if _, ok := closeData[id]; ok {
				found := ""
				if buffered[i] {
					found = "Leak on select with possible buffered partner:\n"
				} else {
					found = "Leak on select with possible unbuffered partner:\n"
				}
				found += "\tselect: " + tID + "\n"
				found += "\tpartner: " + closeData[id].TID
				logging.Result(found, logging.CRITICAL)
				foundPartner = true
			}
		}
	}

	if !foundPartner {
		println("No partner found")
		for i, id := range ids {
			// add all select operations to leaking Channels,
			leakingChannels[id] = append(leakingChannels[id], VectorClockTID2{id, vc, tID, opTypes[i], tPre, buffered[i], true})
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
	found += "\tlast: " + mostRecentAcquireTotal[id].TID + "\n"
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
