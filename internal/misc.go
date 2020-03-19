package internal

import (
	"fmt"
	"strings"
)

// GetLastArgument returns the last argument
func GetLastArgument(params []interface{}) (interface{}, bool) {
	if i := len(params); i > 0 {
		return params[i-1], true
	}
	return nil, false
}

// GetLastStringArgument returns the last argument
func GetLastStringArgument(params []string) (string, bool) {
	if i := len(params); i > 0 {
		return params[i-1], true
	}
	return "", false
}

// GetLastIntArgument returns the last argument
func GetLastIntArgument(params []int) (int, bool) {
	if i := len(params); i > 0 {
		return params[i-1], true
	}
	return 0, false
}

func MakeURL(base, url string, a ...interface{}) string {
	url = fmt.Sprintf(url, a...)
	if base == "" {
		return url
	}
	if url == "" {
		return base
	}
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(url, "/")
}
