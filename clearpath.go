package hit

import (
	"fmt"

	"strings"

	"github.com/google/go-cmp/cmp"
)

type pathPart struct {
	Func      string
	Arguments []interface{}
}

type clearPath []pathPart

func (cleanPath clearPath) Push(name string, arguments []interface{}) clearPath {
	if arguments == nil {
		arguments = []interface{}{}
	}
	return append(cleanPath, pathPart{
		Func:      name,
		Arguments: arguments,
	})
}

func newClearPath(name string, arguments []interface{}) clearPath {
	if arguments == nil {
		arguments = []interface{}{}
	}
	return []pathPart{
		{
			Func:      name,
			Arguments: arguments,
		},
	}
}

func funcComparer(x, y func(Hit)) bool {
	// return runtime.FuncForPC(reflect.ValueOf(x).Pointer()).Name() == runtime.FuncForPC(reflect.ValueOf(y).Pointer()).Name()
	return fmt.Sprintf("%p", x) == fmt.Sprintf("%p", y)
}

func (cleanPath clearPath) Contains(needle clearPath) bool {
	haySize := len(cleanPath)
	needleSize := len(needle)
	if needleSize > haySize {
		return false
	}
	haySize = needleSize

	for i := 0; i < haySize; i++ {
		if cleanPath[i].Func != needle[i].Func {
			return false
		}

		hayArgSize := len(cleanPath[i].Arguments)
		needleArgSize := len(needle[i].Arguments)

		if hayArgSize > needleArgSize {
			hayArgSize = needleArgSize
		}
		if !cmp.Equal(cleanPath[i].Arguments[:hayArgSize], needle[i].Arguments, cmp.Comparer(funcComparer)) {
			return false
		}
	}
	return true
}

func (cleanPath clearPath) String() string {
	parts := make([]string, len(cleanPath))
	for i := range cleanPath {
		parts[i] = cleanPath[i].Func
	}
	return strings.Join(parts, ".")
}

func (cleanPath clearPath) CallString() string {
	parts := make([]string, len(cleanPath))
	for i := range cleanPath {
		args := make([]string, len(cleanPath[i].Arguments))
		for i, argument := range cleanPath[i].Arguments {
			args[i] = fmt.Sprintf("%#v", argument)
		}
		parts[i] = fmt.Sprintf("%s(%s)", cleanPath[i].Func, strings.Join(args, ", "))
	}
	return strings.Join(parts, ".")
}
