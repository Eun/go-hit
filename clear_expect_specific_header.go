package hit

// IClearExpectSpecificHeader provides a clear functionality to remove previous steps from running in the Expect().Header(...) and Expect().Headers().Get(...) scope
type IClearExpectSpecificHeader interface {
	IStep

	// Contains removes all previous Expect().Header(...).Contains() steps.
	//
	// If you specify an argument it will only remove the Expect().Header(...).Contains() steps matching that argument.
	//
	// Examples:
	//           Clear().Expect().Header("Content-Type").Contains()              // will remove all Expect().Header("Content-Type").Contains() steps
	//           Clear().Expect().Header("Content-Type").Contains("json")        // will remove all Expect().Header("Content-Type").Contains("json") steps
	//           Clear().Expect().Headers().Get("Content-Type").Contains()       // will remove all Expect().Headers().Get("Content-Type").Contains() steps
	//           Clear().Expect().Headers().Get("Content-Type").Contains("json") // will remove all Expect().Headers().Get("Content-Type").Contains("json") steps
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect().Header("Content-Type").Contains("json")
	//               Clear().Expect().Header("Content-Type").Contains()
	//               Expect().Header("Content-Type").Equal("application/json")
	//           )
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect().Headers().Get("Content-Type").Contains("json")
	//               Clear().Expect().Headers().Get("Content-Type").Contains("json")
	//               Expect().Headers().Get("Content-Type").Equal("application/json")
	//           )
	Contains(...string) IStep

	// NotContains removes all previous Expect().Header(...).NotContains() steps.
	//
	// If you specify an argument it will only remove the Expect().Header(...).NotContains() steps matching that argument.
	//
	// Examples:
	//           Clear().Expect().Header("Content-Type").NotContains()              // will remove all Expect().Header("Content-Type").NotContains() steps
	//           Clear().Expect().Header("Content-Type").NotContains("json")        // will remove all Expect().Header("Content-Type").NotContains("json") steps
	//           Clear().Expect().Headers().Get("Content-Type").NotContains()       // will remove all Expect().Headers().Get("Content-Type").NotContains() steps
	//           Clear().Expect().Headers().Get("Content-Type").NotContains("json") // will remove all Expect().Headers().Get("Content-Type").NotContains("json") steps
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect().Header("Content-Type").NotContains("json")
	//               Clear().Expect().Header("Content-Type").NotContains()
	//               Expect().Header("Content-Type").Contains("xml")
	//           )
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect().Headers().Get("Content-Type").NotContains("json")
	//               Clear().Expect().Headers().Get("Content-Type").NotContains("json")
	//               Expect().Headers().Get("Content-Type").Contains("xml")
	//           )
	NotContains(...string) IStep

	// OneOf removes all previous Expect().Header(...).OneOf() steps.
	//
	// If you specify an argument it will only remove the Expect().Header(...).OneOf() steps matching that argument.
	//
	// Examples:
	//           Clear().Expect().Header("Content-Type").OneOf()              // will remove all Expect().Header("Content-Type").OneOf() steps
	//           Clear().Expect().Header("Content-Type").OneOf("json")        // will remove all Expect().Header("Content-Type").OneOf("json") steps
	//           Clear().Expect().Headers().Get("Content-Type").OneOf()       // will remove all Expect().Headers().Get("Content-Type").OneOf() steps
	//           Clear().Expect().Headers().Get("Content-Type").OneOf("json") // will remove all Expect().Headers().Get("Content-Type").OneOf("json") steps
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect().Header("Content-Type").OneOf("application/json", "application/xml")
	//               Clear().Expect().Header("Content-Type").OneOf()
	//               Expect().Header("Content-Type").Equal("application/json")
	//           )
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect().Header().Get("Content-Type").OneOf("application/json", "application/xml")
	//               Clear().Expect().Header().Get("Content-Type").OneOf("application/json", "application/xml")
	//               Expect().Header().Get("Content-Type").Equal("application/json")
	//           )
	OneOf(...string) IStep

	// NotOneOf removes all previous Expect().Header(...).NotOneOf() steps.
	//
	// If you specify an argument it will only remove the Expect().Header(...).NotOneOf() steps matching that argument.
	//
	// Examples:
	//           Clear().Expect().Header("Content-Type").NotOneOf()              // will remove all Expect().Header("Content-Type").NotOneOf() steps
	//           Clear().Expect().Header("Content-Type").NotOneOf("json")        // will remove all Expect().Header("Content-Type").NotOneOf("json") steps
	//           Clear().Expect().Headers().Get("Content-Type").NotOneOf()       // will remove all Expect().Headers().Get("Content-Type").NotOneOf() steps
	//           Clear().Expect().Headers().Get("Content-Type").NotOneOf("json") // will remove all Expect().Headers().Get("Content-Type").NotOneOf("json") steps
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect().Header("Content-Type").NotOneOf("application/json", "application/xml")
	//               Clear().Expect().Header("Content-Type").NotOneOf()
	//               Expect().Header("Content-Type").NotOneOf("application/xml")
	//           )
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect().Header().Get("Content-Type").NotOneOf("application/json", "application/xml")
	//               Clear().Expect().Header().Get("Content-Type").NotOneOf("application/json", "application/xml")
	//               Expect().Header().Get("Content-Type").NotOneOf("application/xml")
	//           )
	NotOneOf(...string) IStep

	// Empty removes all previous Expect().Header(...).Empty() steps.
	//
	// Examples:
	//           Clear().Expect().Header("Content-Type").Empty()              // will remove all Expect().Header("Content-Type").Empty() steps
	//           Clear().Expect().Headers().Get("Content-Type").Empty()       // will remove all Expect().Headers().Get("Content-Type").Empty() steps
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect().Header("Content-Type").Empty()
	//               Clear().Expect().Header("Content-Type").Empty()
	//               Expect().Header("Content-Type").Equal("application/json")
	//           )
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect().Headers().Get("Content-Type").Empty()
	//               Clear().Expect().Headers().Get("Content-Type").Empty()
	//               Expect().Headers().Get("Content-Type").Equal("application/json")
	//           )
	Empty() IStep

	// Len removes all previous Expect().Header(...).Len() steps.
	//
	// If you specify an argument it will only remove the Expect().Header(...).Len() steps matching that argument.
	//
	// Examples:
	//           Clear().Expect().Header("Content-Type").Len()              // will remove all Expect().Header("Content-Type").Len() steps
	//           Clear().Expect().Header("Content-Type").Len(10)            // will remove all Expect().Header("Content-Type").Len(10) steps
	//           Clear().Expect().Headers().Get("Content-Type").Len()       // will remove all Expect().Headers().Get("Content-Type").Len() steps
	//           Clear().Expect().Headers().Get("Content-Type").Len(10)     // will remove all Expect().Headers().Get("Content-Type").Len(10) steps
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect().Header("Content-Type").Len(10)
	//               Clear().Expect().Header("Content-Type").Len()
	//               Expect().Header("Content-Type").Len(20)
	//           )
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect().Headers().Get("Content-Type").Len(10)
	//               Clear().Expect().Headers().Get("Content-Type").Len(10)
	//               Expect().Headers().Get("Content-Type").Len(20)
	//           )
	Len(...int) IStep

	// Equal removes all previous Expect().Header(...).Equal() steps.
	//
	// If you specify an argument it will only remove the Expect().Header(...).Equal() steps matching that argument.
	//
	// Examples:
	//           Clear().Expect().Header("Content-Type").Equal()                          // will remove all Expect().Header("Content-Type").Equal() steps
	//           Clear().Expect().Header("Content-Type").Equal("application/json")        // will remove all Expect().Header("Content-Type").Equal("application/json") steps
	//           Clear().Expect().Headers().Get("Content-Type").Equal()                   // will remove all Expect().Headers().Get("Content-Type").Equal() steps
	//           Clear().Expect().Headers().Get("Content-Type").Equal("application/json") // will remove all Expect().Headers().Get("Content-Type").Equal("application/json") steps
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect().Header("Content-Type").Equal("application/xml")
	//               Clear().Expect().Header("Content-Type").Equal()
	//               Expect().Header("Content-Type").Equal("application/json")
	//           )
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect().Headers().Get("Content-Type").Equal("application/xml")
	//               Clear().Expect().Headers().Get("Content-Type").Equal("application/xml")
	//               Expect().Headers().Get("Content-Type").Equal("application/json")
	//           )
	Equal(...string) IStep

	// NotEqual removes all previous Expect().Header(...).NotEqual() steps.
	//
	// If you specify an argument it will only remove the Expect().Header(...).NotEqual() steps matching that argument.
	//
	// Examples:
	//           Clear().Expect().Header("Content-Type").NotEqual()                          // will remove all Expect().Header("Content-Type").NotEqual() steps
	//           Clear().Expect().Header("Content-Type").NotEqual("application/json")        // will remove all Expect().Header("Content-Type").NotEqual("application/json") steps
	//           Clear().Expect().Headers().Get("Content-Type").NotEqual()                   // will remove all Expect().Headers().Get("Content-Type").NotEqual() steps
	//           Clear().Expect().Headers().Get("Content-Type").NotEqual("application/json") // will remove all Expect().Headers().Get("Content-Type").NotEqual("application/json") steps
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect().Header("Content-Type").NotEqual("application/json")
	//               Clear().Expect().Header("Content-Type").NotEqual()
	//               Expect().Header("Content-Type").NotEqual("application/xml")
	//           )
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect().Headers().Get("Content-Type").NotEqual("application/json")
	//               Clear().Expect().Headers().Get("Content-Type").NotEqual("application/json")
	//               Expect().Headers().Get("Content-Type").NotEqual("application/xml")
	//           )
	NotEqual(...string) IStep
}
type clearExpectSpecificHeader struct {
	expect    IClearExpect
	cleanPath clearPath
}

func newClearExpectSpecificHeader(expect IClearExpect, cleanPath clearPath) IClearExpectSpecificHeader {
	return &clearExpectSpecificHeader{
		expect:    expect,
		cleanPath: cleanPath,
	}
}

func (hdr *clearExpectSpecificHeader) when() StepTime {
	return CleanStep
}

func (hdr *clearExpectSpecificHeader) exec(hit Hit) error {
	// this runs if we called Clear().Expect().Header("...") or Clear().Expect().Headers().Get("...")
	removeSteps(hit, hdr.cleanPath)
	return nil
}

func (hdr *clearExpectSpecificHeader) clearPath() clearPath {
	return hdr.cleanPath
}

func (hdr *clearExpectSpecificHeader) Contains(v ...string) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(hdr.cleanPath.Push("Contains", args))
}

func (hdr *clearExpectSpecificHeader) NotContains(v ...string) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(hdr.cleanPath.Push("NotContains", args))
}

func (hdr *clearExpectSpecificHeader) OneOf(v ...string) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(hdr.cleanPath.Push("OneOf", args))
}

func (hdr *clearExpectSpecificHeader) NotOneOf(v ...string) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(hdr.cleanPath.Push("NotOneOf", args))
}

func (hdr *clearExpectSpecificHeader) Empty() IStep {
	return removeStep(hdr.cleanPath.Push("Empty", nil))
}

func (hdr *clearExpectSpecificHeader) Len(v ...int) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(hdr.cleanPath.Push("Len", args))
}

func (hdr *clearExpectSpecificHeader) Equal(v ...string) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(hdr.cleanPath.Push("Equal", args))
}

func (hdr *clearExpectSpecificHeader) NotEqual(v ...string) IStep {
	args := make([]interface{}, len(v))
	for i := range v {
		args[i] = v[i]
	}
	return removeStep(hdr.cleanPath.Push("NotEqual", args))
}
