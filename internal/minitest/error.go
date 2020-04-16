package minitest

import (
	"fmt"
	"reflect"

	"strings"

	"errors"

	"github.com/Eun/go-hit/internal/minitest/contains"
	"github.com/Eun/go-hit/internal/misc"
	"github.com/google/go-cmp/cmp"
)

type ReturnError struct{}

func makeError(err string, customMessageAndArgs ...interface{}) error {
	var sb strings.Builder

	if detail := formatMessage(customMessageAndArgs); detail != "" {
		fmt.Fprintln(&sb, detail)
	}

	if err != "" {
		sb.WriteString(err)
	}

	return errors.New(sb.String())
}

func (ReturnError) Errorf(messageAndArgs ...interface{}) error {
	return makeError(formatMessage(messageAndArgs))
}

func (ReturnError) NoError(err error, customMessageAndArgs ...interface{}) error {
	if err != nil {
		return makeError(err.Error(), customMessageAndArgs...)
	}
	return nil
}

func (ReturnError) Equal(expected, actual interface{}, customMessageAndArgs ...interface{}) error {
	if !cmp.Equal(expected, actual) {
		return makeError(stringJoin("\n", "Not equal", actualExpectedDiff(actual, expected)), customMessageAndArgs...)
	}
	return nil
}

func (ReturnError) NotEqual(expected, actual interface{}, customMessageAndArgs ...interface{}) error {
	if cmp.Equal(expected, actual) {
		return makeError(stringJoin("\n", fmt.Sprintf("should not be %s", PrintValue(actual))), customMessageAndArgs...)
	}
	return nil
}

func (ReturnError) Contains(object interface{}, value interface{}, customMessageAndArgs ...interface{}) error {
	if !contains.Contains(object, value) {
		return makeError(fmt.Sprintf(`%s does not contain %s`, PrintValue(object), PrintValue(value)), customMessageAndArgs...)
	}
	return nil
}

func (ReturnError) NotContains(object interface{}, value interface{}, customMessageAndArgs ...interface{}) error {
	if contains.Contains(object, value) {
		return makeError(fmt.Sprintf(`%s should not contain %s`, PrintValue(object), PrintValue(value)), customMessageAndArgs...)
	}
	return nil
}

func (ReturnError) Empty(object interface{}, customMessageAndArgs ...interface{}) error {
	v := misc.GetValue(object)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		l := v.Len()
		if l != 0 {
			return makeError(fmt.Sprintf(`%s should be empty, but has %d item(s)`, PrintValue(object), l), customMessageAndArgs...)
		}
		return nil
	default:
		return makeError(fmt.Sprintf("called Len() on %s", PrintValue(object)))
	}
}

func (ReturnError) Len(object interface{}, length int, customMessageAndArgs ...interface{}) error {
	v := misc.GetValue(object)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		l := v.Len()
		if l != length {
			return makeError(fmt.Sprintf(`%s should have %d item(s), but has %d`, PrintValue(object), length, l), customMessageAndArgs...)
		}
		return nil
	default:
		return makeError(fmt.Sprintf("called Len() on %s", PrintValue(object)))
	}
}

func (ReturnError) True(value bool, customMessageAndArgs ...interface{}) error {
	if !value {
		return makeError(`Expected bool to be true but is false`, customMessageAndArgs...)
	}
	return nil
}

func (ReturnError) False(value bool, customMessageAndArgs ...interface{}) error {
	if value {
		return makeError(`Expected bool to be false but is true`, customMessageAndArgs...)
	}
	return nil
}
