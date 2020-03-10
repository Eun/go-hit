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
			},
		},
		{"with function path",
			runtime.Frame{
				PC:       0,
				Func:     nil,
				Function: "github.com/Eun/go-hit/errortrace.(*defaultInstance).Trace.func1",
				File:     "somefile.go",
				Line:     10,
				Entry:    0,
			},
			Call{
				PackageName:  "github.com/Eun/go-hit/errortrace",
				FunctionPath: "(*defaultInstance).Trace",
				FunctionName: "func1",
				File:         "somefile.go",
				Line:         10,
			},
		},
		{"full",
			runtime.Frame{
				PC:       0,
				Func:     nil,
				Function: "github.com/Eun/go-hit/errortrace.Trace",
				File:     "somefile.go",
				Line:     10,
				Entry:    0,
			},
			Call{
				PackageName:  "github.com/Eun/go-hit/errortrace",
				FunctionName: "Trace",
				File:         "somefile.go",
				Line:         10,
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
			},
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.Name, func(t *testing.T) {
			require.Equal(t, test.ExpectedCall, makeCall(test.Frame))
		})
	}
}

func TestCall_FullName(t *testing.T) {
	t.Run("", func(t *testing.T) {
		c := Call{PackageName: "package", FunctionName: "Func"}
		require.Equal(t, "package.Func", c.FullName)
	})

	t.Run("", func(t *testing.T) {
		c := Call{PackageName: "package", FunctionName: ""}
		require.Equal(t, "package.", c.FullName)
	})
	t.Run("", func(t *testing.T) {
		c := Call{PackageName: "", FunctionName: "Func"}
		require.Equal(t, "Func", c.FullName)
	})
}

// func TestErrorTrace(t *testing.T) {
// 	t.Run("", func(t *testing.T) {
// 		tm, err := New(0)
// 		require.NoError(t, err)
// 		et := tm.Prepare()
//
// 		var wg sync.WaitGroup
// 		wg.Add(1)
// 		go func() {
// 			fmt.Println(et.Format("Hello World", "Some Error"))
// 			wg.Done()
// 		}()
//
// 		wg.Wait()
// 	})
// }
