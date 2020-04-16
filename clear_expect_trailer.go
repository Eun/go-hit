package hit

import (
	"github.com/Eun/go-hit/errortrace"
)

// IClearExpectTrailer provides a clear functionality to remove previous steps from running in the Expect().Trailer(...)  scope
type IClearExpectTrailer interface {
	IStep

	// Contains removes all previous Expect().Trailer(...).Contains() steps.
	//
	// If you specify an argument it will only remove the Expect().Trailer(...).Contains() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Trailer().Contains()                     // will remove all Expect().Trailer().Contains() steps
	//     Clear().Expect().Trailer().Contains("Content-Type")       // will remove all Expect().Trailer().Contains("Content-Type") steps
	//     Clear().Expect().Trailer("Content-Type").Contains()       // will remove all Expect().Trailer("Content-Type").Contains() steps
	//     Clear().Expect().Trailer("Content-Type").Contains("json") // will remove all Expect().Trailer("Content-Type").Contains("json") steps
	//
	// Examples:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Trailer().Contains("Content-Type")
	//         Clear().Expect().Trailer().Contains()
	//         Expect().Trailer().Contains("Set-Cookie")
	//     )
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Trailer("Content-Type").Contains("json")
	//         Clear().Expect().Trailer("Content-Type").Contains()
	//         Expect().Trailer("Content-Type").Equal("application/json")
	//     )
	Contains(value ...interface{}) IStep

	// NotContains removes all previous Expect().Trailer(...).NotContains() steps.
	//
	// If you specify an argument it will only remove the Expect().Trailer(...).NotContains() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Trailer().NotContains()                     // will remove all Expect().Trailer().NotContains() steps
	//     Clear().Expect().Trailer().NotContains("Content-Type")       // will remove all Expect().Trailer().NotContains("Content-Type") steps
	//     Clear().Expect().Trailer("Content-Type").NotContains()       // will remove all Expect().Trailer("Content-Type").NotContains() steps
	//     Clear().Expect().Trailer("Content-Type").NotContains("json") // will remove all Expect().Trailer("Content-Type").NotContains("json") steps
	//
	// Examples:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Trailer().NotContains("Content-Type")
	//         Clear().Expect().Trailer().NotContains()
	//         Expect().Trailer().Contains("Set-Cookie")
	//     )
	//
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Trailer("Content-Type").NotContains("json")
	//         Clear().Expect().Trailer("Content-Type").NotContains()
	//         Expect().Trailer("Content-Type").Contains("xml")
	//     )
	NotContains(value ...interface{}) IStep

	// OneOf removes all previous Expect().Trailer(...).OneOf() steps.
	//
	// If you specify an argument it will only remove the Expect().Trailer(...).OneOf() steps matching that argument.
	//
	// USage:
	//     Clear().Expect().Trailer("Content-Type").OneOf()       // will remove all Expect().Trailer("Content-Type").OneOf() steps
	//     Clear().Expect().Trailer("Content-Type").OneOf("json") // will remove all Expect().Trailer("Content-Type").OneOf("json") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Trailer("Content-Type").OneOf("application/json", "application/xml")
	//         Clear().Expect().Trailer("Content-Type").OneOf()
	//         Expect().Trailer("Content-Type").Equal("application/json")
	//     )
	OneOf(values ...interface{}) IStep

	// NotOneOf removes all previous Expect().Trailer(...).NotOneOf() steps.
	//
	// If you specify an argument it will only remove the Expect().Trailer(...).NotOneOf() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Trailer("Content-Type").NotOneOf()       // will remove all Expect().Trailer("Content-Type").NotOneOf() steps
	//     Clear().Expect().Trailer("Content-Type").NotOneOf("json") // will remove all Expect().Trailer("Content-Type").NotOneOf("json") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Trailer("Content-Type").NotOneOf("application/json", "application/xml")
	//         Clear().Expect().Trailer("Content-Type").NotOneOf()
	//         Expect().Trailer("Content-Type").NotOneOf("application/xml")
	//     )
	NotOneOf(values ...interface{}) IStep

	// Empty removes all previous Expect().Trailer(...).Empty() steps.
	//
	// Usage:
	//     Clear().Expect().Trailer().Empty()               // will remove all Expect().Trailer().Empty() steps
	//     Clear().Expect().Trailer("Content-Type").Empty() // will remove all Expect().Trailer("Content-Type").Empty() steps
	//
	// Examples:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Trailer().Empty()
	//         Clear().Expect().Trailer().Empty()
	//         Expect().Trailer().Contains("Content-Type")
	//     )
	//
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Trailer("Content-Type").Empty()
	//         Clear().Expect().Trailer("Content-Type").Empty()
	//         Expect().Trailer("Content-Type").Equal("application/json")
	//     )
	Empty() IStep

	// Len removes all previous Expect().Trailer(...).Len() steps.
	//
	// If you specify an argument it will only remove the Expect().Trailer(...).Len() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Trailer().Len()                 // will remove all Expect().Trailer().Len() steps
	//     Clear().Expect().Trailer().Len(10)               // will remove all Expect().Trailer().Len(10) steps
	//     Clear().Expect().Trailer("Content-Type").Len()   // will remove all Expect().Trailer("Content-Type").Len() steps
	//     Clear().Expect().Trailer("Content-Type").Len(10) // will remove all Expect().Trailer("Content-Type").Len(10) steps
	//
	// Examples:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Trailer().Len(1)
	//         Clear().Expect().Trailer().Len()
	//         Expect().Trailer().Len(2)
	//     )
	//
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Trailer("Content-Type").Len(10)
	//         Clear().Expect().Trailer("Content-Type").Len()
	//         Expect().Trailer("Content-Type").Len(20)
	//     )
	Len(size ...int) IStep

	// Equal removes all previous Expect().Trailer(...).Equal() steps.
	//
	// If you specify an argument it will only remove the Expect().Trailer(...).Equal() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Trailer().Equal()                                                           // will remove all Expect().Trailer().Equal() steps
	//     Clear().Expect().Trailer().Equal(map[string]interface{}{"Content-Type": "application/json"}) // will remove all Expect().Trailer().Equal(map[string]interface{}{"Content-Type": "application/json"}) steps
	//     Clear().Expect().Trailer("Content-Type").Equal()                                             // will remove all Expect().Trailer("Content-Type").Equal() steps
	//     Clear().Expect().Trailer("Content-Type").Equal("application/json")                           // will remove all Expect().Trailer("Content-Type").Equal("application/json") steps
	//
	// Examples:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Trailer().Equal(map[string]interface{}{"Content-Type": "application/xml"})
	//         Clear().Expect().Trailer().Equal()
	//         Expect().Trailer().Equal(map[string]interface{}{"Content-Type": "application/json"})
	//     )
	//
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Trailer("Content-Type").Equal("application/xml")
	//         Clear().Expect().Trailer("Content-Type").Equal("application/xml")
	//         Expect().Trailer("Content-Type").Equal("application/json")
	//     )
	Equal(value ...interface{}) IStep

	// NotEqual removes all previous Expect().Trailer(...).NotEqual() steps.
	//
	// If you specify an argument it will only remove the Expect().Trailer(...).NotEqual() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Trailer().NotEqual()                                                           // will remove all Expect().Trailer().NotEqual() steps
	//     Clear().Expect().Trailer().NotEqual(map[string]interface{}{"Content-Type": "application/json"}) // will remove all Expect().Trailer().NotEqual(map[string]interface{}{"Content-Type": "application/json"}) steps
	//     Clear().Expect().Trailer("Content-Type").NotEqual()                                             // will remove all Expect().Trailer("Content-Type").NotEqual() steps
	//     Clear().Expect().Trailer("Content-Type").NotEqual("application/json")                           // will remove all Expect().Trailer("Content-Type").NotEqual("application/json") steps
	//
	// Examples:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Trailer().NotEqual(map[string]interface{}{"Content-Type": "application/json"})
	//         Clear().Expect().Trailer().NotEqual()
	//         Expect().Trailer().NotEqual(map[string]interface{}{"Content-Type": "application/xml"})
	//     )
	//
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Trailer("Content-Type").NotEqual("application/json")
	//         Clear().Expect().Trailer("Content-Type").NotEqual("application/json")
	//         Expect().Trailer("Content-Type").NotEqual("application/xml")
	//     )
	NotEqual(value ...interface{}) IStep
}
type clearExpectTrailer struct {
	cleanPath clearPath
	trace     *errortrace.ErrorTrace
}

func newClearExpectTrailer(cleanPath clearPath) IClearExpectTrailer {
	return &clearExpectTrailer{
		cleanPath: cleanPath,
		trace:     ett.Prepare(),
	}
}

func (hdr *clearExpectTrailer) when() StepTime {
	return CleanStep
}

func (hdr *clearExpectTrailer) exec(hit Hit) error {
	// this runs if we called Clear().Expect().Trailer("...")
	if err := removeSteps(hit, hdr.clearPath()); err != nil {
		return hdr.trace.Format(hit.Description(), err.Error())
	}
	return nil
}

func (hdr *clearExpectTrailer) clearPath() clearPath {
	return hdr.cleanPath
}

func (hdr *clearExpectTrailer) Contains(value ...interface{}) IStep {
	return removeStep(hdr.clearPath().Push("Contains", value))
}

func (hdr *clearExpectTrailer) NotContains(value ...interface{}) IStep {
	return removeStep(hdr.clearPath().Push("NotContains", value))
}

func (hdr *clearExpectTrailer) OneOf(values ...interface{}) IStep {
	return removeStep(hdr.clearPath().Push("OneOf", values))
}

func (hdr *clearExpectTrailer) NotOneOf(values ...interface{}) IStep {
	return removeStep(hdr.clearPath().Push("NotOneOf", values))
}

func (hdr *clearExpectTrailer) Empty() IStep {
	return removeStep(hdr.clearPath().Push("Empty", nil))
}

func (hdr *clearExpectTrailer) Len(size ...int) IStep {
	args := make([]interface{}, len(size))
	for i := range size {
		args[i] = size[i]
	}
	return removeStep(hdr.clearPath().Push("Len", args))
}

func (hdr *clearExpectTrailer) Equal(value ...interface{}) IStep {
	return removeStep(hdr.clearPath().Push("Equal", value))
}

func (hdr *clearExpectTrailer) NotEqual(value ...interface{}) IStep {
	return removeStep(hdr.clearPath().Push("NotEqual", value))
}
