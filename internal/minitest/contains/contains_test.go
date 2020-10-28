package contains

import (
	"testing"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		haystack interface{}
		needle   interface{}
		want     bool
		wantErr  bool
	}{
		// string
		{
			"string contains string",
			"Hello World 123",
			"Hello",
			true,
			false,
		},
		{
			"string contains int",
			"Hello World 123",
			123,
			true,
			false,
		},
		{
			"string not contains int",
			"Hello World 123",
			456,
			false,
			false,
		},
		{
			"string not contains struct",
			"Hello World 123",
			struct{}{},
			false,
			false,
		},

		// slice
		{
			"slice contains string",
			[]string{"Hello", "World", "123"},
			"Hello",
			true,
			false,
		},
		{
			"slice contains int",
			[]string{"Hello", "World", "123"},
			123,
			true,
			false,
		},
		// {
		// 	"slice contains slice",
		// 	[]string{"Hello", "World", "123"},
		// 	[]string{"Hello"},
		// 	true,
		// 	false,
		// },
		{
			"slice not contains int",
			[]string{"Hello", "World", "123"},
			456,
			false,
			false,
		},
		{
			"slice not contains struct",
			[]string{"Hello", "World", "123"},
			struct{}{},
			false,
			false,
		},

		// map
		{
			"map contains string",
			map[string]interface{}{"Name": "Joe", "Id": 10},
			"Name",
			true,
			false,
		},
		{
			"map not contains struct",
			map[string]interface{}{"Name": "Joe", "Id": 10},
			struct{}{},
			false,
			false,
		},
		// not required for now
		// {
		// 	"map contains map",
		// 	map[string]interface{}{"Name": "Joe", "Id": 10},
		// 	map[string]interface{}{"Name": "Joe"},
		// 	true,
		// 	nil,
		// },

		// struct
		{
			"struct contains string",
			struct{ Name string }{Name: "Joe"},
			"Name",
			true,
			false,
		},
		// not required for now
		// {
		// 	"struct contains struct",
		// 	struct {
		// 		Name string
		// 		Id   int
		// 	}{Name: "Joe", Id: 10},
		// 	struct{ Name string }{Name: "Joe"},
		// 	true,
		// 	nil,
		// },
		{
			"struct not contains struct",
			struct{ Name string }{Name: "Joe"},
			struct{}{},
			false,
			false,
		},
		// edge cases
		{
			"nil haystack",
			nil,
			0,
			false,
			false,
		},
		{
			"nil needle",
			0,
			nil,
			false,
			false,
		},
		{
			"nil haystack and needle",
			nil,
			nil,
			true,
			false,
		},
		{
			"unknown haystack type",
			1,
			0,
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Contains(tt.haystack, tt.needle)
			if (err != nil) != tt.wantErr {
				t.Errorf("Contains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Contains() got = %v, want %v", got, tt.want)
			}
		})
	}
}
