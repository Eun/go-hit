package hit

import (
	"github.com/Eun/go-hit/httpbody"
)

// IStoreBodyJSON defines the functions that can be used to store data from the http request/response body
// (in JSON format).
type IStoreBodyJSON interface {
	IStoreStep

	// JQ runs an jq expression on the JSON body the result can than be stored afterwards
	//
	// Example:
	//     // given the following body: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	//     var name string
	//     MustDo(
	//         Get("https://example.com/json"),
	//         Store().Response().Body().JSON().JQ(".Name").In(&name), // store "Joe" in name
	//     )
	JQ(expression ...string) IStoreStep
}

type storeBodyJSON struct {
	mode storeBodyMode
}

func newStoreBodyJSON(mode storeBodyMode) IStoreBodyJSON {
	return &storeBodyJSON{
		mode: mode,
	}
}

func (s *storeBodyJSON) body(hit Hit) *httpbody.HTTPBody {
	if s.mode == storeBodyRequest {
		return hit.Request().Body()
	}
	return hit.Response().Body()
}

func (s *storeBodyJSON) JQ(expression ...string) IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return s.body(hit).JSON().JQ(v, expression...)
	})
}

func (s *storeBodyJSON) In(v interface{}) IStep {
	return newStoreInStep(func(hit Hit, v interface{}) error {
		return s.body(hit).JSON().Decode(v)
	}, v)
}
