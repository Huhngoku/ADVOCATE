// DEDEGO-FILE-START

package runtime

/*
 * Get a string representation of an uint32
 * Args:
 * 	n: int to convert
 * Return:
 * 	string representation of the int
 */
func uint64ToString(n uint64) string {
	if n < 10 {
		return string(n + '0')
	} else {
		return uint64ToString(n/10) + string(n%10+'0')

	}
}

// DEDEGO-FILE-END
