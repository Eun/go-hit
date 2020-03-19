# hit [![Travis](https://img.shields.io/travis/Eun/go-hit.svg)](https://travis-ci.org/Eun/go-hit) [![Codecov](https://img.shields.io/codecov/c/github/Eun/go-hit.svg)](https://codecov.io/gh/Eun/go-hit) [![GoDoc](https://godoc.org/github.com/Eun/go-hit?status.svg)](https://godoc.org/github.com/Eun/go-hit) [![go-report](https://goreportcard.com/badge/github.com/Eun/go-hit)](https://goreportcard.com/report/github.com/Eun/go-hit)
hit is an **h**ttp **i**ntegration **t**est framework written in golang.

It is designed to be flexible as possible, but to keep a simple to use interface for developers.

So lets get started!

> go get -u github.com/Eun/go-hit

```go
import (
	"net/http"
	. "github.com/Eun/go-hit"
)

func TestHttpBin(t *testing.T) {
    Test(t,
        Description("Post to httpbin.org"),
        Get("https://httpbin.org/post"),
        Expect().Status(http.StatusMethodNotAllowed),
    )
}
``` 

## expect, expect, expect, ....
```go
Test(t,
    Get("https://httpbin.org/post"),
    Expect().Status(http.StatusMethodNotAllowed),
    Expect().Headers().Contains("Content-Type"),
    Expect().Body().Contains("Method Not Allowed"),
)
``` 

## Sending some data
```go
Test(t,
    Post("https://httpbin.org/post"),
    Send().Body("Hello HttpBin"),
    Expect().Status(http.StatusOK),
    Expect().Body().Contains("Hello HttpBin"), 
)
``` 

### Sending and expecting JSON
```go
Test(t,
    Post("https://httpbin.org/post"),
    Send().Headers().Set("Content-Type", "application/json"),
    Send().Body().JSON(map[string][]string{"Foo": []string{"Bar", "Baz"}}),
    Expect().Status(http.StatusOK),
    Expect().Body().JSON().Equal("json.Foo.1", "Baz"),
)
``` 

## Problems? `Debug`!
```go
Test(
    Post(t, "https://httpbin.org/post"),
    Debug(),
)
```

## Twisted!
Although the following is hard to read it is possible to do!
```go
Test(t,
    Post("https://httpbin.org/post"),
    Expect().Status(200),
    Send("Hello World"),
    Expect().Body().Contains("Hello"),
)
```

## Custom Send and Expects
```go
Test(t,
    Get("https://httpbin.org/get"),
    Send().Custom(func(hit Hit) {
        hit.Request().Body().SetStringf("Hello %s", "World")
    }),
    Expect().Custom(func(hit Hit) {
        if len(hit.Response().Body().String()) <= 0 {
            t.FailNow()
        }
    }),
    Custom(AfterExpectStep, func(Hit) {
        fmt.Println("everything done")
    }),
)
```

## Templates / Multiuse
```go
template := []IStep{
    Post("https://httpbin.org/post"),
    Send().Headers().Set("Content-Type", "application/json"),
    Expect().Headers("Content-Type").Equal("application/json"),
}
Test(t,
    append(template,
        Send().Body().JSON("Hello World"),
    )...,
)

Test(t,
    append(template,
        Send().Body().JSON("Hello Universe"),
    )...,
)
```
More examples can be found in the `examples directory`

### Changelog
#### 0.4.0
* Fixed a double run bug in `CombineSteps` (#3)
* Better `Clear` functionality, you can now clear any previous step by prepending `Clear()`,
eg. 
  ```go
  Do(
      Get("https://example.com"),
      Expect().Body("Hello World"),
      Clear().Expect(),                        // remove all previous Expect() steps
      // Clear().Expect().Body("Hello World"), // remove only the Expect().Body("Hello World") step
  )
  
  // also works for CombineSteps
  Do(
      Post("https://example.com"),        
      CombineSteps(
          Send().Body("Hello World"),
          Expect().Body("Hello World"),
      ),
      Clear().Expect(),
  )
  ```
* Simplified `Expect().Header()` use, no more `Expect().Headers()`,
everything can now be done in the `Expect().Header()` function.
* More documentation and examples
* `hit.Do` and `hit.MustDo` for inline steps.
* Removal of inline steps (use `hit.Do` and `hit.MustDo`)
  ```go
  Do(
      Get("https://example.com"),
      Expect().Custon(func (hit Hit) {
          // Expect("Hello World") is invalid now
          // you must use MustDo() or Do()
          hit.MustDo(
              Expect("Hello World"),
          )
      }),
  )
  ```
* `hit.InsertSteps` to insert steps during runtime