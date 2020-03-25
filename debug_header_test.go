package hit_test

import (
	"bytes"
	"encoding/json"
	"testing"

	. "github.com/Eun/go-hit"

	"io/ioutil"

	"strings"

	"github.com/Eun/go-hit/expr"
	"github.com/lunixbochs/vtclean"
	"github.com/stretchr/testify/require"
)

func TestDebugRequestHeader(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("without expression", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Header("X-Header1", "Foo"),
			Send().Header("X-Header2", "Bar"),
			Stdout(buf),
			Debug().Request().Header(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Len(t, m, 2)
		require.Equal(t, "Foo", expr.MustGetValue(m, "X-Header1"))
		require.Equal(t, "Bar", expr.MustGetValue(m, "X-Header2"))
	})

	t.Run("with expression", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		Test(t,
			Post(s.URL),
			Send().Header("X-Header1", "Foo"),
			Send().Header("X-Header2", "Bar"),
			Stdout(buf),
			Debug().Request().Header("X-Header1"),
		)

		b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
		require.NoError(t, err)
		require.Equal(t, `"Foo"`, strings.TrimSpace(string(b)))
	})

	t.Run("clear", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		Test(t,
			Post(s.URL),
			Send().Header("X-Header1", "Foo"),
			Send().Header("X-Header2", "Bar"),
			Expect("Hello World"),
			Stdout(buf),
			Debug().Request().Header(),
			Clear().Expect(),
		)
	})
}

func TestDebugResponseHeader(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("without expression", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Header("X-Header1", "Foo"),
			Send().Header("X-Header2", "Bar"),
			Stdout(buf),
			Debug().Response().Header(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Equal(t, "Foo", expr.MustGetValue(m, "X-Header1"))
		require.Equal(t, "Bar", expr.MustGetValue(m, "X-Header2"))
	})

	t.Run("with expression", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		Test(t,
			Post(s.URL),
			Send().Header("X-Header1", "Foo"),
			Send().Header("X-Header2", "Bar"),
			Stdout(buf),
			Debug().Response().Header("X-Header1"),
		)

		b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
		require.NoError(t, err)
		require.Equal(t, `"Foo"`, strings.TrimSpace(string(b)))
	})

	t.Run("clear", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		Test(t,
			Post(s.URL),
			Send().Header("X-Header1", "Foo"),
			Send().Header("X-Header2", "Bar"),
			Expect("Hello World"),
			Stdout(buf),
			Debug().Response().Header(),
			Clear().Expect(),
		)
	})
}
