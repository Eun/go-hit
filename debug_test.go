package hit_test

import (
	"bytes"
	"encoding/json"
	"testing"

	. "github.com/Eun/go-hit"

	"github.com/Eun/go-hit/expr"
	"github.com/lunixbochs/vtclean"
	"github.com/stretchr/testify/require"
)

func TestDebug(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("no json decode", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send().Body("Hello World"),
			Debug(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

		require.NotNil(t, expr.MustGetValue(m, "Request"))
		require.Equal(t, "Hello World", expr.MustGetValue(m, "Request.Body"))
		require.NotNil(t, expr.MustGetValue(m, "Response"))
	})

	t.Run("json decode", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send().Body().JSON([]int{1, 2, 3}),
			Debug(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

		require.NotNil(t, expr.MustGetValue(m, "Request"))
		require.Equal(t, []interface{}{1.0, 2.0, 3.0}, expr.MustGetValue(m, "Request.Body"))
		require.NotNil(t, expr.MustGetValue(m, "Response"))
	})

	t.Run("debug without body", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Stdout(buf),
			Debug(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

		require.NotNil(t, expr.MustGetValue(m, "Request"))
		require.Equal(t, "", expr.MustGetValue(m, "Request.Body"))
		require.NotNil(t, expr.MustGetValue(m, "Response"))
	})

	t.Run("debug with expression", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send().Body("Hello World"),
			Debug("Request"),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Equal(t, "Hello World", expr.MustGetValue(m, "Body"))
	})

	t.Run("debug request", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send().Body("Hello World"),
			Debug().Request(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Equal(t, "Hello World", expr.MustGetValue(m, "Body"))
	})

	t.Run("debug request", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send().Body("Hello World"),
			Debug().Response(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Equal(t, "200 OK", expr.MustGetValue(m, "Status"))
	})

	t.Run("debug in custom", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send().Body("Hello World"),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.MustDo(Debug().Request())
			}),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Equal(t, "Hello World", expr.MustGetValue(m, "Body"))
	})

	t.Run("Time", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send().Body("Hello World"),
			Expect().Body("Hello World"),
			Debug().Time(),
			Clear().Expect(),
		)
	})

	t.Run("Debug with Clear", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send().Body("Hello World"),
			Expect().Body("Hello World"),
			Debug(),
			Clear().Expect(),
		)
	})
}
