// Package errortrace provides a method to track function stacktrace and populate it in case of error.
package errortrace

import (
	"fmt"
	"runtime"
	"strings"

	"reflect"

	"github.com/gookit/color"

	"github.com/Eun/go-hit/internal/minitest"
)

// Template is used as a factory to create a new ErrorTrace when you need it, use New() to create an Template.
type Template struct {
	ignore []string
}

// ErrorTrace represents an ErrorTrace, it includes the parent trace (if called with Prepare).
type ErrorTrace struct {
	inheritedTrace bool
	inheritedPC    []uintptr
	template       *Template
	description    string
	error          error
	ctx            string
}

// New can be used to create an Template with the desired parameters.
func New(ignore ...string) *Template {
	return &Template{
		ignore: ignore,
	}
}

// IgnoreFunc can be used to ignore the specified function in a package for the trace.
func IgnoreFunc(fn interface{}) string {
	v := reflect.ValueOf(fn)
	if !v.IsValid() {
		panic("function is is not valid")
	}
	return strings.TrimSuffix(runtime.FuncForPC(v.Pointer()).Name(), "-fm")
}

// IgnorePackage can be used to ignore all functions in a package for the trace.
func IgnorePackage(fn interface{}) string {
	return makeCall(&runtime.Frame{
		PC:       0,
		Func:     nil,
		Function: IgnoreFunc(fn),
		File:     "",
		Line:     0,
		Entry:    0,
	}).PackageName
}

// IgnoreStruct can be used to ignore all struct functions in the trace.
func IgnoreStruct(fn interface{}) string {
	c := makeCall(&runtime.Frame{
		PC:       0,
		Func:     nil,
		Function: IgnoreFunc(fn),
		File:     "",
		Line:     0,
		Entry:    0,
	})
	return strings.Join([]string{c.PackageName, c.FunctionPath}, ".")
}

// Prepare collects the current trace, if later a Panic will be called the collected trace will be included
// in the error trace.
func (t *Template) Prepare() *ErrorTrace {
	var et ErrorTrace
	et.inheritedTrace = true
	// skip  runtime.Callers(), ett.currentTraceCalls(...), ett.(*ErrorTrace).Prepare(...), Call to this Function
	//nolint:gomnd
	et.inheritedPC = currentTraceCalls(4)
	et.template = t
	return &et
}

// Error generates an ErrorTrace on the current position.
func (t *Template) Error(description string, err error, ctx string) error {
	et := t.Prepare()
	et.SetDescription(description)
	et.SetError(err)
	et.SetContext(ctx)
	return et
}

func (et *ErrorTrace) isIncluded(call *Call) bool {
	if call.FunctionName == "" {
		return false
	}

	for _, f := range et.template.ignore {
		if strings.HasPrefix(call.FullName, f) {
			return false
		}
	}

	return true
}

func (et *ErrorTrace) filterTraceCalls(calls []Call) []Call {
	var filtered []Call

	for i := 0; i < len(calls); i++ {
		if et.isIncluded(&calls[i]) {
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
		calls = append(calls, makeCall(&frame))
		if !more {
			return calls
		}
	}
}

func (et *ErrorTrace) formatStack(calls []Call) string {
	if len(calls) == 0 {
		return "<nil>"
	}
	var sb strings.Builder

	if n := len(calls); n > 0 {
		cl := color.New(color.FgBlue, color.OpUnderscore)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%s(...)\n\t", calls[i].FullName)
			sb.WriteString(cl.Sprintf("%s:%d", calls[i].File, calls[i].Line))
			fmt.Fprintln(&sb)
		}
	}
	return sb.String()
}

// SetDescription sets the description that should be printed for this error.
func (et *ErrorTrace) SetDescription(description string) {
	et.description = description
}

// SetError sets the error.
func (et *ErrorTrace) SetError(err error) {
	et.error = err
}

// SetContext sets the context information that should be printed for this error.
func (et *ErrorTrace) SetContext(ctx string) {
	et.ctx = ctx
}

// Error returns the string representation for the error. It includes the Stacktrace, Description and Context.
func (et *ErrorTrace) Error() string {
	// collect the current trace
	// skip  runtime.Callers(), ett.currentTraceCalls(...), ett.(*ErrorTrace).Format(...), call to this function
	//nolint:gomnd
	traceCalls := resolveTraceCalls(currentTraceCalls(4))

	// if we have a inherited trace resolve the items and filter the current trace
	if et.inheritedTrace {
		// resolve the inherited calls
		traceCalls = append(resolveTraceCalls(et.inheritedPC), traceCalls...)
	}

	// filter
	traceCalls = et.filterTraceCalls(traceCalls)

	// print
	var sb strings.Builder
	if et.description != "" {
		sb.WriteString(minitest.Format("Description:", et.description, color.FgBlue))
	}
	sb.WriteString(minitest.Format("Error:      ", et.error.Error(), color.FgRed))
	sb.WriteString(minitest.Format("Error Trace:", et.formatStack(traceCalls)))
	if et.ctx != "" {
		sb.WriteString(minitest.Format("Context:    ", et.ctx))
	}

	return sb.String()
}

// ErrorText returns only the error text for the error.
func (et *ErrorTrace) ErrorText() string {
	return et.error.Error()
}

// Implement xerrors

// Is implements the xerrors interface so we can use the xerrors.Is() function.
func (et *ErrorTrace) Is(err error) bool {
	return et.error == err
}

// Unwrap implements the xerrors.Wrapper interface.
func (et *ErrorTrace) Unwrap() error {
	return et.error
}
