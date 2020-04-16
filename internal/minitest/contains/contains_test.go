package contains

import (
	"testing"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		heystack interface{}
		needle   interface{}
		want     bool
	}{
		// string
		{
			"string contains string",
			"Hello World 123",
			"Hello",
			true,
		},
		{
			"string contains int",
			"Hello World 123",
			123,
			true,
		},
		{
			"string not contains int",
			"Hello World 123",
			456,
			false,
		},
		{
			"string not contains struct",
			"Hello World 123",
			struct{}{},
			false,
		},

		// slice
		{
			"slice contains string",
			[]string{"Hello", "World", "123"},
			"Hello",
			true,
		},
		{
			"slice contains int",
			[]string{"Hello", "World", "123"},
			123,
			true,
		},
		{
			"slice not contains int",
			[]string{"Hello", "World", "123"},
			456,
			false,
		},
		{
			"slice not contains struct",
			[]string{"Hello", "World", "123"},
			struct{}{},
			false,
		},

		// map
		{
			"map contains string",
			map[string]interface{}{"Name": "Joe", "Id": 10},
			"Name",
			true,
		},
		{
			"map not contains struct",
			map[string]interface{}{"Name": "Joe", "Id": 10},
			struct{}{},
			false,
		},
		// not required for now
		// {
		// 	"map contains map",
		// 	map[string]interface{}{"Name": "Joe", "Id": 10},
		// 	map[string]interface{}{"Name": "Joe"},
		// 	true,
		// },

		// struct
		{
			"struct contains string",
			struct{ Name string }{Name: "Joe"},
			"Name",
			true,
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
		// },
		{
			"struct not contains struct",
			struct{ Name string }{Name: "Joe"},
			struct{}{},
			false,
		},
		// edge cases
		{
			"nil haystack",
			nil,
			0,
			false,
		},
		{
			"nil needle",
			0,
			nil,
			false,
		},
		{
			"nil haystack and needle",
			nil,
			nil,
			true,
		},
		{
			"unknown haystack type",
			1,
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.heystack, tt.needle); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
