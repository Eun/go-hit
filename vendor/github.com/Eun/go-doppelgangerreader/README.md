# DoppelgangerReader [![Travis](https://img.shields.io/travis/Eun/go-doppelgangerreader.svg)](https://travis-ci.org/Eun/go-doppelgangerreader) [![Codecov](https://img.shields.io/codecov/c/github/Eun/go-doppelgangerreader.svg)](https://codecov.io/gh/Eun/go-doppelgangerreader) [![GoDoc](https://godoc.org/github.com/Eun/go-doppelgangerreader?status.svg)](https://godoc.org/github.com/Eun/go-doppelgangerreader) [![go-report](https://goreportcard.com/badge/github.com/Eun/go-doppelgangerreader)](https://goreportcard.com/report/github.com/Eun/go-doppelgangerreader)
DoppelgangerReader provides a way to read one `io.Reader` multiple times.


```go
package main

import (
	"fmt"
	"bytes"
	"io/ioutil"
	"github.com/Eun/go-doppelgangerreader"
)

func main() {
	reader := bytes.NewBufferString("Hello World")
	factory := doppelgangerreader.NewFactory(reader)
	defer factory.Close()

	d1 := factory.NewDoppelganger()
	defer d1.Close()

	fmt.Println(ioutil.ReadAll(d1))

	d2 := factory.NewDoppelganger()
	defer d2.Close()

	fmt.Println(ioutil.ReadAll(d2))
}
```