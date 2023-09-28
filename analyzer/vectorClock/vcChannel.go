package vectorClock

import "analyzer/debug"

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

/*
 * Update and calculate the vector clocks given a send on a buffered channel.
 * Args:
 * 	rout (int): the route of the sender
 * 	id (int): the id of the sender
 * 	size (int): buffer size
 * 	vc (map[int]VectorClock): the current vector clocks
 * Returns:
 * 	the vector clock of the send
 */
func Send(rout int, id int, size int, vc map[int]VectorClock) VectorClock {
	newBufferedVCs(id, size, vc[rout].size)

	count := bufferedVCsCount[id]
	if count > size || bufferedVCs[id][count].occupied {
		debug.Log("Write to occupied buffer position or to big count", debug.ERROR)
	}

	v := bufferedVCs[id][count].vc
	vc[rout].Sync(v)
	bufferedVCs[id][count] = bufferedVC{true, vc[rout].Copy()}
	return vc[rout].Inc(rout).Copy()
}

/*
 * Update and calculate the vector clocks given a receive on a buffered channel.
 * Args:
 * 	rout (int): the route of the sender
 * 	id (int): the id of the sender
 * 	size (int): buffer size
 * 	vc (map[int]VectorClock): the current vector clocks
 * Returns:
 * 	the vector clock of the receive
 */
func Recv(rout int, id int, size int, vc map[int]VectorClock) VectorClock {
	newBufferedVCs(id, size, vc[rout].size)

	v := bufferedVCs[id][0].vc
	if !bufferedVCs[id][0].occupied {
		debug.Log("Read from unoccupied buffer position", debug.ERROR)
	}
	vc[rout].Sync(v)
	bufferedVCs[id] = bufferedVCs[id][1:]
	bufferedVCs[id] = append(bufferedVCs[id], bufferedVC{false, vc[rout].Copy()})

	return vc[rout].Inc(rout).Copy()
}

/*
 * Update and calculate the vector clocks given a close on a channel.
 * Args:
 * 	rout (int): the route of the sender
 * 	id (int): the id of the sender
 * 	vc (map[int]VectorClock): the current vector clocks
 * Returns:
 * 	the vector clock of the close
 */
func Close(rout int, id int, vc map[int]VectorClock) VectorClock {
	closeVC[id] = vc[rout].Copy()
	return vc[rout].Inc(rout).Copy()
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
	vc[rout].Sync(closeVC[id])
	return vc[rout].Inc(rout).Copy()
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
