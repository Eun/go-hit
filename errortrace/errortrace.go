package errortrace

import (
	"runtime"
)

type ErrorTrace struct {
	Panic Panicer
}

// Prepare collects the current trace, if later a Panic will be called the collected trace will be included
// in the error trace
func Prepare() *ErrorTrace {
	var et ErrorTrace
	et.Panic.inheritedTrace = true
	et.Panic.inheritedPC = currentTraceCalls(4)
	return &et
}

// Panic panics directly with the current full stack
var Panic = Panicer{}

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
