// +build !generate_numeric

// ⚠️⚠️⚠️ This file was autogenerated by generators/expect/numeric ⚠️⚠️⚠️ //
package hit

import minitest "github.com/Eun/go-hit/internal/minitest"

// IExpectUint16 provides assertions for the uint16 type.
type IExpectUint16 interface {
	// Equal expects the uint16 to be equal to the specified value.
	Equal(value uint16) IStep

	// NotEqual expects the uint16 to be not equal to the specified value.
	NotEqual(value ...uint16) IStep

	// OneOf expects the uint16 to be equal to one of the specified values.
	OneOf(values ...uint16) IStep

	// NotOneOf expects the uint16 to be not equal to one of the specified values.
	NotOneOf(values ...uint16) IStep

	// GreaterThan expects the uint16 to be not greater than the specified value.
	GreaterThan(value uint16) IStep

	// LessThan expects the uint16 to be less than the specified value.
	LessThan(value uint16) IStep

	// GreaterOrEqualThan expects the uint16 to be greater or equal than the specified value.
	GreaterOrEqualThan(value uint16) IStep

	// LessOrEqualThan expects the uint16 to be less or equal than the specified value.
	LessOrEqualThan(value uint16) IStep

	// Between expects the uint16 to be between the specified min and max value (inclusive, min >= uint16 >= max).
	Between(min, max uint16) IStep

	// NotBetween expects the uint16 to be not between the specified min and max value (inclusive, min >= uint16 >= max).
	NotBetween(min, max uint16) IStep
}
type expectUint16ValueCallback func(hit Hit) uint16
type expectUint16 struct {
	cp            callPath
	valueCallback expectUint16ValueCallback
}

func newExpectUint16(cp callPath, valueCallback expectUint16ValueCallback) IExpectUint16 {
	return &expectUint16{cp: cp, valueCallback: valueCallback}
}
func (v *expectUint16) Equal(value uint16) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("Equal", []interface{}{value}), Exec: func(hit *hitImpl) error {
		return minitest.Error.Equal(v.valueCallback(hit), value)
	}}
}
func (v *expectUint16) NotEqual(values ...uint16) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("NotEqual", uint16SliceToInterfaceSlice(values)), Exec: func(hit *hitImpl) error {
		return minitest.Error.NotEqual(v.valueCallback(hit), uint16SliceToInterfaceSlice(values)...)
	}}
}
func (v *expectUint16) OneOf(values ...uint16) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("OneOf", uint16SliceToInterfaceSlice(values)), Exec: func(hit *hitImpl) error {
		return minitest.Error.OneOf(v.valueCallback(hit), uint16SliceToInterfaceSlice(values)...)
	}}
}
func (v *expectUint16) NotOneOf(values ...uint16) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("NotOneOf", uint16SliceToInterfaceSlice(values)), Exec: func(hit *hitImpl) error {
		return minitest.Error.NotOneOf(v.valueCallback(hit), uint16SliceToInterfaceSlice(values)...)
	}}
}
func (v *expectUint16) GreaterThan(value uint16) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("GreaterThan", []interface{}{value}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l <= value {
			return minitest.Error.Errorf("expected %d to be greater than %d", l, value)
		}
		return nil
	}}
}
func (v *expectUint16) LessThan(value uint16) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("LessThan", []interface{}{value}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l >= value {
			return minitest.Error.Errorf("expected %d to be less than %d", l, value)
		}
		return nil
	}}
}
func (v *expectUint16) GreaterOrEqualThan(value uint16) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("GreaterOrEqualThan", []interface{}{value}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l < value {
			return minitest.Error.Errorf("expected %d to be greater or equal than %d", l, value)
		}
		return nil
	}}
}
func (v *expectUint16) LessOrEqualThan(value uint16) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("LessOrEqualThan", []interface{}{value}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l > value {
			return minitest.Error.Errorf("expected %d to be less or equal than %d", l, value)
		}
		return nil
	}}
}
func (v *expectUint16) Between(min, max uint16) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("Between", []interface{}{min, max}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l < min || l > max {
			return minitest.Error.Errorf("expected %d to be between %d and %d", l, min, max)
		}
		return nil
	}}
}
func (v *expectUint16) NotBetween(min, max uint16) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("NotBetween", []interface{}{min, max}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l >= min && l <= max {
			return minitest.Error.Errorf("expected %d not to be between %d and %d", l, min, max)
		}
		return nil
	}}
}
