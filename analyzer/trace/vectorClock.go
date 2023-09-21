package trace

/*
 * vectorClock is a vector clock
 * Fields:
 *   size (int): The size of the vector clock
 *   clock ([]int): The vector clock
 */
type vectorClock struct {
	size  int
	clock []int
}

type happensBefore int

const (
	Before happensBefore = iota
	After
	Concurrent
	None
)

/*
 * Create a new vector clock
 * Args:
 *   size (int): The size of the vector clock
 * Returns:
 *   (vectorClock): The new vector clock
 */
func newVectorClock(size int) vectorClock {
	return vectorClock{
		size:  size,
		clock: make([]int, size),
	}
}

/*
 * Increment the vector clock at the given position
 * Args:
 *   pos (int): The position to increment
 */
func (vc *vectorClock) inc(pos int) {
	vc.clock[pos]++
}

/*
 * Update the vector clock given a received vector clock by taking the
 * element wise maximum of the two vector clocks
 * Args:
 *   rec (vectorClock): The received vector clock
 */
func (vc *vectorClock) sync(rec *vectorClock) {
	for i := 0; i < vc.size; i++ {
		if vc.clock[i] < rec.clock[i] {
			vc.clock[i] = rec.clock[i]
		}
	}
}

/*
 * Get the happens before relation between two operations given there
 * vector clocks
 * Args:
 *   vc1 (vectorClock): The first vector clock
 *   vc2 (vectorClock): The second vector clock
 * Returns:
 *   happensBefore: The happens before relation between the two vector clocks
 */
func getHappensBefore(pre1 *vectorClock, post1 *vectorClock,
	pre2 *vectorClock, post2 *vectorClock) happensBefore {
	isCausePre1 := isCause(pre1, pre2)
	isCausePre2 := isCause(pre2, pre1)

	isCausePre := None
	if isCausePre1 {
		isCausePre = Before
	} else if isCausePre2 {
		isCausePre = After
	} else {
		return Concurrent
	}

	isCausePost1 := isCause(post1, post2)
	isCausePost2 := isCause(post2, post1)

	isCausePost := None
	if isCausePost1 {
		isCausePost = Before
	} else if isCausePost2 {
		isCausePost = After
	} else {
		return Concurrent
	}

	if isCausePre == isCausePost {
		return isCausePre
	}
	return Concurrent
}

/*
 * Check if vc1 is a cause of vc2
 * Args:
 *   vc1 (vectorClock): The first vector clock
 *   vc2 (vectorClock): The second vector clock
 * Returns:
 *   bool: True if vc1 is a cause of vc2, false otherwise
 */
func isCause(vc1 *vectorClock, vc2 *vectorClock) bool {
	atLeastOneSmaller := false
	for i := 0; i < vc1.size; i++ {
		if vc1.clock[i] > vc2.clock[i] {
			return false
		} else if vc1.clock[i] < vc2.clock[i] {
			atLeastOneSmaller = true
		}
	}
	return atLeastOneSmaller
}
