package main

import (
	"bytes"
	"errors"
	"log"
	"reflect"

	"github.com/otto-eng/go-hit/generators/helpers"

	"fmt"
	"strings"

	"unicode"

	"io"

	"github.com/dave/jennifer/jen"
	"golang.org/x/xerrors"

	"github.com/otto-eng/go-hit"
)

// maximum recursion to generate (e.g JQ().JQ().JQ()...)
const maxDepth = 5

var hitStepType = reflect.TypeOf((*hit.IStep)(nil)).Elem()

func getDefaultValueRepresentation(t reflect.Type, isVariadic bool) string {
	v := reflect.Zero(t)

	switch v.Interface().(type) {
	case string:
		return `"Foo-Bar"`
	case []string:
		if isVariadic {
			return `"Foo", "Bar"`
		}
		return `[]string{"Foo", "Bar"}`
	case uint, uint8, uint16, uint32, uint64:
		return `0x2`
	case int, int8, int16, int32, int64:
		return `2`
	case []uint, []uint8, []uint16, []uint32, []uint64:
		if isVariadic {
			return `0x1, 0x2`
		}
		return fmt.Sprintf(`%s{0x1, 0x2}`, v.Type().String())
	case []int, []int8, []int16, []int32, []int64:
		if isVariadic {
			return `1, 2`
		}
		return fmt.Sprintf(`%s{1, 2}`, v.Type().String())
	case bool:
		return `false`
	case float32, float64:
		return `1.000000`
	case []float32, []float64:
		if isVariadic {
			return `1.000000, 2.000000`
		}
		return fmt.Sprintf(`%s{1.000000, 2.000000}`, v.Type().String())
	case hit.Callback:
		return `func(hit Hit)error{return nil}`
	case []interface{}:
		if isVariadic {
			return `"Foo", "Baz"`
		}
		return `[]interface{}{"Foo", "Baz"}`
	default:
		if t.Implements(reflect.TypeOf((*io.Reader)(nil)).Elem()) {
			return `bytes.NewReader(nil)`
		}
		if v.Kind() == reflect.Interface {
			return `"Foo-Taz"`
		}
		panic(fmt.Errorf("got no sample for %s", v.Type()))
	}
}

func getSampleValueRepresentation(t reflect.Type, isVariadic bool) string {
	v := reflect.Zero(t)

	switch v.Interface().(type) {
	case string:
		return `"Hello-World"`
	case []string:
		if isVariadic {
			return `"Hello", "World"`
		}
		return `[]string{"Hello", "World"}`
	case uint, uint8, uint16, uint32, uint64:
		return `0x3`
	case int, int8, int16, int32, int64:
		return `3`
	case []uint, []uint8, []uint16, []uint32, []uint64:
		if isVariadic {
			return `0x3, 0x4`
		}
		return fmt.Sprintf(`%s{0x3, 0x4}`, v.Type().String())
	case []int, []int8, []int16, []int32, []int64:
		if isVariadic {
			return `3, 4`
		}
		return fmt.Sprintf(`%s{3, 4}`, v.Type().String())
	case bool:
		return `true`
	case float32, float64:
		return `3.000000`
	case []float32, []float64:
		if isVariadic {
			return `3.000000, 4.000000`
		}
		return fmt.Sprintf(`%s{3.000000, 4.000000}`, v.Type().String())
	case hit.Callback:
		return `func(hit Hit)error{panic("Err")}`
	case []interface{}:
		if isVariadic {
			return `"Hello", "Earth"`
		}
		return `[]interface{}{"Hello", "Earth"}`
	default:
		if t.Implements(reflect.TypeOf((*io.Reader)(nil)).Elem()) {
			return `bytes.NewReader([]byte{1, 2, 3})`
		}
		if v.Kind() == reflect.Interface {
			return `"Hello-Universe"`
		}
		panic(fmt.Errorf("got no sample for %s", v.Type()))
	}
}

type call struct {
	Name       string
	Args       []jen.Code
	SampleArgs []jen.Code
}

type callPath []call

func (c callPath) Join(sep string) string {
	s := make([]string, len(c))
	for i := 0; i < len(c); i++ {
		s[i] = c[i].Name
	}
	return strings.Join(s, sep)
}

type generateOpts struct {
	CallPath callPath
	Func     reflect.Value
	FuncOut  reflect.Value
}

func generateForStruct(f *jen.File, options generateOpts) error {
	rfn := options.Func
	rfnStruct := options.FuncOut

	if rfn.Kind() != reflect.Func {
		return xerrors.Errorf("not a func was a %s", rfn.Kind().String())
	}

	var defaultStmtLine *jen.Statement
	var sampleStmtLine *jen.Statement

	var genericLine *jen.Statement
	clearGenericLine := jen.Id("Clear").Call()
	clearDefaultLine := jen.Id("Clear").Call()

	for _, s := range options.CallPath {
		if genericLine == nil {
			genericLine = jen.Id(s.Name).Call()
		} else {
			genericLine = genericLine.Dot(s.Name).Call()
		}
		clearGenericLine = clearGenericLine.Dot(s.Name).Call()
		clearDefaultLine = clearDefaultLine.Dot(s.Name).Call(s.Args...)
		if defaultStmtLine == nil {
			defaultStmtLine = jen.Id(s.Name).Call(s.Args...)
			sampleStmtLine = jen.Id(s.Name).Call(s.SampleArgs...)
		} else {
			defaultStmtLine = defaultStmtLine.Dot(s.Name).Call(s.Args...)
			sampleStmtLine = sampleStmtLine.Dot(s.Name).Call(s.SampleArgs...)
		}
	}

	specificTest := true

	lastArg := make([]reflect.Type, rfn.Type().NumIn())
	for i := 0; i < len(lastArg); i++ {
		lastArg[i] = rfn.Type().In(i)
	}

	if rfn.Type().Out(0) != hitStepType {
		specificTest = false
		stmtValue := rfnStruct

	next:
		if stmtValue.NumMethod() == 0 {
			return xerrors.Errorf("no functions present in %s", stmtValue.Type().String())
		}
		m := stmtValue.Type().Method(0)
		stmtValue = stmtValue.Method(0)

		args, jenArgs, sampleArgs := getArgs(stmtValue.Type())
		defaultStmtLine = defaultStmtLine.Dot(m.Name).Call(jenArgs...)
		sampleStmtLine = sampleStmtLine.Dot(m.Name).Call(sampleArgs...)
		if stmtValue.Type().IsVariadic() {
			stmtValue = stmtValue.CallSlice(args)[0].Elem()
		} else {
			stmtValue = stmtValue.Call(args)[0].Elem()
		}
		if m.Type.Out(0) != hitStepType {
			goto next
		}
	}

	if s := len(lastArg); s > 0 {
		// filter out functions with hit.Callback as parameter
		if lastArg[s-1] == reflect.TypeOf(new(hit.Callback)).Elem() {
			return nil
		}

		// filter out functions with io.Reader as parameter
		if lastArg[s-1] == reflect.TypeOf(new(io.Reader)).Elem() {
			return nil
		}
	}

	var genericLineRepresentation bytes.Buffer
	if err := genericLine.Render(&genericLineRepresentation); err != nil {
		return err
	}

	f.Func().
		Id(fmt.Sprintf("TestGenClear_Generic_%s", options.CallPath.Join(""))).
		Params(jen.Id("t").Op("*").Qual("testing", "T")).
		Block(
			jen.Id("s").Op(":=").Id("EchoServer").Call(),
			jen.Defer().Id("s").Dot("Close").Call(),
			jen.Var().Id("steps").Id("[]IStep"),
			jen.Id("ExpectError").Call(
				jen.Id("t"),
				jen.Line().Id("Do").Call(
					jen.Line().Id("Post").Call(jen.Id("s").Dot("URL")),
					jen.Line().Add(defaultStmtLine),
					jen.Line().Add(sampleStmtLine),
					jen.Line().Id("storeSteps").Call(jen.Op("&").Id("steps")),
					jen.Line().Add(clearGenericLine),
					jen.Line().Id("expectSteps").Call(jen.Id("t"), jen.Op("&").Id("steps"), jen.Id("2")),
				),
				jen.Line().Id("PtrStr").Call(jen.Lit("TestOK")),
				jen.Line(),
			),

			jen.Comment("test unable to find a step"),
			jen.Id("ExpectError").Call(
				jen.Id("t"),
				jen.Line().Id("Do").Call(
					jen.Line().Id("Post").Call(jen.Id("s").Dot("URL")),
					jen.Line().Add(clearGenericLine),
				),
				jen.Line().Id("PtrStr").Call(jen.Lit("unable to find a step with "+genericLineRepresentation.String())),
				jen.Line().Id("PtrStr").Call(jen.Lit("got these steps:")),
				jen.Line().Id("PtrStr").Call(jen.Lit("Post()")),
				jen.Line(),
			),
		)

	if specificTest {
		f.Func().
			Id(fmt.Sprintf("TestGenClear_Specific_%s", options.CallPath.Join(""))).
			Params(jen.Id("t").Op("*").Qual("testing", "T")).
			Block(
				jen.Id("s").Op(":=").Id("EchoServer").Call(),
				jen.Defer().Id("s").Dot("Close").Call(),
				jen.Var().Id("steps").Id("[]IStep"),
				jen.Id("ExpectError").Call(
					jen.Id("t"),
					jen.Line().Id("Do").Call(
						jen.Line().Id("Post").Call(jen.Id("s").Dot("URL")),
						jen.Line().Add(defaultStmtLine),
						jen.Line().Add(sampleStmtLine),
						jen.Line().Id("storeSteps").Call(jen.Op("&").Id("steps")),
						jen.Line().Add(clearDefaultLine),
						jen.Line().Id("expectSteps").Call(jen.Id("t"), jen.Op("&").Id("steps"), jen.Id("1")),
					),
					jen.Line().Id("PtrStr").Call(jen.Lit("TestOK")),
					jen.Line(),
				),
			)
	}

	return nil
}

func isCallPathExported(v callPath) bool {
	i := len(v)
	if i == 0 {
		return false
	}
	return unicode.IsUpper([]rune(v[i-1].Name)[0])
}

func getArgs(t reflect.Type) (args []reflect.Value, argsJen, sampleArgs []jen.Code) {
	size := t.NumIn()
	args = make([]reflect.Value, size)
	argsJen = make([]jen.Code, size)
	sampleArgs = make([]jen.Code, size)
	for i := 0; i < size; i++ {
		args[i] = reflect.Zero(t.In(i))
		argsJen[i] = jen.Op(getDefaultValueRepresentation(t.In(i), t.IsVariadic() && i == size-1))
		sampleArgs[i] = jen.Op(getSampleValueRepresentation(t.In(i), t.IsVariadic() && i == size-1))
	}
	return args, argsJen, sampleArgs
}

func collectForFunc(cp callPath, fv reflect.Value, depth int) ([]generateOpts, error) {
	if depth > maxDepth {
		return nil, nil
	}
	if fv.Kind() != reflect.Func {
		return nil, fmt.Errorf("not a func was a %s", fv.Kind().String())
	}

	if fv.Type().NumOut() != 1 {
		return nil, errors.New("not one out")
	}

	if !isCallPathExported(cp) {
		return nil, errors.New("not exported")
	}

	self := generateOpts{
		Func: fv,
	}

	args, argsRepresentation, sampleArgs := getArgs(fv.Type())
	if fv.Type().IsVariadic() {
		self.FuncOut = fv.CallSlice(args)[0].Elem()
	} else {
		self.FuncOut = fv.Call(args)[0].Elem()
	}

	self.CallPath = make([]call, len(cp))
	for i := 0; i < len(cp); i++ {
		self.CallPath[i].Name = cp[i].Name
		self.CallPath[i].Args = make([]jen.Code, len(cp[i].Args))
		copy(self.CallPath[i].Args, cp[i].Args)
		self.CallPath[i].SampleArgs = make([]jen.Code, len(cp[i].SampleArgs))
		copy(self.CallPath[i].SampleArgs, cp[i].SampleArgs)
	}
	self.CallPath[len(cp)-1].Args = argsRepresentation
	self.CallPath[len(cp)-1].SampleArgs = sampleArgs

	if !self.FuncOut.IsValid() {
		return nil, errors.New("invalid result")
	}

	result := []generateOpts{self}
	result = append(result, collectForStruct(cp, self.FuncOut, depth+1)...)
	return result, nil
}

func collectForStruct(cp callPath, sv reflect.Value, depth int) []generateOpts {
	var result []generateOpts

	for i := 0; i < sv.NumMethod(); i++ {
		m := sv.Method(i)
		_, argsRepresentation, sampleArgs := getArgs(sv.Method(i).Type())
		sub, err := collectForFunc(append(cp, call{
			Name:       sv.Type().Method(i).Name,
			Args:       argsRepresentation,
			SampleArgs: sampleArgs,
		}), m, depth)
		if err != nil {
			continue
		}

		result = append(result, sub...)
	}
	return result
}

func generateClearSend(f *jen.File) {
	generations, err := collectForFunc([]call{{"Send", nil, nil}}, reflect.ValueOf(hit.Send), 0)
	if err != nil {
		panic(err)
	}
	for _, g := range generations {
		if err := generateForStruct(f, g); err != nil {
			panic(err)
		}
	}
}

func generateClearExpect(f *jen.File) {
	generations, err := collectForFunc([]call{{"Expect", nil, nil}}, reflect.ValueOf(hit.Expect), 0)
	if err != nil {
		panic(err)
	}
	for _, g := range generations {
		if err := generateForStruct(f, g); err != nil {
			panic(err)
		}
	}
}

func main() {
	f := jen.NewFile("hit_test")

	f.Op(`import . "github.com/otto-eng/go-hit"`)
	f.Op(`import "github.com/stretchr/testify/require"`)
	f.Op(`import "errors"`)

	f.Comment("⚠️⚠️⚠️ This file was autogenerated by generators/clear/tests ⚠️⚠️⚠️ //")

	// helper func

	f.Op(`
func storeSteps(steps *[]IStep) IStep {
	return Custom(StepTime(CleanStep), func(hit Hit) error {
		*steps = hit.Steps()
		return nil
	})
}
func expectSteps(t *testing.T, expectSteps *[]IStep, removedStepsCount int) IStep {
	return Custom(BeforeExpectStep, func(hit Hit) error {
		require.Len(t, hit.Steps(), len(*expectSteps)-removedStepsCount)
		return errors.New("TestOK")
	})
}
`)

	generateClearSend(f)
	generateClearExpect(f)

	if err := helpers.WriteJenFile("clear_gen_test.go", f); err != nil {
		log.Fatal(err)
	}
}
