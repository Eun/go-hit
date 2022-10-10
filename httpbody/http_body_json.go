package httpbody

import (
	"bytes"
	"io"

	"github.com/itchyny/gojq"
	jsoniter "github.com/json-iterator/go"
	"github.com/tomwright/dasel"
	daselstorage "github.com/tomwright/dasel/storage"
	"golang.org/x/xerrors"
)

// HTTPBodyJSON provides JSON functions for the HTTPBody.
type HTTPBodyJSON struct { //nolint:revive //ignore type name will be used as httpbody.HTTPBodyJSON by other packages
	body    func() io.ReadCloser
	setBody func([]byte)
}

func newHTTPBodyJSON(body func() io.ReadCloser, setBody func([]byte)) *HTTPBodyJSON {
	return &HTTPBodyJSON{
		body:    body,
		setBody: setBody,
	}
}

// Decode decodes the body as JSON.
func (jsn *HTTPBodyJSON) Decode(container interface{}) error {
	return json.NewDecoder(jsn.body()).Decode(&container)
}

// MustDecode decodes the body as JSON, evaluates the expression and puts the result into the container
// it will panic if something goes wrong.
func (jsn *HTTPBodyJSON) MustDecode(container interface{}) {
	if err := jsn.Decode(container); err != nil {
		panic(err)
	}
}

// Set sets the body to the specified json data.
func (jsn *HTTPBodyJSON) Set(data interface{}) error {
	if jsn.setBody == nil {
		return xerrors.New("setBody is nil")
	}
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	jsn.setBody(buf)
	return nil
}

type jsonInputIter struct {
	dec *jsoniter.Decoder
	buf *bytes.Buffer
	err error
}

func newJSONInputIter(r io.Reader) *jsonInputIter {
	buf := new(bytes.Buffer)
	return &jsonInputIter{dec: json.NewDecoder(io.TeeReader(r, buf)), buf: buf}
}

func (i *jsonInputIter) Next() (interface{}, bool) {
	if i.err != nil {
		return nil, false
	}
	var v interface{}
	if err := i.dec.Decode(&v); err != nil {
		if err == io.EOF {
			i.err = err
			return err, false
		}
		i.err = err
		return i.err, true
	}
	if i.buf.Len() >= 256*1024 {
		i.buf.Reset()
	}
	return v, true
}

func (i *jsonInputIter) Close() error {
	i.err = io.EOF
	return nil
}

// JQ runs an jq expression on the JSON body the result will be stored into container.
func (jsn *HTTPBodyJSON) JQ(container interface{}, expression ...string) error {
	var iter gojq.Iter
	jsonIter := newJSONInputIter(jsn.body())
	defer func() {
		_ = jsonIter.Close()
	}()
	iter = jsonIter

	for _, e := range expression {
		query, err := gojq.Parse(e)
		if err != nil {
			return err
		}

		code, err := gojq.Compile(query, gojq.WithInputIter(iter))
		if err != nil {
			return err
		}
		v, hasNext := iter.Next()
		if err, ok := v.(error); ok {
			return err
		}
		if !hasNext {
			return xerrors.New("no source iter")
		}
		iter = code.Run(v)
	}

	// read first item
	v, hasNext := iter.Next()
	if err, ok := v.(error); ok {
		return err
	}
	// first item is missing
	if !hasNext {
		return json.NewDecoder(bytes.NewReader([]byte("null"))).Decode(container)
	}

	// remember the first item for later
	col := []interface{}{
		v,
		nil, // initialize the second item so we can save a resize later
	}

	// read second item
	v, hasNext = iter.Next()
	if err, ok := v.(error); ok {
		return err
	}
	// no second item, exit with the first item
	if !hasNext {
		// let the json encode encode and decode the item
		// so we can support the `json:"name"` tag
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(col[0]); err != nil {
			return err
		}

		return json.NewDecoder(&buf).Decode(container)
	}
	// set the second item
	col[1] = v

	for {
		v, hasNext = iter.Next()
		if err, ok := v.(error); ok {
			return err
		}
		if !hasNext {
			// let the json encode encode and decode the item
			// so we can support the `json:"name"` tag
			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(col); err != nil {
				return err
			}
			return json.NewDecoder(&buf).Decode(container)
		}
		col = append(col, v)
	}
}

// MustJQ runs an jq expression on the JSON body the result will be stored into container, if an error occurs it will
// panic.
func (jsn *HTTPBodyJSON) MustJQ(container interface{}, expression ...string) {
	if err := jsn.JQ(container, expression...); err != nil {
		panic(err)
	}
}

// Dasel runs an dasel expression on the JSON body the result will be stored into container.
func (jsn *HTTPBodyJSON) Dasel(container interface{}, expression ...string) error {
	value, err := daselstorage.Load(&daselstorage.JSONParser{}, jsn.body())
	if err != nil {
		return err
	}

	if value == nil {
		return io.EOF
	}

	node := dasel.New(value)

	if !node.Value.IsValid() {
		return json.NewDecoder(bytes.NewReader([]byte("null"))).Decode(container)
	}

	for _, e := range expression {
		nodes, err := node.QueryMultiple(e)
		if err != nil {
			if !xerrors.Is(err, dasel.ValueNotFound{}) {
				return err
			}
			return json.NewDecoder(bytes.NewReader([]byte("null"))).Decode(container)
		}

		switch n := len(nodes); n {
		case 0:
			return json.NewDecoder(bytes.NewReader([]byte("null"))).Decode(container)
		case 1:
			if !nodes[0].Value.IsValid() {
				return json.NewDecoder(bytes.NewReader([]byte("null"))).Decode(container)
			}
			node = dasel.New(nodes[0].InterfaceValue())
		default:
			values := make([]interface{}, n)
			for i := range nodes {
				values[i] = nodes[i].InterfaceValue()
			}
			node = dasel.New(values)
		}
	}

	// item is invalid
	if !node.Value.IsValid() {
		return json.NewDecoder(bytes.NewReader([]byte("null"))).Decode(container)
	}

	switch n := len(node.NextMultiple); n {
	case 0:
		// let the json encoder encode and decode the item
		// so we can support the `json:"name"` tag
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(node.InterfaceValue()); err != nil {
			return err
		}

		return json.NewDecoder(&buf).Decode(container)
	default:
		col := make([]interface{}, 1+n)
		col[0] = node.InterfaceValue()
		// remember the first item for later
		for i := 0; i < n; i++ {
			col[i+1] = node.NextMultiple[i].InterfaceValue()
			node = node.Next
		}
		// let the json encoder encode and decode the item
		// so we can support the `json:"name"` tag
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(col); err != nil {
			return err
		}

		return json.NewDecoder(&buf).Decode(container)
	}

}

// MustDasel runs an dasel expression on the JSON body the result will be stored into container, if an error occurs it will
// panic.
func (jsn *HTTPBodyJSON) MustDasel(container interface{}, expression ...string) {
	if err := jsn.Dasel(container, expression...); err != nil {
		panic(err)
	}
}
