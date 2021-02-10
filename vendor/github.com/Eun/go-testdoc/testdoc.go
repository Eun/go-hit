// Package testdoc allows you to test your code your code documentation examples during normal test time.
package testdoc

import (
	"go/ast"
	"go/parser"

	"github.com/traefik/yaegi/interp"

	yaegi_template "github.com/Eun/yaegi-template"

	"github.com/pkg/errors"

	"go/token"
	"strings"

	"github.com/hashicorp/go-multierror"
)

//go:generate /usr/bin/env -S docker run -v ${PWD}:/app/ -w /app --rm eunts/minigo:latest -t -o README.md ./internal/readme_template.md

func testDocForFile(fileName string, pos int, doc string, options *Options) error {
	doc = strings.TrimSpace(doc)
	if doc == "" {
		return nil
	}
	d, err := ParseDoc(doc)
	if err != nil {
		return errors.Wrapf(err, "unable to parse doc for %s:%d", fileName, pos)
	}
	var result error
	for field, code := range d.Fields {
		for _, section := range options.Sections {
			if strings.EqualFold(field, section) {
				if err := testDocForSection(options, code); err != nil {
					result = multierror.Append(result, errors.Wrapf(err, "%s:%d", fileName, pos))
				}
				break
			}
		}
	}
	return result
}

func testDocForSection(options *Options, code string) error {
	t, err := yaegi_template.New(*options.Options, options.Symbols...)
	if err != nil {
		return errors.Wrapf(err, "unable to create template")
	}
	t.StartTokens = []rune{}
	t.EndTokens = []rune{}

	if err := t.ParseString(code); err != nil {
		return errors.Wrap(err, "unable to parse code")
	}

	for _, i := range options.Imports {
		if err := t.Import(yaegi_template.Import(i)); err != nil {
			return errors.Wrapf(err, "unable to import %v", i)
		}
	}
	if _, err := t.Exec(nil, nil); err != nil {
		return errors.Wrap(err, "exec failed")
	}
	return nil
}

// Import represents an import that should be evaluated.
type Import yaegi_template.Import

// Options options for TestCodeDocumentation.
type Options struct {
	Path           string
	PkgName        string
	Sections       []string
	Symbols        []interp.Exports
	Options        *interp.Options
	Imports        []Import
	IncludePrivate bool
}

// TestCodeDocumentation can be used to test documentation for functions, structs, interfaces...
func TestCodeDocumentation(t interface {
	Fatalf(string, ...interface{})
	Fatal(...interface{})
}, options *Options) {
	if options == nil {
		t.Fatalf("options cannot be nil")
		return
	}
	if options.Symbols == nil {
		options.Symbols = yaegi_template.DefaultSymbols()
	}
	if options.Options == nil {
		opts := yaegi_template.DefaultOptions()
		options.Options = &opts
	}

	set := token.NewFileSet()
	pkgs, err := parser.ParseDir(set, options.Path, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("failed to parse package: %w", err)
		return
	}

	pkg, ok := pkgs[options.PkgName]
	if !ok {
		t.Fatalf("unable to find %q in %v", options.PkgName, pkgs)
		return
	}

	var result error
	for fileName, f := range pkg.Files {
		var docEntries []docEntry
		var errorSlice []error

		if f.Doc != nil && strings.TrimSpace(f.Doc.Text()) != "" {
			docEntries = append(docEntries, docEntry{
				text:   f.Doc.Text(),
				lineno: set.Position(f.Doc.Pos()).Line,
			})
		}

		ast.Inspect(f, astInspectorFunc(set, &docEntries, &errorSlice, options))

		result = multierror.Append(result, errorSlice...)

		for _, doc := range docEntries {
			if err := testDocForFile(fileName, doc.lineno, doc.text, options); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	if result == nil {
		return
	}
	if merr, ok := result.(*multierror.Error); ok {
		switch merr.Len() {
		case 0:
			return
		case 1:
			t.Fatal(merr.Errors[0])
			return
		}
	}
	t.Fatal(result)
}

type docEntry struct {
	text   string
	lineno int
}

//nolint:dupl,prealloc
func astInspectorFunc(set *token.FileSet, docEntries *[]docEntry, errorSlice *[]error, options *Options) func(n ast.Node) bool {
	return func(n ast.Node) bool {
		var entries []docEntry
		var errs []error
		switch v := n.(type) {
		// function
		case *ast.FuncDecl:
			entries, errs = inspectFuncDecl(set, v, options)
		case *ast.GenDecl:
			entries, errs = inspectGenDecl(set, v, options)
		}
		*docEntries = append(*docEntries, entries...)
		*errorSlice = append(*errorSlice, errs...)
		return true
	}
}

//nolint:dupl,prealloc,unparam
func inspectFuncDecl(set *token.FileSet, decl *ast.FuncDecl, options *Options) ([]docEntry, []error) {
	if !options.IncludePrivate && !decl.Name.IsExported() {
		return nil, nil
	}
	if decl.Doc == nil || strings.TrimSpace(decl.Doc.Text()) == "" {
		return nil, nil
	}

	return []docEntry{
			{
				text:   decl.Doc.Text(),
				lineno: set.Position(decl.Doc.Pos()).Line,
			},
		},
		nil
}

//nolint:dupl,prealloc
func inspectGenDecl(set *token.FileSet, decl *ast.GenDecl, options *Options) ([]docEntry, []error) {
	var entries []docEntry
	var errs []error

	for _, spec := range decl.Specs {
		var subEntries []docEntry
		var subErrors []error
		switch v := spec.(type) {
		case *ast.TypeSpec:
			subEntries, errs = inspectTypeSpec(set, decl, v, options)
		case *ast.ValueSpec:
			subEntries, errs = inspectValueSpec(set, decl, v, options)
		}
		entries = append(entries, subEntries...)
		errs = append(errs, subErrors...)
	}
	return entries, errs
}

//nolint:dupl,prealloc
func inspectTypeSpec(set *token.FileSet, parent *ast.GenDecl, spec *ast.TypeSpec, options *Options) ([]docEntry, []error) {
	var entries []docEntry
	var errs []error

	if options.IncludePrivate || spec.Name.IsExported() {
		if parent.Doc != nil && strings.TrimSpace(parent.Doc.Text()) != "" {
			entries = append(entries, docEntry{
				text:   parent.Doc.Text(),
				lineno: set.Position(parent.Doc.Pos()).Line,
			})
		}
		if spec.Doc != nil && strings.TrimSpace(spec.Doc.Text()) != "" {
			entries = append(entries, docEntry{
				text:   spec.Doc.Text(),
				lineno: set.Position(spec.Doc.Pos()).Line,
			})
		}
	}
	var subEntries []docEntry
	switch v := spec.Type.(type) {
	case *ast.InterfaceType:
		subEntries, errs = inspectInterfaceType(set, v, options)
	case *ast.StructType:
		subEntries, errs = inspectStructType(set, v, options)
	}

	entries = append(entries, subEntries...)
	return entries, errs
}

//nolint:dupl,prealloc
func inspectValueSpec(set *token.FileSet, parent *ast.GenDecl, spec *ast.ValueSpec, options *Options) ([]docEntry, []error) {
	var entries []docEntry
	var errs []error

	if options.IncludePrivate || isOneExported(spec.Names...) {
		if parent.Doc != nil && strings.TrimSpace(parent.Doc.Text()) != "" {
			entries = append(entries, docEntry{
				text:   parent.Doc.Text(),
				lineno: set.Position(parent.Doc.Pos()).Line,
			})
		}
		if spec.Doc != nil && strings.TrimSpace(spec.Doc.Text()) != "" {
			entries = append(entries, docEntry{
				text:   spec.Doc.Text(),
				lineno: set.Position(spec.Doc.Pos()).Line,
			})
		}
	}
	var subEntries []docEntry
	switch v := spec.Type.(type) {
	case *ast.InterfaceType:
		subEntries, errs = inspectInterfaceType(set, v, options)
	case *ast.StructType:
		subEntries, errs = inspectStructType(set, v, options)
	}

	entries = append(entries, subEntries...)
	return entries, errs
}

//nolint:dupl,prealloc
func inspectInterfaceType(set *token.FileSet, i *ast.InterfaceType, options *Options) ([]docEntry, []error) {
	if i.Methods.List == nil {
		return nil, nil
	}
	var entries []docEntry
	var errs []error
	for _, m := range i.Methods.List {
		if !options.IncludePrivate && !isOneExported(m.Names...) {
			continue
		}

		if m.Doc == nil || strings.TrimSpace(m.Doc.Text()) == "" {
			continue
		}

		entries = append(entries, docEntry{
			text:   m.Doc.Text(),
			lineno: set.Position(m.Doc.Pos()).Line,
		})
	}
	return entries, errs
}

//nolint:dupl,prealloc
func inspectStructType(set *token.FileSet, s *ast.StructType, options *Options) ([]docEntry, []error) {
	if s.Fields.List == nil {
		return nil, nil
	}
	var entries []docEntry
	var errs []error
	for _, f := range s.Fields.List {
		if !options.IncludePrivate && !isOneExported(f.Names...) {
			continue
		}

		if f.Doc == nil || strings.TrimSpace(f.Doc.Text()) == "" {
			continue
		}
		entries = append(entries, docEntry{
			text:   f.Doc.Text(),
			lineno: set.Position(f.Doc.Pos()).Line,
		})
	}
	return entries, errs
}

func isOneExported(idents ...*ast.Ident) bool {
	for _, ident := range idents {
		if ident.IsExported() {
			return true
		}
	}
	return false
}
