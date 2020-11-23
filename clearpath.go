package hit

import (
	"fmt"

	"strings"

	"io"

	"github.com/google/go-cmp/cmp"
)

type pathPart struct {
	Func      string
	Arguments []interface{}
}

type callPath []pathPart

func (cleanPath callPath) Push(name string, arguments []interface{}) callPath {
	if arguments == nil {
		arguments = []interface{}{}
	}
	return append(cleanPath, pathPart{
		Func:      name,
		Arguments: arguments,
	})
}

func newCallPath(name string, arguments []interface{}) callPath {
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
	return fmt.Sprintf("%p", x) == fmt.Sprintf("%p", y)
}

func ioReaderComparer(x, y io.Reader) bool {
	return fmt.Sprintf("%p", x) == fmt.Sprintf("%p", y)
}

func (cleanPath callPath) Contains(needle callPath) bool {
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
		if !cmp.Equal(cleanPath[i].Arguments[:hayArgSize], needle[i].Arguments,
			cmp.Comparer(funcComparer),
			cmp.Comparer(ioReaderComparer),
		) {
			return false
		}
	}
	return true
}

func (cleanPath callPath) Equal(b callPath) bool {
	aSize := len(cleanPath)
	bSize := len(b)
	if bSize != aSize {
		return false
	}
	for i := 0; i < aSize; i++ {
		if cleanPath[i].Func != b[i].Func {
			return false
		}

		aArgSize := len(cleanPath[i].Arguments)
		bArgSize := len(b[i].Arguments)

		if aArgSize > bArgSize {
			aArgSize = bArgSize
		}
		if !cmp.Equal(cleanPath[i].Arguments[:aArgSize], b[i].Arguments,
			cmp.Comparer(funcComparer),
			cmp.Comparer(ioReaderComparer),
		) {
			return false
		}
	}
	return true
}

func (cleanPath callPath) CallString(withArguments bool) string {
	parts := make([]string, len(cleanPath))
	for i := range cleanPath {
		var args []string
		if withArguments {
			args = make([]string, len(cleanPath[i].Arguments))
			for i, argument := range cleanPath[i].Arguments {
				switch v := argument.(type) {
				case []byte:
					values := make([]string, len(v))
					for j, b := range v {
						values[j] = fmt.Sprintf("%#v", b)
					}
					args[i] = fmt.Sprintf("[]uint8{%s}", strings.Join(values, ", "))
				case float32, float64:
					args[i] = fmt.Sprintf("%#f", argument)
				default:
					args[i] = fmt.Sprintf("%#v", argument)
				}
			}
		}
		parts[i] = fmt.Sprintf("%s(%s)", cleanPath[i].Func, strings.Join(args, ", "))
	}
	return strings.Join(parts, ".")
}
