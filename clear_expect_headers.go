package hit

type IClearExpectHeaders interface {
	IStep
	// Contains removes all Expect().Headers().Contains() steps
	Contains() IStep

	// NotContains removes all Expect().Headers().NotContains() steps
	NotContains() IStep

	// Empty removes all Expect().Headers().Empty() steps
	Empty() IStep

	// Len removes all Expect().Headers().Len() steps
	Len() IStep

	// Equal removes all Expect().Headers().Equal() steps
	Equal() IStep

	// NotEqual removes all Expect().Headers().NotEqual() steps
	NotEqual() IStep

	// Get removes all Expect().Headers().Get() steps, if you specify the header it will only remove
	// the Expect().Headers().Get() steps with the matching header.
	// Examples:
	//           Expect().Headers().Get()                    // will remove all Get() steps
	//           Expect().Headers().Get("Content-Type")      // will only remove Get("Content-Type") steps
	Get(header ...string) IClearExpectSpecificHeader
}

type clearExpectHeaders struct {
	expect    IClearExpect
	cleanPath CleanPath
}

func newClearExpectHeaders(expect IClearExpect, cleanPath CleanPath) IClearExpectHeaders {
	return &clearExpectHeaders{
		expect:    expect,
		cleanPath: cleanPath,
	}
}

// implement IStep interface to make sure we can call just Expect().Headers()

func (hdr *clearExpectHeaders) When() StepTime {
	return CleanStep
}

// Exec contains the logic for Clear().Expect().Headers(), it will remove all Expect().Headers() and Expect().Headers().* Steps
func (hdr *clearExpectHeaders) Exec(hit Hit) error {
	removeSteps(hit, hdr.cleanPath)
	return nil
}

func (hdr *clearExpectHeaders) CleanPath() CleanPath {
	return hdr.cleanPath
}

// NotEqual removes all Expect().Headers().Contains() steps
func (hdr *clearExpectHeaders) Contains() IStep {
	return removeStep(hdr.cleanPath.Push("Contains", nil))
}

// NotEqual removes all Expect().Headers().NotContains() steps
func (hdr *clearExpectHeaders) NotContains() IStep {
	return removeStep(hdr.cleanPath.Push("NotContains", nil))
}

// NotEqual removes all Expect().Headers().Empty() steps
func (hdr *clearExpectHeaders) Empty() IStep {
	return removeStep(hdr.cleanPath.Push("Empty", nil))
}

// NotEqual removes all Expect().Headers().Len() steps
func (hdr *clearExpectHeaders) Len() IStep {
	return removeStep(hdr.cleanPath.Push("Len", nil))
}

// NotEqual removes all Expect().Headers().Equal() steps
func (hdr *clearExpectHeaders) Equal() IStep {
	return removeStep(hdr.cleanPath.Push("Equal", nil))
}

// NotEqual removes all Expect().Headers().NotEqual() steps
func (hdr *clearExpectHeaders) NotEqual() IStep {
	return removeStep(hdr.cleanPath.Push("NotEqual", nil))
}

// Get removes all Expect().Headers().Get() steps, if you specify the header it will only remove
// the Expect().Headers().Get() steps with the matching header.
// Examples:
//           Expect().Headers().Get()                    // will remove all Get() steps
//           Expect().Headers().Get("Content-Type")      // will only remove Get("Content-Type") steps
func (hdr *clearExpectHeaders) Get(header ...string) IClearExpectSpecificHeader {
	args := make([]interface{}, len(header))
	for i := range header {
		args[i] = header[i]
	}
	return newClearExpectSpecificHeader(hdr.expect, hdr.cleanPath.Push("Get", args))
}
