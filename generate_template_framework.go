// +build -ignore

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"strings"

	"sort"

	"golang.org/x/tools/imports"
)

func main() {
	var sb strings.Builder
	fmt.Fprintln(&sb, `// +build -ignore`)
	fmt.Fprintln(&sb)
	fmt.Fprintln(&sb, `// You can use this file as an template to build your own framework. Just change / add the functions you need.`)
	fmt.Fprintln(&sb, `// See also examples/extensibility`)
	fmt.Fprintln(&sb)
	fmt.Fprintln(&sb, `package main`)
	fmt.Fprintln(&sb, `import "github.com/Eun/go-hit"`)

	set := token.NewFileSet()
	pkgs, err := parser.ParseDir(set, ".", nil, parser.ParseComments)
	if err != nil {
		log.Fatal("Failed to parse package:", err)
	}

	types := collectTypes(pkgs, "hit")

	decls := getDecls(pkgs, "hit")

	for i := 0; i < len(decls); i++ {
		d := decls[i]
		fn, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}
		// skip over receiver functions
		if fn.Recv != nil {
			continue
		}
		// skip over unexported functions
		if !fn.Name.IsExported() {
			continue
		}
		if strings.HasPrefix(fn.Name.Name, "Test") && len(fn.Name.Name) > 4 {
			continue
		}

		if strings.HasPrefix(fn.Name.Name, "Benchmark") && len(fn.Name.Name) > 9 {
			continue
		}

		if fn.Doc != nil {
			for i := 0; i < len(fn.Doc.List); i++ {
				fmt.Fprintln(&sb, fn.Doc.List[i].Text)
			}
		} else {
			fmt.Fprintln(&sb)
		}

		fmt.Fprintf(&sb, "func %s", fn.Name)

		fmt.Fprint(&sb, buildParams(fn, types))
		fmt.Fprint(&sb, buildReturns(fn, types))

		fmt.Fprintln(&sb, "{")

		if fn.Type.Results != nil {
			fmt.Fprint(&sb, "return")
		}

		fmt.Fprintf(&sb, " hit.%s", fn.Name)
		fmt.Fprintln(&sb, buildCall(fn, types))
		fmt.Fprintln(&sb, "}")
	}

	buf, err := imports.Process("", []byte(sb.String()), nil)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("template_framework.go", buf, 0)
	if err != nil {
		log.Fatal(err)
	}
}

func collectTypes(pkgs map[string]*ast.Package, limitToPackage string) (types []string) {
	// keep track of all types we use
	decls := getDecls(pkgs, limitToPackage)
	for _, d := range decls {
		decl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, s := range decl.Specs {
			typeSpec, ok := s.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if !typeSpec.Name.IsExported() {
				continue
			}
			types = append(types, typeSpec.Name.String())
		}
	}

	return types
}

func buildParams(fn *ast.FuncDecl, types []string) string {
	if fn.Type.Params == nil {
		return "()"
	}
	size := len(fn.Type.Params.List)
	if size <= 0 {
		return "()"
	}
	var sb strings.Builder

	set := token.NewFileSet()

	params := make([]string, size)
	for i := 0; i < size; i++ {
		sb.Reset()
		if err := printer.Fprint(&sb, set, fn.Type.Params.List[i].Type); err != nil {
			log.Fatal(err)
		}
		params[i] = sb.String()
		params[i] = expandType(params[i], "hit", types)

		if fn.Type.Params.List[i].Names != nil {
			names := make([]string, len(fn.Type.Params.List[i].Names))
			for j, name := range fn.Type.Params.List[i].Names {
				names[j] = name.String()
			}
			params[i] = fmt.Sprintf("%s %s", strings.Join(names, ", "), params[i])
		}
	}

	sb.Reset()
	sb.WriteString(`(`)
	sb.WriteString(strings.Join(params, ","))
	sb.WriteString(`)`)
	return sb.String()
}

func buildReturns(fn *ast.FuncDecl, types []string) string {
	if fn.Type.Results == nil {
		return ""
	}
	var sb strings.Builder
	size := len(fn.Type.Results.List)
	if size <= 0 {
		return ""
	}

	set := token.NewFileSet()

	returns := make([]string, size)
	for i := 0; i < size; i++ {
		sb.Reset()
		if err := printer.Fprint(&sb, set, fn.Type.Results.List[i].Type); err != nil {
			log.Fatal(err)
		}
		returns[i] = sb.String()
		returns[i] = expandType(returns[i], "hit", types)
		if fn.Type.Results.List[i].Names != nil && len(fn.Type.Results.List[i].Names) >= 1 {
			returns[i] = fmt.Sprintf("%s %s", fn.Type.Results.List[i].Names[0].Name, returns[i])
		}
	}

	sb.Reset()

	if size > 1 {
		sb.WriteString(`(`)
	}
	sb.WriteString(strings.Join(returns, ","))
	if size > 1 {
		sb.WriteString(`)`)
	}
	return sb.String()
}

func buildCall(fn *ast.FuncDecl, types []string) string {
	if fn.Type.Params == nil {
		return "()"
	}
	var sb strings.Builder
	size := len(fn.Type.Params.List)
	if size <= 0 {
		return "()"
	}

	params := make([]string, size)
	for i := 0; i < size; i++ {
		if fn.Type.Params.List[i].Names == nil || len(fn.Type.Params.List[i].Names) <= 0 {
			log.Fatalf("missing parameter name for %s", fn.Name.String())
		}

		names := make([]string, len(fn.Type.Params.List[i].Names))
		for j, name := range fn.Type.Params.List[i].Names {
			names[j] = name.String()
		}
		params[i] = strings.Join(names, ", ")
		if _, isEllipsis := fn.Type.Params.List[i].Type.(*ast.Ellipsis); isEllipsis {
			params[i] = params[i] + "..."
		}
	}

	sb.WriteString(`(`)
	sb.WriteString(strings.Join(params, ","))
	sb.WriteString(`)`)
	return sb.String()
}

func expandType(typ string, pkg string, types []string) string {
	ellipsis := ""
	if strings.HasPrefix(typ, "...") {
		typ = typ[3:]
		ellipsis = "..."
	}
	for _, t := range types {
		if t == typ {
			return fmt.Sprintf("%s%s.%s", ellipsis, pkg, typ)
		}
	}

	if typ == "func(Hit)" {
		typ = "func(hit.Hit)"
	}

	return ellipsis + typ
}

func getDecls(pkgs map[string]*ast.Package, pkgName string) (decl []ast.Decl) {
	for _, pkg := range pkgs {
		if pkg.Name != pkgName {
			continue
		}
		for _, f := range pkg.Files {
			decl = append(decl, f.Decls...)
		}
	}

	sort.Slice(decl, func(i, j int) bool {
		return strings.Compare(fmt.Sprintf("%s", decl[i]), fmt.Sprintf("%s", decl[j])) < 0
	})

	return decl
}
