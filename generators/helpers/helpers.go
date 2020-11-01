// Package helpers contains some helper functions for the generators.
package helpers

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"

	"github.com/dave/jennifer/jen"
	"golang.org/x/tools/imports"
	"golang.org/x/xerrors"
)

// WriteJenFile writes a generated jen.File to the specified filename, it pretty prints the file and runs goimports
// afterwards.
func WriteJenFile(fileName string, file *jen.File) error {
	var buf bytes.Buffer
	if err := file.Render(&buf); err != nil {
		return err
	}

	// pretty print
	set := token.NewFileSet()
	astFile, err := parser.ParseFile(set, "", buf.String(), parser.ParseComments|parser.AllErrors)
	if err != nil {
		return xerrors.Errorf("Failed to parse package: %w", err)
	}

	buf.Reset()
	if err = printer.Fprint(&buf, set, astFile); err != nil {
		return xerrors.Errorf("Failed to parse package: %w", err)
	}

	data, err := imports.Process("", buf.Bytes(), nil)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(fileName, data, 0600); err != nil {
		return xerrors.Errorf("unable to write file %s: %w", fileName, err)
	}

	fmt.Printf("Created %s\n", fileName)
	return nil
}
