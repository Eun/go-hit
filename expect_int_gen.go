// +build !generate_numeric

package hit

import minitest "github.com/Eun/go-hit/internal/minitest"

// ⚠️⚠️⚠️ This file was autogenerated by generators/expect/numeric ⚠️⚠️⚠️ //

// IExpectInt provides assertions for the int type.
type IExpectInt interface {
	// Equal expects the int to be equal to the specified value.
	Equal(value int) IStep

	// NotEqual expects the int to be not equal to the specified value.
	NotEqual(value ...int) IStep

	// OneOf expects the int to be equal to one of the specified values.
	OneOf(values ...int) IStep

	// NotOneOf expects the int to be not equal to one of the specified values.
	NotOneOf(values ...int) IStep

	// GreaterThan expects the int to be not greater than the specified value.
	GreaterThan(value int) IStep

	// LessThan expects the int to be less than the specified value.
	LessThan(value int) IStep

	// GreaterOrEqualThan expects the int to be greater or equal than the specified value.
	GreaterOrEqualThan(value int) IStep

	// LessOrEqualThan expects the int to be less or equal than the specified value.
	LessOrEqualThan(value int) IStep

	// Between expects the int to be between the specified min and max value (inclusive, min >= int >= max).
	Between(min, max int) IStep

	// NotBetween expects the int to be not between the specified min and max value (inclusive, min >= int >= max).
	NotBetween(min, max int) IStep
}
type expectIntValueCallback func(hit Hit) int
type expectInt struct {
	cp            callPath
	valueCallback expectIntValueCallback
}

func newExpectInt(cp callPath, valueCallback expectIntValueCallback) IExpectInt {
	return &expectInt{cp: cp, valueCallback: valueCallback}
}
func (v *expectInt) Equal(value int) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("Equal", []interface{}{value}), Exec: func(hit *hitImpl) error {
		return minitest.Equal(v.valueCallback(hit), value)
	}}
}
func (v *expectInt) NotEqual(values ...int) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("NotEqual", intSliceToInterfaceSlice(values)), Exec: func(hit *hitImpl) error {
		return minitest.NotEqual(v.valueCallback(hit), intSliceToInterfaceSlice(values)...)
	}}
}
func (v *expectInt) OneOf(values ...int) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("OneOf", intSliceToInterfaceSlice(values)), Exec: func(hit *hitImpl) error {
		return minitest.OneOf(v.valueCallback(hit), intSliceToInterfaceSlice(values)...)
	}}
}
func (v *expectInt) NotOneOf(values ...int) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("NotOneOf", intSliceToInterfaceSlice(values)), Exec: func(hit *hitImpl) error {
		return minitest.NotOneOf(v.valueCallback(hit), intSliceToInterfaceSlice(values)...)
	}}
}
func (v *expectInt) GreaterThan(value int) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("GreaterThan", []interface{}{value}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l <= value {
			return minitest.Errorf("expected %d to be greater than %d", l, value)
		}
		return nil
	}}
}
func (v *expectInt) LessThan(value int) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("LessThan", []interface{}{value}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l >= value {
			return minitest.Errorf("expected %d to be less than %d", l, value)
		}
		return nil
	}}
}
func (v *expectInt) GreaterOrEqualThan(value int) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("GreaterOrEqualThan", []interface{}{value}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l < value {
			return minitest.Errorf("expected %d to be greater or equal than %d", l, value)
		}
		return nil
	}}
}
func (v *expectInt) LessOrEqualThan(value int) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("LessOrEqualThan", []interface{}{value}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l > value {
			return minitest.Errorf("expected %d to be less or equal than %d", l, value)
		}
		return nil
	}}
}
func (v *expectInt) Between(min, max int) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("Between", []interface{}{min, max}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l < min || l > max {
			return minitest.Errorf("expected %d to be between %d and %d", l, min, max)
		}
		return nil
	}}
}
func (v *expectInt) NotBetween(min, max int) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("NotBetween", []interface{}{min, max}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l >= min && l <= max {
			return minitest.Errorf("expected %d not to be between %d and %d", l, min, max)
		}
		return nil
	}}
}
