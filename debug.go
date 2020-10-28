package hit

import (
	"fmt"
	"time"

	"io"

	"reflect"

	"os"

	"github.com/gookit/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/pretty"

	"github.com/Eun/go-hit/errortrace"
)

// IDebug provides a debug functionality for the Request and Response.
type IDebug interface {
	IStep
	// Time prints the current Time
	Time() IStep
	// Request prints the Request
	//
	// Usage:
	//     Debug().Request()                           // print the whole Request
	//     Debug().Request().Body()                    // print the body
	//     Debug().Request().Body().JSON().JQ(".Name") // print the Name value of the JSON body
	//     Debug().Request().Headers()                  // print all headers
	//     Debug().Request().Headers("Content-Type")    // print the Content-Type header
	Request() IDebugRequest

	// Response prints the Response
	//
	// Usage:
	//     Debug().Response()                           // print the whole Response
	//     Debug().Response().Body()                    // print only the body
	//     Debug().Response().Body().JSON().JQ(".Name") // print the Name value of the JSON body
	//     Debug().Response().Headers()                  // print all headers
	//     Debug().Response().Headers("Content-Type")    // print the Content-Type header
	Response() IDebugResponse
}

type debug struct {
	cp callPath
	w  io.Writer
}

func newDebug(cp callPath, w io.Writer) IDebug {
	return &debug{
		cp: cp,
		w:  w,
	}
}

func (d *debug) out(*hitImpl) io.Writer {
	if d.w == nil {
		return os.Stdout
	}
	return d.w
}
func (*debug) trace() *errortrace.ErrorTrace {
	return nil
}
func (*debug) when() StepTime {
	return BeforeExpectStep
}

func (*debug) getMap(header map[string][]string) map[string]string {
	m := make(map[string]string)
	for key := range header {
		// skip empty fields
		if s := header[key]; len(s) > 0 && s[0] != "" {
			m[key] = s[0]
		}
	}
	return m
}

func (d *debug) print(out io.Writer, v interface{}) error {
	if v == nil {
		_, err := io.WriteString(out, "<nil>")
		return err
	}

	// if v is struct, map, or slice, use json
	// if not print directly

	rv := reflect.TypeOf(v)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct && rv.Kind() != reflect.Slice && rv.Kind() != reflect.Map {
		_, err := fmt.Fprintf(out, "%v", v)
		return err
	}

	bytes, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(v)
	if err != nil {
		return err
	}

	bytes = pretty.Pretty(bytes)
	if color.IsSupportColor() {
		bytes = pretty.Color(bytes, nil)
	}
	_, err = out.Write(bytes)
	return err
}

func (d *debug) exec(hit *hitImpl) error {
	type M map[string]interface{}

	m := M{
		"Time": time.Now().String(),
	}

	if hit.Request() != nil {
		m["Request"] = newDebugRequest(d.cp, d).data(hit)
	}

	if hit.Response() != nil {
		m["Response"] = newDebugResponse(d.cp, d).data(hit)
	}

	return d.print(d.out(hit), m)
}

func (d *debug) callPath() callPath {
	return d.cp
}

func (d *debug) Time() IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     BeforeExpectStep,
		CallPath: nil,
		Exec: func(hit *hitImpl) error {
			_, err := io.WriteString(d.out(hit), time.Now().String())
			return err
		},
	}
}

func (d *debug) Request() IDebugRequest {
	return newDebugRequest(d.cp.Push("Request", nil), d)
}

func (d *debug) Response() IDebugResponse {
	return newDebugResponse(d.cp.Push("Response", nil), d)
}
