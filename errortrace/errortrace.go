package errortrace

import (
	"fmt"
	"runtime"
	"strings"

	"reflect"

	"github.com/Eun/go-hit/internal/minitest"
	"github.com/gookit/color"
	"golang.org/x/xerrors"
)

type ErrorTraceTemplate struct {
	ignore []string
}

type ErrorTrace struct {
	inheritedTrace bool
	inheritedPC    []uintptr
	template       *ErrorTraceTemplate
}

type ErrorTraceError string

func (e ErrorTraceError) Error() string {
	return string(e)
}

func New(ignore ...string) *ErrorTraceTemplate {
	return &ErrorTraceTemplate{
		ignore: ignore,
	}
}

func IgnoreFunc(fn interface{}) string {
	v := reflect.ValueOf(fn)
	if !v.IsValid() {
		panic(xerrors.New("function is is not valid"))
	}
	return strings.TrimSuffix(runtime.FuncForPC(v.Pointer()).Name(), "-fm")
}

func IgnorePackage(fn interface{}) string {
	return makeCall(runtime.Frame{
		PC:       0,
		Func:     nil,
		Function: IgnoreFunc(fn),
		File:     "",
		Line:     0,
		Entry:    0,
	}).PackageName
}

func IgnoreStruct(fn interface{}) string {
	c := makeCall(runtime.Frame{
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
// in the error trace
func (t *ErrorTraceTemplate) Prepare() *ErrorTrace {
	var et ErrorTrace
	et.inheritedTrace = true
	// skip  runtime.Callers(), errortrace.currentTraceCalls(...), errortrace.(*ErrorTrace).Prepare(...), Call to this Function
	et.inheritedPC = currentTraceCalls(4)
	et.template = t
	return &et
}

func (t *ErrorTraceTemplate) Format(description, errText string) ErrorTraceError {
	return t.Prepare().Format(description, errText)
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
			exits := false
			for _, call := range filtered {
				if call.FullName == calls[i].FullName {
					exits = true
					break
				}
			}
			if !exits {
				filtered = append(filtered, calls[i])
			}
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

func (et *ErrorTrace) formatStack(calls []Call) string {
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

func (et *ErrorTrace) Format(description, errText string) ErrorTraceError {
	// collect the current trace
	// skip  runtime.Callers(), errortrace.currentTraceCalls(...), errortrace.(*ErrorTrace).Format(...), call to this function
	traceCalls := resolveTraceCalls(currentTraceCalls(4))

	// if we have a inherited trace resolve the items and filter the current trace
	if et.inheritedTrace {
		// resolve the inherited calls

		traceCalls = append(traceCalls, resolveTraceCalls(et.inheritedPC)...)
	}

	// filter
	traceCalls = et.filterTraceCalls(traceCalls)

	// print
	var sb strings.Builder
	if description != "" {
		sb.WriteString(minitest.Format("Description:", description, color.FgBlue))
	}
	sb.WriteString(minitest.Format("Error:      ", errText, color.FgRed))
	sb.WriteString(minitest.Format("Error Trace:", et.formatStack(traceCalls)))
	return ErrorTraceError(sb.String())
}
