package hit_test

import (
	"bytes"
	"encoding/json"
	"testing"

	. "github.com/Eun/go-hit"

	"io/ioutil"

	"strings"

	"github.com/lunixbochs/vtclean"
	"github.com/stretchr/testify/require"
)

func TestDebugRequestHeader(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("without header", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Headers("X-Header1").Add("Foo"),
			Send().Headers("X-Header2").Add("Bar"),
			Fdebug(buf).Request().Headers(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Len(t, m, 2)
		require.Equal(t, "Foo", m["X-Header1"])
		require.Equal(t, "Bar", m["X-Header2"])
	})

	t.Run("with header", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		Test(t,
			Post(s.URL),
			Send().Headers("X-Header1").Add("Foo"),
			Send().Headers("X-Header2").Add("Bar"),
			Fdebug(buf).Request().Headers("X-Header1"),
		)

		b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
		require.NoError(t, err)
		require.Equal(t, "Foo", strings.TrimSpace(string(b)))
	})

	t.Run("clear", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		Test(t,
			Post(s.URL),
			Send().Headers("X-Header1").Add("Foo"),
			Send().Headers("X-Header2").Add("Bar"),
			Expect().Body().String().Equal("Hello World"),
			Fdebug(buf).Request().Headers(),
			Clear().Expect(),
		)
	})
}

func TestDebugResponseHeader(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("without header", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Headers("X-Header1").Add("Foo"),
			Send().Headers("X-Header2").Add("Bar"),
			Fdebug(buf).Response().Headers(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Equal(t, "Foo", m["X-Header1"])
		require.Equal(t, "Bar", m["X-Header2"])
	})

	t.Run("with header", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		Test(t,
			Post(s.URL),
			Send().Headers("X-Header1").Add("Foo"),
			Send().Headers("X-Header2").Add("Bar"),
			Fdebug(buf).Response().Headers("X-Header1"),
		)

		b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
		require.NoError(t, err)
		require.Equal(t, "Foo", strings.TrimSpace(string(b)))
	})

	t.Run("clear", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		Test(t,
			Post(s.URL),
			Send().Headers("X-Header1").Add("Foo"),
			Send().Headers("X-Header2").Add("Bar"),
			Expect().Body().String().Equal("Hello World"),
			Fdebug(buf).Response().Headers(),
			Clear().Expect(),
		)
	})
}
