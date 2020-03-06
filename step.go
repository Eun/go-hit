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
	CleanStep
)

type IStep interface {
	When() StepTime
	CleanPath() CleanPath
	Exec(Hit) error
}

type Step struct {
	When      StepTime
	CleanPath CleanPath
	Exec      Callback
}

type hitStep struct {
	et        *errortrace.ErrorTrace
	call      Callback
	w         StepTime
	cleanPath CleanPath
}

func (step *hitStep) Exec(h Hit) (err error) {
	if step.call == nil {
		return nil
	}
	defer func() {
		r := recover()
		if r != nil {
			var ok bool
			err, ok = r.(errortrace.ErrorTraceError)
			if !ok {
				err = step.et.Format(h.Description(), fmt.Sprint(r))
			}
		}
	}()
	step.call(h)
	return err
}

func (step *hitStep) When() StepTime {
	return step.w
}

func (step *hitStep) CleanPath() CleanPath {
	return step.cleanPath
}

func Custom(when StepTime, exec Callback) IStep {
	return custom(Step{
		When:      when,
		CleanPath: NewCleanPath("Custom", []interface{}{when, exec}),
		Exec:      exec,
	})
}

func custom(step Step) IStep {
	s := &hitStep{
		et:        errortrace.Prepare(),
		w:         step.When,
		cleanPath: step.CleanPath,
		call:      step.Exec,
	}

	// if we found no instance we are running as a Main Step, e.g. Do(SomeStep)
	// this means we do not execute anything, we just append to the chain
	if getContext() == nil {
		return s
	}

	// we have an instance, this means we are running inside a function, e.g.Do(Custom(func(Hit){SomeStep}))
	panic("This is unsupported, use RunSteps()")
	return s
}

type wrapStep struct {
	IStep
	CleanPath CleanPath
}

func wrap(step wrapStep) IStep {
	return &hitStep{
		et:        errortrace.Prepare(),
		w:         step.When(),
		cleanPath: step.CleanPath,
		call: func(hit Hit) {
			if err := step.Exec(hit); err != nil {
				panic(err)
			}
		},
	}
}
