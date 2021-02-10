package testdoc

import (
	"bufio"
	"bytes"
	"strings"
)

// Doc represents an documentation block.
type Doc struct {
	Description string
	Fields      map[string]string
}

// ParseDoc parses an documentation block and returns an instance to Doc.
// It will parse Sections separated by a colon and populate them in the Doc.Fields field.
// Example:
// Take following input:
// Lorem ipsum dolor sit amet, consectetur adipiscing elit.
//
// Field1:
//     Hello World
//
// Field2:
//    Good Bye World
// Will return
// Doc{
//    Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
//    Fields: map[string]string{
//        "Field1": "Hello World",
//        "Field2": "Good Bye World",
//    },
// }.
func ParseDoc(s string) (*Doc, error) {
	var d Doc

	setField := func(fieldName, value string) {
		if fieldName == "" {
			d.Description = strings.TrimSpace(value)
			return
		}
		if d.Fields == nil {
			d.Fields = make(map[string]string)
		}
		d.Fields[fieldName] = strings.TrimSpace(value)
	}

	scanner := bufio.NewScanner(bytes.NewReader([]byte(s)))

	var previousText string
	var sb strings.Builder
	var currentField string

	for scanner.Scan() {
		txt := scanner.Text()

		if strings.HasSuffix(txt, ":") && previousText == "" {
			setField(currentField, sb.String())
			sb.Reset()
			currentField = txt[:len(txt)-1]
			continue
		}
		sb.WriteString(strings.TrimSpace(txt))
		sb.WriteRune('\n')
		previousText = txt
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if sb.Len() > 0 {
		setField(currentField, sb.String())
	}

	return &d, nil
}
