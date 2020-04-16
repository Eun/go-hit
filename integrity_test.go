package hit_test

import (
	"errors"
	"testing"

	"reflect"

	"sort"
	"strings"

	"github.com/Eun/go-hit"
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
	ImplementsAllFunctionsOf(t, hit.Send(), hit.Clear().Send())
	ImplementsAllFunctionsOf(t, hit.Send().Body(), hit.Clear().Send().Body())
}

func TestClearExpect(t *testing.T) {
	ImplementsAllFunctionsOf(t, hit.Expect(), hit.Clear().Expect())
	ImplementsAllFunctionsOf(t, hit.Expect().Body(), hit.Clear().Expect().Body())
	ImplementsAllFunctionsOf(t, hit.Expect().Body().JSON(), hit.Clear().Expect().Body().JSON())
	ImplementsAllFunctionsOf(t, hit.Expect().Header(), hit.Clear().Expect().Header())
	ImplementsAllFunctionsOf(t, hit.Expect().Trailer(), hit.Clear().Expect().Trailer())
	ImplementsAllFunctionsOf(t, hit.Expect().Status(), hit.Clear().Expect().Status())
}
