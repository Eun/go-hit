package hit

import "github.com/Eun/go-hit/internal/converter"

type IStoreURL interface {
	IStoreStep
	// Scheme stores the URL's scheme
	Scheme() IStoreStep

	// Scheme stores the URL's opaque status
	Opaque() IStoreStep

	// Scheme stores the URL's UserInfo
	User() IStoreUserInfo

	// Host stores the URL's host
	Host() IStoreStep

	// Path stores the URL's path
	Path() IStoreStep

	// RawPath stores the URL's RawPath value
	RawPath() IStoreStep

	// ForceQuery stores the URL's ForceQuery value
	ForceQuery() IStoreStep

	// RawQuery stores the URL's RawQuery value
	RawQuery() IStoreStep

	// Fragment stores the URL's fragment value
	Fragment() IStoreStep
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

func (s *storeURL) Path() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.Path, v)
	})
}

func (s *storeURL) RawPath() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.RawPath, v)
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
