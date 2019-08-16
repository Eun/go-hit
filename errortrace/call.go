package errortrace

import (
	"runtime"
	"strings"
)

type Call struct {
	PackageName  string
	FunctionPath string
	FunctionName string
	File         string
	Line         int
	PC           uintptr
	Entry        uintptr
}

func (c *Call) FullName() string {
	var sb strings.Builder
	if c.PackageName != "" {
		sb.WriteString(c.PackageName)
		sb.WriteRune('.')
	}
	if c.FunctionPath != "" {
		sb.WriteString(c.FunctionPath)
		sb.WriteRune('.')
	}
	sb.WriteString(c.FunctionName)
	return sb.String()
}

func makeCall(frame runtime.Frame) Call {
	// find the last slash
	lastSlash := strings.LastIndexFunc(frame.Function, func(r rune) bool {
		return r == '/'
	})
	if lastSlash <= -1 {
		lastSlash = 0
	}

	call := Call{
		File:  frame.File,
		Line:  frame.Line,
		PC:    frame.PC,
		Entry: frame.Entry,
	}

	// the first dot after the slash ends the package name
	dot := strings.IndexRune(frame.Function[lastSlash:], '.')
	if dot < 0 {
		// no dot means no package
		call.FunctionName = frame.Function
	} else {
		dot += lastSlash
		call.PackageName = frame.Function[:dot]
		call.FunctionName = strings.TrimLeft(frame.Function[dot:], ".")
	}

	parts := strings.FieldsFunc(call.FunctionName, func(r rune) bool {
		return r == '.'
	})

	size := len(parts)
	if size <= 1 {
		return call
	}
	size--

	call.FunctionPath = strings.Join(parts[:size], ".")
	call.FunctionName = parts[size]
	return call
}
