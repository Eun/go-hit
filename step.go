package hit

import (
	"fmt"

	"github.com/Eun/go-hit/errortrace"
)

type StepTime uint8

const (
	CombineStep StepTime = iota + 1
	CleanStep
	BeforeSendStep
	SendStep
	AfterSendStep
	BeforeExpectStep
	ExpectStep
	AfterExpectStep
)

func (s StepTime) String() string {
	switch s {
	case CleanStep:
		return "CleanStep"
	case CombineStep:
		return "CombineStep"
	case BeforeSendStep:
		return "BeforeSendStep"
	case SendStep:
		return "SendStep"
	case AfterSendStep:
		return "AfterSendStep"
	case BeforeExpectStep:
		return "BeforeExpectStep"
	case ExpectStep:
		return "ExpectStep"
	case AfterExpectStep:
		return "AfterExpectStep"
	}
	return ""
}

type IStep interface {
	when() StepTime
	clearPath() clearPath
	exec(Hit) error
}

type hitStep struct {
	Trace     *errortrace.ErrorTrace
	When      StepTime
	ClearPath clearPath
	Exec      func(hit Hit) error
}

func (step *hitStep) exec(h Hit) (err error) {
	if step.Exec == nil {
		return nil
	}
	defer func() {
		r := recover()
		if r != nil {
			var ok bool
			err, ok = r.(errortrace.ErrorTraceError)
			if !ok {
				err = step.Trace.Format(h.Description(), fmt.Sprint(r))
			}
		}
	}()
	err = step.Exec(h)
	return err
}

func (step *hitStep) when() StepTime {
	return step.When
}

func (step *hitStep) clearPath() clearPath {
	return step.ClearPath
}
