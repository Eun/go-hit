package expr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetValue(t *testing.T) {
	heyStack := map[string]interface{}{
		"Name":   "Joe",
		"UserID": 10,
		"Roles":  []string{"Admin", "User"},
		"Details": map[string]interface{}{
			"Surname": "Doe",
			"Email":   "joe@example.com",
		},
		"Company": struct {
			ID   int
			Name string
		}{
			1,
			"Wood Inc",
		},
	}

	tests := []struct {
		Name       string
		Heystack   interface{}
		Expression string
		Options    options
		Ok         bool
		Value      interface{}
		ErrorText  string
	}{
		{
			"All",
			heyStack,
			"",
			options{},
			true,
			heyStack,
			"",
		},
		{
			"String",
			heyStack,
			"Name",
			options{},
			true,
			"Joe",
			"",
		},
		{
			"Int",
			heyStack,
			"UserID",
			options{},
			true,
			10,
			"",
		},
		{
			"Slice",
			heyStack,
			"Roles.0",
			options{},
			true,
			"Admin",
			"",
		},
		{
			"Slice (close to bounds)",
			heyStack,
			"Roles.1",
			options{},
			true,
			"User",
			"",
		},
		{
			"Slice (out of bounds)",
			heyStack,
			"Roles.-1",
			options{},
			false,
			nil,
			"",
		},
		{
			"Slice (out of bounds)",
			heyStack,
			"Roles.10",
			options{},
			false,
			nil,
			"",
		},
		{
			"Slice (Content)",
			heyStack,
			"Roles.Admin",
			options{},
			true,
			"Admin",
			"",
		},
		{
			"Map",
			heyStack,
			"Details.Surname",
			options{},
			true,
			"Doe",
			"",
		},
		{
			"Map (Wrong Case)",
			heyStack,
			"Details.sUrName",
			options{},
			false,
			nil,
			"",
		},
		{
			"Map (Wrong Case)",
			heyStack,
			"Details.sUrName",
			options{IgnoreCase},
			true,
			"Doe",
			"",
		},
		{
			"Map First Object",
			heyStack,
			"Details.0",
			options{},
			true,
			"joe@example.com",
			"",
		},
		{
			"Map Second Object",
			heyStack,
			"Details.1",
			options{},
			true,
			"Doe",
			"",
		},
		{
			"Map (out of bounds)",
			heyStack,
			"Details.2",
			options{},
			false,
			nil,
			"",
		},
		{
			"Map (not existing key)",
			heyStack,
			"Details.Address",
			options{},
			false,
			nil,
			"",
		},
		{
			"Struct",
			heyStack,
			"Company.Name",
			options{},
			true,
			"Wood Inc",
			"",
		},
		{
			"Struct (Wrong Case)",
			heyStack,
			"Company.nAmE",
			options{},
			false,
			nil,
			"",
		},
		{
			"Struct (Wrong Case)",
			heyStack,
			"Company.nAmE",
			options{IgnoreCase},
			true,
			"Wood Inc",
			"",
		},
		{
			"Struct First Object",
			heyStack,
			"Company.0",
			options{},
			true,
			1,
			"",
		},
		{
			"Struct Second Object",
			heyStack,
			"Company.1",
			options{},
			true,
			"Wood Inc",
			"",
		},
		{
			"Struct (out of bounds)",
			heyStack,
			"Company.2",
			options{},
			false,
			nil,
			"",
		},
		{
			"Struct (not existing key)",
			heyStack,
			"Company.Address",
			options{},
			false,
			nil,
			"",
		},
		{
			"Nil",
			nil,
			"",
			options{},
			false,
			nil,
			"",
		},
		{
			"Nil interface{}",
			map[string]interface{}{"Foo": nil},
			"Foo",
			options{},
			true,
			nil,
			"",
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.Name, func(t *testing.T) {
			v, ok, err := GetValue(test.Heystack, test.Expression, test.Options...)
			if test.ErrorText == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, test.ErrorText)
			}
			require.Equal(t, test.Ok, ok)
			require.Equal(t, test.Value, v)
		})
	}
}

func TestGetValue_EmbeddedStruct(t *testing.T) {
	type Company struct {
		Name string
		ID   int
	}

	type User struct {
		Name string
		Company
	}

	require.Equal(t, "Joe", MustGetValue(User{Name: "Joe", Company: Company{Name: "Wood Inc", ID: 10}}, "Name"))
	require.Equal(t, "Wood Inc", MustGetValue(User{Name: "Joe", Company: Company{Name: "Wood Inc", ID: 10}}, "Company.Name"))
}

func TestInvalidType(t *testing.T) {
	_, _, err := GetValue("Hello World", "0")
	require.EqualError(t, err, "string cannot be used with the expression 0")
}

func TestMustGetValue(t *testing.T) {
	t.Run("", func(t *testing.T) {
		require.PanicsWithValue(t, "string cannot be used with the expression 0", func() {
			_ = MustGetValue("Hello World", "0")
		})
	})
	t.Run("", func(t *testing.T) {
		require.PanicsWithValue(t, "unable to find value with Name in map[string]string{}", func() {
			_ = MustGetValue(map[string]string{}, "Name")
		})
	})
	t.Run("", func(t *testing.T) {
		require.Equal(t, "Joe", MustGetValue(map[string]string{"Name": "Joe"}, "Name"))
	})
}
