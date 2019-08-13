package hit

import (
	"fmt"
)

type PanicT struct{}

func (PanicT) Errorf(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}
func (PanicT) FailNow() {
	panic("FailNow")
}

// getLastArgument returns the last argument
func getLastArgument(params []interface{}) (interface{}, bool) {
	if i := len(params); i > 0 {
		return params[i-1], true
	}
	return nil, false
}
