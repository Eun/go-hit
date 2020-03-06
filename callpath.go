package hit

import (
	"fmt"

	"github.com/Eun/go-hit/internal"
)

type CleanPath []string

func callPathFormat(name string, arguments ...interface{}) string {
	arg, _ := internal.GetLastArgument(arguments)
	return fmt.Sprintf("%s(%v)", name, arg)
}

func (p CleanPath) Push(name string, arguments ...interface{}) CleanPath {
	return append(p, callPathFormat(name, arguments))
}

func NewCleanPath(name string, arguments ...interface{}) CleanPath {
	return []string{callPathFormat(name, arguments)}
}
