package errortrace

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeCall(t *testing.T) {
	tests := []struct {
		Name         string
		Frame        runtime.Frame
		ExpectedCall Call
	}{
		{"short",
			runtime.Frame{
				PC:       0,
				Func:     nil,
				Function: "packagename.Test",
				File:     "somefile.go",
				Line:     10,
				Entry:    0,
			},
			Call{
				PackageName:  "packagename",
				FunctionName: "Test",
				File:         "somefile.go",
				Line:         10,
				FullName:     "packagename.Test",
			},
		},
		{"with function path",
			runtime.Frame{
				PC:       0,
				Func:     nil,
				Function: "github.com/otto-eng/go-hit/ett.(*defaultInstance).Trace.func1",
				File:     "somefile.go",
				Line:     10,
				Entry:    0,
			},
			Call{
				PackageName:  "github.com/otto-eng/go-hit/ett",
				FunctionPath: "(*defaultInstance).Trace",
				FunctionName: "func1",
				File:         "somefile.go",
				Line:         10,
				FullName:     "github.com/otto-eng/go-hit/ett.(*defaultInstance).Trace.func1",
			},
		},
		{"full",
			runtime.Frame{
				PC:       0,
				Func:     nil,
				Function: "github.com/otto-eng/go-hit/ett.Trace",
				File:     "somefile.go",
				Line:     10,
				Entry:    0,
			},
			Call{
				PackageName:  "github.com/otto-eng/go-hit/ett",
				FunctionName: "Trace",
				File:         "somefile.go",
				Line:         10,
				FullName:     "github.com/otto-eng/go-hit/ett.Trace",
			},
		},
		{"no package",
			runtime.Frame{
				PC:       0,
				Func:     nil,
				Function: "Trace",
				File:     "somefile.go",
				Line:     10,
				Entry:    0,
			},
			Call{
				PackageName:  "",
				FunctionName: "Trace",
				File:         "somefile.go",
				Line:         10,
				FullName:     "Trace",
			},
		},
		{"no package",
			runtime.Frame{
				PC:       0,
				Func:     nil,
				Function: ".Trace",
				File:     "somefile.go",
				Line:     10,
				Entry:    0,
			},
			Call{
				PackageName:  "",
				FunctionName: "Trace",
				File:         "somefile.go",
				Line:         10,
				FullName:     "Trace",
			},
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.Name, func(t *testing.T) {
			require.Equal(t, test.ExpectedCall, makeCall(&test.Frame))
		})
	}
}
