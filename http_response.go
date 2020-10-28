package hit

import (
	"net/http"

	"github.com/Eun/go-hit/httpbody"
)

// HTTPResponse contains the http.Response and methods to extract/set data for the body.
type HTTPResponse struct {
	Hit Hit
	*http.Response
	body *httpbody.HTTPBody
}

func newHTTPResponse(hit Hit, response *http.Response) *HTTPResponse {
	r := &HTTPResponse{
		Hit:      hit,
		Response: response,
		body:     httpbody.NewHTTPBody(response.Body, response.Header),
	}
	r.body.SetReader(response.Body)
	return r
}

// Body provides methods for accessing the http body.
func (r *HTTPResponse) Body() *httpbody.HTTPBody {
	return r.body
}
