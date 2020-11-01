// +build !generate

package hit

import errortrace "github.com/Eun/go-hit/errortrace"

// ⚠️⚠️⚠️ This file was autogenerated by generators/clear/clear ⚠️⚠️⚠️ //

// IClearExpectString provides methods to clear steps.
type IClearExpectString interface {
	IStep
	// Contains clears all matching Contains steps
	Contains(value ...string) IStep
	// Equal clears all matching Equal steps
	Equal(value ...string) IStep
	// Len clears all matching Len steps
	Len() IClearExpectInt
	// NotContains clears all matching NotContains steps
	NotContains(value ...string) IStep
	// NotEqual clears all matching NotEqual steps
	NotEqual(value ...string) IStep
	// NotOneOf clears all matching NotOneOf steps
	NotOneOf(value ...string) IStep
	// OneOf clears all matching OneOf steps
	OneOf(value ...string) IStep
}
type clearExpectString struct {
	cp callPath
	tr *errortrace.ErrorTrace
}

func newClearExpectString(cp callPath) IClearExpectString {
	return &clearExpectString{cp: cp, tr: ett.Prepare()}
}
func (v *clearExpectString) trace() *errortrace.ErrorTrace {
	return v.tr
}
func (*clearExpectString) when() StepTime {
	return CleanStep
}
func (v *clearExpectString) callPath() callPath {
	return v.cp
}
func (v *clearExpectString) exec(hit *hitImpl) error {
	if err := removeSteps(hit, v.callPath()); err != nil {
		return err
	}
	return nil
}
func (v *clearExpectString) Contains(value ...string) IStep {
	return removeStep(v.callPath().Push("Contains", stringSliceToInterfaceSlice(value)))
}
func (v *clearExpectString) Equal(value ...string) IStep {
	return removeStep(v.callPath().Push("Equal", stringSliceToInterfaceSlice(value)))
}
func (v *clearExpectString) Len() IClearExpectInt {
	return newClearExpectInt(v.callPath().Push("Len", nil))
}
func (v *clearExpectString) NotContains(value ...string) IStep {
	return removeStep(v.callPath().Push("NotContains", stringSliceToInterfaceSlice(value)))
}
func (v *clearExpectString) NotEqual(value ...string) IStep {
	return removeStep(v.callPath().Push("NotEqual", stringSliceToInterfaceSlice(value)))
}
func (v *clearExpectString) NotOneOf(value ...string) IStep {
	return removeStep(v.callPath().Push("NotOneOf", stringSliceToInterfaceSlice(value)))
}
func (v *clearExpectString) OneOf(value ...string) IStep {
	return removeStep(v.callPath().Push("OneOf", stringSliceToInterfaceSlice(value)))
}
