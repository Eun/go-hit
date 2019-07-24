package hit

import (
	"net/http"
)

type HTTPResponse struct {
	Hit Hit
	*http.Response
	body *HTTPBody
}

func newHTTPResponse(hit Hit, response *http.Response) *HTTPResponse {
	r := &HTTPResponse{
		Hit:      hit,
		Response: response,
		body:     newHTTPBody(hit, response.Body),
	}
	r.body.SetReader(response.Body)
	return r
}

func (r *HTTPResponse) Body() *HTTPBody {
	return r.body
}
