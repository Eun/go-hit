package hit

// IClearExpectHeader provides a clear functionality to remove previous steps from running in the Expect().Header(...)  scope
type IClearExpectHeader interface {
	IStep

	// Contains removes all previous Expect().Header(...).Contains() steps.
	//
	// If you specify an argument it will only remove the Expect().Header(...).Contains() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Header().Contains()                     // will remove all Expect().Header().Contains() steps
	//     Clear().Expect().Header().Contains("Content-Type")       // will remove all Expect().Header().Contains("Content-Type") steps
	//     Clear().Expect().Header("Content-Type").Contains()       // will remove all Expect().Header("Content-Type").Contains() steps
	//     Clear().Expect().Header("Content-Type").Contains("json") // will remove all Expect().Header("Content-Type").Contains("json") steps
	//
	// Examples:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Header().Contains("Content-Type")
	//         Clear().Expect().Header().Contains()
	//         Expect().Header().Contains("Set-Cookie")
	//     )
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Header("Content-Type").Contains("json")
	//         Clear().Expect().Header("Content-Type").Contains()
	//         Expect().Header("Content-Type").Equal("application/json")
	//     )
	Contains(value ...interface{}) IStep

	// NotContains removes all previous Expect().Header(...).NotContains() steps.
	//
	// If you specify an argument it will only remove the Expect().Header(...).NotContains() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Header().NotContains()                     // will remove all Expect().Header().NotContains() steps
	//     Clear().Expect().Header().NotContains("Content-Type")       // will remove all Expect().Header().NotContains("Content-Type") steps
	//     Clear().Expect().Header("Content-Type").NotContains()       // will remove all Expect().Header("Content-Type").NotContains() steps
	//     Clear().Expect().Header("Content-Type").NotContains("json") // will remove all Expect().Header("Content-Type").NotContains("json") steps
	//
	// Examples:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Header().NotContains("Content-Type")
	//         Clear().Expect().Header().NotContains()
	//         Expect().Header().Contains("Set-Cookie")
	//     )
	//
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Header("Content-Type").NotContains("json")
	//         Clear().Expect().Header("Content-Type").NotContains()
	//         Expect().Header("Content-Type").Contains("xml")
	//     )
	NotContains(value ...interface{}) IStep

	// OneOf removes all previous Expect().Header(...).OneOf() steps.
	//
	// If you specify an argument it will only remove the Expect().Header(...).OneOf() steps matching that argument.
	//
	// USage:
	//     Clear().Expect().Header("Content-Type").OneOf()       // will remove all Expect().Header("Content-Type").OneOf() steps
	//     Clear().Expect().Header("Content-Type").OneOf("json") // will remove all Expect().Header("Content-Type").OneOf("json") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Header("Content-Type").OneOf("application/json", "application/xml")
	//         Clear().Expect().Header("Content-Type").OneOf()
	//         Expect().Header("Content-Type").Equal("application/json")
	//     )
	OneOf(values ...interface{}) IStep

	// NotOneOf removes all previous Expect().Header(...).NotOneOf() steps.
	//
	// If you specify an argument it will only remove the Expect().Header(...).NotOneOf() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Header("Content-Type").NotOneOf()       // will remove all Expect().Header("Content-Type").NotOneOf() steps
	//     Clear().Expect().Header("Content-Type").NotOneOf("json") // will remove all Expect().Header("Content-Type").NotOneOf("json") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Header("Content-Type").NotOneOf("application/json", "application/xml")
	//         Clear().Expect().Header("Content-Type").NotOneOf()
	//         Expect().Header("Content-Type").NotOneOf("application/xml")
	//     )
	NotOneOf(values ...interface{}) IStep

	// Empty removes all previous Expect().Header(...).Empty() steps.
	//
	// Usage:
	//     Clear().Expect().Header().Empty()               // will remove all Expect().Header().Empty() steps
	//     Clear().Expect().Header("Content-Type").Empty() // will remove all Expect().Header("Content-Type").Empty() steps
	//
	// Examples:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Header().Empty()
	//         Clear().Expect().Header().Empty()
	//         Expect().Header().Contains("Content-Type")
	//     )
	//
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Header("Content-Type").Empty()
	//         Clear().Expect().Header("Content-Type").Empty()
	//         Expect().Header("Content-Type").Equal("application/json")
	//     )
	Empty() IStep

	// Len removes all previous Expect().Header(...).Len() steps.
	//
	// If you specify an argument it will only remove the Expect().Header(...).Len() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Header().Len()                 // will remove all Expect().Header().Len() steps
	//     Clear().Expect().Header().Len(10)               // will remove all Expect().Header().Len(10) steps
	//     Clear().Expect().Header("Content-Type").Len()   // will remove all Expect().Header("Content-Type").Len() steps
	//     Clear().Expect().Header("Content-Type").Len(10) // will remove all Expect().Header("Content-Type").Len(10) steps
	//
	// Examples:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Header().Len(1)
	//         Clear().Expect().Header().Len()
	//         Expect().Header().Len(2)
	//     )
	//
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Header("Content-Type").Len(10)
	//         Clear().Expect().Header("Content-Type").Len()
	//         Expect().Header("Content-Type").Len(20)
	//     )
	Len(size ...int) IStep

	// Equal removes all previous Expect().Header(...).Equal() steps.
	//
	// If you specify an argument it will only remove the Expect().Header(...).Equal() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Header().Equal()                                                           // will remove all Expect().Header().Equal() steps
	//     Clear().Expect().Header().Equal(map[string]interface{}{"Content-Type": "application/json"}) // will remove all Expect().Header().Equal(map[string]interface{}{"Content-Type": "application/json"}) steps
	//     Clear().Expect().Header("Content-Type").Equal()                                             // will remove all Expect().Header("Content-Type").Equal() steps
	//     Clear().Expect().Header("Content-Type").Equal("application/json")                           // will remove all Expect().Header("Content-Type").Equal("application/json") steps
	//
	// Examples:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Header().Equal(map[string]interface{}{"Content-Type": "application/xml"})
	//         Clear().Expect().Header().Equal()
	//         Expect().Header().Equal(map[string]interface{}{"Content-Type": "application/json"})
	//     )
	//
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Header("Content-Type").Equal("application/xml")
	//         Clear().Expect().Header("Content-Type").Equal("application/xml")
	//         Expect().Header("Content-Type").Equal("application/json")
	//     )
	Equal(value ...interface{}) IStep

	// NotEqual removes all previous Expect().Header(...).NotEqual() steps.
	//
	// If you specify an argument it will only remove the Expect().Header(...).NotEqual() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Header().NotEqual()                                                           // will remove all Expect().Header().NotEqual() steps
	//     Clear().Expect().Header().NotEqual(map[string]interface{}{"Content-Type": "application/json"}) // will remove all Expect().Header().NotEqual(map[string]interface{}{"Content-Type": "application/json"}) steps
	//     Clear().Expect().Header("Content-Type").NotEqual()                                             // will remove all Expect().Header("Content-Type").NotEqual() steps
	//     Clear().Expect().Header("Content-Type").NotEqual("application/json")                           // will remove all Expect().Header("Content-Type").NotEqual("application/json") steps
	//
	// Examples:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Header().NotEqual(map[string]interface{}{"Content-Type": "application/json"})
	//         Clear().Expect().Header().NotEqual()
	//         Expect().Header().NotEqual(map[string]interface{}{"Content-Type": "application/xml"})
	//     )
	//
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Header("Content-Type").NotEqual("application/json")
	//         Clear().Expect().Header("Content-Type").NotEqual("application/json")
	//         Expect().Header("Content-Type").NotEqual("application/xml")
	//     )
	NotEqual(value ...interface{}) IStep
}
type clearExpectHeader struct {
	cleanPath clearPath
}

func newClearExpectHeader(cleanPath clearPath) IClearExpectHeader {
	return &clearExpectHeader{
		cleanPath: cleanPath,
	}
}

func (hdr *clearExpectHeader) when() StepTime {
	return CleanStep
}

func (hdr *clearExpectHeader) exec(hit Hit) error {
	// this runs if we called Clear().Expect().Header("...") or Clear().Expect().Headers().Get("...")
	removeSteps(hit, hdr.clearPath())
	return nil
}

func (hdr *clearExpectHeader) clearPath() clearPath {
	return hdr.cleanPath
}

func (hdr *clearExpectHeader) Contains(value ...interface{}) IStep {
	return removeStep(hdr.cleanPath.Push("Contains", value))
}

func (hdr *clearExpectHeader) NotContains(value ...interface{}) IStep {
	return removeStep(hdr.cleanPath.Push("NotContains", value))
}

func (hdr *clearExpectHeader) OneOf(values ...interface{}) IStep {
	return removeStep(hdr.cleanPath.Push("OneOf", values))
}

func (hdr *clearExpectHeader) NotOneOf(values ...interface{}) IStep {
	return removeStep(hdr.cleanPath.Push("NotOneOf", values))
}

func (hdr *clearExpectHeader) Empty() IStep {
	return removeStep(hdr.cleanPath.Push("Empty", nil))
}

func (hdr *clearExpectHeader) Len(size ...int) IStep {
	args := make([]interface{}, len(size))
	for i := range size {
		args[i] = size[i]
	}
	return removeStep(hdr.cleanPath.Push("Len", args))
}

func (hdr *clearExpectHeader) Equal(value ...interface{}) IStep {
	return removeStep(hdr.cleanPath.Push("Equal", value))
}

func (hdr *clearExpectHeader) NotEqual(value ...interface{}) IStep {
	return removeStep(hdr.cleanPath.Push("NotEqual", value))
}
