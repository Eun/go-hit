package hit

import (
	"fmt"

	"github.com/Eun/go-hit/errortrace"
)

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
	exec(Hit) error
	when() StepTime
}

type hitStep struct {
	call Callback
	et   *errortrace.ErrorTrace
	w    StepTime
}

func (step hitStep) exec(h Hit) (err error) {
	if step.call == nil {
		return nil
	}
	defer func() {
		r := recover()
		if r != nil {
			err = step.et.Format(h.Description(), fmt.Sprint(r))
		}
	}()
	step.call(h)
	return err
}

func (step hitStep) when() StepTime {
	return step.w
}

// Custom calls a custom Step on the specified execution time
func Custom(when StepTime, exec Callback) IStep {
	return hitStep{
		et:   errortrace.Prepare(),
		w:    when,
		call: exec,
	}
}
