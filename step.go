package hit

import (
	"bytes"
	"fmt"

	"golang.org/x/xerrors"

	"github.com/Eun/go-hit/errortrace"
)

// StepTime defines when a step should be run.
type StepTime uint8

const (
	// combineStep is a special step that runs before everything else and is used exclusively for the function
	// CombineSteps().
	combineStep StepTime = iota + 1

	// cleanStep is a step that runs during the clean step phase.
	cleanStep

	// requestCreateStep will be run before the actual http request will be created.
	requestCreateStep

	// BeforeSendStep runs before the Send() steps.
	BeforeSendStep

	// SendStep runs during the Send() steps.
	SendStep

	// AfterSendStep runs after the Send() steps, note that this is still before the actual sending process.
	AfterSendStep

	// BeforeExpectStep runs before the Expect() steps (this is after we got the data from the server).
	BeforeExpectStep

	// ExpectStep runs during the Expect() steps.
	ExpectStep

	// AfterExpectStep runs after the Expect() steps.
	AfterExpectStep
)

// String represents the string representation of StepTime.
func (s StepTime) String() string {
	switch s {
	case combineStep:
		return "combineStep"
	case cleanStep:
		return "cleanStep"
	case requestCreateStep:
		return "requestCreateStep"
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
	default:
		return ""
	}
}

// IStep defines a hit step.
type IStep interface {
	trace() *errortrace.ErrorTrace
	when() StepTime
	callPath() callPath
	exec(instance *hitImpl) error
}

type hitStep struct {
	Trace    *errortrace.ErrorTrace
	When     StepTime
	CallPath callPath
	Exec     func(hit *hitImpl) error
}

func (step *hitStep) trace() *errortrace.ErrorTrace {
	return step.Trace
}

func (step *hitStep) when() StepTime {
	return step.When
}

func (step *hitStep) callPath() callPath {
	return step.CallPath
}

func (step *hitStep) exec(hit *hitImpl) (err error) {
	if step.Exec == nil {
		return nil
	}
	return step.Exec(hit)
}

func execStep(hit *hitImpl, step IStep) (err *Error) {
	setError := func(r interface{}) {
		if r == nil {
			return
		}

		setMeta := func() {
			step.trace().SetDescription(hit.Description())
			var b bytes.Buffer
			if newDebug(step.callPath().Push("Debug", nil), &b).exec(hit) == nil {
				step.trace().SetContext(b.String())
			}
		}

		switch v := r.(type) {
		// this is already a hit.Error
		case *Error:
			err = v
		case error:
			step.trace().SetError(v)
			setMeta()
			err = &Error{
				callPath: step.callPath(),
				et:       step.trace(),
			}
		default:
			step.trace().SetError(xerrors.New(fmt.Sprint(r)))
			setMeta()
			err = &Error{
				callPath: step.callPath(),
				et:       step.trace(),
			}
		}
	}

	defer func() {
		setError(recover())
	}()
	setError(step.exec(hit))
	return err
}

// StepCallPath returns the representation of the step that is passed in.
func StepCallPath(step IStep, withArguments bool) string {
	return step.callPath().CallString(withArguments)
}
