package hit

import (
	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/internal"
	"golang.org/x/xerrors"
)

type IClearSendBody interface {
	IStep
	// JSON removes all previous Send().Body().JSON() steps.
	//
	// If you specify an argument it will only remove the Send().Body().JSON() steps matching that argument.
	//
	// Usage:
	//     Clear().Send().Body().JSON()                                              // will remove all Send().Body().JSON() steps
	//     Clear().Send().Body().JSON(map[string]interface{}{"Name": "Joe"})         // will remove all Send().Body().JSON(map[string]interface{}{"Name": "Joe"}) steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Send().Body().JSON(map[string]interface{}{"Name": "Joe"}),
	//         Clear().Send().Body().JSON(),
	//         Send().Body().JSON(map[string]interface{}{"Name": "Alice"}),
	//     )
	JSON(...interface{}) IStep

	// Interface removes all previous Send().Body().Interface() steps.
	//
	// If you specify an argument it will only remove the Send().Body().Interface() steps matching that argument.
	//
	// Usage:
	//     Clear().Send().Body().Interface()              // will remove all Send().Body().Interface() steps
	//     Clear().Send().Body().Interface("Hello World") // will remove all Send().Body().Interface("Hello World") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Send().Body().Interface("Hello Earth"),
	//         Clear().Send().Body().Interface(),
	//         Send().Body().Interface("Hello World"),
	//     )
	Interface(...interface{}) IStep
}

type clearSendBody struct {
	cleanPath clearPath
	trace     *errortrace.ErrorTrace
}

func newClearSendBody(clearPath clearPath, params []interface{}) IClearSendBody {
	if _, ok := internal.GetLastArgument(params); ok {
		// this runs if we called Clear().Send().Body(something)
		return finalClearSendBody{removeStep(clearPath)}
	}

	return &clearSendBody{
		cleanPath: clearPath,
		trace:     ett.Prepare(),
	}
}

func (*clearSendBody) when() StepTime {
	return SendStep
}

func (body *clearSendBody) exec(hit Hit) error {
	// this runs if we called Clear().Send().Body()
	if err := removeSteps(hit, body.clearPath()); err != nil {
		return body.trace.Format(hit.Description(), err.Error())
	}
	return nil
}

func (body *clearSendBody) clearPath() clearPath {
	return body.cleanPath
}

func (body *clearSendBody) JSON(data ...interface{}) IStep {
	return removeStep(body.cleanPath.Push("JSON", data))
}

func (body *clearSendBody) Interface(data ...interface{}) IStep {
	return removeStep(body.cleanPath.Push("Interface", data))
}

type finalClearSendBody struct {
	IStep
}

func (finalClearSendBody) JSON(...interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Clear().Send().Body() not with Clear().Send().Body(value)")
		},
	}
}

func (finalClearSendBody) Interface(...interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New("only usable with Clear().Send().Body() not with Clear().Send().Body(value)")
		},
	}
}
