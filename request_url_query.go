package hit

import (
	"net/url"

	"github.com/Eun/go-hit/internal/converter"
)

// IRequestURLQuery provides methods to send header/trailer.
type IRequestURLQuery interface {
	// Add adds the specified value to the url query.
	//
	// Usage:
	//     Request().URL().Query("page").Add(1)
	Add(value ...interface{}) IStep
}

type requestURLQueryValueCallback func(hit Hit) (*string, url.Values)

type requestURLQuery struct {
	cleanPath     callPath
	valueCallback requestURLQueryValueCallback
	name          string
}

func newRequestURLQuery(clearPath callPath, valueCallback requestURLQueryValueCallback, name string) IRequestURLQuery {
	return &requestURLQuery{
		cleanPath:     clearPath,
		valueCallback: valueCallback,
		name:          name,
	}
}

func (v *requestURLQuery) Add(values ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: v.cleanPath.Push("Add", values),
		Exec: func(hit *hitImpl) error {
			for _, value := range values {
				var s string
				if err := converter.Convert(value, &s); err != nil {
					return err
				}
				rawQuery, queryValues := v.valueCallback(hit)
				queryValues.Add(v.name, s)
				*rawQuery = queryValues.Encode()
			}
			return nil
		},
	}
}
