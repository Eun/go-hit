// Package yaegi_template is a package that provides a templating engine using yeagi.
// yaegi is a golang interpreter and can be used to run go code inside an go application.
//
// Example usage:
//    package main
//
//    import (
//        "os"
//
//        "github.com/Eun/yaegi-template"
//    )
//
//    func main() {
//        template := yaegi_template.MustNew(yaegi_template.DefaultOptions(), yaegi_template.DefaultSymbols()...)
//        template.MustParseString(`
//    <html>
//    <$
//        import "time"
//        func GreetUser(name string) {
//            fmt.Printf("Hello %s, it is %s", name, time.Now().Format(time.Kitchen))
//        }
//    $>
//
//    <p>
//    <$
//        if context.LoggedIn {
//            GreetUser(context.UserName)
//        }
//    $>
//    </p>
//    </html>
//    `)
//
//        type Context struct {
//            LoggedIn bool
//            UserName string
//        }
//
//        template.MustExec(os.Stdout, &Context{
//            LoggedIn: true,
//            UserName: "Joe Doe",
//        })
//    }
package yaegi_template

import (
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/traefik/yaegi/stdlib"

	"reflect"

	"bytes"

	"fmt"

	"sync"

	"go/parser"
	"go/scanner"
	"go/token"

	"go/ast"

	"github.com/traefik/yaegi/interp"

	"github.com/Eun/yaegi-template/codebuffer"
)

// Template represents a template.
type Template struct {
	options        interp.Options
	use            interp.Exports
	imports        importSymbols
	templateReader io.Reader
	StartTokens    []rune
	EndTokens      []rune
	interp         *interp.Interpreter
	outputBuffer   *outputBuffer
	codeBuffer     *codebuffer.CodeBuffer
	mu             sync.Mutex
}

// DefaultOptions return the default options for the New and MustNew functions.
func DefaultOptions() interp.Options {
	return interp.Options{
		GoPath:    os.Getenv("GOPATH"),
		BuildTags: nil,
		Stdin:     nil,
		Stdout:    nil,
		Stderr:    nil,
	}
}

// DefaultSymbols return the default symbols for the New and MustNew functions.
func DefaultSymbols() []interp.Exports {
	return []interp.Exports{stdlib.Symbols}
}

// New creates a new Template that can be used in a later time.
func New(
	options interp.Options, //nolint:gocritic // disable hugeParam: options is heavy
	use ...interp.Exports) (*Template, error) {
	t := &Template{
		options:     options,
		use:         mergeExports(use...),
		StartTokens: []rune("<$"),
		EndTokens:   []rune("$>"),
	}
	return t, nil
}

// MustNew is like New, except it panics on failure.
func MustNew(
	options interp.Options, //nolint:gocritic // disable hugeParam: options is heavy
	use ...interp.Exports) *Template {
	t, err := New(options, use...)
	if err != nil {
		panic(err.Error())
	}
	return t
}

// Parse parses the specified reader, after success it is possible to call Exec() on the template.
func (t *Template) Parse(reader io.Reader) error {
	if err := t.LazyParse(reader); err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	it, err := t.codeBuffer.Iterator()
	if err != nil {
		return err
	}

	// parse everything now
	for it.Next() {
	}
	return it.Error()
}

// MustParse is like Parse, except it panics on failure.
func (t *Template) MustParse(r io.Reader) *Template {
	if err := t.Parse(r); err != nil {
		panic(err.Error())
	}
	return t
}

// LazyParse parses the specified reader during usage of Exec().
func (t *Template) LazyParse(reader io.Reader) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	// maybe in the future we parse the template here
	// for now we don't
	t.templateReader = reader

	t.outputBuffer = newOutputBuffer(true)
	t.codeBuffer = codebuffer.New(reader, t.StartTokens, t.EndTokens)
	t.options.Stdout = t.outputBuffer

	t.interp = interp.New(t.options)

	// if we already have some uses
	// use them
	if len(t.use) != 0 {
		if err := t.interp.Use(t.use); err != nil {
			return errors.Wrap(err, "unable to use")
		}
	}

	// if we already have some imports
	// import them
	if len(t.imports) != 0 {
		if _, err := t.safeEval(t.imports.ImportBlock()); err != nil {
			return err
		}
	}

	return nil
}

// MustLazyParse is like LazyParse, except it panics on failure.
func (t *Template) MustLazyParse(r io.Reader) *Template {
	if err := t.LazyParse(r); err != nil {
		panic(err.Error())
	}
	return t
}

// ParseString parses the specified string, after success it is possible to call Exec() on the template.
func (t *Template) ParseString(s string) error {
	return t.Parse(bytes.NewReader([]byte(s)))
}

// ParseBytes parses the specified byte slice, after success it is possible to call Exec() on the template.
func (t *Template) ParseBytes(b []byte) error {
	return t.Parse(bytes.NewReader(b))
}

// MustParseString is like ParseString, except it panics on failure.
func (t *Template) MustParseString(s string) *Template {
	if err := t.ParseString(s); err != nil {
		panic(err.Error())
	}
	return t
}

// MustParseBytes is like ParseBytes, except it panics on failure.
func (t *Template) MustParseBytes(b []byte) error {
	if err := t.ParseBytes(b); err != nil {
		panic(err.Error())
	}
	return nil
}

// Exec executes the template, and writes the output to the specified writer.
func (t *Template) Exec(writer io.Writer, context interface{}) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.codeBuffer == nil {
		return 0, errors.New("template was never parsed")
	}

	it, err := t.codeBuffer.Iterator()
	if err != nil {
		return 0, err
	}

	total := 0

	for it.Next() {
		part := it.Value()
		switch part.Type {
		case codebuffer.CodePartType:
			n, err := t.execCode(string(part.Content), writer, context)
			if err != nil {
				return total, err
			}
			if n > 0 {
				total += n
			}
		case codebuffer.TextPartType:
			var n int
			if writer != nil {
				n, err = writer.Write(part.Content)
			}
			if err != nil {
				return total, err
			}
			if n > 0 {
				total += n
			}
		}
	}

	return total, it.Error()
}

// MustExec is like Exec, except it panics on failure.
func (t *Template) MustExec(writer io.Writer, context interface{}) {
	if _, err := t.Exec(writer, context); err != nil {
		panic(err.Error())
	}
}

func (t *Template) execCode(code string, out io.Writer, context interface{}) (int, error) {
	if err := t.evalImports(&code); err != nil {
		return 0, err
	}
	if context != nil {
		// do we need to
		err := t.interp.Use(interp.Exports{
			"internal/internal": map[string]reflect.Value{
				"context": reflect.ValueOf(context),
			},
		})
		if err != nil {
			return 0, errors.Wrapf(err, "unable to use context")
		}

		// always reimport internal
		if _, err := t.safeEval(`import . "internal"`); err != nil {
			return 0, err
		}
	}

	// make sure the buffer is empty
	t.outputBuffer.DiscardWrites(false)
	res, err := t.safeEval(code)
	if err != nil {
		return 0, err
	}

	if t.outputBuffer.Length() == 0 {
		// implicit write
		fmt.Fprint(t.outputBuffer, printValue(res))
	}
	var n int
	if out != nil {
		n, err = out.Write(t.outputBuffer.Bytes())
	}
	t.outputBuffer.DiscardWrites(true)
	t.outputBuffer.Reset()
	return n, err
}

func (t *Template) safeEval(code string) (res reflect.Value, err error) {
	if strings.TrimSpace(code) == "" {
		return reflect.Value{}, nil
	}

	defer func() {
		e := recover()
		if e == nil {
			return
		}
		switch v := e.(type) {
		case error:
			err = v
		default:
			err = fmt.Errorf("%v", v)
		}
	}()

	res, err = t.interp.Eval(code)
	if err != nil {
		return res, err
	}
	return res, err
}

func printValue(v reflect.Value) string {
	if !v.IsValid() || !v.CanInterface() {
		return ""
	}

	switch x := v.Interface().(type) {
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprint(x)
	case string:
		return x
	default:
		return ""
	}
}

// evalImports finds all "import" lines evaluates them and removes them from the code.
func (t *Template) evalImports(code *string) error {
	var ok bool
	ok, err := t.hasPackage(*code)
	if err != nil {
		return err
	}
	var c string
	if !ok {
		c = "package main\n" + *code
	} else {
		c = *code
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", c, parser.ImportsOnly)
	if err != nil {
		return err
	}

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.IMPORT {
			continue
		}

		syms := make(importSymbols, 0, len(genDecl.Specs))
		for _, spec := range genDecl.Specs {
			importSpec, ok := spec.(*ast.ImportSpec)
			if !ok {
				continue
			}

			sym := Import{
				Name: "",
				Path: strings.TrimFunc(importSpec.Path.Value, func(r rune) bool {
					return r == '`' || r == '"'
				}),
			}

			if importSpec.Name != nil {
				sym.Name = importSpec.Name.Name
			}

			syms = append(syms, sym)
		}

		if err := t.Import(syms...); err != nil {
			return err
		}

		pos := int(genDecl.Pos()) - 1
		end := int(genDecl.End()) - 1
		c = c[:pos] + strings.Repeat(" ", end-pos) + c[end:]
	}

	// remove package main\n
	*code = c[f.Name.End():]

	return nil
}

// hasPackage returns true when the code has a 'package' line.
func (*Template) hasPackage(s string) (bool, error) {
	_, err := parser.ParseFile(token.NewFileSet(), "", s, parser.PackageClauseOnly)
	if err != nil {
		errList, ok := err.(scanner.ErrorList)
		if !ok {
			return false, err
		}
		if len(errList) == 0 {
			return false, err
		}
		if !strings.HasPrefix(errList[0].Msg, fmt.Sprintf("expected '%s', found", token.PACKAGE)) {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

// Import imports the specified imports to the interpreter.
func (t *Template) Import(imports ...Import) error {
	var symbolsToImport importSymbols
	for _, symbol := range imports {
		if !t.imports.Contains(symbol) {
			symbolsToImport = append(symbolsToImport, symbol)
		}
	}

	if len(symbolsToImport) == 0 {
		return nil
	}

	if t.interp != nil { // if we have an interpreter, import right now
		if _, err := t.safeEval(symbolsToImport.ImportBlock()); err != nil {
			return err
		}
	}
	t.imports = append(t.imports, symbolsToImport...)
	return nil
}

// MustImport is like Import, except it panics on failure.
func (t *Template) MustImport(imports ...Import) *Template {
	if err := t.Import(imports...); err != nil {
		panic(err)
	}
	return t
}

// Use loads binary runtime symbols in the interpreter context so
// they can be used in interpreted code.
func (t *Template) Use(values ...interp.Exports) error {
	return t.useExports(mergeExports(values...))
}

func (t *Template) useExports(values interp.Exports) error {
	if len(values) == 0 {
		return nil
	}

	t.use = mergeExports(t.use, values)
	// if we have an interpreter, use right now
	if t.interp != nil {
		if err := t.interp.Use(t.use); err != nil {
			return errors.Wrap(err, "unable to use")
		}
	}
	return nil
}

// MustUse is like Use, except it panics on failure.
func (t *Template) MustUse(values ...interp.Exports) *Template {
	if err := t.Use(values...); err != nil {
		panic(err)
	}
	return t
}

func mergeExports(values ...interp.Exports) interp.Exports {
	result := make(map[string]*map[string]reflect.Value)
	for i := range values {
		for packageName, funcMap := range values[i] {
			existingFuncMap, ok := result[packageName]
			if !ok {
				m := make(map[string]reflect.Value)
				existingFuncMap = &m
				result[packageName] = existingFuncMap
			}
			for funcName, funcReference := range funcMap {
				(*existingFuncMap)[funcName] = funcReference
			}
		}
	}
	r := make(interp.Exports, len(result))
	for s, m := range result {
		r[s] = *m
	}
	return r
}
