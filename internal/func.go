package internal

import "reflect"

func GetGenericFunc(v interface{}) reflect.Value {
	r := GetValue(v)
	if !r.IsValid() {
		return r
	}
	if r.Kind() != reflect.Func {
		return reflect.Value{} // return an invalid func
	}
	return r
}

func CallGenericFunc(f reflect.Value) {
	t := f.Type()
	in := make([]reflect.Value, t.NumIn())
	for i := range in {
		in[i] = reflect.Zero(t.In(i))
	}
	f.Call(in)
}
