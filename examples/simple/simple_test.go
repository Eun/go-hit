package simple_test

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"testing"

	"golang.org/x/xerrors"

	"github.com/stretchr/testify/require"

	_ "github.com/Eun/go-hit/doctest/implicit"
	. "github.com/otto-eng/go-hit"
)

func TestHead(t *testing.T) {
	Test(t,
		Head("https://httpbin.org/status/300"),
		Expect().Status().Equal(300),
	)
}

func TestPost(t *testing.T) {
	Test(t,
		Post("https://httpbin.org/post"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(map[string]interface{}{"Foo": "Bar"}),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().JQ(".json.Foo").Equal("Bar"),
	)
}

func TestPatch(t *testing.T) {
	Test(t,
		Patch("http://%s/%s", "127.0.0.1:8081/v1", "user"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("Authorization").Add("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI4ZTMzNGMxZC0yYzdjLTRhYTAtYTM2Mi1iNWViNWIwMGJkMmQiLCJleHAiOjE2NTMyODcwMTh9.6oD8LDhftnQ9iIdqo6XSnZWMQCitc5nxTyx0q5cv2wc"),
		Send().Body().JSON(map[string]interface{}{"firstName": "paultest"}),
		Expect().Status().Equal(http.StatusNoContent),
	)
}

func TestStatusCode(t *testing.T) {
	Test(t,
		Head("http://google.com"),
		Expect().Custom(func(e Hit) error {
			if e.Response().StatusCode > 400 {
				// hit will catch errors
				// so feel free to panic here
				return xerrors.New("Expected StatusCode to be less than 400")
			}
			return nil
		}),
	)
}

func TestCookie(t *testing.T) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Jar: jar,
	}

	steps := []IStep{
		HTTPClient(client),
		BaseURL("https://httpbin.org/cookies"),
		Expect().Status().Equal(http.StatusOK),
	}

	Test(t, append(steps, Get("/set/CookieA/Value123"))...)

	Test(t, append(steps,
		Get(""),
		Expect().Body().JSON().JQ(".cookies.CookieA").Equal("Value123"),
	)...)
}

func TestStore(t *testing.T) {
	var name string
	var roles []string
	Test(t,
		Post("https://httpbin.org/post"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Roles": []string{"Admin", "Developer"}}),
		Expect().Status().Equal(http.StatusOK),
		Store().Response().Body().JSON().JQ(".json.Name").In(&name),
		Store().Response().Body().JSON().JQ(".json.Roles").In(&roles),
	)
	require.Equal(t, "Joe", name)
	require.Equal(t, []string{"Admin", "Developer"}, roles)
}
