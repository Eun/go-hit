package doctest

import (
	"reflect"
	"testing"
)

func TestParseDoc(t *testing.T) {
	tests := []struct {
		Input    string
		Expected *Doc
	}{
		{
			Input: `Line #1
Line #2

Usage:
    Start Usage

    End Usage

Example:
    Start Example

    End Example
`,
			Expected: &Doc{
				Description: "Line #1\nLine #2",
				Fields: map[string]string{
					"Usage":   "Start Usage\n\nEnd Usage",
					"Example": "Start Example\n\nEnd Example",
				},
			},
		},

		{
			Input: `Description

Example:
    Example
`,
			Expected: &Doc{
				Description: "Description",
				Fields: map[string]string{
					"Example": "Example",
				},
			},
		},
	}

	for _, test := range tests {
		actual, err := ParseDoc(test.Input)
		if err != nil {
			t.Error(err)
		}

		if actual.Description != test.Expected.Description {
			t.Fatalf("\nexpected: %s\ngot:       %s", test.Expected.Description, actual.Description)
		}

		if !reflect.DeepEqual(actual.Fields, test.Expected.Fields) {
			t.Fatalf("\nexpected: %#v\ngot:      %#v", test.Expected, actual)
		}
	}
}
