// +build !generate

// ⚠️⚠️⚠️ This file was autogenerated by generators/clear/clear ⚠️⚠️⚠️ //
package hit

import errortrace "github.com/Eun/go-hit/errortrace"

// IClearExpectBodyJSONJQ provides methods to clear steps.
type IClearExpectBodyJSONJQ interface {
	IStep
	// Contains clears all matching Contains steps
	Contains(value ...interface{}) IStep
	// Equal clears all matching Equal steps
	Equal(value ...interface{}) IStep
	// JQ clears all matching JQ steps
	JQ(value ...string) IClearExpectBodyJSONJQ
	// Len clears all matching Len steps
	Len() IClearExpectInt
	// NotContains clears all matching NotContains steps
	NotContains(value ...interface{}) IStep
	// NotEqual clears all matching NotEqual steps
	NotEqual(value ...interface{}) IStep
}
type clearExpectBodyJSONJQ struct {
	cp callPath
	tr *errortrace.ErrorTrace
}

func newClearExpectBodyJSONJQ(cp callPath) IClearExpectBodyJSONJQ {
	return &clearExpectBodyJSONJQ{cp: cp, tr: ett.Prepare()}
}
func (v *clearExpectBodyJSONJQ) trace() *errortrace.ErrorTrace {
	return v.tr
}
func (*clearExpectBodyJSONJQ) when() StepTime {
	return CleanStep
}
func (v *clearExpectBodyJSONJQ) callPath() callPath {
	return v.cp
}
func (v *clearExpectBodyJSONJQ) exec(hit *hitImpl) error {
	if err := removeSteps(hit, v.callPath()); err != nil {
		return err
	}
	return nil
}
func (v *clearExpectBodyJSONJQ) Contains(value ...interface{}) IStep {
	return removeStep(v.callPath().Push("Contains", value))
}
func (v *clearExpectBodyJSONJQ) Equal(value ...interface{}) IStep {
	return removeStep(v.callPath().Push("Equal", value))
}
func (v *clearExpectBodyJSONJQ) JQ(value ...string) IClearExpectBodyJSONJQ {
	return newClearExpectBodyJSONJQ(v.callPath().Push("JQ", stringSliceToInterfaceSlice(value)))
}
func (v *clearExpectBodyJSONJQ) Len() IClearExpectInt {
	return newClearExpectInt(v.callPath().Push("Len", nil))
}
func (v *clearExpectBodyJSONJQ) NotContains(value ...interface{}) IStep {
	return removeStep(v.callPath().Push("NotContains", value))
}
func (v *clearExpectBodyJSONJQ) NotEqual(value ...interface{}) IStep {
	return removeStep(v.callPath().Push("NotEqual", value))
}
