package hit

import (
	"net/url"
)

// IRequestURLUser provides methods to set the username and/or password for the request.
type IRequestURLUser interface {
	// Username sets the username for the request.
	Username(name string) IStep

	// Password sets the password for the request.
	Password(password string) IStep
}

type requestURLUser struct {
	cleanPath callPath
}

func newRequestURLUser(clearPath callPath) IRequestURLUser {
	return &requestURLUser{
		cleanPath: clearPath,
	}
}

func (v *requestURLUser) Username(name string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: v.cleanPath.Push("Username", []interface{}{name}),
		Exec: func(hit *hitImpl) error {
			password, ok := hit.request.URL.User.Password()
			if ok {
				hit.request.URL.User = url.UserPassword(name, password)
				return nil
			}
			hit.request.URL.User = url.User(name)
			return nil
		},
	}
}

func (v *requestURLUser) Password(password string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: v.cleanPath.Push("Password", []interface{}{password}),
		Exec: func(hit *hitImpl) error {
			hit.request.URL.User = url.UserPassword(hit.request.URL.User.Username(), password)
			return nil
		},
	}
}
