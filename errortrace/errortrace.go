package errortrace

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/Eun/go-hit/internal"
	"github.com/google/go-cmp/cmp"
	"github.com/gookit/color"
	"github.com/k0kubun/pp"
	"github.com/lunixbochs/vtclean"
)

type ErrorTrace struct {
	inheritedTrace bool
	inheritedPC    []uintptr
}

type ErrorTraceError string

func (e ErrorTraceError) Error() string {
	return string(e)
}

var defaultErrorTrace ErrorTrace

// Prepare collects the current trace, if later a Panic will be called the collected trace will be included
// in the error trace
func Prepare() *ErrorTrace {
	var et ErrorTrace
	et.inheritedTrace = true
	et.inheritedPC = currentTraceCalls(4)
	return &et
}

func isIncluded(call *Call) bool {
	if call.FunctionName == "" {
		return false
	}
	if call.PackageName != "github.com/Eun/go-hit" {
		return false
	}

	switch call.FunctionName {
	case "Do", "Custom", "runExpectCalls", "runSendCalls":
		return false
	}

	return true
}

func filterTraceCalls(calls []Call) []Call {
	var filtered []Call

	for i := 0; i < len(calls); i++ {
		if isIncluded(&calls[i]) {
			filtered = append(filtered, calls[i])
		}
	}
	return filtered
}

func currentTraceCalls(skip int) []uintptr {
	var pc [16]uintptr
	var calls []uintptr
	for index, n := skip, 0; ; index += n {
		n = runtime.Callers(index, pc[:])
		if n <= 0 {
			break
		}
		calls = append(calls, pc[:n]...)
	}
	return calls
}

func resolveTraceCalls(pc []uintptr) []Call {
	var calls []Call
	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()
		calls = append(calls, makeCall(frame))
		if !more {
			return calls
		}
	}
}

func stringJoin(seperator string, a ...string) string {
	return strings.Join(a, seperator)
}

func (p *ErrorTrace) formatMessage(customMessageAndArgs []interface{}) string {
	if len(customMessageAndArgs) <= 0 {
		return ""
	}
	s, ok := customMessageAndArgs[0].(string)
	if !ok {
		p.panicNow(p.formatError("expected custom message to be a string"))
	}
	return strings.TrimSpace(fmt.Sprintf(s, customMessageAndArgs[1:]...))
}

func (p *ErrorTrace) formatStack(calls []Call) string {
	var sb strings.Builder

	if n := len(calls); n > 0 {
		cl := color.New(color.FgBlue, color.OpUnderscore)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%s(...)\n\t", calls[i].FullName())
			sb.WriteString(cl.Sprintf("%s:%d", calls[i].File, calls[i].Line))
			fmt.Fprintln(&sb)
		}
	}
	return sb.String()
}

func (p *ErrorTrace) formatError(errText string, customMessageAndArgs ...interface{}) string {
	// if helper, ok := t.(TestHelper); ok {
	// 	helper.ErrorMessage(errText, p.formatMessage(customMessageAndArgs))
	// 	return errors.New(errText)
	// }

	// collect the current trace
	traceCalls := resolveTraceCalls(currentTraceCalls(4))

	// if we have a inherited trace resolve the items and filter the current trace
	if p.inheritedTrace {
		// filter
		traceCalls = filterTraceCalls(traceCalls)

		// resolve the inherited calls
		traceCalls = append(traceCalls, resolveTraceCalls(p.inheritedPC)...)
	}

	// print
	var sb strings.Builder
	if detail := p.formatMessage(customMessageAndArgs); detail != "" {
		fmt.Fprintln(&sb, detail)
	}
	sb.WriteString(errText)
	errText = sb.String()
	sb.Reset()

	fmt.Fprint(&sb, format("Error:      ", errText, color.FgRed))
	fmt.Fprint(&sb, format("Error Trace:", p.formatStack(traceCalls)))
	return sb.String()
}

func (p *ErrorTrace) panicNow(errText string) {
	panic(ErrorTraceError(errText))
}

func (p *ErrorTrace) actualExpectedDiff(actual, expected interface{}) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "expected:\t%+v\n", expected)
	fmt.Fprintf(&sb, "actual:  \t%+v\n", actual)
	if diff := cmp.Diff(expected, actual); diff != "" {
		fmt.Fprintf(&sb, format("diff:    ", trimLeftSpaces(diff)))
	}
	return sb.String()
}

func (p *ErrorTrace) FailNow(err error, customMessageAndArgs ...interface{}) {
	if err != nil {
		p.panicNow(p.formatError(err.Error(), customMessageAndArgs...))
	}
	p.panicNow(p.formatError("", customMessageAndArgs...))
}

func (p *ErrorTrace) Errorf(messageAndArgs ...interface{}) {
	p.panicNow(p.formatError(p.formatMessage(messageAndArgs)))
}

func (p *ErrorTrace) NoError(err error, customMessageAndArgs ...interface{}) {
	if err != nil {
		p.panicNow(p.formatError(err.Error(), customMessageAndArgs...))
	}
}

func (p *ErrorTrace) Equal(expected, actual interface{}, customMessageAndArgs ...interface{}) {
	if !cmp.Equal(expected, actual) {
		p.panicNow(p.formatError(stringJoin("\n", "Not equal", p.actualExpectedDiff(actual, expected)), customMessageAndArgs...))
	}
}

func (p *ErrorTrace) Contains(object interface{}, contains interface{}, customMessageAndArgs ...interface{}) {
	if !internal.Contains(object, contains) {
		p.panicNow(p.formatError(fmt.Sprintf(`%s does not contain %s`, vtclean.Clean(pp.Sprint(object), false), vtclean.Clean(pp.Sprint(contains), false)), customMessageAndArgs...))
	}
}

func (p *ErrorTrace) Empty(object interface{}, customMessageAndArgs ...interface{}) {
	v := internal.GetValue(object)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		l := v.Len()
		if l != 0 {
			p.panicNow(p.formatError(fmt.Sprintf(`%s should be empty, but has %d item(s)`, vtclean.Clean(pp.Sprint(object), false), l), customMessageAndArgs...))
		}
	default:
		p.panicNow(p.formatError(fmt.Sprintf("called Len() on %s", vtclean.Clean(pp.Sprint(object), false))))
	}
}

func (p *ErrorTrace) Len(object interface{}, length int, customMessageAndArgs ...interface{}) {
	v := internal.GetValue(object)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		l := v.Len()
		if l != length {
			p.panicNow(p.formatError(fmt.Sprintf(`%s should have %d item(s), but has %d`, vtclean.Clean(pp.Sprint(object), false), length, l), customMessageAndArgs...))
		}
	default:
		p.panicNow(p.formatError(fmt.Sprintf("called Len() on %s", vtclean.Clean(pp.Sprint(object), false))))
	}
}

func (p *ErrorTrace) True(value bool, customMessageAndArgs ...interface{}) {
	if !value {
		p.panicNow(p.formatError(`Expected bool to be true but is false`, customMessageAndArgs...))
	}
}

func (p *ErrorTrace) False(value bool, customMessageAndArgs ...interface{}) {
	if value {
		p.panicNow(p.formatError(`Expected bool to be false but is true`, customMessageAndArgs...))
	}
}

func FailNow(err error, customMessageAndArgs ...interface{}) {
	defaultErrorTrace.FailNow(err, customMessageAndArgs...)
}

func Errorf(messageAndArgs ...interface{}) {
	defaultErrorTrace.Errorf(messageAndArgs...)
}

func NoError(err error, customMessageAndArgs ...interface{}) {
	defaultErrorTrace.NoError(err, customMessageAndArgs...)
}

func Equal(expected, actual interface{}, customMessageAndArgs ...interface{}) {
	defaultErrorTrace.Equal(expected, actual, customMessageAndArgs...)
}

func Contains(object interface{}, contains interface{}, customMessageAndArgs ...interface{}) {
	defaultErrorTrace.Contains(object, contains, customMessageAndArgs...)
}

func Empty(object interface{}, customMessageAndArgs ...interface{}) {
	defaultErrorTrace.Empty(object, customMessageAndArgs...)
}

func Len(object interface{}, length int, customMessageAndArgs ...interface{}) {
	defaultErrorTrace.Len(object, length, customMessageAndArgs...)
}

func True(value bool, customMessageAndArgs ...interface{}) {
	defaultErrorTrace.True(value, customMessageAndArgs...)
}

func False(value bool, customMessageAndArgs ...interface{}) {
	defaultErrorTrace.False(value, customMessageAndArgs...)
}

func FormatError(errText string, customMessageAndArgs ...interface{}) error {
	return ErrorTraceError(defaultErrorTrace.formatError(errText, customMessageAndArgs...))
}
