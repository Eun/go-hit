//go:build !generate
// +build !generate

package hit

import errortrace "github.com/Eun/go-hit/errortrace"

// ⚠️⚠️⚠️ This file was autogenerated by generators/clear/clear ⚠️⚠️⚠️ //

// IClearExpectHeaderValue provides methods to clear steps.
type IClearExpectHeaderValue interface {
	IStep
	// Contains clears all matching Contains steps
	Contains(value ...interface{}) IStep
	// Empty clears all matching Empty steps
	Empty() IStep
	// Equal clears all matching Equal steps
	Equal(value ...interface{}) IStep
	// Len clears all matching Len steps
	Len() IClearExpectInt
	// NotContains clears all matching NotContains steps
	NotContains(value ...interface{}) IStep
	// NotEmpty clears all matching NotEmpty steps
	NotEmpty() IStep
	// NotEqual clears all matching NotEqual steps
	NotEqual(value ...interface{}) IStep
	// NotOneOf clears all matching NotOneOf steps
	NotOneOf(value ...interface{}) IStep
	// OneOf clears all matching OneOf steps
	OneOf(value ...interface{}) IStep
}
type clearExpectHeaderValue struct {
	cp callPath
	tr *errortrace.ErrorTrace
}

func newClearExpectHeaderValue(cp callPath) IClearExpectHeaderValue {
	return &clearExpectHeaderValue{cp: cp, tr: ett.Prepare()}
}
func (v *clearExpectHeaderValue) trace() *errortrace.ErrorTrace {
	return v.tr
}
func (*clearExpectHeaderValue) when() StepTime {
	return cleanStep
}
func (v *clearExpectHeaderValue) callPath() callPath {
	return v.cp
}
func (v *clearExpectHeaderValue) exec(hit *hitImpl) error {
	if err := removeSteps(hit, v.callPath()); err != nil {
		return err
	}
	return nil
}
func (v *clearExpectHeaderValue) Contains(value ...interface{}) IStep {
	return removeStep(v.callPath().Push("Contains", value))
}
func (v *clearExpectHeaderValue) Empty() IStep {
	return removeStep(v.callPath().Push("Empty", nil))
}
func (v *clearExpectHeaderValue) Equal(value ...interface{}) IStep {
	return removeStep(v.callPath().Push("Equal", value))
}
func (v *clearExpectHeaderValue) Len() IClearExpectInt {
	return newClearExpectInt(v.callPath().Push("Len", nil))
}
func (v *clearExpectHeaderValue) NotContains(value ...interface{}) IStep {
	return removeStep(v.callPath().Push("NotContains", value))
}
func (v *clearExpectHeaderValue) NotEmpty() IStep {
	return removeStep(v.callPath().Push("NotEmpty", nil))
}
func (v *clearExpectHeaderValue) NotEqual(value ...interface{}) IStep {
	return removeStep(v.callPath().Push("NotEqual", value))
}
func (v *clearExpectHeaderValue) NotOneOf(value ...interface{}) IStep {
	return removeStep(v.callPath().Push("NotOneOf", value))
}
func (v *clearExpectHeaderValue) OneOf(value ...interface{}) IStep {
	return removeStep(v.callPath().Push("OneOf", value))
}
