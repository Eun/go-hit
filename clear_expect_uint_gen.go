// +build !generate

// ⚠️⚠️⚠️ This file was autogenerated by generators/clear/clear ⚠️⚠️⚠️ //
package hit

import errortrace "github.com/Eun/go-hit/errortrace"

// IClearExpectUint provides methods to clear steps.
type IClearExpectUint interface {
	IStep
	// Between clears all matching Between steps
	Between(value ...uint) IStep
	// Equal clears all matching Equal steps
	Equal(value ...uint) IStep
	// GreaterOrEqualThan clears all matching GreaterOrEqualThan steps
	GreaterOrEqualThan(value ...uint) IStep
	// GreaterThan clears all matching GreaterThan steps
	GreaterThan(value ...uint) IStep
	// LessOrEqualThan clears all matching LessOrEqualThan steps
	LessOrEqualThan(value ...uint) IStep
	// LessThan clears all matching LessThan steps
	LessThan(value ...uint) IStep
	// NotBetween clears all matching NotBetween steps
	NotBetween(value ...uint) IStep
	// NotEqual clears all matching NotEqual steps
	NotEqual(value ...uint) IStep
	// NotOneOf clears all matching NotOneOf steps
	NotOneOf(value ...uint) IStep
	// OneOf clears all matching OneOf steps
	OneOf(value ...uint) IStep
}
type clearExpectUint struct {
	cp callPath
	tr *errortrace.ErrorTrace
}

func newClearExpectUint(cp callPath) IClearExpectUint {
	return &clearExpectUint{cp: cp, tr: ett.Prepare()}
}
func (v *clearExpectUint) trace() *errortrace.ErrorTrace {
	return v.tr
}
func (*clearExpectUint) when() StepTime {
	return CleanStep
}
func (v *clearExpectUint) callPath() callPath {
	return v.cp
}
func (v *clearExpectUint) exec(hit *hitImpl) error {
	if err := removeSteps(hit, v.callPath()); err != nil {
		return err
	}
	return nil
}
func (v *clearExpectUint) Between(value ...uint) IStep {
	return removeStep(v.callPath().Push("Between", uintSliceToInterfaceSlice(value)))
}
func (v *clearExpectUint) Equal(value ...uint) IStep {
	return removeStep(v.callPath().Push("Equal", uintSliceToInterfaceSlice(value)))
}
func (v *clearExpectUint) GreaterOrEqualThan(value ...uint) IStep {
	return removeStep(v.callPath().Push("GreaterOrEqualThan", uintSliceToInterfaceSlice(value)))
}
func (v *clearExpectUint) GreaterThan(value ...uint) IStep {
	return removeStep(v.callPath().Push("GreaterThan", uintSliceToInterfaceSlice(value)))
}
func (v *clearExpectUint) LessOrEqualThan(value ...uint) IStep {
	return removeStep(v.callPath().Push("LessOrEqualThan", uintSliceToInterfaceSlice(value)))
}
func (v *clearExpectUint) LessThan(value ...uint) IStep {
	return removeStep(v.callPath().Push("LessThan", uintSliceToInterfaceSlice(value)))
}
func (v *clearExpectUint) NotBetween(value ...uint) IStep {
	return removeStep(v.callPath().Push("NotBetween", uintSliceToInterfaceSlice(value)))
}
func (v *clearExpectUint) NotEqual(value ...uint) IStep {
	return removeStep(v.callPath().Push("NotEqual", uintSliceToInterfaceSlice(value)))
}
func (v *clearExpectUint) NotOneOf(value ...uint) IStep {
	return removeStep(v.callPath().Push("NotOneOf", uintSliceToInterfaceSlice(value)))
}
func (v *clearExpectUint) OneOf(value ...uint) IStep {
	return removeStep(v.callPath().Push("OneOf", uintSliceToInterfaceSlice(value)))
}
