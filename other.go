package hit

type PanicT struct{}

func (PanicT) Error(args ...interface{}) {
	panic(args)
}
