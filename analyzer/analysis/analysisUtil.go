package analysis

import (
	"errors"
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
	spilt1 := strings.Split(tid, "@")

	if len(spilt1) != 2 {
		return "", 0, 0, errors.New("TID not correct: no @")
	}

	split2 := strings.Split(spilt1[0], ":")
	if len(split2) != 2 {
		return "", 0, 0, errors.New("TID not correct: no :")
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
