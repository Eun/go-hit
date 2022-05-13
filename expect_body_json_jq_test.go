package hit_test

import (
	"testing"

	. "github.com/otto-eng/go-hit"
)

func TestExpectBodyJSONJQ_Equal(t *testing.T) {
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
			Expect().Body().JSON().JQ(".Name").Equal("Joe"),
		)
	})
	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().JQ(".Roles").Equal([]interface{}{"Admin", "User"}),
		)
	})

	t.Run("slice of string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().JQ(".Roles").Equal([]string{"Admin", "User"}),
		)
	})

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().JQ(".Details").Equal(map[string]interface{}{"Surname": "Doe", "Email": "joe@example.com"}),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().JQ(".UserID").Equal(10),
		)
	})

	t.Run("struct", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().JQ(".Company").Equal(struct {
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
			Expect().Body().JSON().JQ(".").Equal(payload),
		)
	})

	t.Run("nil", func(t *testing.T) {
		t.Run("equal", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Expect().Body().JSON().JQ(".NotExistent").Equal(nil),
			)
		})

		t.Run("nil in expect", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Expect().Body().JSON().JQ(".UserID").Equal(nil),
				),
				PtrStr("not equal"), nil, nil, nil, nil, nil,
			)
		})

		t.Run("nil in response", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Expect().Body().JSON().JQ(".NotExistent").Equal("Hello World"),
				),
				PtrStr("not equal"), nil, nil, nil, nil, nil,
			)
		})
	})

	t.Run("multi jq", func(t *testing.T) {
		t.Run("param", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Expect().Body().JSON().JQ(".Details", ".Email").Equal("joe@example.com"),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Expect().Body().JSON().JQ(".Details", ".NotExistent").Equal("Hello World"),
				),
				PtrStr("not equal"), nil, nil, nil, nil, nil,
			)
		})
		t.Run("func", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Expect().Body().JSON().JQ(".Details").JQ(".Email").Equal("joe@example.com"),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Expect().Body().JSON().JQ(".Details").JQ(".NotExistent").Equal("Hello World"),
				),
				PtrStr("not equal"), nil, nil, nil, nil, nil,
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
			Expect().Body().JSON().JQ(`.[] | select(.role == "Admin") | .name`).Equal([]string{"Joe", "Bob"}),
		)
	})
}

func TestExpectBodyJSONJQ_NotEqual(t *testing.T) {
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
			Expect().Body().JSON().JQ(".Name").NotEqual("Alice"),
		)
	})
	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().JQ(".Roles").NotEqual([]interface{}{"Admin", "Developer"}),
		)
	})

	t.Run("slice of string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().JQ(".Roles").NotEqual([]string{"Admin", "Developer"}),
		)
	})

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().JQ(".Details").NotEqual(map[string]interface{}{"Surname": "Joe", "Email": "joe@example.com"}),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().JQ(".UserID").NotEqual(5),
		)
	})

	t.Run("struct", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().JQ(".Company").NotEqual(struct {
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
				Expect().Body().JSON().JQ(".").NotEqual(payload),
			),
			PtrStr("should not be map[string]interface {}{"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		)
	})

	t.Run("nil", func(t *testing.T) {
		t.Run("equal", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Expect().Body().JSON().JQ(".NotExistent").NotEqual(nil),
				),
				PtrStr("should not be nil"),
			)
		})

		t.Run("nil in expect", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Expect().Body().JSON().JQ(".UserID").NotEqual(nil),
			)
		})

		t.Run("nil in response", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Expect().Body().JSON().JQ(".NotExistent").NotEqual("Hello World"),
			)
		})
	})
}

func TestExpectBodyJSONJQ_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	Test(t,
		Post(s.URL),
		Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Id": 10}),
		Expect().Body().JSON().JQ(".Name").Contains("Joe"),
	)
}

func TestExpectBodyJSONJQ_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	Test(t,
		Post(s.URL),
		Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Id": 10}),
		Expect().Body().JSON().JQ(".Name").NotContains("A"),
	)
}

func TestExpectBodyJSONJQ_NoJSON(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	ExpectError(t, Do(
		Head(s.URL),
		Expect().Body().JSON().JQ(".").Equal(""),
	),
		PtrStr(`EOF`),
	)
}

func TestExpectBodyJSONJQ_Len(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String(`["Hello", "World"]`),
		Expect().Body().JSON().JQ(".").Len().Equal(2),
	)

	Test(t,
		Post(s.URL),
		Send().Body().String(`{"Name":"Joe", "ID": 10}`),
		Expect().Body().JSON().JQ(".").Len().Equal(2),
	)

	Test(t,
		Post(s.URL),
		Send().Body().String(`{"Name":"Joe", "ID": 10}`),
		Expect().Body().JSON().JQ(".Name").Len().Equal(3),
	)

	Test(t,
		Post(s.URL),
		Send().Body().String(`{"Name":"Joe", "ID": 10}`),
		Expect().Body().JSON().JQ(".Group").Len().Equal(0),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String(`"Hello World"`),
			Expect().Body().JSON().JQ(".").Len().Equal(10),
		),
		PtrStr("not equal"), PtrStr("expected: 10"), PtrStr("actual: 11"), nil, nil, nil, nil,
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String(`10`),
			Expect().Body().JSON().JQ(".").Len().Equal(10),
		),
		PtrStr("cannot get len for 10"),
	)
}
