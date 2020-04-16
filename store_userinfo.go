package hit

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
	return newStoreInStep(func(hit Hit) (interface{}, error) {
		return hit.Request().URL.User, nil
	}, v)
}

func (s *storeUserInfo) Username() IStoreStep {
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Request().URL.User.Username(), nil
	})
}

func (s *storeUserInfo) Password() IStoreStep {
	return newStoreStep(func(hit Hit) (interface{}, error) {
		pass, _ := hit.Request().URL.User.Password()
		return pass, nil
	})
}

func (s *storeUserInfo) String() IStoreStep {
	return newStoreStep(func(hit Hit) (interface{}, error) {
		return hit.Request().URL.User.String(), nil
	})
}
