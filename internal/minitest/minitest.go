package minitest

import (
	"fmt"

	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/k0kubun/pp"
	"github.com/lunixbochs/vtclean"
)

func stringJoin(separator string, a ...string) string {
	return strings.Join(a, separator)
}

func formatMessage(customMessageAndArgs []interface{}) string {
	if len(customMessageAndArgs) == 0 {
		return ""
	}
	s, ok := customMessageAndArgs[0].(string)
	if !ok {
		return "expected custom message to be a string"
	}
	return strings.TrimSpace(fmt.Sprintf(s, customMessageAndArgs[1:]...))
}

func actualExpectedDiff(actual, expected interface{}) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "expected:\t%s\n", PrintValue(expected))
	fmt.Fprintf(&sb, "actual:  \t%s\n", PrintValue(actual))
	if diff := cmp.Diff(expected, actual); diff != "" {
		fmt.Fprint(&sb, Format("diff:    ", trimLeftSpaces(diff)))
	}
	return sb.String()
}
func PrintValue(v interface{}) string {
	return vtclean.Clean(pp.Sprint(v), false)
}

var Panic PanicNow
var Error ReturnError
