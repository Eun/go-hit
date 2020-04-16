package hit

import (
	"github.com/Eun/go-hit/expr"
	"github.com/Eun/go-hit/internal/misc"
)

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
	TransferEncoding(expression ...string) IStoreStep

	// Host stores the Request's Host
	Host() IStoreStep

	// Header stores the Request's Header(s)
	//
	// If you omit the argument it will store all the headers
	//
	// Usage:
	//     var headers http.Header
	//     Store().Request().Header().In(&headers)
	//     var contentType string
	//     Store().Request().Header("Content-Type").In(&contentType)
	Header(headerName ...string) IStoreStep

	// Trailer stores the Request's Trailer(s)
	//
	// If you omit the argument it will store all the trailers
	//
	// Usage:
	//     var trailers http.Header
	//     Store().Request().Trailer().In(&trailers)
	//     var contentType string
	//     Store().Request().Trailer("Content-Type").In(&contentType)
	Trailer(trailerName ...string) IStoreStep

	// Body stores the Request's Body
	//
	// given the following body: { "ID": 10, "Name": "Joe", "Roles": ["Admin", "User"] }
	// Usage:
	//     var body string
	//     Store().Request().Body().In(&body) // store the whole body as string
	//     var name string
	//     Store().Request().Body().JSON("Name").In(&name) // store "Joe" in name
	Body(expression ...string) IStoreBody
}

type storeRequest struct{}

func newStoreRequest() *storeRequest {
	return &storeRequest{}
}

func (d *storeRequest) Method() IStoreStep {
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Request().Method, nil
	})
}

func (d *storeRequest) URL() IStoreURL {
	return newStoreURL()
}

func (d *storeRequest) Proto() IStoreStep {
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Request().Proto, nil
	})
}

func (d *storeRequest) ProtoMajor() IStoreStep {
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Request().ProtoMajor, nil
	})
}

func (d *storeRequest) ProtoMinor() IStoreStep {
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Request().ProtoMinor, nil
	})
}

func (d *storeRequest) ContentLength() IStoreStep {
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Request().ContentLength, nil
	})
}

func (d *storeRequest) TransferEncoding(expression ...string) IStoreStep {
	if e, ok := misc.GetLastStringArgument(expression); ok {
		return newStoreStep(func(hit Hit) (interface{}, error) {
			v, _, err := expr.GetValue(hit.Request().TransferEncoding, e, expr.IgnoreCase)
			return v, err
		})
	}
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Request().TransferEncoding, nil
	})
}

func (d *storeRequest) Host() IStoreStep {
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Request().Host, nil
	})
}

func (d *storeRequest) Header(headerName ...string) IStoreStep {
	if header, ok := misc.GetLastStringArgument(headerName); ok {
		return newStoreStep(func(hit Hit) (interface{}, error) {
			return hit.Request().Header.Get(header), nil
		})
	}
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Request().Header, nil
	})
}

func (d *storeRequest) Trailer(trailerName ...string) IStoreStep {
	if trailer, ok := misc.GetLastStringArgument(trailerName); ok {
		return newStoreStep(func(hit Hit) (interface{}, error) {
			return hit.Request().Trailer.Get(trailer), nil
		})
	}
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Request().Trailer, nil
	})
}

func (d *storeRequest) Body(expression ...string) IStoreBody {
	return newStoreBody(storeBodyRequest, expression)
}
