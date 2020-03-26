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

func TestDebugRequest_Method(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Stdout(buf),
		Debug().Request().Method(),
	)

	b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
	require.NoError(t, err)
	require.Equal(t, `POST`, strings.TrimSpace(string(b)))
}

func TestDebugRequest_URL(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Stdout(buf),
		Debug().Request().URL("Host"),
	)

	b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
	require.NoError(t, err)
	require.Equal(t, s.Listener.Addr().String(), strings.TrimSpace(string(b)))
}

func TestDebugRequest_Proto(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Stdout(buf),
		Debug().Request().Proto(),
	)

	b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
	require.NoError(t, err)
	require.Equal(t, "HTTP/1.1", strings.TrimSpace(string(b)))
}

func TestDebugRequest_ProtoMajor(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Stdout(buf),
		Debug().Request().ProtoMajor(),
	)

	b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
	require.NoError(t, err)
	require.Equal(t, "1", strings.TrimSpace(string(b)))
}

func TestDebugRequest_ProtoMinor(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Stdout(buf),
		Debug().Request().ProtoMinor(),
	)

	b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
	require.NoError(t, err)
	require.Equal(t, "1", strings.TrimSpace(string(b)))
}

func TestDebugRequest_ContentLength(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Stdout(buf),
		Send("Hello World"),
		Debug().Request().ContentLength(),
	)

	b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
	require.NoError(t, err)
	require.Equal(t, "11", strings.TrimSpace(string(b)))
}

func TestDebugRequest_TransferEncoding(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Stdout(buf),
		Send("Hello World"),
		Debug().Request().TransferEncoding(),
	)

	b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
	require.NoError(t, err)
	require.Equal(t, "[]", strings.TrimSpace(string(b)))
}

func TestDebugRequest_Host(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Stdout(buf),
		Send("Hello World"),
		Send().Header("Host", "example.com"),
		Debug().Request().Host(),
	)

	b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
	require.NoError(t, err)
	require.Equal(t, "example.com", strings.TrimSpace(string(b)))
}

func TestDebugRequest_Body(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("no json decode", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send("Hello World"),
			Debug().Request().Body(),
		)

		b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
		require.NoError(t, err)
		require.Equal(t, `Hello World`, strings.TrimSpace(string(b)))
	})

	t.Run("json decode", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Stdout(buf),
			Send([]int{1, 2, 3}),
			Debug().Request().Body(),
		)

		var m []interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

		require.Equal(t, []interface{}{1.0, 2.0, 3.0}, expr.MustGetValue(m, "."))
	})
}

func TestDebugRequest_Header(t *testing.T) {
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
			Send().Body("Hello"),
			Send().Header("X-Header1", "Foo"),
			Send().Header("X-Header2", "Bar"),
			Stdout(buf),
			Debug().Request().Header("X-Header1"),
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
			Expect("Hello World"),
			Stdout(buf),
			Debug().Request().Header(),
			Clear().Expect(),
		)
	})
}

func TestDebugRequest_Trailer(t *testing.T) {
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
			Debug().Request().Trailer(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Len(t, m, 2)
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
			Debug().Request().Trailer("X-Trailer1"),
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
			Expect("Hello World"),
			Stdout(buf),
			Debug().Request().Trailer(),
			Clear().Expect(),
		)
	})
}
