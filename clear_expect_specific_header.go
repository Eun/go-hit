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
	hit       Hit
	cleanPath CleanPath
}

func newClearExpectSpecificHeader(expect IClearExpect, hit Hit, cleanPath CleanPath) IClearExpectSpecificHeader {
	return &clearExpectSpecificHeader{
		expect:    expect,
		hit:       hit,
		cleanPath: cleanPath,
	}
}

// implement IStep interface to make sure we can call just Expect().Header()

func (hdr *clearExpectSpecificHeader) When() StepTime {
	return CleanStep
}

// Exec contains the logic for Clear().Expect().Header() and/or Clear().Expect().Headers(...).Get(...),
// it will remove all Expect().Header(), Expect().Header().* and/or Expect().Header(...).Get(...), Expect().Header(...).Get(...).* Steps
func (hdr *clearExpectSpecificHeader) Exec(hit Hit) error {
	removeSteps(hit, hdr.cleanPath)
	return nil
}

func (hdr *clearExpectSpecificHeader) CleanPath() CleanPath {
	return hdr.cleanPath
}

// Contains removes all Expect().Header(..).Contains() steps
func (hdr *clearExpectSpecificHeader) Contains() IStep {
	return removeStep(hdr.hit, hdr.cleanPath.Push("Contains"))
}

// NotContains removes all Expect().Header(..).NotContains() steps
func (hdr *clearExpectSpecificHeader) NotContains() IStep {
	return removeStep(hdr.hit, hdr.cleanPath.Push("NotContains"))
}

// OneOf removes all Expect().Header(..).OneOf() steps
func (hdr *clearExpectSpecificHeader) OneOf() IStep {
	return removeStep(hdr.hit, hdr.cleanPath.Push("OneOf"))
}

// NotOneOf removes all Expect().Header(..).NotOneOf() steps
func (hdr *clearExpectSpecificHeader) NotOneOf() IStep {
	return removeStep(hdr.hit, hdr.cleanPath.Push("NotOneOf"))
}

// Empty removes all Expect().Header(..).Empty() steps
func (hdr *clearExpectSpecificHeader) Empty() IStep {
	return removeStep(hdr.hit, hdr.cleanPath.Push("Empty"))
}

// Len removes all Expect().Header(..).Len() steps
func (hdr *clearExpectSpecificHeader) Len() IStep {
	return removeStep(hdr.hit, hdr.cleanPath.Push("Len"))
}

// Equal removes all Expect().Header(..).Equal() steps
func (hdr *clearExpectSpecificHeader) Equal() IStep {
	return removeStep(hdr.hit, hdr.cleanPath.Push("Equal"))
}

// NotEqual removes all Expect().Header(..).NotEqual() steps
func (hdr *clearExpectSpecificHeader) NotEqual() IStep {
	return removeStep(hdr.hit, hdr.cleanPath.Push("NotEqual"))
}
