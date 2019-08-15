package hit

type StepTime uint8

const (
	BeforeSendStep StepTime = iota + 1
	SendStep
	AfterSendStep
	BeforeExpectStep
	ExpectStep
	AfterExpectStep
	CleanStep StepTime = 0x80
)

type IStep interface {
	exec(Hit)
	when() StepTime
}

type hitStep struct {
	e func(Hit)
	w StepTime
}

func (c hitStep) exec(h Hit) {
	c.e(h)
}

func (c hitStep) when() StepTime {
	return c.w
}

// Custom calls a custom Step on the specified execution time
func Custom(when StepTime, exec func(Hit)) IStep {
	return hitStep{
		w: when,
		e: exec,
	}
}
