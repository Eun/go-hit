package hit_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/Eun/go-hit"
)

func TestRequestURL_Set(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	u, err := url.Parse(s.URL)
	require.NoError(t, err)

	var v string
	Test(t,
		Get(""),
		RequestURL().Set(u),
		Expect().Status().Equal(http.StatusOK),
		Store().Request().URL().Scheme().In(&v),
	)
	require.Equal(t, "http", v)
}

func TestRequestURL_Scheme(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var v string
	Test(t,
		Get(s.URL),
		RequestURL().Scheme("http"),
		Expect().Status().Equal(http.StatusOK),
		Store().Request().URL().Scheme().In(&v),
	)
	require.Equal(t, "http", v)
}

func TestRequestURL_Host(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	u, err := url.Parse(s.URL)
	require.NoError(t, err)

	Test(t,
		Get(s.URL),
		RequestURL().Host(u.Host),
		Expect().Status().Equal(http.StatusOK),
	)
}

func TestRequestURL_Path(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/foo", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Get(s.URL),
		RequestURL().Path("/foo"),
		Expect().Status().Equal(http.StatusOK),
	)
}

func TestRequestURL_RawPath(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/foo/bar", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Get(s.URL),
		RequestURL().RawPath("/foo%2fbar"),
		Expect().Status().Equal(http.StatusOK),
	)
}

func TestRequestURL_ForceQuery(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var v string
	Test(t,
		Get(s.URL),
		RequestURL().ForceQuery(true),
		Expect().Status().Equal(http.StatusOK),
		Store().Request().URL().String().In(&v),
	)
	require.Equal(t, s.URL+"?", v)
}

func TestRequestURL_Query(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		require.NoError(t, json.NewEncoder(writer).Encode(request.URL.Query()))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	t.Run("single value", func(t *testing.T) {
		Test(t,
			Get(s.URL),
			RequestURL().Query("username").Add("joe"),
			RequestURL().Query("page").Add(1),
			Expect().Status().Equal(http.StatusOK),
			Expect().Body().JSON().Equal(map[string][]string{
				"username": {"joe"},
				"page":     {"1"},
			}),
		)
	})

	t.Run("multiple values", func(t *testing.T) {
		Test(t,
			Get(s.URL),
			RequestURL().Query("usernames").Add("joe"),
			RequestURL().Query("usernames").Add("alice"),
			RequestURL().Query("pages").Add(1, 2),
			Expect().Status().Equal(http.StatusOK),
			Expect().Body().JSON().Equal(map[string][]string{
				"usernames": {"joe", "alice"},
				"pages":     {"1", "2"},
			}),
		)
	})
}

func TestRequestURL_RawQuery(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		require.NoError(t, json.NewEncoder(writer).Encode(request.URL.Query()))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Get(s.URL),
		RequestURL().RawQuery("x=1&y=2"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().Equal(map[string][]string{
			"x": {"1"},
			"y": {"2"},
		}),
	)
}

func TestRequestURL_Fragment(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	var v string
	Test(t,
		Get(s.URL),
		RequestURL().Fragment("anchor"),
		Expect().Status().Equal(http.StatusOK),
		Store().Request().URL().Fragment().In(&v),
	)
	require.Equal(t, "anchor", v)
}

func TestRequestURL_User(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		auth := request.Header.Get("Authorization")
		require.NotEmpty(t, auth)
		require.True(t, strings.HasPrefix(auth, "Basic "))
		userAndPassword, err := base64.StdEncoding.DecodeString(auth[6:])
		require.NoError(t, err)

		writer.WriteHeader(http.StatusOK)
		writer.Write(userAndPassword)
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	t.Run("only username", func(t *testing.T) {
		Test(t,
			Get(s.URL),
			RequestURL().User().Username("joe"),
			Expect().Status().Equal(http.StatusOK),
			Expect().Body().String().Equal("joe:"),
		)
	})

	t.Run("only password", func(t *testing.T) {
		Test(t,
			Get(s.URL),
			RequestURL().User().Password("secret"),
			Expect().Status().Equal(http.StatusOK),
			Expect().Body().String().Equal(":secret"),
		)
	})

	t.Run("username and password", func(t *testing.T) {
		Test(t,
			Get(s.URL),
			RequestURL().User().Username("joe"),
			RequestURL().User().Password("secret"),
			Expect().Status().Equal(http.StatusOK),
			Expect().Body().String().Equal("joe:secret"),
		)
	})
}
