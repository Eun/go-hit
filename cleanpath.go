package hit

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

type pathPart struct {
	Func      string
	Arguments []interface{}
}

type CleanPath []pathPart

func (p CleanPath) Push(name string, arguments []interface{}) CleanPath {
	return append(p, pathPart{
		Func:      name,
		Arguments: arguments,
	})
}

func NewCleanPath(name string, arguments []interface{}) CleanPath {
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

func (haystack CleanPath) Contains(needle CleanPath) bool {
	haySize := len(haystack)
	needleSize := len(needle)
	if needleSize > haySize {
		return false
	}
	haySize = needleSize

	for i := 0; i < haySize; i++ {
		if haystack[i].Func != needle[i].Func {
			return false
		}

		if len(needle[i].Arguments) > 0 {
			if !cmp.Equal(haystack[i].Arguments, needle[i].Arguments, cmp.Comparer(funcComparer)) {
				return false
			}
		}
	}
	return true
}
