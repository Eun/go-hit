package yaegi_template

import (
	"fmt"
	"strings"
)

type importSymbols []Import

func (is importSymbols) Contains(symbol Import) bool {
	for _, s := range is {
		if s.Equals(symbol) {
			return true
		}
	}
	return false
}

func (is importSymbols) ImportBlock() string {
	switch len(is) {
	case 0:
		return ""
	case 1:
		return "import " + is[0].importLine()
	default:
		var sb strings.Builder
		sb.WriteString("import (\n")
		for _, symbol := range is {
			sb.WriteString(symbol.importLine())
			sb.WriteRune('\n')
		}
		sb.WriteString(")")
		return sb.String()
	}
}

// Import represents an import that should be evaluated.
type Import struct {
	Name string
	Path string
}

// Equals returns true if the specified import is equal to this import.
func (v Import) Equals(i Import) bool {
	return v.Name == i.Name && strings.EqualFold(v.Path, i.Path)
}

func (v Import) importLine() string {
	if v.Name != "" {
		return fmt.Sprintf("%s %q", v.Name, v.Path)
	}
	return fmt.Sprintf("%q", v.Path)
}
