package hit_test

import (
	"testing"

	"github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestExpectBodyJSON_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body(`"Hello World"`).
			Expect().Body().JSON("Hello World").
			Do()
	})
	t.Run("slice", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body(`["A","B"]`).
			Expect().Body().JSON([]interface{}{"A", "B"}).
			Do()
	})

	t.Run("slice of string", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body(`["A","B"]`).
			Expect().Body().JSON([]string{"A", "B"}).
			Do()
	})

	t.Run("object", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body(`{"A":"1","B":"2"}`).
			Expect().Body().JSON(map[string]interface{}{"A": "1", "B": "2"}).
			Do()
	})

	t.Run("object", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body(`{"A":"1","B":"2"}`).
			Expect().Body().JSON(map[string]string{"A": "1", "B": "2"}).
			Do()
	})

	t.Run("int", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body(`8`).
			Expect().Body().JSON(8).
			Do()
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
			hit.Post(t, s.URL).
				Send().Body(`{"Name":"Joe", "ID": 10}`).
				Expect().Body().JSON(user).
				Do()
		})

		// ptr
		t.Run("", func(t *testing.T) {
			hit.Post(t, s.URL).
				Send().Body(`{"Name":"Joe", "ID": 10}`).
				Expect().Body().JSON(&user).
				Do()
		})

		// double ptr
		puser := &user
		t.Run("", func(t *testing.T) {
			hit.Post(t, s.URL).
				Send().Body(`{"Name":"Joe", "ID": 10}`).
				Expect().Body().JSON(&puser).
				Do()
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
		hit.Post(t, s.URL).
			Expect().Body().JSON().Equal("Name", "Joe").
			Do()
	})
	t.Run("slice", func(t *testing.T) {
		hit.Post(t, s.URL).
			Expect().Body().JSON().Equal("Roles", []interface{}{"Admin", "User"}).
			Do()
	})

	t.Run("slice of string", func(t *testing.T) {
		hit.Post(t, s.URL).
			Expect().Body().JSON().Equal("Roles", []string{"Admin", "User"}).
			Do()
	})

	t.Run("object", func(t *testing.T) {
		hit.Post(t, s.URL).
			Expect().Body().JSON().Equal("Details", map[string]interface{}{"Surname": "Doe", "Email": "joe@example.com"}).
			Do()
	})

	t.Run("int", func(t *testing.T) {
		hit.Post(t, s.URL).
			Expect().Body().JSON().Equal("UserID", 10).
			Do()
	})

	t.Run("struct", func(t *testing.T) {
		hit.Post(t, s.URL).
			Expect().Body().JSON().Equal("Company", struct {
			ID   int
			Name string
		}{
			1,
			"Wood Inc",
		}).Do()
	})

	t.Run("nil", func(t *testing.T) {
		t.Run("equal", func(t *testing.T) {
			hit.Post(t, s.URL).
				Expect().Body().JSON().Equal("NotExistent", nil).
				Do()
		})

		t.Run("nil in expect", func(t *testing.T) {
			require.Panics(t, func() {
				hit.Post(NewPanicWithMessage(t, PtrStr("Not equal"), nil, nil, nil, nil, nil), s.URL).
					Expect().Body().JSON().Equal("UserID", nil).
					Do()
			})
		})

		t.Run("nil in response", func(t *testing.T) {
			require.Panics(t, func() {
				hit.Post(NewPanicWithMessage(t, PtrStr("Not equal"), nil, nil, nil, nil, nil), s.URL).
					Expect().Body().JSON().Equal("NotExistent", "Hello World").
					Do()
			})
		})
	})
}

func TestExpectBodyJSON_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("object", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body(`{"Name":"Joe", "ID": 10}`).
			Expect().Body().JSON().Contains("", "Name").
			Do()
	})

	t.Run("slice", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body(`[1, 2, 3]`).
			Expect().Body().JSON().Contains("", 2).
			Do()
	})

	t.Run("string", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body(`"Hello World"`).
			Expect().Body().JSON().Contains("", "W").
			Do()
	})

	t.Run("not contains", func(t *testing.T) {
		require.Panics(t, func() {
			hit.Post(NewPanicWithMessage(t, PtrStr(`"Hello World" does not contain "Bye"`)), s.URL).
				Send().Body(`"Hello World"`).
				Expect().Body().JSON().Contains("", "Bye").
				Do()
		})
	})
	t.Run("nil contains", func(t *testing.T) {
		hit.Post(t, s.URL).
			Send().Body(`null`).
			Expect().Body().JSON().Contains("", nil).
			Do()
	})
}

func TestExpectBodyJSON_NilResponse(t *testing.T) {
	s := PrintJSONServer(nil)
	defer s.Close()

	t.Run("", func(t *testing.T) {
		hit.Post(t, s.URL).
			Expect().Body().JSON(nil).
			Do()
	})

	t.Run("", func(t *testing.T) {
		hit.Post(t, s.URL).
			Expect().Body().JSON().Equal("", nil).
			Do()
	})
}

func TestExpectBodyJSON_NoJSON(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	require.Panics(t, func() {
		hit.Head(NewPanicWithMessage(t, PtrStr(`EOF`)), s.URL).
			Expect().Body().JSON().Equal("", "").
			Do()
	})
}

func TestExpectBodyJSON_GetAs(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	type User struct {
		ID   int
		Name string
	}

	hit.Post(t, s.URL).
		Send(User{10, "Joe"}).
		Expect(func(h hit.Hit) {
			var user User
			h.Response().Body().JSON().GetAs(&user)
			require.Equal(t, User{10, "Joe"}, user)
		}).
		Do()
}
