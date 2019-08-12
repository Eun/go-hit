# hit [![Travis](https://img.shields.io/travis/Eun/go-hit.svg)](https://travis-ci.org/Eun/go-hit) [![Codecov](https://img.shields.io/codecov/c/github/Eun/go-hit.svg)](https://codecov.io/gh/Eun/go-hit) [![GoDoc](https://godoc.org/github.com/Eun/go-hit?status.svg)](https://godoc.org/github.com/Eun/go-hit) [![go-report](https://goreportcard.com/badge/github.com/Eun/go-hit)](https://goreportcard.com/report/github.com/Eun/go-hit)
hit is an **h**ttp **i**ntegration **t**est framework written in golang.

It is designed to be flexible as possible, but to keep a simple to use interface for developers.

So lets get started!

> go get -u github.com/Eun/go-hit

```go
import (
	"net/http"
	"github.com/Eun/go-hit"
)

func TestHttpBin(t *testing.T) {
	hit.Get(t, "https://httpbin.org/post").
		Expect().Status(http.StatusMethodNotAllowed).
		Do()
}
``` 

## expect, expect, expect, ....
```go
hit.Get(t, "https://httpbin.org/post").
	Expect().Status(http.StatusMethodNotAllowed).
	Expect().Headers().Contains("Content-Type").
	Expect().Body().Contains("Method Not Allowed").
	Do()
}
``` 

## Sending some data
```go
hit.Post(t, "https://httpbin.org/post").
	Send().Body("Hello HttpBin").
	Expect().Status(http.StatusOK).
	Expect().Body().Contains("Hello HttpBin").
	Do()
``` 

### Sending and expecting JSON
```go
hit.Post(t, "https://httpbin.org/post").
	Send().Headers().Set("Content-Type", "application/json").
	Send().Body().JSON(map[string][]string{"Foo": []string{"Bar", "Baz"}}).
	Expect().Status(http.StatusOK).
	Expect().Body().JSON().Equal("json.Foo.1", "Baz").
	Do()
``` 

## Problems? `Debug`!
```go
hit.Post(t, "https://httpbin.org/post").
	Debug().
	Do().
	Debug()
```

## Twisted!
Although the following is hard to read it is possible to do!
```go
hit.Post(t, "https://httpbin.org/post").
	Expect().Status(200).
	Send("Hello World").
	Do().
	Expect().Body().Contains("Hello")
```

## Custom Send and Expects
```go
hit.Get(t, "https://httpbin.org/get").
	Send().Custom(func(hit hit.Hit) {
		hit.Request().Body().SetStringf("Hello %s", "World")
	}).
	Expect().Custom(func(hit hit.Hit) {
		if len(hit.Response().Body().String()) <= 0 {
			t.FailNow()
		}
	}).
	Do()
```

## Templates
```go
template := hit.Post(t, "https://httpbin.org/post").
	Send().Headers().Set("Content-Type", "application/json").
	Expect().Headers("Content-Type").Equal("application/json")

template.Copy().
	Send().Body().JSON("Hello World").
	Do()

template.Copy().
	Send().Body().JSON("Hello Universe").
	Do()
```

## BaseURL Template
```go
template := hit.New(t, "https://httpbin.org").

template.Copy().Get("/get").
	Do()

template.Copy().Post("/post").
	Do()
```

More examples in the `example_*.go` files