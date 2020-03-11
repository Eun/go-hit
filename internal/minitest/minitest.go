package minitest

import (
	"fmt"
	"reflect"

	"strings"

	"github.com/Eun/go-hit/internal"
	"github.com/google/go-cmp/cmp"
	"github.com/k0kubun/pp"
	"github.com/lunixbochs/vtclean"
)

func panicNow(err string, customMessageAndArgs ...interface{}) {
	var sb strings.Builder

	if detail := formatMessage(customMessageAndArgs); detail != "" {
		fmt.Fprintln(&sb, detail)
	}

	if err != "" {
		sb.WriteString(err)
	}

	panic(sb.String())
}

func stringJoin(seperator string, a ...string) string {
	return strings.Join(a, seperator)
}

func formatMessage(customMessageAndArgs []interface{}) string {
	if len(customMessageAndArgs) <= 0 {
		return ""
	}
	s, ok := customMessageAndArgs[0].(string)
	if !ok {
		panicNow("expected custom message to be a string")
	}
	return strings.TrimSpace(fmt.Sprintf(s, customMessageAndArgs[1:]...))
}

func actualExpectedDiff(actual, expected interface{}) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "expected:\t%s\n", PrintValue(expected))
	fmt.Fprintf(&sb, "actual:  \t%s\n", PrintValue(actual))
	if diff := cmp.Diff(expected, actual); diff != "" {
		fmt.Fprintf(&sb, Format("diff:    ", trimLeftSpaces(diff)))
	}
	return sb.String()
}

func FailNow(err error, customMessageAndArgs ...interface{}) {
	if err != nil {
		panicNow(err.Error(), customMessageAndArgs...)
	}
	panicNow("", customMessageAndArgs...)
}

func Errorf(messageAndArgs ...interface{}) {
	panicNow(formatMessage(messageAndArgs))
}

func NoError(err error, customMessageAndArgs ...interface{}) {
	if err != nil {
		panicNow(err.Error(), customMessageAndArgs...)
	}
}

func Equal(expected, actual interface{}, customMessageAndArgs ...interface{}) {
	if !cmp.Equal(expected, actual) {
		panicNow(stringJoin("\n", "Not equal", actualExpectedDiff(actual, expected)), customMessageAndArgs...)
	}
}

func NotEqual(expected, actual interface{}, customMessageAndArgs ...interface{}) {
	if cmp.Equal(expected, actual) {
		panicNow(stringJoin("\n", fmt.Sprintf("should not be %s", PrintValue(actual))), customMessageAndArgs...)
	}
}

func Contains(object interface{}, contains interface{}, customMessageAndArgs ...interface{}) {
	if !internal.Contains(object, contains) {
		panicNow(fmt.Sprintf(`%s does not contain %s`, PrintValue(object), PrintValue(contains)), customMessageAndArgs...)
	}
}

func NotContains(object interface{}, contains interface{}, customMessageAndArgs ...interface{}) {
	if internal.Contains(object, contains) {
		panicNow(fmt.Sprintf(`%s should not contain %s`, PrintValue(object), PrintValue(contains)), customMessageAndArgs...)
	}
}

func Empty(object interface{}, customMessageAndArgs ...interface{}) {
	v := internal.GetValue(object)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		l := v.Len()
		if l != 0 {
			panicNow(fmt.Sprintf(`%s should be empty, but has %d item(s)`, PrintValue(object), l), customMessageAndArgs...)
		}
	default:
		panicNow(fmt.Sprintf("called Len() on %s", PrintValue(object)))
	}
}

func Len(object interface{}, length int, customMessageAndArgs ...interface{}) {
	v := internal.GetValue(object)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		l := v.Len()
		if l != length {
			panicNow(fmt.Sprintf(`%s should have %d item(s), but has %d`, PrintValue(object), length, l), customMessageAndArgs...)
		}
	default:
		panicNow(fmt.Sprintf("called Len() on %s", PrintValue(object)))
	}
}

func True(value bool, customMessageAndArgs ...interface{}) {
	if !value {
		panicNow(`Expected bool to be true but is false`, customMessageAndArgs...)
	}
}

func False(value bool, customMessageAndArgs ...interface{}) {
	if value {
		panicNow(`Expected bool to be false but is true`, customMessageAndArgs...)
	}
}

func PrintValue(v interface{}) string {
	return vtclean.Clean(pp.Sprint(v), false)
}
