package errortrace

import (
	"fmt"
	"strings"

	"io"
	"os"

	"reflect"

	"github.com/Eun/go-hit/internal"
	"github.com/google/go-cmp/cmp"
	"github.com/gookit/color"
	"github.com/k0kubun/pp"
	"github.com/lunixbochs/vtclean"
)

var ErrorOut io.Writer = os.Stderr

type Panicer struct {
	inheritedTrace bool
	inheritedPC    []uintptr
}

type TestHelper interface {
	ErrorMessage(msg string, detail string)
}

func (p *Panicer) formatStack(calls []Call) string {
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

func (p *Panicer) panicNow(t TestingT, errText string, customMessageAndArgs ...interface{}) {
	if helper, ok := t.(TestHelper); ok {
		helper.ErrorMessage(errText, p.formatMessage(t, customMessageAndArgs))
		t.FailNow()
		return
	}

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
	if detail := p.formatMessage(t, customMessageAndArgs); detail != "" {
		fmt.Fprintln(&sb, detail)
	}
	sb.WriteString(errText)
	errText = sb.String()
	sb.Reset()

	fmt.Fprint(&sb, format("Error:      ", errText, color.FgRed))
	fmt.Fprint(&sb, format("Error Trace:", p.formatStack(traceCalls)))
	fmt.Fprintln(ErrorOut, sb.String())
	t.FailNow()
}

func (p *Panicer) actualExpectedDiff(actual, expected interface{}) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "expected:\t%+v\n", expected)
	fmt.Fprintf(&sb, "actual:  \t%+v\n", actual)
	if diff := cmp.Diff(expected, actual); diff != "" {
		fmt.Fprintf(&sb, format("diff:    ", trimLeftSpaces(diff)))
	}
	return sb.String()
}

func (p *Panicer) FailNow(t TestingT, err error, customMessageAndArgs ...interface{}) {
	if err != nil {
		p.panicNow(t, err.Error(), customMessageAndArgs...)
		return
	}
	p.panicNow(t, "", customMessageAndArgs...)
}

func (p *Panicer) Errorf(t TestingT, messageAndArgs ...interface{}) {
	p.panicNow(t, p.formatMessage(t, messageAndArgs))
}

func (p *Panicer) NoError(t TestingT, err error, customMessageAndArgs ...interface{}) {
	if err != nil {
		p.panicNow(t, err.Error(), customMessageAndArgs...)
	}
}

func (p *Panicer) Equal(t TestingT, expected, actual interface{}, customMessageAndArgs ...interface{}) {
	if !cmp.Equal(expected, actual) {
		p.panicNow(t, stringJoin("\n", "Not equal", p.actualExpectedDiff(actual, expected)), customMessageAndArgs...)
	}
}

func (p *Panicer) Contains(t TestingT, object interface{}, contains interface{}, customMessageAndArgs ...interface{}) {
	if !internal.Contains(object, contains) {
		p.panicNow(t, fmt.Sprintf(`%s does not contain %s`, vtclean.Clean(pp.Sprint(object), false), vtclean.Clean(pp.Sprint(contains), false)), customMessageAndArgs...)
	}
}

func (p *Panicer) Empty(t TestingT, object interface{}, customMessageAndArgs ...interface{}) {
	v := internal.GetValue(object)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		l := v.Len()
		if l != 0 {
			p.panicNow(t, fmt.Sprintf(`%s should be empty, but has %d item(s)`, vtclean.Clean(pp.Sprint(object), false), l), customMessageAndArgs...)
		}
	default:
		p.panicNow(t, fmt.Sprintf("called Len() on %s", vtclean.Clean(pp.Sprint(object), false)))
	}
}

func (p *Panicer) Len(t TestingT, object interface{}, length int, customMessageAndArgs ...interface{}) {
	v := internal.GetValue(object)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		l := v.Len()
		if l != length {
			p.panicNow(t, fmt.Sprintf(`%s should have %d item(s), but has %d`, vtclean.Clean(pp.Sprint(object), false), length, l), customMessageAndArgs...)
		}
	default:
		p.panicNow(t, fmt.Sprintf("called Len() on %s", vtclean.Clean(pp.Sprint(object), false)))
	}
}

func (p *Panicer) True(t TestingT, value bool, customMessageAndArgs ...interface{}) {
	if !value {
		p.panicNow(t, `Expected bool to be true but is false`, customMessageAndArgs...)
	}
}

func (p *Panicer) False(t TestingT, value bool, customMessageAndArgs ...interface{}) {
	if value {
		p.panicNow(t, `Expected bool to be false but is true`, customMessageAndArgs...)
	}
}

func (p *Panicer) formatMessage(t TestingT, customMessageAndArgs []interface{}) string {
	if len(customMessageAndArgs) <= 0 {
		return ""
	}
	s, ok := customMessageAndArgs[0].(string)
	if !ok {
		p.panicNow(t, "expected custom message to be a string")
	}
	return strings.TrimSpace(fmt.Sprintf(s, customMessageAndArgs[1:]...))
}

func stringJoin(seperator string, a ...string) string {
	return strings.Join(a, seperator)
}
