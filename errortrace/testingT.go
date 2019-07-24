package errortrace

type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
}
