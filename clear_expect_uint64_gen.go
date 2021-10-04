//go:build !generate
// +build !generate

package hit

import errortrace "github.com/Eun/go-hit/errortrace"

// ⚠️⚠️⚠️ This file was autogenerated by generators/clear/clear ⚠️⚠️⚠️ //

// IClearExpectUint64 provides methods to clear steps.
type IClearExpectUint64 interface {
	IStep
	// Between clears all matching Between steps
	Between(value ...uint64) IStep
	// Equal clears all matching Equal steps
	Equal(value ...uint64) IStep
	// GreaterOrEqualThan clears all matching GreaterOrEqualThan steps
	GreaterOrEqualThan(value ...uint64) IStep
	// GreaterThan clears all matching GreaterThan steps
	GreaterThan(value ...uint64) IStep
	// LessOrEqualThan clears all matching LessOrEqualThan steps
	LessOrEqualThan(value ...uint64) IStep
	// LessThan clears all matching LessThan steps
	LessThan(value ...uint64) IStep
	// NotBetween clears all matching NotBetween steps
	NotBetween(value ...uint64) IStep
	// NotEqual clears all matching NotEqual steps
	NotEqual(value ...uint64) IStep
	// NotOneOf clears all matching NotOneOf steps
	NotOneOf(value ...uint64) IStep
	// OneOf clears all matching OneOf steps
	OneOf(value ...uint64) IStep
}
type clearExpectUint64 struct {
	cp callPath
	tr *errortrace.ErrorTrace
}

func newClearExpectUint64(cp callPath) IClearExpectUint64 {
	return &clearExpectUint64{cp: cp, tr: ett.Prepare()}
}
func (v *clearExpectUint64) trace() *errortrace.ErrorTrace {
	return v.tr
}
func (*clearExpectUint64) when() StepTime {
	return cleanStep
}
func (v *clearExpectUint64) callPath() callPath {
	return v.cp
}
func (v *clearExpectUint64) exec(hit *hitImpl) error {
	if err := removeSteps(hit, v.callPath()); err != nil {
		return err
	}
	return nil
}
func (v *clearExpectUint64) Between(value ...uint64) IStep {
	return removeStep(v.callPath().Push("Between", uint64SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint64) Equal(value ...uint64) IStep {
	return removeStep(v.callPath().Push("Equal", uint64SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint64) GreaterOrEqualThan(value ...uint64) IStep {
	return removeStep(v.callPath().Push("GreaterOrEqualThan", uint64SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint64) GreaterThan(value ...uint64) IStep {
	return removeStep(v.callPath().Push("GreaterThan", uint64SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint64) LessOrEqualThan(value ...uint64) IStep {
	return removeStep(v.callPath().Push("LessOrEqualThan", uint64SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint64) LessThan(value ...uint64) IStep {
	return removeStep(v.callPath().Push("LessThan", uint64SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint64) NotBetween(value ...uint64) IStep {
	return removeStep(v.callPath().Push("NotBetween", uint64SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint64) NotEqual(value ...uint64) IStep {
	return removeStep(v.callPath().Push("NotEqual", uint64SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint64) NotOneOf(value ...uint64) IStep {
	return removeStep(v.callPath().Push("NotOneOf", uint64SliceToInterfaceSlice(value)))
}
func (v *clearExpectUint64) OneOf(value ...uint64) IStep {
	return removeStep(v.callPath().Push("OneOf", uint64SliceToInterfaceSlice(value)))
}
