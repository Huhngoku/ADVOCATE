package utils

/*
* Check if a slice Contains an element
* Args:
*   s: slice to check
*   e: element to check
 */
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
