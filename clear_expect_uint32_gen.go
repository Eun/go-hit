// +build !generate

// ⚠️⚠️⚠️ This file was autogenerated by generators/clear/clear ⚠️⚠️⚠️ //
package hit

import errortrace "github.com/Eun/go-hit/errortrace"

// IClearExpectUint32 provides methods to clear steps.
type IClearExpectUint32 interface {
	IStep
	// Between clears all matching Between steps
	Between(value ...uint32) IStep
	// Equal clears all matching Equal steps
	Equal(value ...uint32) IStep
	// GreaterOrEqualThan clears all matching GreaterOrEqualThan steps
	GreaterOrEqualThan(value ...uint32) IStep
	// GreaterThan clears all matching GreaterThan steps
	GreaterThan(value ...uint32) IStep
	// LessOrEqualThan clears all matching LessOrEqualThan steps
	LessOrEqualThan(value ...uint32) IStep
	// LessThan clears all matching LessThan steps
	LessThan(value ...uint32) IStep
	// NotBetween clears all matching NotBetween steps
	NotBetween(value ...uint32) IStep
	// NotEqual clears all matching NotEqual steps
	NotEqual(value ...uint32) IStep
	// NotOneOf clears all matching NotOneOf steps
	NotOneOf(value ...uint32) IStep
	// OneOf clears all matching OneOf steps
	OneOf(value ...uint32) IStep
}
type clearExpectUint32 struct {
	cp callPath
	tr *errortrace.ErrorTrace
}

func newClearExpectUint32(cp callPath) IClearExpectUint32 {
	return &clearExpectUint32{cp: cp, tr: ett.Prepare()}
}
func (v *clearExpectUint32) trace() *errortrace.ErrorTrace {
	return v.tr
}
func (*clearExpectUint32) when() StepTime {
	return CleanStep
}
func (v *clearExpectUint32) callPath() callPath {
	return v.cp
}
func (v *clearExpectUint32) exec(hit *hitImpl) error {
	if err := removeSteps(hit, v.callPath()); err != nil {
		return err
	}
	return nil
}
func (v *clearExpectUint32) Between(value ...uint32) IStep {
	return removeStep(v.callPath().Push("Between", uint32SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint32) Equal(value ...uint32) IStep {
	return removeStep(v.callPath().Push("Equal", uint32SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint32) GreaterOrEqualThan(value ...uint32) IStep {
	return removeStep(v.callPath().Push("GreaterOrEqualThan", uint32SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint32) GreaterThan(value ...uint32) IStep {
	return removeStep(v.callPath().Push("GreaterThan", uint32SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint32) LessOrEqualThan(value ...uint32) IStep {
	return removeStep(v.callPath().Push("LessOrEqualThan", uint32SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint32) LessThan(value ...uint32) IStep {
	return removeStep(v.callPath().Push("LessThan", uint32SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint32) NotBetween(value ...uint32) IStep {
	return removeStep(v.callPath().Push("NotBetween", uint32SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint32) NotEqual(value ...uint32) IStep {
	return removeStep(v.callPath().Push("NotEqual", uint32SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint32) NotOneOf(value ...uint32) IStep {
	return removeStep(v.callPath().Push("NotOneOf", uint32SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint32) OneOf(value ...uint32) IStep {
	return removeStep(v.callPath().Push("OneOf", uint32SliceToInterfaceSlice(value)))
}
