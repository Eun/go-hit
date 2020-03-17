package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestClearExpectBodyJSON_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().JSON(map[string]interface{}{"Name": "Joe"}),
				Expect().Body().JSON().Equal("Name", "Alice"),
				Clear().Expect().Body().JSON(),
				Expect().Body().JSON().Equal("Name", "Bob"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Bob"`), PtrStr(`actual: "Joe"`), nil, nil, nil, nil,
		)
	})

	t.Run("specific only first parameter", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Surname": "Doe"}),
				Expect().Body().JSON().Equal("Name", "Alice"),
				Expect().Body().JSON().Equal("Name", "Bob"),
				Expect().Body().JSON().Equal("Surname", "Hunt"),
				Clear().Expect().Body().JSON().Equal("Name"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Hunt"`), PtrStr(`actual: "Doe"`), nil, nil, nil, nil,
		)
	})

	t.Run("specific (all)", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().JSON(map[string]interface{}{"Name": "Joe"}),
				Expect().Body().JSON().Equal("Name", "Alice"),
				Expect().Body().JSON().Equal("Name", "Bob"),
				Clear().Expect().Body().JSON().Equal("Name", "Alice"),
			),
			PtrStr("Not equal"), PtrStr(`expected: "Bob"`), PtrStr(`actual: "Joe"`), nil, nil, nil, nil,
		)
	})
}

func TestClearExpectBodyJSON_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().JSON(map[string]interface{}{"Name": "Joe"}),
			Expect().Body().JSON().NotEqual("Name", "Joe"),
			Clear().Expect().Body().JSON(),
			Expect().Body().JSON().NotEqual("Name", "Alice"),
		)
	})

	t.Run("specific only first parameter", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Surname": "Doe"}),
			Expect().Body().JSON().NotEqual("Name", "Joe"),
			Expect().Body().JSON().NotEqual("Surname", "Hunt"),
			Clear().Expect().Body().JSON().NotEqual("Name"),
		)
	})

	t.Run("specific (all)", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().JSON(map[string]interface{}{"Name": "Joe"}),
			Expect().Body().JSON().NotEqual("Name", "Joe"),
			Expect().Body().JSON().NotEqual("Name", "Bob"),
			Clear().Expect().Body().JSON().NotEqual("Name", "Joe"),
		)
	})
}

func TestClearExpectBodyJSON_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().JSON(map[string]interface{}{"Name": "Joe"}),
				Expect().Body().JSON().Contains("Name", "Alice"),
				Clear().Expect().Body().JSON(),
				Expect().Body().JSON().Contains("Name", "Bob"),
			),
			PtrStr(`"Joe" does not contain "Bob"`),
		)
	})

	t.Run("specific only first parameter", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Surname": "Doe"}),
				Expect().Body().JSON().Contains("Name", "Alice"),
				Expect().Body().JSON().Contains("Name", "Bob"),
				Expect().Body().JSON().Contains("Surname", "Hunt"),
				Clear().Expect().Body().JSON().Contains("Name"),
			),
			PtrStr(`"Doe" does not contain "Hunt"`),
		)
	})

	t.Run("specific (all)", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().JSON(map[string]interface{}{"Name": "Joe"}),
				Expect().Body().JSON().Contains("Name", "Alice"),
				Expect().Body().JSON().Contains("Name", "Bob"),
				Clear().Expect().Body().JSON().Contains("Name", "Alice"),
			),
			PtrStr(`"Joe" does not contain "Bob"`),
		)
	})
}

func TestClearExpectBodyJSON_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("all", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().JSON(map[string]interface{}{"Name": "Joe"}),
			Expect().Body().JSON().NotContains("Name", "Joe"),
			Clear().Expect().Body().JSON(),
			Expect().Body().JSON().NotContains("Name", "Alice"),
		)
	})

	t.Run("specific only first parameter", func(t *testing.T) {
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Surname": "Doe"}),
				Expect().Body().JSON().NotContains("Name", "Alice"),
				Expect().Body().JSON().NotContains("Name", "Bob"),
				Expect().Body().JSON().NotContains("Surname", "Doe"),
				Clear().Expect().Body().JSON().NotContains("Name"),
			),
			PtrStr(`"Doe" does contain "Doe"`),
		)
	})

	t.Run("specific (all)", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().JSON(map[string]interface{}{"Name": "Joe"}),
			Expect().Body().JSON().NotContains("Name", "Joe"),
			Expect().Body().JSON().NotContains("Name", "Bob"),
			Clear().Expect().Body().JSON().NotContains("Name", "Joe"),
		)
	})
}

func TestClearExpectBodyJSONFinal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("Clear().Expect().Body().JSON(value).Equal()", func(t *testing.T) {
		require.PanicsWithValue(t, "only usable with Clear().Expect().Body().JSON() not with Clear().Expect().Body().JSON(value)", func() {
			Do(Clear().Expect().Body().JSON("").Equal())
		})
	})
	t.Run("Clear().Expect().Body().JSON(value).NotEqual()", func(t *testing.T) {
		require.PanicsWithValue(t, "only usable with Clear().Expect().Body().JSON() not with Clear().Expect().Body().JSON(value)", func() {
			Do(Clear().Expect().Body().JSON("").NotEqual())
		})
	})
	t.Run("Clear().Expect().Body().JSON(value).Contains()", func(t *testing.T) {
		require.PanicsWithValue(t, "only usable with Clear().Expect().Body().JSON() not with Clear().Expect().Body().JSON(value)", func() {
			Do(Clear().Expect().Body().JSON("").Contains())
		})
	})
	t.Run("Clear().Expect().Body().JSON(value).NotContains()", func(t *testing.T) {
		require.PanicsWithValue(t, "only usable with Clear().Expect().Body().JSON() not with Clear().Expect().Body().JSON(value)", func() {
			Do(Clear().Expect().Body().JSON("").NotContains())
		})
	})
}
