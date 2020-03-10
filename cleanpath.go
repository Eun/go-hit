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
	return append(cleanPath, pathPart{
		Func:      name,
		Arguments: arguments,
	})
}

func newClearPath(name string, arguments []interface{}) clearPath {
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

		if len(needle[i].Arguments) > 0 {
			if !cmp.Equal(cleanPath[i].Arguments, needle[i].Arguments, cmp.Comparer(funcComparer)) {
				return false
			}
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
