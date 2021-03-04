package hit_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/Eun/go-hit"
)

func TestExpectBodyJSON_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`"Hello World"`),
			Expect().Body().JSON().Equal("Hello World"),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`"Hello Universe"`),
				Expect().Body().JSON().Equal("Hello World"),
			),
			PtrStr("not equal"), nil, nil, nil, nil, nil, nil,
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`"Hello Universe"`),
				Expect().Custom(func(hit Hit) error {
					hit.MustDo(Expect().Body().JSON().Equal("Hello World"))
					return nil
				}),
			),
			PtrStr("not equal"), nil, nil, nil, nil, nil, nil,
		)
	})
	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`["A","B"]`),
			Expect().Body().JSON().Equal([]interface{}{"A", "B"}),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`["A","B","C"]`),
				Expect().Body().JSON().Equal([]interface{}{"A", "B"}),
			),
			PtrStr("not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`"Hello World"`),
				Expect().Body().JSON().Equal([]interface{}{"Hello World"}),
			),
			PtrStr("not equal"), nil, nil, nil, nil, nil, nil, nil, nil,
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`["A","B","C"]`),
				Expect().Body().JSON().Equal("Hello World"),
			),
			PtrStr("not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		)
	})

	t.Run("slice of string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`["A","B"]`),
			Expect().Body().JSON().Equal([]string{"A", "B"}),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`["A","B","C"]`),
				Expect().Body().JSON().Equal([]string{"A", "B"}),
			),
			PtrStr("not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		)
	})

	t.Run("map interface", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`{"A":"1","B":"2"}`),
			Expect().Body().JSON().Equal(map[string]interface{}{"A": "1", "B": "2"}),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`{"A":"1","B":"2","C":"3"}`),
				Expect().Body().JSON().Equal(map[string]interface{}{"A": "1", "B": "2"}),
			),
			PtrStr("not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		)
	})

	t.Run("map string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`{"A":"1","B":"2"}`),
			Expect().Body().JSON().Equal(map[string]string{"A": "1", "B": "2"}),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`{"A":"1","B":"2","C": "3"}`),
				Expect().Body().JSON().Equal(map[string]string{"A": "1", "B": "2"}),
			),
			PtrStr("not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String("8"),
			Expect().Body().JSON().Equal(8),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String("6"),
				Expect().Body().JSON().Equal(8),
			),
			PtrStr("not equal"), nil, nil, nil, nil, nil, nil,
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
				Send().Body().String(`{"Name":"Joe", "ID": 10}`),
				Expect().Body().JSON().Equal(user),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body().String(`{"Name":"Joe", "ID": 11}`),
					Expect().Body().JSON().Equal(user),
				),
				PtrStr("not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			)
		})

		t.Run("not all fields", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body().String(`{"Name":"Joe"}`),
					Expect().Body().JSON().Equal(user),
				),
				PtrStr("not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			)
		})

		// ptr
		t.Run("", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Send().Body().String(`{"Name":"Joe", "ID": 10}`),
				Expect().Body().JSON().Equal(&user),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body().String(`{"Name":"Joe", "ID": 11}`),
					Expect().Body().JSON().Equal(&user),
				),
				PtrStr("not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			)
		})

		// double ptr
		puser := &user
		t.Run("", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Send().Body().String(`{"Name":"Joe", "ID": 10}`),
				Expect().Body().JSON().Equal(&puser),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body().String(`{"Name":"Joe", "ID": 11}`),
					Expect().Body().JSON().Equal(&puser),
				),
				PtrStr("not equal"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
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
				Send().Body().String(`[{"Name":"Wood Works", "ID": 1}, {"Name":"Steel Mechanix", "ID": 10, "Address": "Steel Road 1"}]`),
				Expect().Body().JSON().Equal(expect),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body().String(`[{"Name":"Wood Works", "ID": 1}, {"Name":"Steel Mechanix", "ID": 10, "Address": "Steel Road 2"}]`),
					Expect().Body().JSON().Equal(expect),
				),
				PtrStr("not equal"), nil, nil, nil,
				nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			)

			Test(t,
				Post(s.URL),
				Send().Body().String(`[{"Name":"Wood Works", "ID": 1}, {"Name":"Steel Mechanix", "ID": 10, "Address": "Steel Road 1"}]`),
				Expect().Body().JSON().Equal(&expect),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body().String(`[{"Name":"Wood Works", "ID": 1}, {"Name":"Steel Mechanix", "ID": 10, "Address": "Steel Road 2"}]`),
					Expect().Body().JSON().Equal(&expect),
				),
				PtrStr("not equal"), nil, nil, nil,
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
			Send().Body().String("null"),
			Expect().Body().JSON().Equal(nil),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String("8"),
				Expect().Body().JSON().Equal(nil),
			),
			PtrStr("not equal"), nil, nil, nil, nil, nil,
		)
	})
}

func TestExpectBodyJSON_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`"Hello World"`),
			Expect().Body().JSON().NotEqual("Hello Universe"),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`"Hello World"`),
				Expect().Body().JSON().NotEqual("Hello World"),
			),
			PtrStr(`should not be "Hello World"`),
		)

		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`"Hello World"`),
				Expect().Custom(func(hit Hit) error {
					hit.MustDo(Expect().Body().JSON().NotEqual("Hello World"))
					return nil
				}),
			),
			PtrStr(`should not be "Hello World"`),
		)
	})
	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`["A","B","C"]`),
			Expect().Body().JSON().NotEqual([]interface{}{"A", "B"}),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`["A","B"]`),
				Expect().Body().JSON().NotEqual([]interface{}{"A", "B"}),
			),
			PtrStr("should not be []interface {}{"), nil, nil, nil,
		)
	})

	t.Run("slice of string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`["A","B"]`),
			Expect().Body().JSON().NotEqual([]string{"A", "B", "C"}),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`["A","B"]`),
				Expect().Body().JSON().NotEqual([]string{"A", "B"}),
			),
			PtrStr("should not be []string{"), nil, nil, nil,
		)
	})

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`{"A":"1","B":"2","C":"3"}`),
			Expect().Body().JSON().NotEqual(map[string]interface{}{"A": "1", "B": "2"}),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`{"A":"1","B":"2"}`),
				Expect().Body().JSON().NotEqual(map[string]interface{}{"A": "1", "B": "2"}),
			),
			PtrStr("should not be map[string]interface {}{"), nil, nil, nil,
		)
	})

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`{"A":"1","B":"2","C": "3"}`),
			Expect().Body().JSON().NotEqual(map[string]string{"A": "1", "B": "2"}),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`{"A":"1","B":"2"}`),
				Expect().Body().JSON().NotEqual(map[string]string{"A": "1", "B": "2"}),
			),
			PtrStr("should not be map[string]string{"), nil, nil, nil,
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String("6"),
			Expect().Body().JSON().NotEqual(8),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String("8"),
				Expect().Body().JSON().NotEqual(8),
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
				Send().Body().String(`{"Name":"Joe", "ID": 11}`),
				Expect().Body().JSON().NotEqual(user),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body().String(`{"Name":"Joe", "ID": 10}`),
					Expect().Body().JSON().NotEqual(user),
				),
				PtrStr("should not be struct { Name string; ID int }{"), nil, nil, nil,
			)
		})

		// ptr
		t.Run("", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Send().Body().String(`{"Name":"Joe", "ID": 11}`),
				Expect().Body().JSON().NotEqual(&user),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body().String(`{"Name":"Joe", "ID": 10}`),
					Expect().Body().JSON().NotEqual(&user),
				),
				PtrStr("should not be &struct { Name string; ID int }{"), nil, nil, nil,
			)
		})

		// double ptr
		puser := &user
		t.Run("", func(t *testing.T) {
			Test(t,
				Post(s.URL),
				Send().Body().String(`{"Name":"Joe", "ID": 11}`),
				Expect().Body().JSON().NotEqual(&puser),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body().String(`{"Name":"Joe", "ID": 10}`),
					Expect().Body().JSON().NotEqual(&puser),
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
				Send().Body().String(`[{"Name":"Wood Works", "ID": 1}, {"Name":"Steel Mechanix", "ID": 10, "Address": "Steel Road 2"}]`),
				Expect().Body().JSON().NotEqual(expect),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body().String(`[{"Name":"Wood Works", "ID": 1}, {"Name":"Steel Mechanix", "ID": 10, "Address": "Steel Road 1"}]`),
					Expect().Body().JSON().NotEqual(expect),
				),
				PtrStr("should not be []hit_test.Company{"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			)

			Test(t,
				Post(s.URL),
				Send().Body().String(`[{"Name":"Wood Works", "ID": 1}, {"Name":"Steel Mechanix", "ID": 10, "Address": "Steel Road 2"}]`),
				Expect().Body().JSON().NotEqual(&expect),
			)
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body().String(`[{"Name":"Wood Works", "ID": 1}, {"Name":"Steel Mechanix", "ID": 10, "Address": "Steel Road 1"}]`),
					Expect().Body().JSON().NotEqual(&expect),
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
			Send().Body().JSON(8),
			Expect().Body().JSON().NotEqual(nil),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String("null"),
				Expect().Body().JSON().NotEqual(nil),
			),
			PtrStr("should not be nil"),
		)
	})
}

func TestExpectBodyJSON_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`{"Name":"Joe", "ID": 10}`),
			Expect().Body().JSON().Contains("Name"),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`{"Name":"Joe", "ID": 10}`),
				Expect().Body().JSON().Contains("Address"),
			),
			PtrStr("map[string]interface {}{"), nil, nil, PtrStr(`} does not contain "Address"`),
		)
	})

	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`[1, 2, 3]`),
			Expect().Body().JSON().Contains(2),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`[1, 2, 3]`),
				Expect().Body().JSON().Contains(4),
			),
			PtrStr("[]interface {}{"), nil, nil, nil, PtrStr("} does not contain 4"),
		)
	})

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`"Hello World"`),
			Expect().Body().JSON().Contains("W"),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`"Hello World"`),
				Expect().Body().JSON().Contains("U"),
			),
			PtrStr(`"Hello World" does not contain "U"`),
		)
	})

	t.Run("nil contains", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String("null"),
			Expect().Body().JSON().Contains(nil),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`"Hello World"`),
				Expect().Body().JSON().Contains(nil),
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
			Send().Body().String(`{"Name":"Joe", "ID": 10}`),
			Expect().Body().JSON().NotContains("Address"),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`{"Name":"Joe", "ID": 10}`),
				Expect().Body().JSON().NotContains("Name"),
			),
			PtrStr("map[string]interface {}{"), nil, nil, PtrStr(`} should not contain "Name"`),
		)
	})

	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`[1, 2, 3]`),
			Expect().Body().JSON().NotContains(4),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`[1, 2, 3]`),
				Expect().Body().JSON().NotContains(2),
			),
			PtrStr("[]interface {}{"), nil, nil, nil, PtrStr("} should not contain 2"),
		)
	})

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`"Hello World"`),
			Expect().Body().JSON().NotContains("U"),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String(`"Hello World"`),
				Expect().Body().JSON().NotContains("W"),
			),
			PtrStr(`"Hello World" should not contain "W"`),
		)
	})

	t.Run("nil contains", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().String(`"Hello World"`),
			Expect().Body().JSON().NotContains(nil),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().String("null"),
				Expect().Body().JSON().NotContains(nil),
			),
			PtrStr(`nil should not contain nil`),
		)
	})
}

func TestExpectBodyJSON_NilResponse(t *testing.T) {
	s := PrintJSONServer(nil)
	defer s.Close()

	t.Run("", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Equal(nil),
		)
	})

	t.Run("", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Expect().Body().JSON().Equal(nil),
		)
	})
}

func TestExpectBodyJSON_NoJSON(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	ExpectError(t, Do(
		Head(s.URL),
		Expect().Body().JSON().Equal(""),
	),
		PtrStr(`EOF`),
	)
}

func TestExpectBodyJSON_Len(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String(`["Hello", "World"]`),
		Expect().Body().JSON().Len().Equal(2),
	)

	Test(t,
		Post(s.URL),
		Send().Body().String(`{"Name":"Joe", "ID": 10}`),
		Expect().Body().JSON().Len().Equal(2),
	)

	Test(t,
		Post(s.URL),
		Send().Body().String(`null`),
		Expect().Body().JSON().Len().Equal(0),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String(`"Hello World"`),
			Expect().Body().JSON().Len().Equal(10),
		),
		PtrStr("not equal"), PtrStr("expected: 10"), PtrStr("actual: 11"), nil, nil, nil, nil,
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body().String(`10`),
			Expect().Body().JSON().Len().Equal(10),
		),
		PtrStr("cannot get len for 10"),
	)
}
