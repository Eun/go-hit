package main

import (
	"errors"
	"reflect"

	"github.com/Eun/go-hit/generators/helpers"

	"fmt"
	"strings"

	"unicode"

	"regexp"

	"github.com/dave/jennifer/jen"
	"golang.org/x/xerrors"

	"github.com/Eun/go-hit"
)

// maximum recursion to generate (e.g JQ().JQ().JQ()...)
const maxDepth = 5

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
var hitStepType = reflect.TypeOf((*hit.IStep)(nil)).Elem()
var interfaceType = reflect.TypeOf((*interface{})(nil)).Elem()

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func isExported(method *reflect.Method) bool {
	return unicode.IsUpper([]rune(method.Name)[0])
}

// isSeparator reports whether the rune could mark a word boundary.
// TODO: update when package unicode captures more of the properties.
// nolint:gocritic,gomnd
func isSeparator(r rune) bool {
	// ASCII alphanumerics and underscore are not separators
	if r <= 0x7F {
		switch {
		case '0' <= r && r <= '9':
			return false
		case 'a' <= r && r <= 'z':
			return false
		case 'A' <= r && r <= 'Z':
			return false
		case r == '_':
			return false
		}
		return true
	}
	// Letters and digits are not separators
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	// Otherwise, all we can do for now is treat spaces as separators.
	return unicode.IsSpace(r)
}

func lowerFirstRune(s string) string {
	// Use a closure here to remember state.
	// Hackish but effective. Depends on Map scanning in order and calling
	// the closure once per rune.
	prev := ' '
	return strings.Map(
		func(r rune) rune {
			if isSeparator(prev) {
				prev = r
				return unicode.ToLower(r)
			}
			prev = r
			return r
		},
		s)
}

func getTypeName(t reflect.Type) string {
	n := t.String()
	if strings.HasPrefix(n, "hit.") {
		return t.Name()
	}
	return n
}

func makeMethodHeader(method *reflect.Method) jen.Code {
	if !isExported(method) {
		return nil
	}

	if method.Type.NumOut() != 1 {
		return nil
	}

	stmt := jen.Id(method.Name)

	var param jen.Code

	if numIn := method.Type.NumIn(); numIn > 1 {
		p := jen.Id("value")
		lastType := method.Type.In(numIn - 1)

		if method.Type.IsVariadic() {
			lastType = lastType.Elem()
		}

		param = p.Id("..." + getTypeName(lastType))
	}

	stmt = stmt.Params(param)

	if method.Type.Out(0) == hitStepType {
		return stmt.Id("IStep")
	}

	return stmt.Id("I" + strings.Title("clear"+method.Type.Out(0).Name()[1:]))
}

func makeSliceToInterfaceFunc(p reflect.Type) string {
	var sb strings.Builder

	var sliceCount int
	for p.Kind() == reflect.Slice {
		p = p.Elem()
		sliceCount++
	}

	t := strings.FieldsFunc(getTypeName(p), func(r rune) bool {
		return r == '.'
	})
	last := t[len(t)-1]

	sb.WriteString(lowerFirstRune(last))

	for i := 0; i < sliceCount; i++ {
		sb.WriteString("Slice")
	}

	sb.WriteString("SliceToInterfaceSlice")
	return sb.String()
}

type generateOpts struct {
	CallPath []string
	Func     reflect.Value
	FuncOut  reflect.Value
}

func generateForStruct(generatedFiles *[]string, options generateOpts) error {
	rfn := options.Func
	rfnStruct := options.FuncOut

	if rfn.Kind() != reflect.Func {
		return xerrors.Errorf("not a func was a %s", rfn.Kind().String())
	}

	f := jen.NewFile("hit")
	f.HeaderComment("+build !generate")

	f.Comment("⚠️⚠️⚠️ This file was autogenerated by generators/clear/clear ⚠️⚠️⚠️ //")

	name := "clear" + options.Func.Type().Out(0).Name()[1:]
	fileName := fmt.Sprintf("%s_gen.go", toSnakeCase(name))

	for _, s := range *generatedFiles {
		if fileName == s {
			return nil
		}
	}

	// create the interface
	methods := make([]jen.Code, rfnStruct.NumMethod())
	for i := 0; i < rfnStruct.Type().NumMethod(); i++ {
		rm := rfnStruct.Type().Method(i)
		m := makeMethodHeader(&rm)
		if m == nil {
			continue
		}
		methods[i] = jen.Comment(fmt.Sprintf("%[1]s clears all matching %[1]s steps", rfnStruct.Type().Method(i).Name)).Line().Add(m)
	}

	f.Line()

	f.Commentf("%s provides methods to clear steps.", "I"+strings.Title(name))

	f.Type().Id("I" + strings.Title(name)).Interface(append([]jen.Code{jen.Id("IStep")}, methods...)...)

	// create struct
	f.Type().Id(name).Struct(
		jen.Id("cp").Id("callPath"),
		jen.Id("tr").Op("*").Qual("github.com/Eun/go-hit/errortrace", "ErrorTrace"),
	)

	// create the constructor
	constructorParams := []jen.Code{
		jen.Id("cp").Id("callPath"),
	}

	numConstructorIn := rfn.Type().NumIn()
	if numConstructorIn > 1 {
		return xerrors.Errorf("func has more than one parameter and this is unsupported")
	}

	f.Func().
		Id(fmt.Sprintf("new%s", strings.Title(name))).
		Params(constructorParams...).
		Id("I" + strings.Title(name)).Block(
		jen.Return(jen.Op("&").Id(name).Values(
			jen.Id("cp").Op(":").Id("cp"),
			jen.Id("tr").Op(":").Id("ett").Dot("Prepare").Call(),
		)),
	)

	// create struct helper functions

	// trace func
	f.Func().Params(jen.Id("v").Op("*").Id(name)).Id("trace").Params().Op("*").Qual("github.com/Eun/go-hit/errortrace", "ErrorTrace").Block(
		jen.Return(jen.Id("v").Dot("tr")),
	)

	// when func
	f.Func().Params(jen.Op("*").Id(name)).Id("when").Params().Id("StepTime").Block(
		jen.Return(jen.Id("cleanStep")),
	)

	// callPath func
	f.Func().Params(jen.Id("v").Op("*").Id(name)).Id("callPath").Params().Id("callPath").Block(
		jen.Return(jen.Id("v").Dot("cp")),
	)

	// exec func
	f.Func().Params(jen.Id("v").Op("*").Id(name)).Id("exec").Params(jen.Id("hit").Op("*").Id("hitImpl")).Error().Block(
		jen.If(jen.Err().Op(":=").Id("removeSteps").Call(jen.Id("hit"), jen.Id("v").Dot("callPath").Call()).Op(";").Err().Op("!=").Nil().Block(
			// jen.Return(jen.Id("v").Dot("trace").Dot("Format").Call(jen.Id("hit").Dot("Description").Call(), jen.Err().Dot("Error").Call())),
			jen.Return(jen.Err()),
		)),
		jen.Return(jen.Nil()),
	)

	// create struct functions
	for i := 0; i < rfnStruct.Type().NumMethod(); i++ {
		rm := rfnStruct.Type().Method(i)
		m := makeMethodHeader(&rm)
		if m == nil {
			continue
		}
		var block jen.Code

		subName := "clear" + rm.Type.Out(0).Name()[1:]

		if rm.Type.Out(0) == hitStepType { //nolint:nestif //ignore if rm.Type.Out(0) == hitStepType` is deeply nested
			clearPathSecondArg := jen.Nil()
			if numIn := rm.Type.NumIn(); numIn > 1 {
				lastType := rm.Type.In(numIn - 1)
				if rm.Type.IsVariadic() {
					lastType = lastType.Elem()
				}

				switch lastType {
				case interfaceType:
					clearPathSecondArg = jen.Id("value")
				default:
					clearPathSecondArg = jen.Id(makeSliceToInterfaceFunc(lastType)).Call(jen.Id("value"))
				}
			}

			block = jen.Return(
				jen.Id("removeStep").Call(
					jen.Id("v").Dot("callPath").Call().Dot("Push").Call(jen.Lit(rm.Name), clearPathSecondArg),
				),
			)
		} else {
			clearPathSecondArg := jen.Nil()
			if numIn := rm.Type.NumIn(); numIn > 1 {
				lastType := rm.Type.In(numIn - 1)
				if rm.Type.IsVariadic() {
					lastType = lastType.Elem()
				}

				if lastType.Kind() == reflect.Interface {
					clearPathSecondArg = jen.Id("value")
				} else {
					clearPathSecondArg = jen.Id(makeSliceToInterfaceFunc(lastType)).Call(jen.Id("value"))
				}
			}

			block = jen.Return(
				jen.Id(fmt.Sprintf("new%s", strings.Title(subName))).Call(
					jen.Id("v").Dot("callPath").Call().Dot("Push").Call(jen.Lit(rm.Name), clearPathSecondArg),
				),
			)
		}

		f.Func().Params(
			jen.Id("v").Op("*").Id(name),
		).Add(m).Block(block)
	}

	if err := helpers.WriteJenFile(fileName, f); err != nil {
		return err
	}
	*generatedFiles = append(*generatedFiles, fileName)
	return nil
}

func isCallPathExported(v []string) bool {
	i := len(v)
	if i == 0 {
		return false
	}
	return unicode.IsUpper([]rune(v[i-1])[0])
}

func collectForFunc(callPath []string, fv reflect.Value, depth int) ([]generateOpts, error) {
	if depth > maxDepth {
		return nil, nil
	}
	if fv.Kind() != reflect.Func {
		return nil, fmt.Errorf("not a func was a %s", fv.Kind().String())
	}

	if fv.Type().NumOut() != 1 {
		return nil, errors.New("not one out")
	}
	if fv.Type().Out(0) == hitStepType {
		return nil, errors.New("is a hit step")
	}

	if !isCallPathExported(callPath) {
		return nil, errors.New("not exported")
	}

	self := generateOpts{
		Func: fv,
	}

	size := fv.Type().NumIn()
	args := make([]reflect.Value, size)
	for i := 0; i < size; i++ {
		args[i] = reflect.Zero(fv.Type().In(i))
	}
	if fv.Type().IsVariadic() {
		self.FuncOut = fv.CallSlice(args)[0].Elem()
	} else {
		self.FuncOut = fv.Call(args)[0].Elem()
	}

	self.CallPath = make([]string, len(callPath))
	copy(self.CallPath, callPath)

	if !self.FuncOut.IsValid() {
		return nil, errors.New("invalid result")
	}

	result := []generateOpts{self}
	result = append(result, collectForStruct(callPath, self.FuncOut, depth+1)...)
	return result, nil
}

func collectForStruct(callPath []string, sv reflect.Value, depth int) []generateOpts {
	var result []generateOpts

	for i := 0; i < sv.NumMethod(); i++ {
		m := sv.Method(i)
		cp := append(callPath, sv.Type().Method(i).Name)
		sub, err := collectForFunc(cp, m, depth)
		if err != nil {
			continue
		}
		result = append(result, sub...)
	}
	return result
}

func generateClearExpect() {
	generations, err := collectForFunc([]string{"Expect"}, reflect.ValueOf(hit.Expect), 0)
	if err != nil {
		panic(err)
	}
	var generatedFiles []string
	for _, g := range generations {
		if err := generateForStruct(&generatedFiles, g); err != nil {
			panic(err)
		}
	}
}

func generateClearSend() {
	generations, err := collectForFunc([]string{"Send"}, reflect.ValueOf(hit.Send), 0)
	if err != nil {
		panic(err)
	}
	var generatedFiles []string
	for _, g := range generations {
		if err := generateForStruct(&generatedFiles, g); err != nil {
			panic(err)
		}
	}
}

func main() {
	generateClearSend()
	generateClearExpect()
}
