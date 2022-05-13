package hit_test

import (
	"bytes"
	"encoding/json"
	"testing"

	. "github.com/otto-eng/go-hit"

	"io/ioutil"

	"strings"

	"github.com/lunixbochs/vtclean"
	"github.com/stretchr/testify/require"
)

func TestDebugResponse_Proto(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	buf := bytes.NewBuffer(nil)

	Test(t,
		Post(s.URL),
		Fdebug(buf).Response().Proto(),
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
		Fdebug(buf).Response().ProtoMajor(),
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
		Fdebug(buf).Response().ProtoMinor(),
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
		Send().Body().String("Hello World"),
		Fdebug(buf).Response().ContentLength(),
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
		Send().Body().String("Hello World"),
		Fdebug(buf).Response().TransferEncoding(),
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
		Send().Body().String("Hello World"),
		Fdebug(buf).Response().Status(),
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
		Send().Body().String("Hello World"),
		Fdebug(buf).Response().StatusCode(),
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
			Send().Body().String("Hello World"),
			Fdebug(buf).Response().Body(),
		)

		b, err := ioutil.ReadAll(vtclean.NewReader(buf, false))
		require.NoError(t, err)
		require.Equal(t, "Hello World", strings.TrimSpace(string(b)))
	})

	t.Run("json decode", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Body().JSON([]int{1, 2, 3}),
			Fdebug(buf).Response().Body(),
		)

		var m []interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))

		require.Equal(t, []interface{}{1.0, 2.0, 3.0}, m)
	})
}

func TestDebugResponse_Header(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("without header", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Body().String("Hello"),
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
			Send().Body().String("Hello"),
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
			Send().Body().String("Hello"),
			Send().Headers("X-Header1").Add("Foo"),
			Send().Headers("X-Header2").Add("Bar"),
			Expect().Body().String().Equal("Hello World"),
			Fdebug(buf).Response().Headers(),
			Clear().Expect(),
		)
	})
}

func TestDebugResponse_Trailer(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("without header", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)

		Test(t,
			Post(s.URL),
			Send().Body().String("Hello"),
			Send().Trailers("X-Trailer1").Add("Foo"),
			Send().Trailers("X-Trailer2").Add("Bar"),
			Fdebug(buf).Response().Trailers(),
		)

		var m map[string]interface{}
		require.NoError(t, json.NewDecoder(vtclean.NewReader(buf, false)).Decode(&m))
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
			Fdebug(buf).Response().Trailers("X-Trailer1"),
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
			Fdebug(buf).Response().Trailers(),
			Clear().Expect(),
		)
	})
}
