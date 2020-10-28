// +build !generate_numeric

// ⚠️⚠️⚠️ This file was autogenerated by generators/expect/numeric ⚠️⚠️⚠️ //
package hit

import minitest "github.com/Eun/go-hit/internal/minitest"

// IExpectInt64 provides assertions for the int64 type.
type IExpectInt64 interface {
	// Equal expects the int64 to be equal to the specified value.
	Equal(value int64) IStep

	// NotEqual expects the int64 to be not equal to the specified value.
	NotEqual(value ...int64) IStep

	// OneOf expects the int64 to be equal to one of the specified values.
	OneOf(values ...int64) IStep

	// NotOneOf expects the int64 to be not equal to one of the specified values.
	NotOneOf(values ...int64) IStep

	// GreaterThan expects the int64 to be not greater than the specified value.
	GreaterThan(value int64) IStep

	// LessThan expects the int64 to be less than the specified value.
	LessThan(value int64) IStep

	// GreaterOrEqualThan expects the int64 to be greater or equal than the specified value.
	GreaterOrEqualThan(value int64) IStep

	// LessOrEqualThan expects the int64 to be less or equal than the specified value.
	LessOrEqualThan(value int64) IStep

	// Between expects the int64 to be between the specified min and max value (inclusive, min >= int64 >= max).
	Between(min, max int64) IStep

	// NotBetween expects the int64 to be not between the specified min and max value (inclusive, min >= int64 >= max).
	NotBetween(min, max int64) IStep
}
type expectInt64ValueCallback func(hit Hit) int64
type expectInt64 struct {
	cp            callPath
	valueCallback expectInt64ValueCallback
}

func newExpectInt64(cp callPath, valueCallback expectInt64ValueCallback) IExpectInt64 {
	return &expectInt64{cp: cp, valueCallback: valueCallback}
}
func (v *expectInt64) Equal(value int64) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("Equal", []interface{}{value}), Exec: func(hit *hitImpl) error {
		return minitest.Error.Equal(v.valueCallback(hit), value)
	}}
}
func (v *expectInt64) NotEqual(values ...int64) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("NotEqual", int64SliceToInterfaceSlice(values)), Exec: func(hit *hitImpl) error {
		return minitest.Error.NotEqual(v.valueCallback(hit), int64SliceToInterfaceSlice(values)...)
	}}
}
func (v *expectInt64) OneOf(values ...int64) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("OneOf", int64SliceToInterfaceSlice(values)), Exec: func(hit *hitImpl) error {
		return minitest.Error.OneOf(v.valueCallback(hit), int64SliceToInterfaceSlice(values)...)
	}}
}
func (v *expectInt64) NotOneOf(values ...int64) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("NotOneOf", int64SliceToInterfaceSlice(values)), Exec: func(hit *hitImpl) error {
		return minitest.Error.NotOneOf(v.valueCallback(hit), int64SliceToInterfaceSlice(values)...)
	}}
}
func (v *expectInt64) GreaterThan(value int64) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("GreaterThan", []interface{}{value}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l <= value {
			return minitest.Error.Errorf("expected %d to be greater than %d", l, value)
		}
		return nil
	}}
}
func (v *expectInt64) LessThan(value int64) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("LessThan", []interface{}{value}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l >= value {
			return minitest.Error.Errorf("expected %d to be less than %d", l, value)
		}
		return nil
	}}
}
func (v *expectInt64) GreaterOrEqualThan(value int64) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("GreaterOrEqualThan", []interface{}{value}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l < value {
			return minitest.Error.Errorf("expected %d to be greater or equal than %d", l, value)
		}
		return nil
	}}
}
func (v *expectInt64) LessOrEqualThan(value int64) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("LessOrEqualThan", []interface{}{value}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l > value {
			return minitest.Error.Errorf("expected %d to be less or equal than %d", l, value)
		}
		return nil
	}}
}
func (v *expectInt64) Between(min, max int64) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("Between", []interface{}{min, max}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l < min || l > max {
			return minitest.Error.Errorf("expected %d to be between %d and %d", l, min, max)
		}
		return nil
	}}
}
func (v *expectInt64) NotBetween(min, max int64) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("NotBetween", []interface{}{min, max}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l >= min && l <= max {
			return minitest.Error.Errorf("expected %d not to be between %d and %d", l, min, max)
		}
		return nil
	}}
}
