package hit

import (
	"github.com/Eun/go-hit/errortrace"
)

// IClearSend provides a clear functionality to remove previous steps from running in the Send() scope
type IClearSend interface {
	IStep
	// Body removes all previous Send().Body() steps and all steps chained to Send().Body() e.g. Send().Body().Interface("Hello World").
	//
	// If you specify an argument it will only remove the Send().Body() steps matching that argument.
	//
	// Usage:
	//     Clear().Send().Body()                      // will remove all Send().Body() steps and all chained steps to Send() e.g. Send().Body("Hello World")
	//     Clear().Send().Body("Hello World")         // will remove all Send().Body("Hello World") steps
	//     Clear().Send().Body().Interface()              // will remove all Send().Body().Interface() steps
	//     Clear().Send().Body().Interface("Hello World") // will remove all Send().Body().Interface("Hello World") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Send().Body("Hello Earth"),
	//         Clear().Send().Body(),
	//         Send().Body("Hello World"),
	//     )
	Body(value ...interface{}) IClearSendBody

	// Header removes all previous Send().Header() steps.
	//
	// If you specify an argument it will only remove the Send().Header() steps matching that argument.
	//
	// Usage:
	//     Clear().Send().Header()                                   // will remove all Send().Header() steps
	//     Clear().Send().Header("Content-Type")                     // will remove all Send().Header("Content-Type", ...) step
	//     Clear().Send().Header("Content-Type", "application/json") // will remove all Send().Header("Content-Type", "application/json") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Send().Header("Content-Type", "application/xml"),
	//         Clear().Send().Header("Content-Type"),
	//         Send().Header("Content-Type", "application/json"),
	//     )
	Header(values ...interface{}) IStep

	// Trailer removes all previous Send().Trailer() steps.
	//
	// If you specify an argument it will only remove the Send().Trailer() steps matching that argument.
	//
	// Usage:
	//     Clear().Send().Trailer()                                   // will remove all Send().Trailer() steps
	//     Clear().Send().Trailer("Content-Type")                     // will remove all Send().Trailer("Content-Type", ...) step
	//     Clear().Send().Trailer("Content-Type", "application/json") // will remove all Send().Trailer("Content-Type", "application/json") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Send().Trailer("Content-Type", "application/xml"),
	//         Clear().Send().Trailer("Content-Type"),
	//         Send().Trailer("Content-Type", "application/json"),
	//     )
	Trailer(values ...interface{}) IStep

	// Custom removes all previous Send().Custom() steps.
	//
	// If you specify an argument it will only remove the Send().Custom() steps matching that argument.
	//
	// Usage:
	//     Clear().Send().Custom(fn) // will remove all Send().Custom(fn) steps
	//     Clear().Send().Custom()   // will remove all Send().Custom() steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Send().Custom(func(hit Hit) {
	//             hit.Request().Body().SetString("Hello Earth")
	//         }),
	//         Clear().Send().Custom(),
	//         Send().Custom(func(hit Hit) {
	//             hit.Request().Body().SetString("Hello World")
	//         }),
	//     )
	Custom(fn ...Callback) IStep
}

type clearSend struct {
	cleanPath clearPath
	trace     *errortrace.ErrorTrace
}

func newClearSend(clearPath clearPath) IClearSend {
	return &clearSend{
		cleanPath: clearPath,
		trace:     ett.Prepare(),
	}
}

func (*clearSend) when() StepTime {
	return CleanStep
}

func (snd *clearSend) exec(hit Hit) error {
	// this runs if we called Clear().Send()
	if err := removeSteps(hit, snd.clearPath()); err != nil {
		return snd.trace.Format(hit.Description(), err.Error())
	}
	return nil
}

func (snd *clearSend) clearPath() clearPath {
	return snd.cleanPath
}

func (snd *clearSend) Body(value ...interface{}) IClearSendBody {
	return newClearSendBody(snd.clearPath().Push("Body", value), value)
}

func (snd *clearSend) Header(values ...interface{}) IStep {
	return removeStep(snd.clearPath().Push("Header", values))
}

func (snd *clearSend) Trailer(values ...interface{}) IStep {
	return removeStep(snd.clearPath().Push("Trailer", values))
}

// custom can be used to send a custom behaviour
func (snd *clearSend) Custom(fn ...Callback) IStep {
	args := make([]interface{}, len(fn))
	for i := range fn {
		args[i] = fn[i]
	}
	return removeStep(snd.clearPath().Push("Custom", args))
}
