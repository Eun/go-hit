package hit

import (
	"encoding/json"

	"github.com/Eun/go-hit/expr"
	"github.com/Eun/go-hit/httpbody"
	"github.com/Eun/go-hit/internal/misc"
)

type IStoreBody interface {
	IStoreStep

	// JSON treats the body as JSON data and stores it
	//
	// given the following body: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	// Usage:
	//     var body map[string]interface{}
	//     Body().JSON().In(&body) // store the whole body as a map
	//     var name string
	//     Store().Request().Body().JSON("Name").In(&name) // store "Joe" in name
	JSON(expression ...string) IStoreStep
}

type storeBodyMode int

const (
	storeBodyRequest storeBodyMode = iota
	storeBodyResponse
)

type storeBody struct {
	mode       storeBodyMode
	expression []string
}

func newStoreBody(mode storeBodyMode, expression []string) IStoreBody {
	return &storeBody{
		mode:       mode,
		expression: expression,
	}
}

func (s *storeBody) body(hit Hit) *httpbody.HttpBody {
	if s.mode == storeBodyRequest {
		return hit.Request().Body()
	}
	return hit.Response().Body()
}

func (s *storeBody) In(v interface{}) IStep {
	if e, ok := misc.GetLastStringArgument(s.expression); ok {
		return newStoreInStep(func(hit Hit) (interface{}, error) {
			v, _, err := expr.GetValue(s.body(hit).GetBestFittingObject(), e, expr.IgnoreCase)
			return v, err
		}, v)
	}
	return newStoreInStep(func(hit Hit) (interface{}, error) {
		return s.body(hit), nil
	}, v)
}

func (s *storeBody) JSON(expression ...string) IStoreStep {
	e, _ := misc.GetLastStringArgument(expression)
	return newStoreStep(func(hit Hit) (interface{}, error) {
		var container interface{}
		if err := json.NewDecoder(s.body(hit).Reader()).Decode(&container); err != nil {
			return nil, err
		}
		v, ok, err := expr.GetValue(container, e, expr.IgnoreCase)
		if err != nil {
			return nil, err
		}
		if !ok {
			v = nil
		}
		return v, nil
	})
}
