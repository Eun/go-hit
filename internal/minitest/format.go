package minitest

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/gookit/color"
)

// Format formats an message by making sure the passed in data is nicely indented.
func Format(message, data string, colors ...color.Color) string {
	lines := strings.FieldsFunc(strings.TrimRightFunc(data, unicode.IsSpace), func(r rune) bool {
		return r == '\n'
	})
	var sb strings.Builder
	runes := []rune(message)
	for i := 0; i < len(runes); i++ {
		if !unicode.IsSpace(runes[i]) {
			runes[i] = ' '
		}
		sb.WriteRune(runes[i])
	}

	prefix := sb.String()
	sb.Reset()

	sb.WriteString(message)

	cl := color.New(colors...)
	for i := 0; i < len(lines); i++ {
		if i > 0 {
			sb.WriteString(prefix)
		}
		sb.WriteString("\t")
		sb.WriteString(cl.Sprint(strings.TrimRightFunc(lines[i], unicode.IsSpace)))
		sb.WriteString("\n")
	}
	return sb.String()
}

func trimLeftSpaces(str string) string {
	lines := strings.FieldsFunc(str, func(r rune) bool {
		return r == '\n'
	})
	var sb strings.Builder
	for i := 0; i < len(lines); i++ {
		if _, err := fmt.Fprintln(&sb, strings.TrimLeftFunc(lines[i], unicode.IsSpace)); err != nil {
			panic(err)
		}
	}
	return sb.String()
}
