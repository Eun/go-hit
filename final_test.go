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

	t.Run("Expect(value).Body(value)", func(t *testing.T) {
		ExpectError(t, Do(Post(s.URL), Send("Hello"), Expect().Body("Bye")), PtrStr("Not equal"), nil, nil, nil, nil, nil, nil)
		require.PanicsWithValue(t, "only usable with Expect() not with Expect(value)", func() {
			Do(Expect("Data").Body("Data"))
		})
	})

	t.Run("Expect(value).Body(value)", func(t *testing.T) {
		ExpectError(t, Do(Post(s.URL), Send("Hello"), Expect().Body().Equal("Bye")), PtrStr("Not equal"), nil, nil, nil, nil, nil, nil)
		require.PanicsWithValue(t, "only usable with Expect() not with Expect(value)", func() {
			Do(Expect("Data").Body("Data").Equal("Data"))
		})
	})

	// t.Run("Expect(value).Interface(value)", func(t *testing.T) {
	// 	ExpectError(t, Do(Post(s.URL), Send("Hello"), Expect().Interface("Bye")), PtrStr("Not equal"), nil, nil, nil, nil, nil, nil)
	// 	ExpectError(t, Do(Expect("Data").Interface("Data")), PtrStr("only usable with Expect() not with Expect(value)"))
	// })
}

// Body(data ...interface{}) IExpectBody
// Interface(interface{}) IStep
// Custom(f Callback) IStep
// Headers() *IExpectHeaders
// Header(name string) *expectSpecificHeader
// Status(code ...int) *expectStatus
// Clear() IStep
