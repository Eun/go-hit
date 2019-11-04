package convert

import (
	"fmt"
	"os"
	"strings"
)

var debug = func(format string, a ...interface{}) {}
var debug2 = func(format string, a ...interface{}) {}

func init() {
	if strings.Contains(os.Getenv("GODEBUG"), "go.convert=1") {
		debug = func(format string, a ...interface{}) {
			fmt.Fprintf(os.Stdout, format, a...)
		}
	}
	if strings.Contains(os.Getenv("GODEBUG"), "go.convert=2") {
		debug = func(format string, a ...interface{}) {
			fmt.Fprintf(os.Stdout, format, a...)
		}
		debug2 = func(format string, a ...interface{}) {
			fmt.Fprintf(os.Stdout, format, a...)
		}
	}
}
