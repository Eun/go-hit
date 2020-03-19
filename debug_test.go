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

func Test_Debug(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("no json decode", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send("Hello World"),
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
			Send([]int{1, 2, 3}),
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
		require.Nil(t, expr.MustGetValue(m, "Request.Body"))
		require.NotNil(t, expr.MustGetValue(m, "Response"))
	})

	t.Run("debug with expression", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send("Hello World"),
			Debug("Request"),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Equal(t, "Hello World", expr.MustGetValue(m, "Body"))
	})

	t.Run("debug in custom", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		// send garbage so Debugs getBody function cannot parse it as json
		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send("Hello World"),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.MustDo(Debug("Request"))
			}),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Equal(t, "Hello World", expr.MustGetValue(m, "Body"))
	})

	t.Run("Debug with Clear", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		// send garbage so Debugs getBody function cannot parse it as json
		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send("Hello World"),
			Expect("Hello World"),
			Debug(),
			Clear().Expect(),
		)
	})
}
