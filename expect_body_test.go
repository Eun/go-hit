package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestExpectBody_Equal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("bytes", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`Hello World`),
			Expect().Body([]byte("Hello World")),
		)
	})

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Body(`Hello World`),
		)
	})

	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`["A","B"]`),
			Expect().Body([]interface{}{"A", "B"}),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("8"),
			Expect().Body(8),
		)
	})
}

func TestExpectBody_NotEqual(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("bytes", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`Hello World`),
			Expect().Body().NotEqual([]byte("Hello Universe")),
		)
	})

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Body().NotEqual(`Hello Universe`),
		)
	})

	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body(`["A","B"]`),
			Expect().Body().NotEqual([]interface{}{"A", "B", "C"}),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("8"),
			Expect().Body().NotEqual(6),
		)
	})
}

func TestExpectBody_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Body().Contains(`World`),
		)
	})

	t.Run("slice", func(t *testing.T) {
		// slice goes to json
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Body().Contains([]string{"Hello World"}),
			),
			PtrStr("invalid character 'H' looking for beginning of value"),
		)
	})
}

func TestExpectBody_NotContains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Body().NotContains(`Universe`),
		)
	})

	t.Run("slice", func(t *testing.T) {
		// slice goes to json
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Body().NotContains([]string{"Hello Universe"}),
			),
			PtrStr("invalid character 'H' looking for beginning of value"),
		)
	})
}

func TestExpectBodyFinal(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("Expect().Body(value).Interface()", func(t *testing.T) {
		require.PanicsWithValue(t, "only usable with Expect().Body() not with Expect().Body(value)", func() {
			Do(Expect().Body("Data").Interface(nil))
		})
	})

	t.Run("Expect().Body(value).JSON()", func(t *testing.T) {
		require.PanicsWithValue(t, "only usable with Expect().Body() not with Expect().Body(value)", func() {
			Do(Expect().Body("Data").JSON(nil))
		})
	})

	t.Run("Expect().Body(value).Equal()", func(t *testing.T) {
		require.PanicsWithValue(t, "only usable with Expect().Body() not with Expect().Body(value)", func() {
			Do(Expect().Body("Data").Equal(nil))
		})
	})

	t.Run("Expect().Body(value).NotEqual()", func(t *testing.T) {
		require.PanicsWithValue(t, "only usable with Expect().Body() not with Expect().Body(value)", func() {
			Do(Expect().Body("Data").NotEqual(nil))
		})
	})

	t.Run("Expect().Body(value).Contains()", func(t *testing.T) {
		require.PanicsWithValue(t, "only usable with Expect().Body() not with Expect().Body(value)", func() {
			Do(Expect().Body("Data").Contains(nil))
		})
	})

	t.Run("Expect().Body(value).NotContains()", func(t *testing.T) {
		require.PanicsWithValue(t, "only usable with Expect().Body() not with Expect().Body(value)", func() {
			Do(Expect().Body("Data").NotContains(nil))
		})
	})
}

func TestExpectBody_WithoutArgument(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Body(),
		),
		PtrStr("unable to run Expect().Body() without an argument or without a chain. Please use Expect().Body(something) or Expect().Body().Something"),
	)
}
