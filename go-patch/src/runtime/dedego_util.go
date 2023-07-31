// DEDEGO-FILE-START

package runtime

/*
 * Get a string representation of an uint64
 * Args:
 * 	n: int to convert
 * Return:
 * 	string representation of the int
 */
func uint64ToString(n uint64) string {
	if n < 10 {
		return string(rune(n + '0'))
	} else {
		return uint64ToString(n/10) + string(rune(n%10+'0'))
	}
}

/*
 * Get a string representation of an int64
 * Args:
 * 	n: int to convert
 * Return:
 * 	string representation of the int
 */
func int64ToString(n int64) string {
	if n < 0 {
		return "-" + int64ToString(-n)
	}
	if n < 10 {
		return string(rune(n + '0'))
	} else {
		return int64ToString(n/10) + string(rune(n%10+'0'))

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
		return "-" + int32ToString(-n)
	}
	if n < 10 {
		return string(rune(n + '0'))
	} else {
		return int32ToString(n/10) + string(rune(n%10+'0'))

	}
}

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
 * Get a string representation of an int
 * Args:
 * 	n: int to convert
 * Return:
 * 	string representation of the int
 */
func intToString(n int) string {
	if n < 0 {
		return "-" + intToString(-n)
	}
	if n < 10 {
		return string(rune(n + '0'))
	} else {
		return intToString(n/10) + string(rune(n%10+'0'))

	}
}

var dedegoCurrentRoutineId uint64
var dedegoCurrentObjectId uint64
var dedegoGlobalCounter uint64

var dedegoCurrentRoutineIdMutex mutex
var dedegoCurrentObjectIdMutex mutex
var dedegoGlobalCounterMutex mutex

/*
 * Get a new id for a routine
 * Return:
 * 	new id
 */
func GetDedegoRoutineId() uint64 {
	lock(&dedegoCurrentRoutineIdMutex)
	defer unlock(&dedegoCurrentRoutineIdMutex)
	dedegoCurrentRoutineId += 1
	return dedegoCurrentRoutineId
}

/*
 * Get a new id for a mutex, channel or waitgroup
 * Return:
 * 	new id
 */
func GetDedegoObjectId() uint64 {
	lock(&dedegoCurrentObjectIdMutex)
	defer unlock(&dedegoCurrentObjectIdMutex)
	dedegoCurrentObjectId += 1
	return dedegoCurrentObjectId
}

/*
 * Get a new counter
 * Return:
 * 	new counter value
 */
func GetDedegoCounter() uint64 {
	lock(&dedegoGlobalCounterMutex)
	defer unlock(&dedegoGlobalCounterMutex)
	dedegoGlobalCounter += 1
	return dedegoGlobalCounter
}

// DEDEGO-FILE-END
