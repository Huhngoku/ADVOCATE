package analysis

import "analyzer/logging"

/*
* CheckForSelectCaseWithoutPartner checks for select cases without a valid
* partner. Call when all elements have been processed.
 */
func CheckForSelectCaseWithoutPartner() {
	for _, cases := range selectCasesSend {
		for _, vcTID := range cases {
			foundPossibleSelectCaseWithoutPartner(vcTID)
		}
	}

	for _, cases := range selectCasesRecv {
		for _, vcTID := range cases {
			foundPossibleSelectCaseWithoutPartner(vcTID)
		}
	}

}

/*
* CheckForSelectCaseWithoutPartner checks for select cases without a valid
* partner. Call when a select statement is encountered.
* Args:
*   cID (int): The channel id
*   tID (string): The trace element id
*   send (bool): Whether the case is a send case
*   buffered (bool): Whether the channel is buffered
*   vc (VectorClock): The vector clock
 */
// TODO: For now only works for unbuffered channels, but will create error for buffered channels
func CheckForSelectCaseWithoutPartnerSelect(cID int, tID string, send bool, buffered bool, vc VectorClock) {
	if send {
		possibleRecv := mostRecentReceive[cID]
		if GetHappensBefore(vc, possibleRecv.vc) == Concurrent {
			return
		}
	} else {
		possibleSend := mostRecentSend[cID]
		if GetHappensBefore(vc, possibleSend.vc) == Concurrent {
			return
		}
	}

	AddUntriggeredSelectCase(cID, tID, send, buffered, vc)
}

/*
* checkForSelectCaseWithoutPartnerChannel checks for select cases without a valid
* partner. Call when a channel operation is encountered.
* Args:
*   cID (int): The channel id
*   tID (string): The trace element id
*   send (bool): Whether the case is a send case
*   buffered (bool): Whether the channel is buffered
*   vc (VectorClock): The vector clock
 */
func checkForSelectCaseWithoutPartnerChannel(cID int, tID string, send bool, buffered bool, vc VectorClock) {
	if send {
		possibleCases := selectCasesRecv[cID]
		for i := 0; i < len(possibleCases); i++ {
			possibleCase := possibleCases[i]
			if GetHappensBefore(vc, possibleCase.vc) == Concurrent {
				selectCasesRecv[cID] = append(possibleCases[:i], possibleCases[i+1:]...)
				i--
			}
		}
	} else {
		possibleCases := selectCasesSend[cID]
		for i := 0; i < len(possibleCases); i++ {
			possibleCase := possibleCases[i]
			if GetHappensBefore(vc, possibleCase.vc) == Concurrent {
				selectCasesSend[cID] = append(possibleCases[:i], possibleCases[i+1:]...)
				i--
			}
		}
	}
}

/*
* AddUntriggeredSelectCase adds an untriggered select case to the analysis
* Args:
*   cId (int): The channel id
*   tID (string): The trace element id
*   send (bool): Whether the case is a send case
*   buffered (bool): Whether the channel is buffered
*   vc (VectorClock): The vector clock
 */
func AddUntriggeredSelectCase(cID int, tID string, send bool, buffered bool,
	vc VectorClock) {
	vcTID := VectorClockTID3{vc, tID, buffered}
	if send {
		if _, ok := selectCasesSend[cID]; !ok {
			selectCasesSend[cID] = make([]VectorClockTID3, 0)
		}

		selectCasesSend[cID] = append(selectCasesSend[cID], vcTID)
	} else {
		if _, ok := selectCasesRecv[cID]; !ok {
			selectCasesRecv[cID] = make([]VectorClockTID3, 0)
		}

		selectCasesRecv[cID] = append(selectCasesRecv[cID], vcTID)
	}

}

/*
 * Log a found possible select case without partner
 * Args:
 *   vcTID (VectorClockTID): The vector clock and trace element id
 */
func foundPossibleSelectCaseWithoutPartner(cvTID VectorClockTID3) {
	msg := "Possible select case without partner:\n"
	msg += "select: " + cvTID.tID + "\n"
	msg += "case: \n" // TODO: add identifier for the case
	logging.Result(msg, logging.WARNING)
}
