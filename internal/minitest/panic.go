package minitest

type PanicNow struct{}

func (PanicNow) FailNow(err error, customMessageAndArgs ...interface{}) {
	if err != nil {
		panic(makeError(err.Error(), customMessageAndArgs...).Error())
	}
	panic(makeError("", customMessageAndArgs...).Error())
}

func (PanicNow) Errorf(messageAndArgs ...interface{}) {
	if err := Error.Errorf(messageAndArgs...); err != nil {
		panic(err.Error())
	}
}

func (PanicNow) NoError(err error, customMessageAndArgs ...interface{}) {
	if err := Error.NoError(err, customMessageAndArgs...); err != nil {
		panic(err.Error())
	}
}

func (PanicNow) Equal(expected, actual interface{}, customMessageAndArgs ...interface{}) {
	if err := Error.Equal(expected, actual, customMessageAndArgs...); err != nil {
		panic(err.Error())
	}
}

func (PanicNow) NotEqual(expected, actual interface{}, customMessageAndArgs ...interface{}) {
	if err := Error.NotEqual(expected, actual, customMessageAndArgs...); err != nil {
		panic(err.Error())
	}
}

func (PanicNow) Contains(object interface{}, value interface{}, customMessageAndArgs ...interface{}) {
	if err := Error.Contains(object, value, customMessageAndArgs...); err != nil {
		panic(err.Error())
	}
}

func (PanicNow) NotContains(object interface{}, value interface{}, customMessageAndArgs ...interface{}) {
	if err := Error.NotContains(object, value, customMessageAndArgs...); err != nil {
		panic(err.Error())
	}
}

func (PanicNow) Empty(object interface{}, customMessageAndArgs ...interface{}) {
	if err := Error.Empty(object, customMessageAndArgs...); err != nil {
		panic(err.Error())
	}
}

func (PanicNow) Len(object interface{}, length int, customMessageAndArgs ...interface{}) {
	if err := Error.Len(object, length, customMessageAndArgs...); err != nil {
		panic(err.Error())
	}
}

func (PanicNow) True(value bool, customMessageAndArgs ...interface{}) {
	if err := Error.True(value, customMessageAndArgs...); err != nil {
		panic(err.Error())
	}
}

func (PanicNow) False(value bool, customMessageAndArgs ...interface{}) {
	if err := Error.False(value, customMessageAndArgs...); err != nil {
		panic(err.Error())
	}
}
