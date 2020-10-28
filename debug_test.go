package hit_test

import (
	"bytes"
	"encoding/json"
	"testing"

	. "github.com/Eun/go-hit"

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
			Send().Body().String("Hello World"),
			Fdebug(buf),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

		require.NotNil(t, m["Request"])
		require.Equal(t, "Hello World", m["Request"].(map[string]interface{})["Body"])
		require.NotNil(t, m["Response"])
	})

	t.Run("json decode", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Body().JSON([]int{1, 2, 3}),
			Fdebug(buf),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

		require.NotNil(t, m["Request"])
		require.Equal(t, []interface{}{1.0, 2.0, 3.0}, m["Request"].(map[string]interface{})["Body"])
		require.NotNil(t, m["Response"])
	})

	t.Run("debug without body", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Fdebug(buf),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

		require.NotNil(t, m["Request"])
		require.Nil(t, m["Request"].(map[string]interface{})["Body"])
		require.NotNil(t, m["Response"])
	})

	t.Run("debug with header", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Body().String("Hello World"),
			Fdebug(buf).Request(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Equal(t, "Hello World", m["Body"])
	})

	t.Run("debug request", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Body().String("Hello World"),
			Fdebug(buf).Request(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Equal(t, "Hello World", m["Body"])
	})

	t.Run("debug response", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Body().String("Hello World"),
			Fdebug(buf).Response(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Equal(t, "200 OK", m["Status"])
	})

	t.Run("debug in custom", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Body().String("Hello World"),
			Custom(BeforeExpectStep, func(hit Hit) {
				hit.MustDo(Fdebug(buf).Request())
			}),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Equal(t, "Hello World", m["Body"])
	})

	t.Run("Time", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().String().Equal("Hello World"),
			Fdebug(buf).Time(),
			Clear().Expect(),
		)
	})

	t.Run("Debug with Clear", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Body().String("Hello World"),
			Expect().Body().String().Equal("Hello World"),
			Fdebug(buf),
			Clear().Expect(),
		)
	})
}
