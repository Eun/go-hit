package hit

import "github.com/Eun/go-hit/internal/converter"

// IStoreUserInfo defines the functions that can be used to store a user part.
type IStoreUserInfo interface {
	IStoreStep
	// Username stores the UserInfo's username
	Username() IStoreStep

	// Password stores the UserInfo's password
	Password() IStoreStep

	// String stores the string representation of UserInfo
	String() IStoreStep
}

type storeUserInfo struct{}

func newStoreUserInfo() IStoreUserInfo {
	return &storeUserInfo{}
}

func (s *storeUserInfo) In(v interface{}) IStep {
	return newStoreInStep(func(hit Hit, container interface{}) error {
		return converter.Convert(hit.Request().URL.User, container)
	}, v)
}

func (s *storeUserInfo) Username() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.User.Username(), v)
	})
}

func (s *storeUserInfo) Password() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		pass, _ := hit.Request().URL.User.Password()
		return converter.Convert(pass, v)
	})
}

func (s *storeUserInfo) String() IStoreStep {
	return newStoreStep(func(hit Hit, v interface{}) error {
		return converter.Convert(hit.Request().URL.User.String(), v)
	})
}
