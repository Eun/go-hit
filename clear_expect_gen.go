// +build !generate

// ⚠️⚠️⚠️ This file was autogenerated by generators/clear/clear ⚠️⚠️⚠️ //
package hit

import errortrace "github.com/Eun/go-hit/errortrace"

// IClearExpect provides methods to clear steps.
type IClearExpect interface {
	IStep
	// Body clears all matching Body steps
	Body() IClearExpectBody
	// Custom clears all matching Custom steps
	Custom(value ...Callback) IStep
	// Headers clears all matching Headers steps
	Headers(value ...string) IClearExpectHeaders
	// Status clears all matching Status steps
	Status() IClearExpectInt64
	// Trailers clears all matching Trailers steps
	Trailers(value ...string) IClearExpectHeaders
}
type clearExpect struct {
	cp callPath
	tr *errortrace.ErrorTrace
}

func newClearExpect(cp callPath) IClearExpect {
	return &clearExpect{cp: cp, tr: ett.Prepare()}
}
func (v *clearExpect) trace() *errortrace.ErrorTrace {
	return v.tr
}
func (*clearExpect) when() StepTime {
	return CleanStep
}
func (v *clearExpect) callPath() callPath {
	return v.cp
}
func (v *clearExpect) exec(hit *hitImpl) error {
	if err := removeSteps(hit, v.callPath()); err != nil {
		return err
	}
	return nil
}
func (v *clearExpect) Body() IClearExpectBody {
	return newClearExpectBody(v.callPath().Push("Body", nil))
}
func (v *clearExpect) Custom(value ...Callback) IStep {
	return removeStep(v.callPath().Push("Custom", callbackSliceToInterfaceSlice(value)))
}
func (v *clearExpect) Headers(value ...string) IClearExpectHeaders {
	return newClearExpectHeaders(v.callPath().Push("Headers", stringSliceToInterfaceSlice(value)))
}
func (v *clearExpect) Status() IClearExpectInt64 {
	return newClearExpectInt64(v.callPath().Push("Status", nil))
}
func (v *clearExpect) Trailers(value ...string) IClearExpectHeaders {
	return newClearExpectHeaders(v.callPath().Push("Trailers", stringSliceToInterfaceSlice(value)))
}
