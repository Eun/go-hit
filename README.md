
# go-hit [![Actions Status](https://github.com/Eun/go-hit/workflows/CI/badge.svg)](https://github.com/Eun/go-hit/actions) [![Coverage Status](https://coveralls.io/repos/github/Eun/go-hit/badge.svg?branch=master)](https://coveralls.io/github/Eun/go-hit?branch=master) [![PkgGoDev](https://img.shields.io/badge/pkg.go.dev-reference-blue)](https://pkg.go.dev/github.com/Eun/go-hit) [![GoDoc](https://godoc.org/github.com/Eun/go-hit?status.svg)](https://godoc.org/github.com/Eun/go-hit) [![go-report](https://goreportcard.com/badge/github.com/Eun/go-hit)](https://goreportcard.com/report/github.com/Eun/go-hit)
hit is an **h**ttp **i**ntegration **t**est framework written in golang.

It is designed to be flexible as possible, but to keep a simple to use interface for developers.

So lets get started!

> go get -u github.com/Eun/go-hit

```go //ignore
package main

import (
    "net/http"
    . "github.com/Eun/go-hit"
)

func main() {
    MustDo(
        Description("Post to httpbin.org"),
        Get("https://httpbin.org/post"),
        Expect().Status().Equal(http.StatusMethodNotAllowed),
        Expect().Body().String().Contains("Method Not Allowed"),
    )
}
```

Or use the `Test()` function:
```go //ignore
package main_test
import (
    "testing"
    "net/http"
    . "github.com/Eun/go-hit"
)

func TestHttpBin(t *testing.T) {
    Test(t,
        Description("Post to httpbin.org"),
        Get("https://httpbin.org/post"),
        Expect().Status().Equal(http.StatusMethodNotAllowed),
        Expect().Body().String().Contains("Method Not Allowed"),
    )
}
``` 

## Expect, Expect, Expect, ....
```go
MustDo(
    Get("https://httpbin.org/post"),
    Expect().Status().Equal(http.StatusMethodNotAllowed),
    Expect().Headers("Content-Type").NotEmpty(),
    Expect().Body().String().Contains("Method Not Allowed"),
)
``` 

## Sending Data
```go
MustDo(
    Post("https://httpbin.org/post"),
    Send().Body().String("Hello HttpBin"),
    Expect().Status().Equal(http.StatusOK),
    Expect().Body().String().Contains("Hello HttpBin"), 
)
``` 


### Sending And Expecting JSON
```go
MustDo(
    Post("https://httpbin.org/post"),
    Send().Headers("Content-Type").Add("application/json"),
    Send().Body().JSON(map[string][]string{"Foo": []string{"Bar", "Baz"}}),
    Expect().Status().Equal(http.StatusOK),
    Expect().Body().JSON().JQ(".json.Foo[1]").Equal("Baz"),
)
``` 


## Storing Data From The Response
```go
var name string
var roles []string
MustDo(
    Post("https://httpbin.org/post"),
    Send().Headers("Content-Type").Add("application/json"),
    Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Roles": []string{"Admin", "Developer"}}),
    Expect().Status().Equal(http.StatusOK),
    Store().Response().Body().JSON().JQ(".json.Name").In(&name),
    Store().Response().Body().JSON().JQ(".json.Roles").In(&roles),
)
fmt.Printf("%s has %d roles\n", name, len(roles))
``` 

## Problems? `Debug`!
```go
MustDo(
    Post("https://httpbin.org/post"),
    Debug(),
    Debug().Response().Body(),
)
```

## Handling Errors
It is possible to handle errors in a custom way.
```go //ignore
func login(username, password string) error {
    err := Do(
         Get("https://httpbin.org/basic-auth/joe/secret"),
         Send().Headers("Authorization").Add("Basic " + base64.StdEncoding.EncodeToString([]byte(username + ":" + password))),
         Expect().Status().Equal(http.StatusOK),
    )
    var hitError *Error
    if errors.As(err, &hitError) {
        if hitError.FailingStepIs(Expect().Status().Equal(http.StatusOK)) {
            return errors.New("login failed")
        }
    }
    return err
}
```

## Build the request url manually
```go
MustDo(
    Request().Method(http.MethodPost),
    Request().URL().Scheme("https"),
    Request().URL().Host("httpbin.org"),
    Request().URL().Path("/post"),
    Request().URL().Query("page").Add(1),
    Expect().Status().Equal(200),
    Send().Body().String("Hello World"),
    Expect().Body().String().Contains("Hello"),
)
```

## Twisted!
Although the following is hard to read it is possible to do!
```go
MustDo(
    Post("https://httpbin.org/post"),
    Expect().Status().Equal(200),
    Send().Body().String("Hello World"),
    Expect().Body().String().Contains("Hello"),
)
```

## Custom Send And Expects
```go
MustDo(
    Get("https://httpbin.org/get"),
    Send().Custom(func(hit Hit) error {
        hit.Request().Body().SetStringf("Hello %s", "World")
        return nil
    }),
    Expect().Custom(func(hit Hit) error {
        if len(hit.Response().Body().MustString()) <= 0 {
            return errors.New("expected the body to be not empty")
        }
        return nil
    }),
    Custom(AfterExpectStep, func(Hit) error {
        fmt.Println("everything done")
        return nil
    }),
)
```

## Templates / Multiuse
```go
template := CombineSteps(
    Post("https://httpbin.org/post"),
    Send().Headers("Content-Type").Add("application/json"),
    Expect().Headers("Content-Type").Equal("application/json"),
)
MustDo(
    template,
    Send().Body().JSON("Hello World"),
)

MustDo(
    template,
    Send().Body().JSON("Hello Universe"),
)
```

## Clean Previous Steps
Sometimes it is necessary to remove some steps that were added before.
```go
template := CombineSteps(
    Get("https://httpbin.org/basic-auth/joe/secret"),
    Expect().Status().Equal(http.StatusOK),
)
MustDo(
    Description("login with correct credentials"),
    template,
    Send().Headers("Authorization").Add("Basic " + base64.StdEncoding.EncodeToString([]byte("joe:secret"))),
)

Test(t,
    Description("login with incorrect credentials"),
    template,
    Clear().Expect().Status(),
    Expect().Status().Equal(http.StatusUnauthorized),
    Send().Headers("Authorization").Add("Basic " + base64.StdEncoding.EncodeToString([]byte("joe:joe"))),
)
```
More examples can be found in the `examples directory`

### Changelog
#### 0.5.0
* Rehaul the api, make things more explicit
* Fix some issues
* `Store()` functionality
* Generate `Clear()` paths
* Infinite `JQ()` functionality
* Test `README.md` and documentation parts

#### 0.4.0
* Fixed a double run bug in `CombineSteps` (#3)
* Better `Clear` functionality, you can now clear any previous step by prepending `Clear()`,
eg. 
  ```go //ignore
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
          Send().Body().String("Hello World"),
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
  ```go //ignore
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
