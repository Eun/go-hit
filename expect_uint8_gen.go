// +build !generate_numeric

// ⚠️⚠️⚠️ This file was autogenerated by generators/expect/numeric ⚠️⚠️⚠️ //
package hit

import minitest "github.com/Eun/go-hit/internal/minitest"

// IExpectUint8 provides assertions for the uint8 type.
type IExpectUint8 interface {
	// Equal expects the uint8 to be equal to the specified value.
	Equal(value uint8) IStep

	// NotEqual expects the uint8 to be not equal to the specified value.
	NotEqual(value ...uint8) IStep

	// OneOf expects the uint8 to be equal to one of the specified values.
	OneOf(values ...uint8) IStep

	// NotOneOf expects the uint8 to be not equal to one of the specified values.
	NotOneOf(values ...uint8) IStep

	// GreaterThan expects the uint8 to be not greater than the specified value.
	GreaterThan(value uint8) IStep

	// LessThan expects the uint8 to be less than the specified value.
	LessThan(value uint8) IStep

	// GreaterOrEqualThan expects the uint8 to be greater or equal than the specified value.
	GreaterOrEqualThan(value uint8) IStep

	// LessOrEqualThan expects the uint8 to be less or equal than the specified value.
	LessOrEqualThan(value uint8) IStep

	// Between expects the uint8 to be between the specified min and max value (inclusive, min >= uint8 >= max).
	Between(min, max uint8) IStep

	// NotBetween expects the uint8 to be not between the specified min and max value (inclusive, min >= uint8 >= max).
	NotBetween(min, max uint8) IStep
}
type expectUint8ValueCallback func(hit Hit) uint8
type expectUint8 struct {
	cp            callPath
	valueCallback expectUint8ValueCallback
}

func newExpectUint8(cp callPath, valueCallback expectUint8ValueCallback) IExpectUint8 {
	return &expectUint8{cp: cp, valueCallback: valueCallback}
}
func (v *expectUint8) Equal(value uint8) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("Equal", []interface{}{value}), Exec: func(hit *hitImpl) error {
		return minitest.Error.Equal(v.valueCallback(hit), value)
	}}
}
func (v *expectUint8) NotEqual(values ...uint8) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("NotEqual", uint8SliceToInterfaceSlice(values)), Exec: func(hit *hitImpl) error {
		return minitest.Error.NotEqual(v.valueCallback(hit), uint8SliceToInterfaceSlice(values)...)
	}}
}
func (v *expectUint8) OneOf(values ...uint8) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("OneOf", uint8SliceToInterfaceSlice(values)), Exec: func(hit *hitImpl) error {
		return minitest.Error.OneOf(v.valueCallback(hit), uint8SliceToInterfaceSlice(values)...)
	}}
}
func (v *expectUint8) NotOneOf(values ...uint8) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("NotOneOf", uint8SliceToInterfaceSlice(values)), Exec: func(hit *hitImpl) error {
		return minitest.Error.NotOneOf(v.valueCallback(hit), uint8SliceToInterfaceSlice(values)...)
	}}
}
func (v *expectUint8) GreaterThan(value uint8) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("GreaterThan", []interface{}{value}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l <= value {
			return minitest.Error.Errorf("expected %d to be greater than %d", l, value)
		}
		return nil
	}}
}
func (v *expectUint8) LessThan(value uint8) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("LessThan", []interface{}{value}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l >= value {
			return minitest.Error.Errorf("expected %d to be less than %d", l, value)
		}
		return nil
	}}
}
func (v *expectUint8) GreaterOrEqualThan(value uint8) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("GreaterOrEqualThan", []interface{}{value}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l < value {
			return minitest.Error.Errorf("expected %d to be greater or equal than %d", l, value)
		}
		return nil
	}}
}
func (v *expectUint8) LessOrEqualThan(value uint8) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("LessOrEqualThan", []interface{}{value}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l > value {
			return minitest.Error.Errorf("expected %d to be less or equal than %d", l, value)
		}
		return nil
	}}
}
func (v *expectUint8) Between(min, max uint8) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("Between", []interface{}{min, max}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l < min || l > max {
			return minitest.Error.Errorf("expected %d to be between %d and %d", l, min, max)
		}
		return nil
	}}
}
func (v *expectUint8) NotBetween(min, max uint8) IStep {
	return &hitStep{Trace: ett.Prepare(), When: ExpectStep, CallPath: v.cp.Push("NotBetween", []interface{}{min, max}), Exec: func(hit *hitImpl) error {
		l := v.valueCallback(hit)
		if l >= min && l <= max {
			return minitest.Error.Errorf("expected %d not to be between %d and %d", l, min, max)
		}
		return nil
	}}
}
