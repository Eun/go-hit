package hit

import (
	"github.com/Eun/go-hit/httpbody"
	"github.com/Eun/go-hit/internal/converter"
)

// ISendFormValues provides methods to send form values.
type ISendFormValues interface {
	// Add adds the specified value to the specified form value.
	//
	// Usage:
	//     Send().Body().FormValues("username").Add("admin")
	Add(value ...interface{}) IStep
}

type sendFormValuesValueCallback func(hit Hit) *httpbody.URLValues

type sendFormValues struct {
	cleanPath     callPath
	valueCallback sendFormValuesValueCallback
	name          string
}

func newSendFormValues(clearPath callPath, valueCallback sendFormValuesValueCallback, name string) ISendFormValues {
	return &sendFormValues{
		cleanPath:     clearPath,
		valueCallback: valueCallback,
		name:          name,
	}
}

func (hdr *sendFormValues) Add(values ...interface{}) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     SendStep,
		CallPath: hdr.cleanPath.Push("Add", values),
		Exec: func(hit *hitImpl) error {
			for _, value := range values {
				var s string
				if err := converter.Convert(value, &s); err != nil {
					return err
				}
				hdr.valueCallback(hit).Add(hdr.name, s)
			}
			return nil
		},
	}
}
