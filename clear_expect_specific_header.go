package hit

type IClearExpectSpecificHeader interface {
	IStep
	// Contains removes all Expect().Header(..).Contains() steps
	Contains() IStep

	// NotContains removes all Expect().Header(..).NotContains() steps
	NotContains() IStep

	// OneOf removes all Expect().Header(..).OneOf() steps
	OneOf() IStep

	// NotOneOf removes all Expect().Header(..).NotOneOf() steps
	NotOneOf() IStep

	// Empty removes all Expect().Header(..).Empty() steps
	Empty() IStep

	// Len removes all Expect().Header(..).Len() steps
	Len() IStep

	// Equal removes all Expect().Header(..).Equal() steps
	Equal() IStep

	// NotEqual removes all Expect().Header(..).NotEqual() steps
	NotEqual() IStep
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

// implement IStep interface to make sure we can call just Expect().Header()

func (hdr *clearExpectSpecificHeader) when() StepTime {
	return CleanStep
}

// exec contains the logic for Clear().Expect().Header() and/or Clear().Expect().Headers(...).Get(...),
// it will remove all Expect().Header(), Expect().Header().* and/or Expect().Header(...).Get(...), Expect().Header(...).Get(...).* Steps
func (hdr *clearExpectSpecificHeader) exec(hit Hit) error {
	removeSteps(hit, hdr.cleanPath)
	return nil
}

func (hdr *clearExpectSpecificHeader) clearPath() clearPath {
	return hdr.cleanPath
}

// Contains removes all Expect().Header(..).Contains() steps
func (hdr *clearExpectSpecificHeader) Contains() IStep {
	return removeStep(hdr.cleanPath.Push("Contains", nil))
}

// NotContains removes all Expect().Header(..).NotContains() steps
func (hdr *clearExpectSpecificHeader) NotContains() IStep {
	return removeStep(hdr.cleanPath.Push("NotContains", nil))
}

// OneOf removes all Expect().Header(..).OneOf() steps
func (hdr *clearExpectSpecificHeader) OneOf() IStep {
	return removeStep(hdr.cleanPath.Push("OneOf", nil))
}

// NotOneOf removes all Expect().Header(..).NotOneOf() steps
func (hdr *clearExpectSpecificHeader) NotOneOf() IStep {
	return removeStep(hdr.cleanPath.Push("NotOneOf", nil))
}

// Empty removes all Expect().Header(..).Empty() steps
func (hdr *clearExpectSpecificHeader) Empty() IStep {
	return removeStep(hdr.cleanPath.Push("Empty", nil))
}

// Len removes all Expect().Header(..).Len() steps
func (hdr *clearExpectSpecificHeader) Len() IStep {
	return removeStep(hdr.cleanPath.Push("Len", nil))
}

// Equal removes all Expect().Header(..).Equal() steps
func (hdr *clearExpectSpecificHeader) Equal() IStep {
	return removeStep(hdr.cleanPath.Push("Equal", nil))
}

// NotEqual removes all Expect().Header(..).NotEqual() steps
func (hdr *clearExpectSpecificHeader) NotEqual() IStep {
	return removeStep(hdr.cleanPath.Push("NotEqual", nil))
}
