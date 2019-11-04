package convert

import (
	"fmt"
	"reflect"
)

type Recipe struct {
	From reflect.Type
	To   reflect.Type
	Func func(c Converter, in reflect.Value, out reflect.Value) error
}

// you can pass in
// func(Converter, INVALUE, *OUTVALUE) (error)

func MakeRecipe(f interface{}) (Recipe, error) {
	v := reflect.ValueOf(f)

	if v.Kind() != reflect.Func {
		return Recipe{}, fmt.Errorf("cannot make an recipe from an %v", v.Kind())
	}
	var r Recipe

	converterInterface := reflect.TypeOf((*Converter)(nil)).Elem()
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()

	if v.Type().NumIn() != 3 {
		return Recipe{}, fmt.Errorf("%s has invalid pattern", v.Type().String())
	}
	if !v.Type().In(0).Implements(converterInterface) {
		return Recipe{}, fmt.Errorf("%s has invalid pattern", v.Type().String())
	}
	r.From = v.Type().In(1)
	r.To = v.Type().In(2)

	r.Func = wrapFunc(v)
	if r.To.Kind() != reflect.Ptr {
		return Recipe{}, fmt.Errorf("%s has invalid pattern", v.Type().String())
	}

	if v.Type().NumOut() != 1 {
		return Recipe{}, fmt.Errorf("%s has invalid pattern", v.Type().String())
	}
	if !v.Type().Out(0).Implements(errorInterface) {
		return Recipe{}, fmt.Errorf("%s has invalid pattern", v.Type().String())
	}

	return r, nil
}

func MakeRecipes(f ...interface{}) ([]Recipe, error) {
	recipes := make([]Recipe, len(f))
	for i, v := range f {
		r, err := MakeRecipe(v)
		if err != nil {
			return nil, err
		}
		recipes[i] = r
	}
	return recipes, nil
}

func MustMakeRecipe(f interface{}) Recipe {
	r, err := MakeRecipe(f)
	if err != nil {
		panic(err)
	}
	return r
}
func MustMakeRecipes(f ...interface{}) []Recipe {
	r, err := MakeRecipes(f...)
	if err != nil {
		panic(err)
	}
	return r
}

func wrapFunc(f reflect.Value) func(c Converter, in reflect.Value, out reflect.Value) error {
	return func(c Converter, in reflect.Value, out reflect.Value) error {
		result := f.Call([]reflect.Value{reflect.ValueOf(c), in, out})
		if !result[0].IsNil() {
			return result[0].Interface().(error)
		}
		return nil
	}
}
