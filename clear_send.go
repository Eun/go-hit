package hit

import (
	"github.com/Eun/go-hit/internal"
)

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

	// Interface removes all previous Send().Interface() steps.
	//
	// If you specify an argument it will only remove the Send().Interface() steps matching that argument.
	//
	// Usage:
	//     Clear().Send().Interface()              // will remove all Send().Interface() steps
	//     Clear().Send().Interface("Hello World") // will remove all Send().Interface("Hello World") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Send().Interface("Hello Earth"),
	//         Clear().Send().Interface(),
	//         Send().Interface("Hello World"),
	//     )
	Interface(value ...interface{}) IStep

	// JSON removes all previous Send().JSON() steps.
	//
	// If you specify an argument it will only remove the Send().JSON() steps matching that argument.
	//
	// Usage:
	//     Clear().Send().JSON()                                      // will remove all Send().JSON() steps
	//     Clear().Send().JSON(map[string]interface{}{"Name": "Joe"}) // will remove all Send().JSON("Hello World") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Send().JSON(map[string]interface{}{"Name": "Joe"}),
	//         Clear().Send().JSON(),
	//         Send().JSON(map[string]interface{}{"Name": "Alice"}),
	//     )
	JSON(value ...interface{}) IStep

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
}

func newClearSend(clearPath clearPath, params []interface{}) IClearSend {
	if _, ok := internal.GetLastArgument(params); ok {
		// this runs if we called Clear().Send(something)
		return finalClearSend{removeStep(clearPath)}
	}
	return &clearSend{
		cleanPath: clearPath,
	}
}

func (*clearSend) when() StepTime {
	return CleanStep
}

func (snd *clearSend) exec(hit Hit) error {
	// this runs if we called Clear().Send()
	removeSteps(hit, snd.cleanPath)
	return nil
}

func (snd *clearSend) clearPath() clearPath {
	return snd.cleanPath
}

func (snd *clearSend) Body(value ...interface{}) IClearSendBody {
	return newClearSendBody(snd.cleanPath.Push("Body", value), value)
}

func (snd *clearSend) Interface(value ...interface{}) IStep {
	return removeStep(snd.cleanPath.Push("Interface", value))
}

// custom can be used to send a custom behaviour
func (snd *clearSend) Custom(fn ...Callback) IStep {
	args := make([]interface{}, len(fn))
	for i := range fn {
		args[i] = fn[i]
	}
	return removeStep(snd.cleanPath.Push("Custom", args))
}

// JSON sets the body to the specific data (shortcut for Body().JSON()
func (snd *clearSend) JSON(value ...interface{}) IStep {
	return removeStep(snd.cleanPath.Push("JSON", value))
}

func (snd *clearSend) Header(values ...interface{}) IStep {
	return removeStep(snd.cleanPath.Push("Header", values))
}

type finalClearSend struct {
	IStep
}

func (finalClearSend) Body(...interface{}) IClearSendBody {
	panic("only usable with Send() not with Send(value)")
}

func (finalClearSend) Custom(...Callback) IStep {
	panic("only usable with Send() not with Send(value)")
}

func (finalClearSend) JSON(...interface{}) IStep {
	panic("only usable with Send() not with Send(value)")
}

func (finalClearSend) Header(...interface{}) IStep {
	panic("only usable with Send() not with Send(value)")
}

func (finalClearSend) Interface(...interface{}) IStep {
	panic("only usable with Send() not with Send(value)")
}
