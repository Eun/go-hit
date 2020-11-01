package minitest

import (
	"fmt"
	"reflect"

	"strings"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/xerrors"

	"github.com/Eun/go-hit/internal/converter"
	"github.com/Eun/go-hit/internal/minitest/contains"
	"github.com/Eun/go-hit/internal/misc"
)

func Errorf(format string, messageAndArgs ...interface{}) error {
	return xerrors.New(strings.TrimSpace(fmt.Sprintf(format, messageAndArgs...)))
}

func NoError(err error) error {
	if err != nil {
		return xerrors.New(err.Error())
	}
	return nil
}

func Equal(object, expected interface{}) error {
	// shortcuts
	if expected == nil && object == nil {
		return nil
	}

	if (expected == nil && object != nil) || (expected != nil && object == nil) {
		return xerrors.New(stringJoin("\n", "not equal", actualExpectedDiff(object, expected)))
	}

	// we might be able to convert this
	compareData := misc.MakeTypeCopy(expected)

	err := converter.Convert(object, &compareData)
	if err == nil {
		object = compareData
	}
	if !cmp.Equal(expected, object) {
		return xerrors.New(stringJoin("\n", "not equal", actualExpectedDiff(object, expected)))
	}
	return nil
}

func NotEqual(object interface{}, values ...interface{}) error {
	for _, value := range values {
		// shortcuts
		if value == nil && object == nil {
			return xerrors.New(stringJoin("\n", fmt.Sprintf("should not be %s", PrintValue(object))))
		}

		if (value == nil && object != nil) || (value != nil && object == nil) {
			return nil
		}
		// we might be able to convert this
		compareData := misc.MakeTypeCopy(value)
		if err := converter.Convert(object, &compareData); err == nil {
			object = compareData
		}
		if cmp.Equal(value, object) {
			return xerrors.New(stringJoin("\n", fmt.Sprintf("should not be %s", PrintValue(object))))
		}
	}

	return nil
}

func Contains(object interface{}, values ...interface{}) error {
	for _, value := range values {
		ok, err := contains.Contains(object, value)
		if err != nil {
			if _, isRecipeNotFoundError := err.(contains.NoRecipeFoundError); isRecipeNotFoundError {
				// is it equal?
				if equalErr := Equal(object, value); equalErr != nil {
					// its not equal
					return xerrors.New(fmt.Sprintf(`%s does not contain %s`, PrintValue(object), PrintValue(value)))
				}
				// its equal, continue
				continue
			}
			return err
		}
		if !ok {
			return xerrors.New(fmt.Sprintf(`%s does not contain %s`, PrintValue(object), PrintValue(value)))
		}
	}
	return nil
}

func NotContains(object interface{}, values ...interface{}) error {
	return NotOneOf(object, values...)
}

func OneOf(object interface{}, values ...interface{}) error {
	for _, value := range values {
		ok, err := contains.Contains(object, value)
		if err != nil {
			if _, isRecipeNotFoundError := err.(contains.NoRecipeFoundError); isRecipeNotFoundError {
				// is it equal?
				if equalErr := NotEqual(object, value); equalErr != nil {
					// its equal so return no error
					return nil
				}
				// its not equal, continue until we find a matching value
				continue
			}
			return err
		}
		if ok {
			return nil
		}
	}
	return xerrors.New(fmt.Sprintf(`%s should be one of %s`, PrintValue(object), PrintValue(values)))
}

func NotOneOf(object interface{}, values ...interface{}) error {
	for _, value := range values {
		ok, err := contains.Contains(object, value)
		if err != nil {
			if _, isRecipeNotFoundError := err.(contains.NoRecipeFoundError); isRecipeNotFoundError {
				// is it equal?
				if equalErr := NotEqual(object, value); equalErr != nil {
					// its equal
					return xerrors.New(fmt.Sprintf(`%s should not contain %s`, PrintValue(object), PrintValue(value)))
				}
				// its not equal
				continue
			}
			return err
		}
		if ok {
			return xerrors.New(fmt.Sprintf(`%s should not contain %s`, PrintValue(object), PrintValue(value)))
		}
	}
	return nil
}

func Empty(object interface{}) error {
	v := misc.GetValue(object)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		l := v.Len()
		if l != 0 {
			return xerrors.New(fmt.Sprintf(`%s should be empty, but has %d item(s)`, PrintValue(object), l))
		}
		return nil
	default:
		return xerrors.New(fmt.Sprintf("called Len() on %s", PrintValue(object)))
	}
}

func NotEmpty(object interface{}) error {
	v := misc.GetValue(object)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		if v.Len() == 0 {
			return xerrors.New(fmt.Sprintf(`%s should be not empty`, PrintValue(object)))
		}
		return nil
	default:
		return xerrors.New(fmt.Sprintf("called Len() on %s", PrintValue(object)))
	}
}

func Len(object interface{}, length int) error {
	v := misc.GetValue(object)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		l := v.Len()
		if l != length {
			return xerrors.New(fmt.Sprintf(`%s should have %d item(s), but has %d`, PrintValue(object), length, l))
		}
		return nil
	default:
		return xerrors.New(fmt.Sprintf("called Len() on %s", PrintValue(object)))
	}
}

func True(value bool) error {
	if !value {
		return xerrors.New(`Expected bool to be true but is false`)
	}
	return nil
}

func False(value bool) error {
	if value {
		return xerrors.New(`Expected bool to be false but is true`)
	}
	return nil
}
