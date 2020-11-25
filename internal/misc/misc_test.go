package misc

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeURL(t *testing.T) {
	tests := []struct {
		Base     string
		URL      string
		Expected string
	}{
		{
			"http://example.com",
			"index.html",
			"http://example.com/index.html",
		},
		{
			"http://example.com/",
			"index.html",
			"http://example.com/index.html",
		},
		{
			"http://example.com",
			"/index.html",
			"http://example.com/index.html",
		},
		{
			"http://example.com",
			"",
			"http://example.com",
		},
		{
			"http://example.com/",
			"",
			"http://example.com/",
		},
		{
			"",
			"index.html",
			"index.html",
		},
		{
			"",
			"/index.html",
			"/index.html",
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run("", func(t *testing.T) {
			if s := MakeURL(test.Base, test.URL); s != test.Expected {
				t.Errorf("expected `%s' got `%s'", test.Expected, s)
			}
		})
	}
}

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

func TestMakeTypeCopy(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}
	tests := []struct {
		Source   interface{}
		Expected interface{}
	}{
		{
			map[string]interface{}{
				"Name": "Joe",
			},
			map[string]interface{}{
				"Name": "",
			},
		},
		{
			map[string]string{
				"Name": "Joe",
			},
			map[string]string{
				"Name": "",
			},
		},
		{
			map[string]interface{}{
				"UserName": map[string]interface{}{
					"Name": "Joe",
				},
			},
			map[string]interface{}{
				"UserName": map[string]interface{}{
					"Name": "",
				},
			},
		},
		{
			map[string]interface{}{
				"Details": map[string]interface{}{
					"UserName": map[string]interface{}{
						"Name": "Joe",
					},
				},
			},
			map[string]interface{}{
				"Details": map[string]interface{}{
					"UserName": map[string]interface{}{
						"Name": "",
					},
				},
			},
		},
		{
			[]interface{}{"Hello", "World"},
			[]interface{}{"", ""},
		},
		{
			[]interface{}{
				map[string]interface{}{
					"UserName": map[string]interface{}{
						"Name": "Joe",
					},
				},
			},
			[]interface{}{
				map[string]interface{}{
					"UserName": map[string]interface{}{
						"Name": "",
					},
				},
			},
		},
		{
			User{
				ID:   10,
				Name: "Joe",
			},
			User{
				ID:   0,
				Name: "",
			},
		},
		{
			[]interface{}{User{10, "Joe"}, User{10, "Joe"}},
			[]interface{}{User{0, ""}, User{0, ""}},
		},
	}
	for _, test := range tests {
		actual := MakeTypeCopy(test.Source)
		if !reflect.DeepEqual(actual, test.Expected) {
			t.Errorf("want %#v; actual %#v", test.Expected, actual)
		}
	}
}
