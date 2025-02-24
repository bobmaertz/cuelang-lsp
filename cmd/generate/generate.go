package main

import (
	"bytes"
	"fmt"
	"strings"
)

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

	for _, extends := range s.Extends {
		if strings.HasPrefix(extends.Name, "_") {
			continue
		}
		fmt.Fprintf(&buf, "\t%s\n", ToTitleCase(extends.Name))
	}

	for _, props := range s.Properties {
		if props.Documentation != "" {
			doc := strings.ReplaceAll(props.Documentation, "\n", " ")
			fmt.Fprintf(&buf, "\t// %s %s\n", ToTitleCase(props.Name), doc)
		}

		typeName := "interface{}"
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
