package hit_test

import (
	"net/http"

	"net/http/cookiejar"

	"github.com/Eun/go-hit"
)

func ExampleHead() {
	hit.Head(hit.PanicT{}, "https://httpbin.org/status/300").
		Expect().Status(300).
		Do()
}

func ExamplePost() {
	hit.Post(hit.PanicT{}, "https://httpbin.org/post").
		Send().Body(map[string]interface{}{"Foo": "Bar"}).
		Do().
		Expect().Status(200).
		Expect().Body().JSON().Equal("json.Foo", "Bar")
}

func Example_statusCode() {
	hit.Head(hit.PanicT{}, "https://google.com").
		Expect().
		Custom(func(e hit.Hit) {
			if e.Response().StatusCode > 400 {
				hit.PanicT{}.FailNow()
			}
		}).
		Do()
}

func Example_cookie() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Jar: jar,
	}

	template := hit.New(hit.PanicT{}).
		SetHTTPClient(client).
		Expect().Status(http.StatusOK)

	template.Copy().
		Get("https://httpbin.org/cookies/set/CookieA/Value123").
		Do()

	template.Copy().
		Get("https://httpbin.org/cookies").
		Expect().Body().JSON().Equal("cookies.CookieA", "Value123").
		Do()
}

func Example_cookie_alternative() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Jar: jar,
	}
	hit.Get(hit.PanicT{}, "https://httpbin.org/cookies/set/CookieA/Value123").
		SetHTTPClient(client).
		Do().
		Expect().Status(http.StatusOK)

	hit.Get(hit.PanicT{}, "https://httpbin.org/cookies").
		SetHTTPClient(client).
		Do().
		Expect().Status(http.StatusOK).
		Expect().Body().JSON().Equal("cookies.CookieA", "Value123")
}
