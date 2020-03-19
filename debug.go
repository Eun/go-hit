package hit

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Eun/go-hit/expr"
	"github.com/Eun/go-hit/internal"
	"github.com/gookit/color"
	"github.com/tidwall/pretty"
)

type debug struct {
	expression []string
}

func newDebug(expression []string) IStep {
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

func (*debug) getHeader(header http.Header) map[string]interface{} {
	m := make(map[string]interface{})
	for key := range header {
		m[key] = header.Get(key)
	}
	return m
}

func (d *debug) exec(hit Hit) error {
	type M map[string]interface{}

	m := M{
		"Time": time.Now().String(),
	}

	if hit.Request() != nil {
		m["Request"] = M{
			"Header":           d.getHeader(hit.Request().Header),
			"Trailer":          d.getHeader(hit.Request().Trailer),
			"Method":           hit.Request().Method,
			"URL":              hit.Request().URL,
			"Proto":            hit.Request().Proto,
			"ProtoMajor":       hit.Request().ProtoMajor,
			"ProtoMinor":       hit.Request().ProtoMinor,
			"ContentLength":    hit.Request().ContentLength,
			"TransferEncoding": hit.Request().TransferEncoding,
			"Host":             hit.Request().Host,
			"Form":             hit.Request().Form,
			"PostForm":         hit.Request().PostForm,
			"MultipartForm":    hit.Request().MultipartForm,
			"RemoteAddr":       hit.Request().RemoteAddr,
			"RequestURI":       hit.Request().RequestURI,
			"Body":             d.getBody(hit.Request().Body()),
		}
	}

	if hit.Response() != nil {
		m["Response"] = M{
			"Header":           d.getHeader(hit.Response().Header),
			"Trailer":          d.getHeader(hit.Response().Trailer),
			"Proto":            hit.Response().Proto,
			"ProtoMajor":       hit.Response().ProtoMajor,
			"ProtoMinor":       hit.Response().ProtoMinor,
			"ContentLength":    hit.Response().ContentLength,
			"TransferEncoding": hit.Response().TransferEncoding,
			"Body":             d.getBody(hit.Response().body),
			"Status":           hit.Response().Status,
			"StatusCode":       hit.Response().StatusCode,
		}
	}

	var v interface{} = m
	if expression, ok := internal.GetLastStringArgument(d.expression); ok {
		v = expr.MustGetValue(m, expression, expr.IgnoreError, expr.IgnoreNotFound)
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

func (*debug) clearPath() clearPath {
	return nil // not clearable
}
