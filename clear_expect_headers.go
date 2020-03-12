package hit

// IClearExpectHeaders provides a clear functionality to remove previous steps from running in the Expect().Headers() scope
type IClearExpectHeaders interface {
	IStep
	// Contains removes all Expect().Headers().Contains() steps.
	//
	// If you specify the expression it will only remove the Expect().Headers().Contains() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Headers().Contains()               // will remove all Expect().Headers().Contains() steps
	//     Clear().Expect().Headers().Contains("Content-Type") // will only remove Expect().Headers().Contains("Content-Type") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Headers().Contains("Content-Type"),
	//         Clear().Expect().Headers().Contains("Content-Type"),
	//     )
	Contains(headerName ...string) IStep

	// NotContains removes all Expect().Headers().NotContains() steps.
	//
	// If you specify the expression it will only remove the Expect().Headers().NotContains() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Headers().NotContains()               // will remove all Expect().Headers().NotContains() steps
	//     Clear().Expect().Headers().NotContains("Content-Type") // will only remove Expect().Headers().NotContains("Content-Type") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Headers().NotContains("Content-Type"),
	//         Clear().Expect().Headers().NotContains("Content-Type"),
	//     )
	NotContains(headerName ...string) IStep

	// Empty removes all Expect().Headers().Empty() steps.
	//
	// Usage:
	//     Clear().Expect().Headers().Empty() // will remove all Empty() steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Headers().Empty(),
	//         Clear().Expect().Headers().Empty(),
	//     )
	Empty() IStep

	// Len removes all Expect().Headers().Len() steps.
	//
	// If you specify the expression it will only remove the Expect().Headers().Len() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Headers().Len()  // will remove all Expect().Headers().Len() steps
	//     Clear().Expect().Headers().Len(6) // will only remove Expect().Headers().Len(6) steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Headers().Len(6),
	//         Clear().Expect().Headers().Len(),
	//         Expect().Headers().Len(10),
	//     )
	Len(size ...int) IStep

	// Equal removes all Expect().Headers().Equal() steps.
	//
	// If you specify the expression it will only remove the Expect().Headers().Equal() steps matching that argument
	//
	// Usage:
	//     Clear().Expect().Headers().Equal()                                                      // will remove all Expect().Headers().Equal() steps
	//     Clear().Expect().Headers().Equal(map[string]string{"Content-Type": "application/json"}) // will only remove Expect().Headers().Equal(map[string]string{"Content-Type": "application/json"}) steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Headers().Equal(map[string]string{"Content-Type": "application/xml"}),
	//         Clear().Expect().Headers().Equal(),
	//         Expect().Headers().Equal(map[string]string{"Content-Type": "application/json"}),
	//     )
	Equal(value ...interface{}) IStep

	// NotEqual removes all Expect().Headers().NotEqual() steps, if you specify the expression it will only remove
	// the Expect().Headers().NotEqual() steps matching that argument
	//
	// Usage:
	//     Clear().Expect().Headers().NotEqual()                                                      // will remove all Expect().Headers().NotEqual() steps
	//     Clear().Expect().Headers().NotEqual(map[string]string{"Content-Type": "application/json"}) // will only remove Expect().Headers().NotEqual(map[string]string{"Content-Type": "application/json"}) steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Headers().NotEqual(map[string]string{"Content-Type": "application/json"}),
	//         Clear().Expect().Headers().NotEqual(),
	//         Expect().Headers().NotEqual(map[string]string{"Content-Type": "application/xml"}),
	//     )
	NotEqual(value ...interface{}) IStep

	// Get removes all Expect().Headers().Get() steps and all steps chained to Expect().Headers().Get() e.g. Expect().Headers().Get("Content-Type").Equal("application/json")
	// if you specify the expression it will only remove
	// the Expect().Headers().Get() steps matching that argument
	//
	// Usage:
	//     Clear().Expect().Headers().Get()               // will remove all Expect().Headers().Get() steps
	//     Clear().Expect().Headers().Get("Content-Type") // will only remove Expect().Headers().Get("Content-Type") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Headers().Get("Content-Type").Contains("application"),
	//         Expect().Headers().Get("Content-Type").OneOf("application/json", "application/xml"),
	//         Clear().Expect().Headers().Get(),
	//         Expect().Headers().Get("Content-Type").Equal("application/json"),
	//     )
	Get(headerName ...string) IClearExpectSpecificHeader
}

type clearExpectHeaders struct {
	expect    IClearExpect
	cleanPath clearPath
}

func newClearExpectHeaders(expect IClearExpect, cleanPath clearPath) IClearExpectHeaders {
	return &clearExpectHeaders{
		expect:    expect,
		cleanPath: cleanPath,
	}
}

func (hdr *clearExpectHeaders) when() StepTime {
	return CleanStep
}

func (hdr *clearExpectHeaders) exec(hit Hit) error {
	// this runs if we called Clear().Expect().Headers()
	removeSteps(hit, hdr.cleanPath)
	return nil
}

func (hdr *clearExpectHeaders) clearPath() clearPath {
	return hdr.cleanPath
}

func (hdr *clearExpectHeaders) Contains(headerName ...string) IStep {
	args := make([]interface{}, len(headerName))
	for i := range headerName {
		args[i] = headerName[i]
	}
	return removeStep(hdr.cleanPath.Push("Contains", args))
}

func (hdr *clearExpectHeaders) NotContains(headerName ...string) IStep {
	args := make([]interface{}, len(headerName))
	for i := range headerName {
		args[i] = headerName[i]
	}
	return removeStep(hdr.cleanPath.Push("NotContains", args))
}

func (hdr *clearExpectHeaders) Empty() IStep {
	return removeStep(hdr.cleanPath.Push("Empty", nil))
}

func (hdr *clearExpectHeaders) Len(size ...int) IStep {
	args := make([]interface{}, len(size))
	for i := range size {
		args[i] = size[i]
	}
	return removeStep(hdr.cleanPath.Push("Len", args))
}

func (hdr *clearExpectHeaders) Equal(value ...interface{}) IStep {
	return removeStep(hdr.cleanPath.Push("Equal", value))
}

func (hdr *clearExpectHeaders) NotEqual(value ...interface{}) IStep {
	return removeStep(hdr.cleanPath.Push("NotEqual", value))
}

func (hdr *clearExpectHeaders) Get(headerName ...string) IClearExpectSpecificHeader {
	args := make([]interface{}, len(headerName))
	for i := range headerName {
		args[i] = headerName[i]
	}
	return newClearExpectSpecificHeader(hdr.expect, hdr.cleanPath.Push("Get", args))
}
