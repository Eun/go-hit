package hit

import (
	"net/http"

	"github.com/Eun/go-hit/httpbody"
)

type HTTPResponse struct {
	Hit Hit
	*http.Response
	body *httpbody.HttpBody
}

func newHTTPResponse(hit Hit, response *http.Response) *HTTPResponse {
	r := &HTTPResponse{
		Hit:      hit,
		Response: response,
		body:     httpbody.NewHttpBody(response.Body, response.Header),
	}
	r.body.SetReader(response.Body)
	return r
}

func (r *HTTPResponse) Body() *httpbody.HttpBody {
	return r.body
}
