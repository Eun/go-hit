package hit_test

import (
	"bytes"
	"encoding/json"
	"testing"

	. "github.com/otto-eng/go-hit"

	"io/ioutil"

	"strings"

	"net/url"

	"github.com/lunixbochs/vtclean"
	"github.com/stretchr/testify/require"
)

func TestDebugRequest_Method(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Fdebug(buf).Request().Method(),
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
		Fdebug(buf).Request().URL(),
	)

	b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
	require.NoError(t, err)

	var u url.URL
	require.NoError(t, json.NewDecoder(bytes.NewReader(b)).Decode(&u))

	require.Equal(t, s.Listener.Addr().String(), u.Host)
}

func TestDebugRequest_Proto(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Fdebug(buf).Request().Proto(),
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
		Fdebug(buf).Request().ProtoMajor(),
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
		Fdebug(buf).Request().ProtoMinor(),
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
		Send().Body().String("Hello World"),
		Fdebug(buf).Request().ContentLength(),
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
		Send().Body().String("Hello World"),
		Fdebug(buf).Request().TransferEncoding(),
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
		Send().Body().String("Hello World"),
		Send().Headers("Host").Add("example.com"),
		Fdebug(buf).Request().Host(),
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
			Send().Body().String("Hello World"),
			Fdebug(buf).Request().Body(),
		)

		b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
		require.NoError(t, err)
		require.Equal(t, `Hello World`, strings.TrimSpace(string(b)))
	})

	t.Run("json decode", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Body().JSON([]int{1, 2, 3}),
			Fdebug(buf).Request().Body(),
		)

		var m []interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

		require.Equal(t, []interface{}{1.0, 2.0, 3.0}, m)
	})
}

func TestDebugRequest_BodyJSON(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Send().Body().JSON([]int{1, 2, 3}),
		Fdebug(buf).Request().Body().JSON(),
	)

	var m []interface{}
	require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

	require.Equal(t, []interface{}{1.0, 2.0, 3.0}, m)
}

func TestDebugRequest_Header(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("without header", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Body().String("Hello"),
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
			Send().Body().String("Hello"),
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
			Send().Body().String("Hello"),
			Send().Headers("X-Header1").Add("Foo"),
			Send().Headers("X-Header2").Add("Bar"),
			Expect().Body().String().Equal("Hello World"),
			Fdebug(buf).Request().Headers(),
			Clear().Expect(),
		)
	})
}

func TestDebugRequest_Trailer(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("without header", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Body().String("Hello"),
			Send().Trailers("X-Trailer1").Add("Foo"),
			Send().Trailers("X-Trailer2").Add("Bar"),
			Fdebug(buf).Request().Trailers(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
		require.Len(t, m, 2)
		require.Equal(t, "Foo", m["X-Trailer1"])
		require.Equal(t, "Bar", m["X-Trailer2"])
	})

	t.Run("with header", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		Test(t,
			Post(s.URL),
			Send().Body().String("Hello"),
			Send().Trailers("X-Trailer1").Add("Foo"),
			Send().Trailers("X-Trailer2").Add("Bar"),
			Fdebug(buf).Request().Trailers("X-Trailer1"),
		)

		b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
		require.NoError(t, err)
		require.Equal(t, "Foo", strings.TrimSpace(string(b)))
	})

	t.Run("clear", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		Test(t,
			Post(s.URL),
			Send().Body().String("Hello"),
			Send().Trailers("X-Trailer1").Add("Foo"),
			Send().Trailers("X-Trailer2").Add("Bar"),
			Expect().Body().String().Equal("Hello World"),
			Fdebug(buf).Request().Trailers(),
			Clear().Expect(),
		)
	})
}
