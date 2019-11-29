package errortrace

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/Eun/go-hit/internal/minitest"
	"github.com/gookit/color"
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
	et.inheritedPC = CurrentTraceCalls(4)
	return &et
}

func isIncluded(call *Call) bool {
	if call.FunctionName == "" {
		return false
	}
	if !strings.HasSuffix(call.PackageName, "github.com/Eun/go-hit") {
		return false
	}

	switch call.FunctionName {
	case "Test", "Do", "Custom", "runSteps", "exec":
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

func CurrentTraceCalls(skip int) []uintptr {
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

func ResolveTraceCalls(pc []uintptr) []Call {
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

func (p *ErrorTrace) Format(description, errText string) ErrorTraceError {
	// collect the current trace
	traceCalls := ResolveTraceCalls(CurrentTraceCalls(4))

	// if we have a inherited trace resolve the items and filter the current trace
	if p.inheritedTrace {
		// filter
		traceCalls = filterTraceCalls(traceCalls)

		// resolve the inherited calls
		traceCalls = append(traceCalls, ResolveTraceCalls(p.inheritedPC)...)
	}

	// print
	var sb strings.Builder
	if description != "" {
		sb.WriteString(minitest.Format("Description:", description, color.FgBlue))
	}
	sb.WriteString(minitest.Format("Error:      ", errText, color.FgRed))
	sb.WriteString(minitest.Format("Error Trace:", p.formatStack(traceCalls)))
	return ErrorTraceError(sb.String())
}

func Format(description, errText string) ErrorTraceError {
	var et ErrorTrace
	return et.Format(description, errText)
}
