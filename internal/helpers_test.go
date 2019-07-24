package internal

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetValue(t *testing.T) {
	ptrString := func(v string) *string {
		return &v
	}
	ptrInterface := func(v interface{}) *interface{} {
		return &v
	}
	ptr := func(v interface{}) interface{} {
		return &v
	}
	tests := []struct {
		Value        interface{}
		ExpectedKind reflect.Kind
	}{
		{
			"Hello",
			reflect.String,
		},
		{
			ptrString("Hello"),
			reflect.String,
		},
		{
			ptrInterface("Hello"),
			reflect.String,
		},
		{
			ptr("Hello"),
			reflect.String,
		},
	}

	for i := range tests {
		test := tests[i]
		v := GetValue(test.Value)
		require.True(t, v.IsValid())
		require.Equal(t, test.ExpectedKind, v.Kind())
	}
}
