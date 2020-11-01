package hit

import (
	"io"
	"io/ioutil"

	"github.com/Eun/go-hit/internal/converter"
)

// IStoreResponse defines the functions that can be used to store data from the http response.
type IStoreResponse interface {
	// Status stores the Response's Status
	Status() IStoreStep

	// Status stores the Response's StatusCode
	StatusCode() IStoreStep

	// Proto stores the Response's Proto
	Proto() IStoreStep

	// ProtoMajor stores the Response's ProtoMajor
	ProtoMajor() IStoreStep

	// ProtoMinor stores the Response's ProtoMinor
	ProtoMinor() IStoreStep

	// ContentLength stores the Response's ContentLength
	ContentLength() IStoreStep

	// TransferEncoding stores the Response's TransferEncoding
	//
	// The argument can be used to narrow down the store path
	TransferEncoding() IStoreStep

	// Header stores the Response's Header(s)
	//
	// If you specify the argument you can directly store the header value
	//
	// Usage:
	//     var headers http.Header
	//     Store().Response().Headers().In(&headers)
	//
	//     var contentType string
	//     Store().Response().Headers("Content-Type").In(&contentType)
	Headers(headerName ...string) IStoreStep

	// Trailer stores the Response's Trailer(s)
	//
	// If you specify the argument you can directly store the trailer value
	//
	// Usage:
	//     var trailers http.Header
	//     Store().Response().Trailers().In(&trailers)
	//
	//     var contentType string
	//     Store().Response().Trailers("Content-Type").In(&contentType)
	Trailers(trailerName ...string) IStoreStep

	// Body stores the Response's Body
	//
	// given the following body: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	// Usage:
	//     var body string
	//     Store().Response().Body().String().In(&body) // store the whole body as string
	//     var name string
	//     Store().Response().Body().JSON().JQ(".Name").In(&name) // store "Joe" in name
	Body() IStoreBody

	// Uncompressed stores the Response's Uncompressed status
	Uncompressed() IStoreStep
}

type storeResponse struct{}

func newStoreResponse() *storeResponse {
	return &storeResponse{}
}

func (d *storeResponse) Status() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Response().Status, v)
	})
}

func (d *storeResponse) StatusCode() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Response().StatusCode, v)
	})
}

func (d *storeResponse) Proto() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Response().Proto, v)
	})
}

func (d *storeResponse) ProtoMajor() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Response().ProtoMajor, v)
	})
}

func (d *storeResponse) ProtoMinor() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Response().ProtoMinor, v)
	})
}

func (d *storeResponse) ContentLength() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Response().ContentLength, v)
	})
}

func (d *storeResponse) TransferEncoding() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Response().TransferEncoding, v)
	})
}

func (d *storeResponse) Headers(headerName ...string) IStoreStep {
	if header, ok := getLastStringArgument(headerName); ok {
		return newStoreStep(func(hit Hit, v interface{}) error {
			return storeStringSlice(hit.Response().Header.Values(header), v)
		})
	}
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Response().Header, v)
	})
}

func (d *storeResponse) Trailers(trailerName ...string) IStoreStep {
	if trailer, ok := getLastStringArgument(trailerName); ok {
		return newStoreStep(func(hit Hit, v interface{}) error {
			// we have to read the body to get the trailers
			_, _ = io.Copy(ioutil.Discard, hit.Response().Body().Reader())
			return storeStringSlice(hit.Response().Trailer.Values(trailer), v)
		})
	}
	return newStoreStep(func(hit Hit, v interface{}) error {
		// we have to read the body to get the trailers
		_, _ = io.Copy(ioutil.Discard, hit.Response().Body().Reader())
		return converter.Convert(hit.Response().Trailer, v)
	})
}

func (d *storeResponse) Body() IStoreBody {
	return newStoreBody(storeBodyResponse)
}

func (d *storeResponse) Uncompressed() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Response().Uncompressed, v)
	})
}
