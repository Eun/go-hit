package simple_test

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"testing"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
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

func TestStatusCode(t *testing.T) {
	Test(t,
		Head("https://google.com"),
		Expect().Custom(func(e Hit) {
			if e.Response().StatusCode > 400 {
				// hit will catch errors
				// so feel free to panic here
				panic("Expected StatusCode to be less than 400")
			}
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
