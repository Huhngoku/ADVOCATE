// DEDEGO-FILE-START

package runtime

import "sync/atomic"

/*
 * Get a string representation of an uint32
 * Args:
 * 	n: int to convert
 * Return:
 * 	string representation of the int
 */
func uint32ToString(n uint32) string {
	if n < 10 {
		return string(rune(n + '0'))
	} else {
		return uint32ToString(n/10) + string(rune(n%10+'0'))

	}
}

/*
 * Get a string representation of an int32
 * Args:
 * 	n: int to convert
 * Return:
 * 	string representation of the int
 */
func int32ToString(n int32) string {
	if n < 0 {
		return "-" + uint32ToString(uint32(-n))
	}
	if n < 10 {
		return string(rune(n + '0'))
	} else {
		return int32ToString(n/10) + string(rune(n%10+'0'))

	}
}

var dedegoCurrentId = uint32(0)

/*
 * Get a new id for a mutex, channel or waitgroup
 * Return:
 * 	new id
 */
func GetNewId() uint32 {
	return atomic.AddUint32(&dedegoCurrentId, 1)
}

// DEDEGO-FILE-END
