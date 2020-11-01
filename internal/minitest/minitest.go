// Package minitest provides some testing functions for the hit package.
package minitest

import (
	"fmt"

	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/k0kubun/pp"
	"github.com/lunixbochs/vtclean"
)

//nolint:unparam
func stringJoin(separator string, a ...string) string {
	return strings.Join(a, separator)
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

// PrintValue prints the specified value in a nice way.
func PrintValue(v interface{}) string {
	return vtclean.Clean(pp.Sprint(v), false)
}
