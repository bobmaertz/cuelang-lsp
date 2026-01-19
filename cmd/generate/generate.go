package main

// This file contains the string-building helpers used by the generator.
// The generator intentionally returns buffers so callers can decide how/when
// to write and format output.

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ToTitleCase converts the first character of a string to uppercase
// and returns the modified string. If the string is empty,
// it returns the empty string
func ToTitleCase(s string) string {
	if s == "" {
		return s
	}

	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size == 1 {
		// Invalid UTF-8; fall back to original string.
		return s
	}

	return string(unicode.ToUpper(r)) + s[size:]
}

// GenerateStructure generates a Go structure definition from the provided
// Structure object.
func GenerateStructure(s Structure) *bytes.Buffer {
	var buf bytes.Buffer

	// Ignore hidden or unexported structures (starts with _)
	if strings.HasPrefix(s.Name, "_") {
		return nil
	}

	start := "type %s struct {\n"
	end := "}\n"
	fmt.Fprintf(&buf, start, s.Name)

	for _, mixin := range s.Mixins {
		if strings.HasPrefix(mixin.Name, "_") {
			continue
		}
		fmt.Fprintf(&buf, "\t%s\n", ToTitleCase(mixin.Name))
	}

	// Check if this type extends another
	for _, extends := range s.Extends {
		if strings.HasPrefix(extends.Name, "_") {
			continue
		}
		fmt.Fprintf(&buf, "\t%s\n", ToTitleCase(extends.Name))
	}

	// Iterate over the properties and generate the necessayr fields
	for _, props := range s.Properties {
		if props.Documentation != "" {
			doc := strings.ReplaceAll(props.Documentation, "\n", " ")
			fmt.Fprintf(&buf, "\t// %s %s\n", ToTitleCase(props.Name), doc)
		}

		// The following code is preparing the "Type" for the field. Arrays
		// and references are handled special. Default will be interface{}
		var typeName string
		if props.Type.Kind == "array" {
			refName := ConvertType(props.Type.Element.Name)
			typeName = fmt.Sprintf("[]%s", refName)
		} else {
			typeName = ConvertType(props.Type.Name)
		}

		if props.Optional != nil && *props.Optional && props.Type.Kind != "array" {
			typeName = fmt.Sprintf("*%s", typeName)
		}

		// TODO: Need json tags for unmarshalling to include omitempty
		fmt.Fprintf(&buf, "\t%s %s\n", ToTitleCase(props.Name), typeName)
	}
	fmt.Fprint(&buf, end)
	return &buf
}

func GenerateNotification(s Notification) *bytes.Buffer {
	var buf bytes.Buffer

	return &buf
}
