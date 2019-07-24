package convert

import (
	"fmt"
	"os"
	"strings"
)

var debug = func(format string, a ...interface{}) {}

func init() {
	if strings.Contains(os.Getenv("GODEBUG"), "go.convert=1") {
		debug = func(format string, a ...interface{}) {
			fmt.Printf(format, a...)
		}
	}
}
