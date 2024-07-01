package analysis

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/*
 * Get the info from a TID
 * Args:
 *   tid (string): The TID
 * Return:
 *   string: the file
 *   int: the line
 *   int: the tPre
 *   error: the error
 */

func infoFromTID(tid string) (string, int, int, error) {
	spilt1 := splitAtLast(tid, "@")

	if len(spilt1) != 2 {
		return "", 0, 0, errors.New(fmt.Sprint("TID not correct: no @: ", tid))
	}

	split2 := strings.Split(spilt1[0], ":")
	if len(split2) != 2 {
		return "", 0, 0, errors.New(fmt.Sprint("TID not correct: no ':': ", tid))
	}

	tPre, err := strconv.Atoi(spilt1[1])
	if err != nil {
		return "", 0, 0, err
	}

	line, err := strconv.Atoi(split2[1])
	if err != nil {
		return "", 0, 0, err
	}

	return split2[0], line, tPre, nil
}

func splitAtLast(str string, sep string) []string {
	i := strings.LastIndex(str, sep)
	if i == -1 {
		return []string{str}
	}
	return []string{str[:i], str[i+1:]}
}
