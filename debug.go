package hit

import (
	"encoding/json"
	"fmt"
	"time"

	"io"

	"reflect"

	"github.com/Eun/go-hit/expr"
	"github.com/Eun/go-hit/internal"
	"github.com/gookit/color"
	"github.com/tidwall/pretty"
)

// IDebug provides a debug functionality for the Request and Response
type IDebug interface {
	IStep
	// Time prints the current Time
	Time() IStep
	// Request prints the Request
	//
	// Usage:
	//     Debug().Request()                        // print the whole Request
	//     Debug().Request().Body()                 // print only the body
	//     Debug().Request().Header()               // print all headers
	//     Debug().Request().Header("Content-Type") // print the Content-Type header
	Request() IDebugRequest

	// Request prints the Response
	//
	// Usage:
	//     Debug().Response()                        // print the whole Response
	//     Debug().Response().Body()                 // print only the body
	//     Debug().Response().Header()               // print all headers
	//     Debug().Response().Header("Content-Type") // print the Content-Type header
	Response() IDebugResponse
}

type debug struct {
	expression []string
}

func newDebug(expression []string) IDebug {
	return &debug{
		expression: expression,
	}
}

func (*debug) when() StepTime {
	return BeforeExpectStep
}

func (*debug) getBody(body *HTTPBody) interface{} {
	reader := body.JSON().body.Reader()
	// if there is a json reader
	if reader != nil {
		var container interface{}
		if err := json.NewDecoder(reader).Decode(&container); err == nil {
			return container
		}
	}
	s := body.String()
	if len(s) == 0 {
		return nil
	}
	return s
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

func (d *debug) printWithExpression(hit Hit, v interface{}, expression []string) error {
	if e, ok := internal.GetLastStringArgument(expression); ok {
		v = expr.MustGetValue(v, e, expr.IgnoreError, expr.IgnoreNotFound)
	}
	return d.print(hit, v)
}

func (*debug) print(hit Hit, v interface{}) error {
	if v == nil {
		_, err := io.WriteString(hit.Stdout(), "<nil>")
		return err
	}

	// if v is struct, map, or slice, use json
	// if not print directly

	rv := reflect.TypeOf(v)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct && rv.Kind() != reflect.Slice && rv.Kind() != reflect.Map {
		_, err := fmt.Fprintf(hit.Stdout(), "%v", v)
		return err
	}

	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	bytes = pretty.Pretty(bytes)
	if color.IsSupportColor() {
		bytes = pretty.Color(bytes, nil)
	}
	_, err = hit.Stdout().Write(bytes)
	return err
}

func (d *debug) exec(hit Hit) error {
	type M map[string]interface{}

	m := M{
		"Time": time.Now().String(),
	}

	if hit.Request() != nil {
		m["Request"] = newDebugRequest().data(hit)
	}

	if hit.Response() != nil {
		m["Response"] = newDebugResponse().data(hit)
	}

	var v interface{} = m
	if e, ok := internal.GetLastStringArgument(d.expression); ok {
		fmt.Println("Warning: Debug(something) is deprecated, use Debug().Something")
		v = expr.MustGetValue(m, e, expr.IgnoreError, expr.IgnoreNotFound)
	}

	return d.print(hit, v)
}

func (*debug) clearPath() clearPath {
	return nil // not clearable
}

func (*debug) Time() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      BeforeExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			_, err := io.WriteString(hit.Stdout(), time.Now().String())
			return err
		},
	}
}

func (*debug) Request() IDebugRequest {
	return newDebugRequest()
}

func (*debug) Response() IDebugResponse {
	return newDebugResponse()
}
