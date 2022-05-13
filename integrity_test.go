// +build !generate

package hit_test

import (
	"errors"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/otto-eng/go-hit"
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

func TestIntegrityClearSend(t *testing.T) {
	ImplementsAllFunctionsOf(t, hit.Send(), hit.Clear().Send())
	ImplementsAllFunctionsOf(t, hit.Send().Body(), hit.Clear().Send().Body())
}

func TestIntegrityClearExpect(t *testing.T) {
	ImplementsAllFunctionsOf(t, hit.Expect(), hit.Clear().Expect())
	ImplementsAllFunctionsOf(t, hit.Expect().Body(), hit.Clear().Expect().Body())
	ImplementsAllFunctionsOf(t, hit.Expect().Body().Int(), hit.Clear().Expect().Body().Int())
	ImplementsAllFunctionsOf(t, hit.Expect().Body().JSON(), hit.Clear().Expect().Body().JSON())
	ImplementsAllFunctionsOf(t, hit.Expect().Body().String(), hit.Clear().Expect().Body().String())
	ImplementsAllFunctionsOf(t, hit.Expect().Headers(""), hit.Clear().Expect().Headers())
	ImplementsAllFunctionsOf(t, hit.Expect().Trailers(""), hit.Clear().Expect().Trailers())
	ImplementsAllFunctionsOf(t, hit.Expect().Status(), hit.Clear().Expect().Status())
}
