# go-testdoc [![Actions Status](https://github.com/Eun/go-testdoc/workflows/CI/badge.svg)](https://github.com/Eun/go-testdoc/actions) [![Coverage Status](https://coveralls.io/repos/github/Eun/go-testdoc/badge.svg?branch=master)](https://coveralls.io/github/Eun/go-testdoc?branch=master) [![PkgGoDev](https://img.shields.io/badge/pkg.go.dev-reference-blue)](https://pkg.go.dev/github.com/Eun/go-testdoc) [![GoDoc](https://godoc.org/github.com/Eun/go-testdoc?status.svg)](https://godoc.org/github.com/Eun/go-testdoc) [![go-report](https://goreportcard.com/badge/github.com/Eun/go-testdoc)](https://goreportcard.com/report/github.com/Eun/go-testdoc)
*go-testdoc* runs your code documentation examples during normal test time.
> go get -u github.com/Eun/go-testdoc

### Example
This example will run the go code inside the `Example` and `Examples` section using [yaegi](https://github.com/traefik/yaegi).
```go
// example.go file
package example

// IsLoggedIn returns true if the current user has access to the service.
//
// Example:
//     if IsLoggedIn() {
//         fmt.Println("You are logged in")
//     }
func IsLoggedIn() bool {
	return true
}

// CurrentUser returns the current username of the logged in user.
//
// Examples:
//     fmt.Println(CurrentUser())
//
//     if IsLoggedIn() {
//         fmt.Println(CurrentUser())
//     }
func CurrentUser() string {
	return "Joe Doe"
}

```
```go
// example_test.go file
package example_test

import (
	"testing"

	"github.com/Eun/go-testdoc"
)

func TestDocumentation(t *testing.T) {
	testdoc.TestCodeDocumentation(t, &testdoc.Options{
		// Test for this folder
		Path: ".",

		// Test the `example` package
		PkgName: "example",

		// Execute code inside the `Example` and `Examples` sections
		Sections: []string{"Example", "Examples"},

		Imports: []testdoc.Import{
			// Import some standard packages we need
			{Name: "", Path: "fmt"},

			// Import the current package so we can call the functions.
			{Name: ".", Path: "./"},
		},
	})
}

```


## See also
[Testing your `README.md` file](https://github.com/Eun/yaegi-template/tree/master/examples/evaluate_readme)