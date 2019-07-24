package expr

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/Eun/go-hit/internal"
)

type expression struct {
	Parts      []string
	CurrentPos int
	Size       int
}

func (e *expression) Current() string {
	return e.Parts[e.CurrentPos]
}

func (e *expression) End() bool {
	return e.CurrentPos >= e.Size
}

func (e *expression) Next() *expression {
	e.CurrentPos++
	return e
}

func (e *expression) String() string {
	return strings.Join(e.Parts, ".")
}

// MustGetValue finds a value in a map/struct/slice. It panics if an error occurred or the value was not found
func MustGetValue(v interface{}, expr string, opts ...Option) interface{} {
	r, found, err := GetValue(v, expr, opts...)
	if err != nil && !options(opts).HasOption(IgnoreError) {
		panic(err.Error())
	}
	if !found && !options(opts).HasOption(IgnoreNotFound) {
		panic(fmt.Sprintf("unable to find value with %s in %#v", expr, v))
	}
	return r
}

// GetValue finds a value in a map/struct/slice, returns the value and true if the value was found
func GetValue(v interface{}, expr string, opts ...Option) (value interface{}, found bool, err error) {
	if v == nil {
		return nil, false, nil
	}
	return getValue(reflect.ValueOf(v), getExpr(expr), opts)
}

func getExpr(expr string) *expression {
	parts := strings.FieldsFunc(strings.TrimSpace(expr), func(r rune) bool {
		return r == '.'
	})
	return &expression{
		Parts:      parts,
		CurrentPos: 0,
		Size:       len(parts),
	}
}

func getValue(v reflect.Value, expr *expression, opts options) (interface{}, bool, error) {
	r := internal.GetElem(v)
	if !r.IsValid() {
		return nil, false, fmt.Errorf("%s cannot be used with the expression %s", v.Type().String(), expr.String())
	}

	if expr.End() {
		if r.IsValid() && r.CanInterface() {
			return r.Interface(), true, nil
		}
		return nil, true, nil
	}

	switch r.Kind() {
	case reflect.Map:
		return getValueFromMap(r, expr, opts)
	case reflect.Struct:
		return getValueFromStruct(r, expr, opts)
	case reflect.Slice:
		return getValueFromSlice(r, expr, opts)
	default:
		return nil, false, fmt.Errorf("%s cannot be used with the expression %s", r.Type().String(), expr.String())
	}
}

func getValueFromMap(m reflect.Value, expr *expression, opts options) (interface{}, bool, error) {
	// if expr is a number
	n, err := strconv.ParseInt(expr.Current(), 0, 64)
	if err == nil {
		return getValueFromMapByIndex(m, expr, int(n), opts)
	}
	return getValueFromMapByName(m, expr, opts)
}

func getValueFromMapByIndex(m reflect.Value, expr *expression, n int, options options) (interface{}, bool, error) {
	if n < 0 || n >= m.Len() {
		return nil, false, nil
	}
	keys := make([]reflect.Value, m.Len())
	for i, key := range m.MapKeys() {
		r := internal.GetElem(key)
		if !r.IsValid() {
			continue
		}
		if r.Kind() != reflect.String {
			continue
		}
		keys[i] = r
	}

	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i].String(), keys[j].String()) < 0
	})
	return getValue(m.MapIndex(keys[n]), expr.Next(), options)
}

func getValueFromMapByName(m reflect.Value, expr *expression, opts options) (interface{}, bool, error) {
	for _, key := range m.MapKeys() {
		r := internal.GetElem(key)
		if !r.IsValid() {
			continue
		}
		if r.Kind() != reflect.String {
			continue
		}
		if !isKeyEqual(r.String(), expr.Current(), opts) {
			continue
		}
		return getValue(m.MapIndex(key), expr.Next(), opts)
	}
	return nil, false, nil
}

func getValueFromStruct(m reflect.Value, expr *expression, opts options) (interface{}, bool, error) {
	// if expr is a number
	n, err := strconv.ParseInt(expr.Current(), 0, 64)
	if err == nil {
		return getValueFromStructByIndex(m, expr, int(n), opts)
	}
	return getValueFromStructByName(m, expr, opts)
}

func getValueFromStructByIndex(m reflect.Value, expr *expression, n int, opts options) (interface{}, bool, error) {
	if n < 0 || n >= m.NumField() {
		return nil, false, nil
	}

	keys := make([]string, m.NumField())
	for i := m.NumField() - 1; i >= 0; i-- {
		keys[i] = m.Type().Field(i).Name
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i], keys[j]) < 0
	})
	return getValue(m.FieldByName(keys[n]), expr.Next(), opts)
}

func getValueFromStructByName(m reflect.Value, expr *expression, opts options) (interface{}, bool, error) {
	for i := m.NumField() - 1; i >= 0; i-- {
		fieldType := m.Type().Field(i)
		if !isKeyEqual(fieldType.Name, expr.Current(), opts) {
			continue
		}
		return getValue(m.Field(i), expr.Next(), opts)
	}
	return nil, false, nil
}
func getValueFromSlice(m reflect.Value, expr *expression, opts options) (interface{}, bool, error) {
	// if expr is a number
	n, err := strconv.ParseInt(expr.Current(), 0, 64)
	if err == nil {
		return getValueFromSliceByIndex(m, expr, int(n), opts)
	}
	return getValueFromSliceByName(m, expr, opts)
}

func getValueFromSliceByIndex(m reflect.Value, expr *expression, n int, opts options) (interface{}, bool, error) {
	if n < 0 || n >= m.Len() {
		return nil, false, nil
	}
	return getValue(m.Index(n), expr.Next(), opts)
}

func getValueFromSliceByName(m reflect.Value, expr *expression, opts options) (interface{}, bool, error) {
	for i := m.Len() - 1; i >= 0; i-- {
		r := internal.GetElem(m.Index(i))
		if !r.IsValid() {
			continue
		}
		if r.Kind() != reflect.String {
			continue
		}

		if !isKeyEqual(r.String(), expr.Current(), opts) {
			continue
		}
		return getValue(m.Index(i), expr.Next(), opts)
	}
	return nil, false, nil
}

func isKeyEqual(a, b string, opts options) bool {
	if opts.HasOption(IgnoreCase) {
		return strings.EqualFold(a, b)
	}
	return a == b
}
