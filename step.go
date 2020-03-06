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
	Exec(Hit) error
	When() StepTime
	CleanPath() CleanPath
}

type Step struct {
	When      StepTime
	CleanPath CleanPath
	Instance  Hit
	Exec      Callback
}

type hitStep struct {
	et   *errortrace.ErrorTrace
	call Callback
	w    StepTime
	path CleanPath
}

func (step *hitStep) Exec(h Hit) (err error) {
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

func (step *hitStep) When() StepTime {
	return step.w
}

func (step *hitStep) CleanPath() CleanPath {
	return step.path
}

func Custom(when StepTime, exec Callback) IStep {
	return custom(Step{
		When:      when,
		CleanPath: NewCleanPath("Custom"),
		Instance:  nil,
		Exec:      exec,
	})
}

func custom(step Step) IStep {
	s := &hitStep{
		et:   errortrace.Prepare(),
		w:    step.When,
		path: step.CleanPath,
		call: step.Exec,
	}

	hit := step.Instance

	if hit == nil {
		hit = getContext()
	}

	// if we found no instance we are running as a Main Step, e.g. Do(SomeStep)
	// this means we do not execute anything, we just append to the chain
	if hit == nil {
		return s
	}

	// we have an instance, this means we are running inside a function, e.g.Do(Custom(func(Hit){SomeStep}))
	_ = s.Exec(hit)
	return s
}
