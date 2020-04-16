package hit

import (
	"errors"
	"testing"

	"reflect"

	"sort"
	"strings"

	"github.com/stretchr/testify/require"
)

func GetFunctions(a interface{}) ([]string, error) {
	if a == nil {
		return nil, errors.New("nil value")
	}
	v := reflect.ValueOf(a)
	if !v.IsValid() {
		return nil, errors.New("invalid value")
	}
	var funcs []string
	for v.IsValid() {
		size := v.Type().NumMethod()
		for i := 0; i < size; i++ {
			funcs = append(funcs, v.Type().Method(i).Name)
		}
		if v.Kind() != reflect.Ptr {
			break
		}
		v = v.Elem()
	}

	sort.Slice(funcs, func(i, j int) bool {
		return strings.Compare(funcs[i], funcs[j]) < 0
	})

	return funcs, nil
}

func ImplementsAllFunctionsOf(t *testing.T, a, b interface{}) {
	funcs1, err := GetFunctions(a)
	if err != nil {
		t.Fatalf("unable to get functions for %#v: %v", a, err)
	}

	funcs2, err := GetFunctions(b)
	if err != nil {
		t.Fatalf("unable to get functions for %#v: %v", b, err)
	}

	require.Equal(t, funcs1, funcs2)
}

func TestClearSend(t *testing.T) {
	ImplementsAllFunctionsOf(t, Send(), Clear().Send())
	ImplementsAllFunctionsOf(t, Send().Body(), Clear().Send().Body())
}

func TestClearExpect(t *testing.T) {
	ImplementsAllFunctionsOf(t, Expect(), Clear().Expect())
	ImplementsAllFunctionsOf(t, Expect().Body(), Clear().Expect().Body())
	ImplementsAllFunctionsOf(t, Expect().Body().JSON(), Clear().Expect().Body().JSON())
	ImplementsAllFunctionsOf(t, Expect().Header(), Clear().Expect().Header())
	ImplementsAllFunctionsOf(t, Expect().Trailer(), Clear().Expect().Trailer())
	ImplementsAllFunctionsOf(t, Expect().Status(), Clear().Expect().Status())
}
