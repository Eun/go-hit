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

type Step interface {
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

func MakeStep(when StepTime, exec func(Hit)) Step {
	return hitStep{
		w: when,
		e: exec,
	}
}
