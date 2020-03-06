package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestExpectBodyJSON_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

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
	return
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
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`["A","B","C"]`),
				Expect().Body().JSON().Equal("", []interface{}{"A", "B"}),
			),
			PtrStr("Not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		)
	})

	t.Run("slice of string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`["A","B"]`),
			Expect().Body().JSON([]string{"A", "B"}),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`["A","B","C"]`),
				Expect().Body().JSON([]string{"A", "B"}),
			),
			PtrStr("Not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		)
	})

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`{"A":"1","B":"2"}`),
			Expect().Body().JSON(map[string]interface{}{"A": "1", "B": "2"}),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`{"A":"1","B":"2","C":"3"}`),
				Expect().Body().JSON(map[string]interface{}{"A": "1", "B": "2"}),
			),
			PtrStr("Not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		)
	})

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`{"A":"1","B":"2"}`),
			Expect().Body().JSON(map[string]string{"A": "1", "B": "2"}),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`{"A":"1","B":"2","C": "3"}`),
				Expect().Body().JSON(map[string]string{"A": "1", "B": "2"}),
			),
			PtrStr("Not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`8`),
			Expect().Body().JSON(8),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`6`),
				Expect().Body().JSON(8),
			),
			PtrStr("Not equal"), nil, nil, nil, nil, nil, nil,
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
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body(`{"Name":"Joe", "ID": 11}`),
					Expect().Body().JSON(user),
				),
				PtrStr("Not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			)
		})

		// ptr
		t.Run("", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Send().Body(`{"Name":"Joe", "ID": 10}`),
				Expect().Body().JSON(&user),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body(`{"Name":"Joe", "ID": 11}`),
					Expect().Body().JSON(&user),
				),
				PtrStr("Not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
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
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body(`{"Name":"Joe", "ID": 11}`),
					Expect().Body().JSON(&puser),
				),
				PtrStr("Not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			)
		})

		t.Run("ptr in struct field", func(t *testing.T) {
			type Company struct {
				ID      int
				Name    string
				Address *string
			}
			expect := []Company{
				{
					ID:      1,
					Name:    "Wood Works",
					Address: nil,
				},
				{
					ID:      10,
					Name:    "Steel Mechanix",
					Address: PtrStr("Steel Road 1"),
				},
			}
			Test(t,
				Post(s.URL),
				Send().Body(`[{"Name":"Wood Works", "ID": 1}, {"Name":"Steel Mechanix", "ID": 10, "Address": "Steel Road 1"}]`),
				Expect().Body().JSON(expect),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body(`[{"Name":"Wood Works", "ID": 1}, {"Name":"Steel Mechanix", "ID": 10, "Address": "Steel Road 2"}]`),
					Expect().Body().JSON(expect),
				),
				PtrStr("Not equal"), nil, nil, nil,
				nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			)

			Test(t,
				Post(s.URL),
				Send().Body(`[{"Name":"Wood Works", "ID": 1}, {"Name":"Steel Mechanix", "ID": 10, "Address": "Steel Road 1"}]`),
				Expect().Body().JSON(&expect),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body(`[{"Name":"Wood Works", "ID": 1}, {"Name":"Steel Mechanix", "ID": 10, "Address": "Steel Road 2"}]`),
					Expect().Body().JSON(&expect),
				),
				PtrStr("Not equal"), nil, nil, nil,
				nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			)

			// test for modification
			require.Equal(t, expect, []Company{
				{
					ID:      1,
					Name:    "Wood Works",
					Address: nil,
				},
				{
					ID:      10,
					Name:    "Steel Mechanix",
					Address: PtrStr("Steel Road 1"),
				},
			})
		})
	})

	t.Run("nil", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(nil),
			Expect().Body().JSON().Equal("", nil),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(8),
				Expect().Body().JSON().Equal("", nil),
			),
			PtrStr("Not equal"), nil, nil, nil, nil, nil,
		)
	})
}

func TestExpectBodyJSON_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`"Hello World"`),
			Expect().Body().JSON().NotEqual("", "Hello Universe"),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`"Hello World"`),
				Expect().Body().JSON().NotEqual("", "Hello World"),
			),
			PtrStr(`should not be "Hello World"`),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`"Hello World"`),
				Expect().Custom(func(hit Hit) {
					hit.Expect().Body().JSON().NotEqual("", "Hello World")
				}),
			),
			PtrStr(`should not be "Hello World"`),
		)
	})
	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`["A","B","C"]`),
			Expect().Body().JSON().NotEqual("", []interface{}{"A", "B"}),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`["A","B"]`),
				Expect().Body().JSON().NotEqual("", []interface{}{"A", "B"}),
			),
			PtrStr("should not be []interface {}{"), nil, nil, nil,
		)
	})

	t.Run("slice of string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`["A","B"]`),
			Expect().Body().JSON().NotEqual("", []string{"A", "B", "C"}),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`["A","B"]`),
				Expect().Body().JSON().NotEqual("", []string{"A", "B"}),
			),
			PtrStr("should not be []string{"), nil, nil, nil,
		)
	})

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`{"A":"1","B":"2","C":"3"}`),
			Expect().Body().JSON().NotEqual("", map[string]interface{}{"A": "1", "B": "2"}),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`{"A":"1","B":"2"}`),
				Expect().Body().JSON().NotEqual("", map[string]interface{}{"A": "1", "B": "2"}),
			),
			PtrStr("should not be map[string]interface {}{"), nil, nil, nil,
		)
	})

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`{"A":"1","B":"2","C": "3"}`),
			Expect().Body().JSON().NotEqual("", map[string]string{"A": "1", "B": "2"}),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`{"A":"1","B":"2"}`),
				Expect().Body().JSON().NotEqual("", map[string]string{"A": "1", "B": "2"}),
			),
			PtrStr("should not be map[string]string{"), nil, nil, nil,
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`6`),
			Expect().Body().JSON().NotEqual("", 8),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`8`),
				Expect().Body().JSON().NotEqual("", 8),
			),
			PtrStr("should not be 8"),
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
				Send().Body(`{"Name":"Joe", "ID": 11}`),
				Expect().Body().JSON().NotEqual("", user),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body(`{"Name":"Joe", "ID": 10}`),
					Expect().Body().JSON().NotEqual("", user),
				),
				PtrStr("should not be struct { Name string; ID int }{"), nil, nil, nil,
			)
		})

		// ptr
		t.Run("", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Send().Body(`{"Name":"Joe", "ID": 11}`),
				Expect().Body().JSON().NotEqual("", &user),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body(`{"Name":"Joe", "ID": 10}`),
					Expect().Body().JSON().NotEqual("", &user),
				),
				PtrStr("should not be &struct { Name string; ID int }{"), nil, nil, nil,
			)
		})

		// double ptr
		puser := &user
		t.Run("", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Send().Body(`{"Name":"Joe", "ID": 11}`),
				Expect().Body().JSON().NotEqual("", &puser),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body(`{"Name":"Joe", "ID": 10}`),
					Expect().Body().JSON().NotEqual("", &puser),
				),
				PtrStr("should not be &&struct { Name string; ID int }{"), nil, nil, nil,
			)
		})

		t.Run("ptr in struct field", func(t *testing.T) {
			type Company struct {
				ID      int
				Name    string
				Address *string
			}
			expect := []Company{
				{
					ID:      1,
					Name:    "Wood Works",
					Address: nil,
				},
				{
					ID:      10,
					Name:    "Steel Mechanix",
					Address: PtrStr("Steel Road 1"),
				},
			}
			Test(t,
				Post(s.URL),
				Send().Body(`[{"Name":"Wood Works", "ID": 1}, {"Name":"Steel Mechanix", "ID": 10, "Address": "Steel Road 2"}]`),
				Expect().Body().JSON().NotEqual("", expect),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body(`[{"Name":"Wood Works", "ID": 1}, {"Name":"Steel Mechanix", "ID": 10, "Address": "Steel Road 1"}]`),
					Expect().Body().JSON().NotEqual("", expect),
				),
				PtrStr("should not be []hit_test.Company{"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			)

			Test(t,
				Post(s.URL),
				Send().Body(`[{"Name":"Wood Works", "ID": 1}, {"Name":"Steel Mechanix", "ID": 10, "Address": "Steel Road 2"}]`),
				Expect().Body().JSON().NotEqual("", &expect),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body(`[{"Name":"Wood Works", "ID": 1}, {"Name":"Steel Mechanix", "ID": 10, "Address": "Steel Road 1"}]`),
					Expect().Body().JSON().NotEqual("", &expect),
				),
				PtrStr("should not be &[]hit_test.Company{"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			)

			// test for modification
			require.Equal(t, expect, []Company{
				{
					ID:      1,
					Name:    "Wood Works",
					Address: nil,
				},
				{
					ID:      10,
					Name:    "Steel Mechanix",
					Address: PtrStr("Steel Road 1"),
				},
			})
		})
	})
	t.Run("nil", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(8),
			Expect().Body().JSON().NotEqual("", nil),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(nil),
				Expect().Body().JSON().NotEqual("", nil),
			),
			PtrStr("should not be nil"),
		)
	})
}

func TestExpectBodyJSON_EqualExpression(t *testing.T) {
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

	t.Run("full payload", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Equal("", payload),
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
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`{"Name":"Joe", "ID": 10}`),
				Expect().Body().JSON().Contains("", "Address"),
			),
			PtrStr("map[string]interface {}{"), nil, nil, PtrStr(`} does not contain "Address"`),
		)
	})

	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`[1, 2, 3]`),
			Expect().Body().JSON().Contains("", 2),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`[1, 2, 3]`),
				Expect().Body().JSON().Contains("", 4),
			),
			PtrStr("[]interface {}{"), nil, nil, nil, PtrStr("} does not contain 4"),
		)
	})

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`"Hello World"`),
			Expect().Body().JSON().Contains("", "W"),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`"Hello World"`),
				Expect().Body().JSON().Contains("", "U"),
			),
			PtrStr(`"Hello World" does not contain "U"`),
		)
	})

	t.Run("nil contains", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`null`),
			Expect().Body().JSON().Contains("", nil),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`"Hello World"`),
				Expect().Body().JSON().Contains("", nil),
			),
			PtrStr(`"Hello World" does not contain nil`),
		)
	})
}

func TestExpectBodyJSON_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`{"Name":"Joe", "ID": 10}`),
			Expect().Body().JSON().NotContains("", "Address"),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`{"Name":"Joe", "ID": 10}`),
				Expect().Body().JSON().NotContains("", "Name"),
			),
			PtrStr("map[string]interface {}{"), nil, nil, PtrStr(`} does contain "Name"`),
		)
	})

	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`[1, 2, 3]`),
			Expect().Body().JSON().NotContains("", 4),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`[1, 2, 3]`),
				Expect().Body().JSON().NotContains("", 2),
			),
			PtrStr("[]interface {}{"), nil, nil, nil, PtrStr("} does contain 2"),
		)
	})

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`"Hello World"`),
			Expect().Body().JSON().NotContains("", "U"),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`"Hello World"`),
				Expect().Body().JSON().NotContains("", "W"),
			),
			PtrStr(`"Hello World" does contain "W"`),
		)
	})

	t.Run("nil contains", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`"Hello World"`),
			Expect().Body().JSON().NotContains("", nil),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body(`null`),
				Expect().Body().JSON().NotContains("", nil),
			),
			PtrStr(`nil does contain nil`),
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
			h.Response().Body().JSON().GetAs("", &user)
			require.Equal(t, User{10, "Joe"}, user)
		}),
	)

	Test(t,
		Post(s.URL),
		Send(User{10, "Joe"}),
		Expect(func(h Hit) {
			var name string
			h.Response().Body().JSON().GetAs("Name", &name)
			require.Equal(t, "Joe", name)
		}),
	)
}
