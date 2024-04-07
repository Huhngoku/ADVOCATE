package analysis

import (
	"analyzer/logging"
)

/*
* CheckForSelectCaseWithoutPartner checks for select cases without a valid
* partner. Call when all elements have been processed.
 */
func CheckForSelectCaseWithoutPartner() {
	// check if not selected cases could be partners
	for i, c1 := range selectCases {
		for j := i + 1; j < len(selectCases); j++ {
			c2 := selectCases[j]

			if c1.partner && c2.partner {
				continue
			}

			if c1.id != c2.id || c1.vcTID.tID == c2.vcTID.tID || c1.send == c2.send {
				continue
			}

			if c2.send { // c1 should be send, c2 should be recv
				c1, c2 = c2, c1
			}

			hb := GetHappensBefore(c1.vcTID.vc, c2.vcTID.vc)
			found := false
			if c1.buffered && (hb == Concurrent || hb == After) {
				found = true
			} else if !c1.buffered && hb == Concurrent {
				found = true
			}

			if found {
				selectCases[i].partner = true
				selectCases[j].partner = true
			}
		}
	}

	// return all cases without a partner
	for _, c := range selectCases {
		if c.partner {
			continue
		}

		logging.Result("Possible select case without partner:\n\tselect: "+
			c.vcTID.tID+"\n\t\n", logging.WARNING)
	}
}

/*
* CheckForSelectCaseWithoutPartnerSelect checks for select cases without a valid
* partner. Call whenever a select is processed.
* Args:
*   ids ([]int): The ids of the channels
*   bufferedInfo ([]bool): The buffer status of the channels
*   sendInfo ([]bool): The send status of the channels
*   vc (VectorClock): The vector clock
*   tID (string): The position of the select in the program
 */
func CheckForSelectCaseWithoutPartnerSelect(ids []int, bufferedInfo []bool,
	sendInfo []bool, vc VectorClock, tID string, chosenIndex int) {
	for i, id := range ids {
		buffered := bufferedInfo[i]
		send := sendInfo[i]

		found := false

		if i == chosenIndex {
			// no need to check if the channel is the chosen case
			found = true
		} else {
			// not select cases
			if send {
				if potentialPartner, ok := mostRecentReceive[id]; ok {
					hb := GetHappensBefore(vc, potentialPartner.vc)
					if buffered && (hb == Concurrent || hb == Before) {
						found = true
					} else if !buffered && hb == Concurrent {
						found = true
					}
				}
			} else { // recv
				if potentialPartner, ok := mostRecentSend[id]; ok {
					hb := GetHappensBefore(vc, potentialPartner.vc)
					if buffered && (hb == Concurrent || hb == After) {
						found = true
					} else if !buffered && hb == Concurrent {
						found = true
					}
				}
			}
		}

		selectCases = append(selectCases,
			allSelectCase{id, VectorClockTID{vc, tID}, send, buffered, found})

	}
}

/*
* CheckForSelectCaseWithoutPartnerChannel checks for select cases without a valid
* partner. Call whenever a channel operation is processed.
* Args:
*   id (int): The id of the channel
*   vc (VectorClock): The vector clock
*   tID (string): The position of the channel operation in the program
*   send (bool): True if the operation is a send
*   buffered (bool): True if the channel is buffered
*   sel (bool): True if the operation is part of a select statement
 */
func CheckForSelectCaseWithoutPartnerChannel(id int, vc VectorClock, tID string,
	send bool, buffered bool) {

	for i, c := range selectCases {
		if c.partner || c.id != id || c.send == send || c.vcTID.tID == tID {
			continue
		}

		hb := GetHappensBefore(vc, c.vcTID.vc)
		found := false
		if send {
			if buffered && (hb == Concurrent || hb == Before) {
				found = true
			} else if !buffered && hb == Concurrent {
				found = true
			}
		} else {
			if buffered && (hb == Concurrent || hb == After) {
				found = true
			} else if !buffered && hb == Concurrent {
				found = true
			}
		}

		if found {
			selectCases[i].partner = true
		}
	}
}

/*
* CheckForSelectCaseWithoutPartnerClose checks for select cases without a valid
* partner. Call whenever a close operation is processed.
* Args:
*   id (int): The id of the channel
*   vc (VectorClock): The vector clock
 */
func CheckForSelectCaseWithoutPartnerClose(id int, vc VectorClock) {
	for i, c := range selectCases {
		if c.partner || c.id != id || c.send {
			continue
		}

		hb := GetHappensBefore(vc, c.vcTID.vc)
		found := false
		if c.buffered && (hb == Concurrent || hb == After) {
			found = true
		} else if !c.buffered && hb == Concurrent {
			found = true
		}

		if found {
			selectCases[i].partner = true
		}
	}
}
