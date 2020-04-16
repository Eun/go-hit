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

func TestDebugResponse_Proto(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Stdout(buf),
		Debug().Response().Proto(),
	)

	b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
	require.NoError(t, err)
	require.Equal(t, "HTTP/1.1", strings.TrimSpace(string(b)))
}

func TestDebugResponse_ProtoMajor(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Stdout(buf),
		Debug().Response().ProtoMajor(),
	)

	b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
	require.NoError(t, err)
	require.Equal(t, "1", strings.TrimSpace(string(b)))
}

func TestDebugResponse_ProtoMinor(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Stdout(buf),
		Debug().Response().ProtoMinor(),
	)

	b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
	require.NoError(t, err)
	require.Equal(t, "1", strings.TrimSpace(string(b)))
}

func TestDebugResponse_ContentLength(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Stdout(buf),
		Send().Body("Hello World"),
		Debug().Response().ContentLength(),
	)

	b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
	require.NoError(t, err)
	require.Equal(t, "11", strings.TrimSpace(string(b)))
}

func TestDebugResponse_TransferEncoding(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Stdout(buf),
		Send().Body("Hello World"),
		Debug().Response().TransferEncoding(),
	)

	b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
	require.NoError(t, err)
	require.Equal(t, "[]", strings.TrimSpace(string(b)))
}

func TestDebugResponse_Status(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Stdout(buf),
		Send().Body("Hello World"),
		Debug().Response().Status(),
	)

	b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
	require.NoError(t, err)
	require.Equal(t, "200 OK", strings.TrimSpace(string(b)))
}

func TestDebugResponse_StatusCode(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Stdout(buf),
		Send().Body("Hello World"),
		Debug().Response().StatusCode(),
	)

	b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
	require.NoError(t, err)
	require.Equal(t, "200", strings.TrimSpace(string(b)))
}

func TestDebugResponse_Body(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("no json decode", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send().Body("Hello World"),
			Debug().Response().Body(),
		)

		b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
		require.NoError(t, err)
		require.Equal(t, "Hello World", strings.TrimSpace(string(b)))
	})

	t.Run("json decode", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send().Body().JSON([]int{1, 2, 3}),
			Debug().Response().Body(),
		)

		var m []interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

		require.Equal(t, []interface{}{1.0, 2.0, 3.0}, expr.MustGetValue(m, "."))
	})
}

func TestDebugResponse_Header(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("without expression", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Body("Hello"),
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
			Send().Body("Hello"),
			Send().Header("X-Header1", "Foo"),
			Send().Header("X-Header2", "Bar"),
			Stdout(buf),
			Debug().Response().Header("X-Header1"),
		)

		b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
		require.NoError(t, err)
		require.Equal(t, "Foo", strings.TrimSpace(string(b)))
	})

	t.Run("clear", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		Test(t,
			Post(s.URL),
			Send().Body("Hello"),
			Send().Header("X-Header1", "Foo"),
			Send().Header("X-Header2", "Bar"),
			Expect().Body("Hello World"),
			Stdout(buf),
			Debug().Response().Header(),
			Clear().Expect(),
		)
	})
}

func TestDebugResponse_Trailer(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("without expression", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Body("Hello"),
			Send().Trailer("X-Trailer1", "Foo"),
			Send().Trailer("X-Trailer2", "Bar"),
			Stdout(buf),
			Debug().Response().Trailer(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Equal(t, "Foo", expr.MustGetValue(m, "X-Trailer1"))
		require.Equal(t, "Bar", expr.MustGetValue(m, "X-Trailer2"))
	})

	t.Run("with expression", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		Test(t,
			Post(s.URL),
			Send().Body("Hello"),
			Send().Trailer("X-Trailer1", "Foo"),
			Send().Trailer("X-Trailer2", "Bar"),
			Stdout(buf),
			Debug().Response().Trailer("X-Trailer1"),
		)

		b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
		require.NoError(t, err)
		require.Equal(t, "Foo", strings.TrimSpace(string(b)))
	})

	t.Run("clear", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		Test(t,
			Post(s.URL),
			Send().Body("Hello"),
			Send().Trailer("X-Trailer1", "Foo"),
			Send().Trailer("X-Trailer2", "Bar"),
			Expect().Body("Hello World"),
			Stdout(buf),
			Debug().Response().Trailer(),
			Clear().Expect(),
		)
	})
}
