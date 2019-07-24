package convert

import (
	"fmt"
	"reflect"
	"strings"
)

type Error struct {
	src              *convertValue
	dst              *convertValue
	underlayingError error
}

func (e *Error) Error() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "unable to convert %s to %s", e.src.getHumanName(), e.dst.getHumanName())
	if e.underlayingError != nil {
		fmt.Fprintf(&sb, ": %s", e.underlayingError.Error())
	}
	return sb.String()
}

type InvalidTypeError struct {
	src              reflect.Value
	underlayingError error
}

func (e *InvalidTypeError) Error() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "unable to get type for %s", e.src.Type().String())
	if e.underlayingError != nil {
		fmt.Fprintf(&sb, ": %s", e.underlayingError.Error())
	}
	return sb.String()
}
