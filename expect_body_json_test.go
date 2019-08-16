package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestExpectBodyJSON_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`"Hello World"`),
			Expect().Body().JSON("Hello World"),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`"Hello Universe"`),
				Expect().Body().JSON("Hello World"),
			),
			PtrStr("Not equal"), nil, nil, nil, nil, nil, nil,
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`"Hello Universe"`),
				Expect().Custom(func(hit Hit) {
					hit.Expect().Body().JSON("Hello World")
				}),
			),
			PtrStr("Not equal"), nil, nil, nil, nil, nil, nil,
		)
	})
	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`["A","B"]`),
			Expect().Body().JSON([]interface{}{"A", "B"}),
		)
	})

	t.Run("slice of string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`["A","B"]`),
			Expect().Body().JSON([]string{"A", "B"}),
		)
	})

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`{"A":"1","B":"2"}`),
			Expect().Body().JSON(map[string]interface{}{"A": "1", "B": "2"}),
		)
	})

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`{"A":"1","B":"2"}`),
			Expect().Body().JSON(map[string]string{"A": "1", "B": "2"}),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`8`),
			Expect().Body().JSON(8),
		)
	})

	t.Run("struct", func(t *testing.T) {
		var user = struct {
			Name string
			ID   int
		}{
			"Joe",
			10,
		}

		t.Run("", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Send().Body(`{"Name":"Joe", "ID": 10}`),
				Expect().Body().JSON(user),
			)
		})

		// ptr
		t.Run("", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Send().Body(`{"Name":"Joe", "ID": 10}`),
				Expect().Body().JSON(&user),
			)
		})

		// double ptr
		puser := &user
		t.Run("", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Send().Body(`{"Name":"Joe", "ID": 10}`),
				Expect().Body().JSON(&puser),
			)
		})
	})
}

func TestExpectBodyJSON_EqualExpression(t *testing.T) {
	s := PrintJSONServer(map[string]interface{}{
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
	})
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Equal("Name", "Joe"),
		)
	})
	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Equal("Roles", []interface{}{"Admin", "User"}),
		)
	})

	t.Run("slice of string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Equal("Roles", []string{"Admin", "User"}),
		)
	})

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Equal("Details", map[string]interface{}{"Surname": "Doe", "Email": "joe@example.com"}),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Equal("UserID", 10),
		)
	})

	t.Run("struct", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Equal("Company", struct {
				ID   int
				Name string
			}{
				1,
				"Wood Inc",
			}),
		)
	})

	t.Run("nil", func(t *testing.T) {
		t.Run("equal", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Expect().Body().JSON().Equal("NotExistent", nil),
			)
		})

		t.Run("nil in expect", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Expect().Body().JSON().Equal("UserID", nil),
				),
				PtrStr("Not equal"), nil, nil, nil, nil, nil,
			)
		})

		t.Run("nil in response", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Expect().Body().JSON().Equal("NotExistent", "Hello World"),
				),
				PtrStr("Not equal"), nil, nil, nil, nil, nil,
			)
		})
	})
}

func TestExpectBodyJSON_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`{"Name":"Joe", "ID": 10}`),
			Expect().Body().JSON().Contains("", "Name"),
		)
	})

	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`[1, 2, 3]`),
			Expect().Body().JSON().Contains("", 2),
		)
	})

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`"Hello World"`),
			Expect().Body().JSON().Contains("", "W"),
		)
	})

	t.Run("not contains", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`"Hello World"`),
				Expect().Body().JSON().Contains("", "Bye"),
			),
			PtrStr(`"Hello World" does not contain "Bye"`),
		)
	})
	t.Run("nil contains", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`null`),
			Expect().Body().JSON().Contains("", nil),
		)
	})
}

func TestExpectBodyJSON_NilResponse(t *testing.T) {
	s := PrintJSONServer(nil)
	defer s.Close()

	t.Run("", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON(nil),
		)
	})

	t.Run("", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Equal("", nil),
		)
	})
}

func TestExpectBodyJSON_NoJSON(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	ExpectError(t, Do(
		Head(s.URL),
		Expect().Body().JSON().Equal("", ""),
	),
		PtrStr(`EOF`),
	)
}

func TestExpectBodyJSON_GetAs(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	type User struct {
		ID   int
		Name string
	}

	Test(t,
		Post(s.URL),
		Send(User{10, "Joe"}),
		Expect(func(h Hit) {
			var user User
			h.Response().Body().JSON().GetAs(&user)
			require.Equal(t, User{10, "Joe"}, user)
		}),
	)
}
