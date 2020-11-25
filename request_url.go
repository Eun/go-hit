package hit

import (
	urlpkg "net/url"
)

// IRequestURL provides methods to set request url parameters.
type IRequestURL interface {
	// Set sets the request url to the specified value.
	Set(url *urlpkg.URL) IStep

	// Scheme sets the request url scheme to the specified value.
	//
	// Usage:
	//     Request().URL().Scheme("https")
	//
	// Example:
	//     MustDo(
	//         Request().Method(http.MethodGet),
	//         Request().URL().Scheme("https"),
	//         Request().URL().Host("example.com"),
	//         Request().URL().Path("/index.html"),
	//     )
	Scheme(scheme string) IStep

	// Host sets the request url host to the specified value.
	//
	// Usage:
	//     Request().URL().Host("example.com")
	//
	// Example:
	//     MustDo(
	//         Request().Method(http.MethodGet),
	//         Request().URL().Host("example.com"),
	//         Request().URL().Path("/index.html"),
	//     )
	Host(host string) IStep

	// Path sets the request url path to the specified value.
	//
	// Usage:
	//     Request().URL().Path("/index.html")
	//
	// Example:
	//     MustDo(
	//         Request().Method(http.MethodGet),
	//         Request().URL().Host("example.com"),
	//         Request().URL().Path("/index.html"),
	//     )
	Path(path string) IStep

	// RawPath sets the request url raw path to the specified value.
	//
	// Usage:
	//     Request().URL().RawPath("/index.html")
	//
	// Example:
	//     MustDo(
	//         Request().Method(http.MethodGet),
	//         Request().URL().Host("example.com"),
	//         Request().URL().RawPath("/index.html"),
	//     )
	RawPath(rawPath string) IStep

	// ForceQuery sets the request url force query to the specified value.
	//
	// Usage:
	//     Request().URL().ForceQuery(true)
	//
	// Example:
	//     MustDo(
	//         Request().Method(http.MethodGet),
	//         Request().URL().Host("example.com"),
	//         Request().URL().Path("/index.html"),
	//         Request().URL().ForceQuery(true),
	//     )
	ForceQuery(forceQuery bool) IStep

	// Query sets the request url query parameters.
	//
	// Usage:
	//     Request().URL().Query("page").Add(1)
	//
	// Example:
	//     MustDo(
	//         Request().Method(http.MethodGet),
	//         Request().URL().Host("example.com"),
	//         Request().URL().Path("/index.html"),
	//         Request().URL().Query("page").Add(1),
	//     )
	Query(name string) IRequestURLQuery

	// RawQuery sets the request url force query to the specified value.
	//
	// Usage:
	//     Request().URL().RawQuery("x=1&y=2")
	//
	// Example:
	//     MustDo(
	//         Request().Method(http.MethodGet),
	//         Request().URL().Host("example.com"),
	//         Request().URL().Path("/index.html"),
	//         Request().URL().RawQuery("x=1&y=2"),
	//     )
	RawQuery(rawQuery string) IStep

	// Fragment sets the request url force query to the specified value.
	//
	// Usage:
	//     Request().URL().Fragment("anchor")
	//
	// Example:
	//     MustDo(
	//         Request().Method(http.MethodGet),
	//         Request().URL().Host("example.com"),
	//         Request().URL().Path("/index.html"),
	//         Request().URL().Fragment("anchor"),
	//     )
	Fragment(fragment string) IStep

	// User sets the request url user and/or password information.
	//
	// Usage:
	//     Request().URL().User().Username("joe")
	//     Request().URL().User().Password("secret")
	//
	// Example:
	//     MustDo(
	//         Request().Method(http.MethodGet),
	//         Request().URL().Host("example.com"),
	//         Request().URL().Path("/index.html"),
	//         Request().URL().User().Username("joe"),
	//         Request().URL().User().Password("secret"),
	//     )
	User() IRequestURLUser
}

type requestURL struct {
	cleanPath callPath
}

func newRequestURL(clearPath callPath) IRequestURL {
	return &requestURL{
		cleanPath: clearPath,
	}
}

func (u *requestURL) Set(url *urlpkg.URL) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: u.cleanPath.Push("Set", []interface{}{url}),
		Exec: func(hit *hitImpl) error {
			hit.request.URL = url
			return nil
		},
	}
}
func (u *requestURL) Scheme(scheme string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: u.cleanPath.Push("Scheme", []interface{}{scheme}),
		Exec: func(hit *hitImpl) error {
			hit.request.URL.Scheme = scheme
			return nil
		},
	}
}

func (u *requestURL) Host(host string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: u.cleanPath.Push("Host", []interface{}{host}),
		Exec: func(hit *hitImpl) error {
			hit.request.URL.Host = host
			return nil
		},
	}
}

func (u *requestURL) Path(path string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: u.cleanPath.Push("Path", []interface{}{path}),
		Exec: func(hit *hitImpl) error {
			hit.request.URL.Path = path
			return nil
		},
	}
}

func (u *requestURL) RawPath(rawPath string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: u.cleanPath.Push("RawPath", []interface{}{rawPath}),
		Exec: func(hit *hitImpl) error {
			var err error
			hit.request.URL.Path, err = urlpkg.PathUnescape(rawPath)
			if err != nil {
				return err
			}
			hit.request.URL.RawPath = rawPath
			return nil
		},
	}
}

func (u *requestURL) ForceQuery(forceQuery bool) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: u.cleanPath.Push("ForceQuery", []interface{}{forceQuery}),
		Exec: func(hit *hitImpl) error {
			hit.request.URL.ForceQuery = forceQuery
			return nil
		},
	}
}

func (u *requestURL) Query(name string) IRequestURLQuery {
	return newRequestURLQuery(u.cleanPath.Push("Query", []interface{}{name}), func(hit Hit) (*string, urlpkg.Values) {
		return &hit.Request().URL.RawQuery, hit.Request().URL.Query()
	}, name)
}

func (u *requestURL) RawQuery(rawQuery string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: u.cleanPath.Push("RawQuery", []interface{}{rawQuery}),
		Exec: func(hit *hitImpl) error {
			hit.request.URL.RawQuery = rawQuery
			return nil
		},
	}
}

func (u *requestURL) Fragment(fragment string) IStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     requestCreateStep,
		CallPath: u.cleanPath.Push("Fragment", []interface{}{fragment}),
		Exec: func(hit *hitImpl) error {
			hit.request.URL.Fragment = fragment
			return nil
		},
	}
}

func (u *requestURL) User() IRequestURLUser {
	return newRequestURLUser(u.cleanPath.Push("User", nil))
}
