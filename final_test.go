package hit_test

import (
	"testing"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestExpectFinal(t *testing.T) {
	s := EchoServer()
	defer s.Close()
	t.Run("Expect(value).Body()", func(t *testing.T) {
		ExpectError(t, Do(Post(s.URL), Send("Hello"), Expect("Bye")), PtrStr("Not equal"), nil, nil, nil, nil, nil, nil)
		require.PanicsWithValue(t, "only usable with Expect() not with Expect(value)", func() {
			Do(Expect("Data").Body())
		})
	})

	t.Run("Expect().Body(value).JSON(value)", func(t *testing.T) {
		ExpectError(t, Do(Post(s.URL), Send(`"Hello"`), Expect().Body().JSON("Bye")), PtrStr("Not equal"), nil, nil, nil, nil, nil, nil)
		require.PanicsWithValue(t, "only usable with Expect().Body() not with Expect().Body(value)", func() {
			Do(Expect().Body("Data").JSON("Data"))
		})
	})

	t.Run("Expect().Body(value).Equal(value)", func(t *testing.T) {
		ExpectError(t, Do(Post(s.URL), Send("Hello"), Expect().Body().Equal("Bye")), PtrStr("Not equal"), nil, nil, nil, nil, nil, nil)
		require.PanicsWithValue(t, "only usable with Expect().Body() not with Expect().Body(value)", func() {
			Do(Expect().Body("Data").Equal("Data"))
		})
	})

	t.Run("Expect().Body(value).NotEqual(value)", func(t *testing.T) {
		ExpectError(t, Do(Post(s.URL), Send("Hello"), Expect().Body().NotEqual("Hello")), PtrStr(`should not be "Hello"`))
		require.PanicsWithValue(t, "only usable with Expect().Body() not with Expect().Body(value)", func() {
			Do(Expect().Body("Data").NotEqual("Data"))
		})
	})

	t.Run("Expect().Body(value).Contains(value)", func(t *testing.T) {
		ExpectError(t, Do(Post(s.URL), Send("Hello"), Expect().Body().Contains("Bye")), PtrStr(`"Hello" does not contain "Bye"`))
		require.PanicsWithValue(t, "only usable with Expect().Body() not with Expect().Body(value)", func() {
			Do(Expect().Body("Data").Contains("Data"))
		})
	})

	t.Run("Expect().Body(value).NotContains(value)", func(t *testing.T) {
		ExpectError(t, Do(Post(s.URL), Send("Hello"), Expect().Body().NotContains("Hello")), PtrStr(`"Hello" does contain "Hello"`))
		require.PanicsWithValue(t, "only usable with Expect().Body() not with Expect().Body(value)", func() {
			Do(Expect().Body("Data").NotContains("Data"))
		})
	})
}
