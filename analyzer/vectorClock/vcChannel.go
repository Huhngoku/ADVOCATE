package vectorClock

/*
 * Update and calculate the vector clocks given a send/receive pair on a unbuffered
 * channel.
 */
func Unbuffered(routSend int, routRecv int, vc map[int]VectorClock) VectorClock {
	vc[routRecv].Sync(vc[routSend])
	vc[routRecv].Inc(routRecv)
	vc[routSend] = vc[routRecv].Copy()
	return vc[routRecv].Copy()
}
