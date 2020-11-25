package hit

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/Eun/go-hit/httpbody"
)

// HTTPRequest contains the http.Request and methods to extract/set data for the body.
type HTTPRequest struct {
	Hit Hit
	*http.Request
	body *httpbody.HTTPBody
}

func newHTTPRequest(hit Hit, req *http.Request) *HTTPRequest {
	newRequest := &http.Request{
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		URL:        new(url.URL),
	}

	if req == nil {
		return &HTTPRequest{
			Hit:     hit,
			Request: newRequest.WithContext(context.Background()),
			body:    httpbody.NewHTTPBody(nil, newRequest.Header),
		}
	}
	u := *req.URL
	newRequest.URL = &u
	newRequest.Method = req.Method

	// copy headers
	if req.Header != nil {
		for k, v := range req.Header {
			newRequest.Header[k] = v
		}
	}
	// copy trailers
	if req.Trailer != nil {
		newRequest.Trailer = make(http.Header)
		for k, v := range req.Trailer {
			newRequest.Trailer[k] = v
		}
	}

	newRequest.Proto = req.Proto
	newRequest.ProtoMajor = req.ProtoMajor
	newRequest.ProtoMinor = req.ProtoMinor
	if req.PostForm != nil {
		newRequest.PostForm = make(url.Values)
		for k, v := range req.PostForm {
			newRequest.PostForm[k] = v
		}
	}
	if req.Form != nil {
		newRequest.Form = make(url.Values)
		for k, v := range req.Form {
			newRequest.Form[k] = v
		}
	}
	if req.MultipartForm != nil {
		newRequest.MultipartForm = new(multipart.Form)
		if req.MultipartForm.Value != nil {
			newRequest.MultipartForm.Value = make(map[string][]string)
			for k, v := range req.MultipartForm.Value {
				newRequest.MultipartForm.Value[k] = v
			}
		}

		if req.MultipartForm.File != nil {
			newRequest.MultipartForm.File = make(map[string][]*multipart.FileHeader)
			for k, v := range req.MultipartForm.File {
				fileHeaders := make([]*multipart.FileHeader, len(v))
				for i, header := range v {
					newHeader := new(multipart.FileHeader)
					*newHeader = *header
					fileHeaders[i] = newHeader
				}
				newRequest.MultipartForm.File[k] = fileHeaders
			}
		}
	}

	newRequest.TransferEncoding = make([]string, len(req.TransferEncoding))
	copy(newRequest.TransferEncoding, req.TransferEncoding)

	body := httpbody.NewHTTPBody(req.Body, newRequest.Header)
	if req.Body != nil {
		req.Body = body.Reader()
	}

	return &HTTPRequest{
		Hit:     hit,
		Request: newRequest.WithContext(req.Context()),
		body:    body,
	}
}

// Body provides methods for accessing the http body.
func (req *HTTPRequest) Body() *httpbody.HTTPBody {
	return req.body
}
