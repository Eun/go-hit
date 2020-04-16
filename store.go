package hit

// IStore provides a store functionality for the Request and Response
type IStore interface {
	// Request prints the Request
	//
	// Usage:
	//     Store().Request().Body().In(&body)                        // store the body
	//     Store().Request().Body().JSON("data").In(&data)           // parse body as json and store the data object into data variable
	//     Store().Request().Header().In(&headers)                   // store all headers
	//     Store().Request().Header("Content-Type").In(&contentType) // store the Content-Type header
	Request() IStoreRequest

	// Request prints the Response
	//
	// Usage:
	//     Store().Response().Body().In(&body)                        // store the body
	//     Store().Response().Body().JSON("data").In(&data)           // parse body as json and store the data object into data variable
	//     Store().Response().Header().In(&headers)                   // store all headers
	//     Store().Response().Header("Content-Type").In(&contentType) // store the Content-Type header
	Response() IStoreResponse
}

type store struct {
}

func newStore() IStore {
	return &store{}
}

func (*store) Request() IStoreRequest {
	return newStoreRequest()
}

func (*store) Response() IStoreResponse {
	return newStoreResponse()
}

type storeFunc func(Hit) (interface{}, error)

type IStoreStep interface {
	In(interface{}) IStep
}

func newStoreStep(f storeFunc) IStoreStep {
	return &storeStep{
		f,
	}
}

type storeStep struct {
	f storeFunc
}

func (s *storeStep) In(v interface{}) IStep {
	return newStoreInStep(s.f, v)
}

func newStoreInStep(f storeFunc, v interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      AfterExpectStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			res, err := f(hit)
			if err != nil {
				return err
			}
			return converter.Convert(res, v)
		},
	}
}
