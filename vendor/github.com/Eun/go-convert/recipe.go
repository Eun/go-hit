package convert

import (
	"fmt"
	"reflect"
)

// Recipe represents a recipe that defines which type can be converted into which type
// and which function should be called to convert this type
type Recipe struct {
	From reflect.Type
	To   reflect.Type
	Func func(c Converter, in reflect.Value, out reflect.Value) error
}

// MakeRecipe makes a recipe for the passed in function.
//
// Note that the functions must match this format:
//     func(Converter, INVALUE, *OUTVALUE) (error)
//
// Example:
//    // convert int to bool
//    MakeRecipe(func(c Converter, in int, out *bool) error {
//        if in == 0 {
//            *out = false
//            return nil
//        }
//        *out = true
//        return nil
//    })
func MakeRecipe(fn interface{}) (Recipe, error) {
	v := reflect.ValueOf(fn)

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

// MakeRecipes makes a recipe slice for the passed in functions.
// see also MakeRecipe
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

// MustMakeRecipe makes a recipe for the passed in functions, but panics on error.
// see also MakeRecipe
func MustMakeRecipe(f interface{}) Recipe {
	r, err := MakeRecipe(f)
	if err != nil {
		panic(err)
	}
	return r
}

// MustMakeRecipes makes a recipe slice for the passed in functions, but panics on error.
// see also MakeRecipe
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
