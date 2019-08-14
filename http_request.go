package hit

import (
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/Eun/go-doppelgangerreader"
)

type HTTPRequest struct {
	Hit Hit
	*http.Request
	body *HTTPBody
}

func newHTTPRequest(hit Hit, req *http.Request) *HTTPRequest {
	u := *req.URL

	newRequest := &http.Request{
		Method: req.Method,
		URL:    &u,
	}
	// copy headers
	if req.Header != nil {
		newRequest.Header = make(http.Header)
		for k, v := range req.Header {
			newRequest.Header[k] = v
		}
	}
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
	for i, v := range req.TransferEncoding {
		newRequest.TransferEncoding[i] = v
	}

	var factory *doppelgangerreader.DoppelgangerFactory
	if req.Body != nil {
		factory = doppelgangerreader.NewFactory(req.Body)
		req.Body = factory.NewDoppelganger()
	}

	return &HTTPRequest{
		Hit:     hit,
		Request: newRequest,
		body: &HTTPBody{
			hit:     hit,
			factory: factory,
		},
	}
}

func (req *HTTPRequest) Body() *HTTPBody {
	return req.body
}
