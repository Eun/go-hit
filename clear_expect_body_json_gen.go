//go:build !generate
// +build !generate

package hit

import errortrace "github.com/Eun/go-hit/errortrace"

// ⚠️⚠️⚠️ This file was autogenerated by generators/clear/clear ⚠️⚠️⚠️ //

// IClearExpectBodyJSON provides methods to clear steps.
type IClearExpectBodyJSON interface {
	IStep
	// Contains clears all matching Contains steps
	Contains(value ...interface{}) IStep
	// Dasel clears all matching Dasel steps
	Dasel(value ...string) IClearExpectBodyJSONDasel
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
type clearExpectBodyJSON struct {
	cp callPath
	tr *errortrace.ErrorTrace
}

func newClearExpectBodyJSON(cp callPath) IClearExpectBodyJSON {
	return &clearExpectBodyJSON{cp: cp, tr: ett.Prepare()}
}
func (v *clearExpectBodyJSON) trace() *errortrace.ErrorTrace {
	return v.tr
}
func (*clearExpectBodyJSON) when() StepTime {
	return cleanStep
}
func (v *clearExpectBodyJSON) callPath() callPath {
	return v.cp
}
func (v *clearExpectBodyJSON) exec(hit *hitImpl) error {
	if err := removeSteps(hit, v.callPath()); err != nil {
		return err
	}
	return nil
}
func (v *clearExpectBodyJSON) Contains(value ...interface{}) IStep {
	return removeStep(v.callPath().Push("Contains", value))
}
func (v *clearExpectBodyJSON) Dasel(value ...string) IClearExpectBodyJSONDasel {
	return newClearExpectBodyJSONDasel(v.callPath().Push("Dasel", stringSliceToInterfaceSlice(value)))
}
func (v *clearExpectBodyJSON) Equal(value ...interface{}) IStep {
	return removeStep(v.callPath().Push("Equal", value))
}
func (v *clearExpectBodyJSON) JQ(value ...string) IClearExpectBodyJSONJQ {
	return newClearExpectBodyJSONJQ(v.callPath().Push("JQ", stringSliceToInterfaceSlice(value)))
}
func (v *clearExpectBodyJSON) Len() IClearExpectInt {
	return newClearExpectInt(v.callPath().Push("Len", nil))
}
func (v *clearExpectBodyJSON) NotContains(value ...interface{}) IStep {
	return removeStep(v.callPath().Push("NotContains", value))
}
func (v *clearExpectBodyJSON) NotEqual(value ...interface{}) IStep {
	return removeStep(v.callPath().Push("NotEqual", value))
}
