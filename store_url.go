package hit

import "github.com/Eun/go-hit/internal/converter"

// IStoreURL defines the functions that can be used to store a URL part.
type IStoreURL interface {
	IStoreStep
	// Scheme stores the URLs scheme
	Scheme() IStoreStep

	// Scheme stores the URLs opaque status
	Opaque() IStoreStep

	// Scheme stores the URLs UserInfo
	User() IStoreUserInfo

	// Host stores the URLs host
	Host() IStoreStep

	// Hostname stores the URLs host, stripping any valid port number if present.
	//
	// If the result is enclosed in square brackets, as literal IPv6 addresses are,
	// the square brackets are removed from the result.
	Hostname() IStoreStep

	// Port stores the port part of the URLs host, without the leading colon.
	//
	// If URLs Host doesn't contain a valid numeric port, Port stores an empty string.
	Port() IStoreStep

	// Path stores the URLs path.
	Path() IStoreStep

	// EscapedPath stores the URLs path.
	EscapedPath() IStoreStep

	// RawPath stores the URLs RawPath value.
	RawPath() IStoreStep

	// Query stores the URLs Query value.
	//
	// Usage:
	//     var values url.Values
	//     Store().Request().URL().Query().In(&values)
	//     var user string
	//     Store().Request().URL().Query("user").In(&user)
	Query(name ...string) IStoreStep

	// ForceQuery stores the URLs ForceQuery value.
	ForceQuery() IStoreStep

	// RawQuery stores the URLs RawQuery value.
	RawQuery() IStoreStep

	// Fragment stores the URLs Fragment value.
	Fragment() IStoreStep

	// IsAbs stores URLs IsAbs value.
	IsAbs() IStoreStep

	// RequestURI stores the URLs RequestURI value.
	RequestURI() IStoreStep

	// String stores the URLs String value.
	String() IStoreStep
}

type storeURL struct{}

func newStoreURL() IStoreURL {
	return &storeURL{}
}

func (s *storeURL) In(v interface{}) IStep {
	return newStoreInStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL, v)
	}, v)
}

func (s *storeURL) Scheme() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.Scheme, v)
	})
}

func (s *storeURL) Opaque() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.Opaque, v)
	})
}

func (s *storeURL) User() IStoreUserInfo {
	return newStoreUserInfo()
}

func (s *storeURL) Host() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.Host, v)
	})
}

func (s *storeURL) Hostname() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.Hostname(), v)
	})
}

func (s *storeURL) Port() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.Port(), v)
	})
}

func (s *storeURL) Path() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.Path, v)
	})
}

func (s *storeURL) EscapedPath() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.EscapedPath(), v)
	})
}

func (s *storeURL) RawPath() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.RawPath, v)
	})
}

func (s *storeURL) Query(name ...string) IStoreStep {
	if header, ok := getLastStringArgument(name); ok {
		return newStoreStep(func(hit Hit, v interface{}) error {
			return storeStringSlice(hit.Request().URL.Query()[header], v)
		})
	}
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.Query(), v)
	})
}

func (s *storeURL) ForceQuery() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.ForceQuery, v)
	})
}

func (s *storeURL) RawQuery() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.RawQuery, v)
	})
}

func (s *storeURL) Fragment() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.Fragment, v)
	})
}

func (s *storeURL) IsAbs() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.IsAbs(), v)
	})
}

func (s *storeURL) RequestURI() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.RequestURI(), v)
	})
}

func (s *storeURL) String() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.String(), v)
	})
}
