package hit

import (
	"github.com/Eun/go-hit/internal/converter"
)

// IStoreRequest defines the functions that can be used to store data from the http request.
type IStoreRequest interface {
	// Method stores the Request's Method
	Method() IStoreStep

	// URL stores the Request's URL
	//
	// The argument can be used to narrow down the store path
	//
	// Usage:
	//     var u url.URL
	//     Store().Request().URL().In(&u)
	//     var path string
	//     Store().Request().URL().Path().In(&path)
	URL() IStoreURL

	// Proto stores the Request's Proto
	Proto() IStoreStep

	// ProtoMajor stores the Request's ProtoMajor
	ProtoMajor() IStoreStep

	// ProtoMinor stores the Request's ProtoMinor
	ProtoMinor() IStoreStep

	// ContentLength stores the Request's ContentLength
	ContentLength() IStoreStep

	// TransferEncoding stores the Request's TransferEncoding
	//
	// The argument can be used to narrow down the store path
	TransferEncoding() IStoreStep

	// Host stores the Request's Host
	Host() IStoreStep

	// Headers stores the Request's Headers(s)
	//
	// If you specify the argument you can directly store the header value
	//
	// Usage:
	//     var headers http.Header
	//     Store().Request().Headers().In(&headers)
	//     var contentType string
	//     Store().Request().Headers("Content-Type").In(&contentType)
	Headers(headerName ...string) IStoreStep

	// Trailers stores the Request's Trailers(s)
	//
	// If you specify the argument you can directly store the trailer value
	//
	// Usage:
	//     var trailers http.Header
	//     Store().Request().Trailers().In(&trailers)
	//     var contentType string
	//     Store().Request().Trailers("Content-Type").In(&contentType)
	Trailers(trailerName ...string) IStoreStep

	// Body stores the Request's Body
	//
	// given the following body: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	// Usage:
	//     var body string
	//     Store().Request().Body().String().In(&body) // store the whole body as string
	//     var name string
	//     Store().Request().Body().JSON().JQ(".Name").In(&name) // store "Joe" in name
	Body() IStoreBody
}

type storeRequest struct{}

func newStoreRequest() *storeRequest {
	return &storeRequest{}
}

func (d *storeRequest) Method() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().Method, v)
	})
}

func (d *storeRequest) URL() IStoreURL {
	return newStoreURL()
}

func (d *storeRequest) Proto() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().Proto, v)
	})
}

func (d *storeRequest) ProtoMajor() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().ProtoMajor, v)
	})
}

func (d *storeRequest) ProtoMinor() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().ProtoMinor, v)
	})
}

func (d *storeRequest) ContentLength() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().ContentLength, v)
	})
}

func (d *storeRequest) TransferEncoding() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().TransferEncoding, v)
	})
}

func (d *storeRequest) Host() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().Host, v)
	})
}

func (d *storeRequest) Headers(headerName ...string) IStoreStep {
	if header, ok := getLastStringArgument(headerName); ok {
		return newStoreStep(func(hit Hit, v interface{}) error {
			return storeStringSlice(hit.Request().Header.Values(header), v)
		})
	}
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().Header, v)
	})
}

func (d *storeRequest) Trailers(trailerName ...string) IStoreStep {
	if trailer, ok := getLastStringArgument(trailerName); ok {
		return newStoreStep(func(hit Hit, v interface{}) error {
			return storeStringSlice(hit.Request().Trailer.Values(trailer), v)
		})
	}
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().Trailer, v)
	})
}

func (d *storeRequest) Body() IStoreBody {
	return newStoreBody(storeBodyRequest)
}
