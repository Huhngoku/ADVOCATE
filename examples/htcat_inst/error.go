package main

import (
	"fmt"
)

type errAssertf struct {
	error
}

func assertErrf(format string, a ...interface{}) errAssertf {
	return errAssertf{fmt.Errorf(format, a...)}
}
