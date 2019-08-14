package hit

import (
	"fmt"
	"strings"
)

type PanicT struct{}

func (PanicT) Error(args ...interface{}) {
	panic(args)
}

// getLastArgument returns the last argument
func getLastArgument(params []interface{}) (interface{}, bool) {
	if i := len(params); i > 0 {
		return params[i-1], true
	}
	return nil, false
}

func makeURL(base, url string, a ...interface{}) string {
	url = fmt.Sprintf(url, a...)
	if base == "" {
		return url
	}
	if url == "" {
		return base
	}
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(url, "/")
}
