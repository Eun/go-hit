package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestExpectBodyJSONDasel_Equal(t *testing.T) {
	payload := map[string]interface{}{
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
		"Bills": []struct {
			ID    int
			Total float64
		}{
			{
				ID:    21,
				Total: 124.23,
			},
			{
				ID:    25,
				Total: 42.55,
			},
		},
	}
	s := PrintJSONServer(payload)
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Dasel(".Name").Equal("Joe"),
		)
	})
	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Dasel(".Roles").Equal([]interface{}{"Admin", "User"}),
		)
	})

	t.Run("slice of string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Dasel(".Roles").Equal([]string{"Admin", "User"}),
		)
	})

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Dasel(".Details").Equal(map[string]interface{}{"Surname": "Doe", "Email": "joe@example.com"}),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Dasel(".UserID").Equal(10),
		)
	})

	t.Run("struct", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Dasel(".Company").Equal(struct {
				ID   int
				Name string
			}{
				1,
				"Wood Inc",
			}),
		)
	})

	t.Run("full payload", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Dasel(".").Equal(payload),
		)
	})

	t.Run("nil", func(t *testing.T) {
		t.Run("equal", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Expect().Body().JSON().Dasel(".NotExistent").Equal(nil),
			)
		})

		t.Run("nil in expect", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Expect().Body().JSON().Dasel(".UserID").Equal(nil),
				),
				PtrStr("not equal"), nil, nil, nil, nil, nil,
			)
		})

		t.Run("nil in response", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Expect().Body().JSON().Dasel(".NotExistent").Equal("Hello World"),
				),
				PtrStr("not equal"), nil, nil, nil, nil, nil,
			)
		})
	})

	t.Run("multi dasel", func(t *testing.T) {
		t.Run("param", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Expect().Body().JSON().Dasel(".Details", ".Email").Equal("joe@example.com"),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Expect().Body().JSON().Dasel(".Details", ".NotExistent").Equal("Hello World"),
				),
				PtrStr("not equal"), nil, nil, nil, nil, nil,
			)
		})
		t.Run("func", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Expect().Body().JSON().Dasel(".Details").Dasel(".Email").Equal("joe@example.com"),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Expect().Body().JSON().Dasel(".Details").Dasel(".NotExistent").Equal("Hello World"),
				),
				PtrStr("not equal"), nil, nil, nil, nil, nil,
			)
		})
		t.Run("array", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Description("first expression returning an array"),
				Expect().Body().JSON().Dasel(".Bills", ".[*].ID").Equal([]int{21, 25}),
			)
			Test(t,
				Post(s.URL),
				Description("first expression returning single json objects"),
				Expect().Body().JSON().Dasel(".Bills.[*]", ".[*].ID").Equal([]int{21, 25}),
			)
		})
	})

	t.Run("stream", func(t *testing.T) {
		se := EchoServer()
		defer se.Close()
		Test(t,
			Post(se.URL),
			Send().Body().JSON([]map[string]interface{}{
				{
					"id":   1,
					"name": "Joe",
					"role": "Admin",
				},
				{
					"id":   2,
					"name": "Alice",
					"role": "Guest",
				},
				{
					"id":   3,
					"name": "Bob",
					"role": "Admin",
				},
			}),
			Expect().Body().JSON().Dasel(`.(?:role=Admin).name`).Equal([]string{"Joe", "Bob"}),
		)
	})

	t.Run("jq", func(t *testing.T) {
		t.Run("func", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Expect().Body().JSON().Dasel(".Details").JQ(".Email").Equal("joe@example.com"),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Expect().Body().JSON().Dasel(".Details").JQ(".NotExistent").Equal("Hello World"),
				),
				PtrStr("not equal"), nil, nil, nil, nil, nil,
			)
		})
	})
}

func TestExpectBodyJSONDasel_NotEqual(t *testing.T) {
	payload := map[string]interface{}{
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
	s := PrintJSONServer(payload)
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Dasel(".Name").NotEqual("Alice"),
		)
	})
	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Dasel(".Roles").NotEqual([]interface{}{"Admin", "Developer"}),
		)
	})

	t.Run("slice of string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Dasel(".Roles").NotEqual([]string{"Admin", "Developer"}),
		)
	})

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Dasel(".Details").NotEqual(map[string]interface{}{"Surname": "Joe", "Email": "joe@example.com"}),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Dasel(".UserID").NotEqual(5),
		)
	})

	t.Run("struct", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Dasel(".Company").NotEqual(struct {
				ID   int
				Name string
			}{
				2,
				"Wood Inc",
			}),
		)
	})

	t.Run("full payload", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Expect().Body().JSON().Dasel(".").NotEqual(payload),
			),
			PtrStr("should not be map[string]interface {}{"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		)
	})

	t.Run("nil", func(t *testing.T) {
		t.Run("equal", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Expect().Body().JSON().Dasel(".NotExistent").NotEqual(nil),
				),
				PtrStr("should not be nil"),
			)
		})

		t.Run("nil in expect", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Expect().Body().JSON().Dasel(".UserID").NotEqual(nil),
			)
		})

		t.Run("nil in response", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Expect().Body().JSON().Dasel(".NotExistent").NotEqual("Hello World"),
			)
		})
	})
}

func TestExpectBodyJSONDasel_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	Test(t,
		Post(s.URL),
		Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Id": 10}),
		Expect().Body().JSON().Dasel(".Name").Contains("Joe"),
	)
}

func TestExpectBodyJSONDasel_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	Test(t,
		Post(s.URL),
		Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Id": 10}),
		Expect().Body().JSON().Dasel(".Name").NotContains("A"),
	)
}

func TestExpectBodyJSONDasel_NoJSON(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	ExpectError(t, Do(
		Head(s.URL),
		Expect().Body().JSON().Dasel(".").Equal(""),
	),
		PtrStr(`EOF`),
	)
}

func TestExpectBodyJSONDasel_Len(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String(`["Hello", "World"]`),
		Expect().Body().JSON().Dasel(".").Len().Equal(2),
	)

	Test(t,
		Post(s.URL),
		Send().Body().String(`{"Name":"Joe", "ID": 10}`),
		Expect().Body().JSON().Dasel(".").Len().Equal(2),
	)

	Test(t,
		Post(s.URL),
		Send().Body().String(`{"Name":"Joe", "ID": 10}`),
		Expect().Body().JSON().Dasel(".Name").Len().Equal(3),
	)

	Test(t,
		Post(s.URL),
		Send().Body().String(`{"Name":"Joe", "ID": 10}`),
		Expect().Body().JSON().Dasel(".Group").Len().Equal(0),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String(`"Hello World"`),
			Expect().Body().JSON().Dasel(".").Len().Equal(10),
		),
		PtrStr("not equal"), PtrStr("expected: 10"), PtrStr("actual: 11"), nil, nil, nil, nil,
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String(`10`),
			Expect().Body().JSON().Dasel(".").Len().Equal(10),
		),
		PtrStr("cannot get len for 10"),
	)
}
