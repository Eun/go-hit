package hit

import (
	"io"

	"github.com/Eun/go-hit/expr"
	"github.com/Eun/go-hit/internal/misc"
)

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
	TransferEncoding(expression ...string) IStoreStep

	// Header stores the Response's Header(s)
	//
	// If you omit the argument it will store all the headers
	//
	// Usage:
	//     var headers http.Header
	//     Store().Response().Header().In(&headers)
	//     var contentType string
	//     Store().Response().Header("Content-Type").In(&contentType)
	Header(headerName ...string) IStoreStep

	// Trailer stores the Response's Trailer(s)
	//
	// If you omit the argument it will store all the trailers
	//
	// Usage:
	//     var trailers http.Header
	//     Store().Response().Trailer().In(&trailers)
	//     var contentType string
	//     Store().Response().Trailer("Content-Type").In(&contentType)
	Trailer(trailerName ...string) IStoreStep

	// Body stores the Response's Body
	//
	// given the following body: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	// Usage:
	//     var body string
	//     Store().Response().Body().In(&body) // store the whole body as string
	//     var name string
	//     Store().Response().Body().JSON("Name").In(&name) // store "Joe" in name
	Body(expression ...string) IStoreBody

	// Uncompressed stores the Response's Uncompressed status
	Uncompressed() IStoreStep
}

type storeResponse struct{}

func newStoreResponse() *storeResponse {
	return &storeResponse{}
}

func (d *storeResponse) Status() IStoreStep {
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Response().Status, nil
	})
}

func (d *storeResponse) StatusCode() IStoreStep {
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Response().StatusCode, nil
	})
}

func (d *storeResponse) Proto() IStoreStep {
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Response().Proto, nil
	})
}

func (d *storeResponse) ProtoMajor() IStoreStep {
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Response().ProtoMajor, nil
	})
}

func (d *storeResponse) ProtoMinor() IStoreStep {
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Response().ProtoMinor, nil
	})
}

func (d *storeResponse) ContentLength() IStoreStep {
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Response().ContentLength, nil
	})
}

func (d *storeResponse) TransferEncoding(expression ...string) IStoreStep {
	if e, ok := misc.GetLastStringArgument(expression); ok {
		return newStoreStep(func(hit Hit) (interface{}, error) {
			v, _, err := expr.GetValue(hit.Response().TransferEncoding, e, expr.IgnoreCase)
			return v, err
		})
	}
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Response().TransferEncoding, nil
	})
}

func (d *storeResponse) Header(headerName ...string) IStoreStep {
	if header, ok := misc.GetLastStringArgument(headerName); ok {
		return newStoreStep(func(hit Hit) (interface{}, error) {
			return hit.Response().Header.Get(header), nil
		})
	}
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Response().Header, nil
	})
}

func (d *storeResponse) Trailer(trailerName ...string) IStoreStep {
	if trailer, ok := misc.GetLastStringArgument(trailerName); ok {
		return newStoreStep(func(hit Hit) (interface{}, error) {
			// we have to read the body to get the trailers
			_, _ = io.Copy(misc.DevNullWriter(), hit.Response().Body().Reader())
			return hit.Response().Trailer.Get(trailer), nil
		})
	}
	return newStoreStep(func(hit Hit) (interface{}, error) {
		// we have to read the body to get the trailers
		_, _ = io.Copy(misc.DevNullWriter(), hit.Response().Body().Reader())
		return hit.Response().Trailer, nil
	})
}

func (d *storeResponse) Body(expression ...string) IStoreBody {
	return newStoreBody(storeBodyResponse, expression)
}

func (d *storeResponse) Uncompressed() IStoreStep {
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Response().Uncompressed, nil
	})
}
