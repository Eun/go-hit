package hit_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/Eun/go-hit"
)

func TestSendBody_String(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String("Hello World"),
		Expect().Body().String().Equal("Hello World"),
	)
}
func TestSendBody_JSON(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("string", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().JSON("Hello World"),
			Expect().Body().String().Equal(`"Hello World"`),
		)
	})
	t.Run("slice", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().JSON([]string{"A", "B"}),
			Expect().Body().String().Equal(`["A","B"]`),
		)
	})

	t.Run("object", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().JSON(map[string]interface{}{"A": "1", "B": "2"}),
			Expect().Body().String().Equal(`{"A":"1","B":"2"}`),
		)
	})

	t.Run("int", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body().JSON(8),
			Expect().Body().String().Equal(`8`),
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

		Test(t,
			Post(s.URL),
			Send().Body().JSON(user),
			Expect().Body().String().Equal(`{"Name":"Joe","ID":10}`),
		)
	})
}

func TestSendBody_XML(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().XML([]string{"A", "B"}),
		Expect().Body().String().Equal(`<string>A</string><string>B</string>`),
	)
}

func TestSendBody_ModifyPreviousBody(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().String("Hello"),
		Send().Custom(func(hit Hit) {
			hit.Request().Body().SetString(fmt.Sprintf("%s World", hit.Request().Body().MustString()))
		}),
		Expect().Body().String().Equal("Hello World"),
	)
}

func TestSendBody_EmptyBody(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Get(s.URL),
		Send().Custom(func(hit Hit) {
			require.Empty(t, hit.Request().Body().MustString())
		}),
	)
}

func TestSendBody_FormValue(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body().FormValues("username").Add("joe"),
		Send().Body().FormValues("password").Add("secret"),
		Expect().Body().String().OneOf("password=secret&username=joe", "username=joe&password=secret"),
	)
}
