package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	defaultInputFile   = "./testdata/metaModel.json"
	defaultPackageName = "main"
	defaultOutputFile  = "test.go"
)

var (
	packageName    string
	inputFileName  string
	outputFileName string
)

func init() {

	flag.StringVar(&inputFileName, "f", defaultInputFile, "location for metamodel file")
	flag.StringVar(&packageName, "p", defaultPackageName, "package name for generated code")
	flag.StringVar(&outputFileName, "o", defaultOutputFile, "location for output file")

	flag.Parse()
	flag.Usage = usage
}

func main() {
	// args := flag.Args()

	// if len(args) < 1 {
	// 	usage()
	// 	return
	// }

	// TOOD: Read in from stdin ..
	b, err := os.ReadFile(inputFileName)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(10)
	}

	// Unmarshall into go represetation
	model := &MetaModel{}
	err = json.Unmarshal(b, model)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(15)
	}

	// Create output file
	file, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(16)
	}
	defer file.Close()

	fileWriter := bufio.NewWriter(file)
	defer func() {
		// Dont forget to flush to file or might lose the info in the buffer
		err = fileWriter.Flush()
		if err != nil {
			os.Stderr.Write([]byte(err.Error()))
			os.Exit(16)
		}
	}()

	h := fmt.Sprintf("package %s\n\n", packageName)
	fmt.Fprint(fileWriter, h)

	// For each structure..
	for _, s := range model.Structures {

		//Ignore hidden or unexported structures (starts with _)
		if strings.HasPrefix(s.Name, "_") {
			continue
		}
		// // TODO: This is a good spot to use a template instead..
		structure := "type %s struct {\n"
		end := "}\n"
		fmt.Fprintf(fileWriter, structure, s.Name)

		// TODO: each field
		for _, p := range s.Properties {
			if p.Type.Name != "" {
				if p.Documentation != "" {
					doc := strings.ReplaceAll(p.Documentation, "\n", " ")
					fmt.Fprintf(fileWriter, "\t // %s %s\n", ToTitleCase(p.Name), doc)
				}
				n := ConvertType(p.Type.Name)
				if p.Optional != nil && *p.Optional {
					n = fmt.Sprintf("*%s", n)
				}
				//TODO: Need json tags for unmarshalling to include omitempty
				fmt.Fprintf(fileWriter, "\t %s %s\n", ToTitleCase(p.Name), n)
			}

			// if p.Type.Kind == "array" {
			// 	fmt.Fprintf(fileWriter, "\t %s []%s\n", ToTitleCase(p.Name), p.Type.Element.Name)
			// }

		}

		fmt.Fprint(fileWriter, end)
	}
	// For each enumeration..
	for _, e := range model.Enumerations {
		start := "type %s %s\n"
		fmt.Fprintf(fileWriter, start, e.Name, ConvertType(e.Type.Name))
	}

	// For each alias..
	for _, t := range model.TypeAliases {
		// TODO: Propertly parse type.
		typ := "interface{}"

		//TODO: Fix SelectionRange self reference - meaning add support for optional fields.
		doc := "// %s %s\n"
		start := "type %s %s\n"

		dv := strings.ReplaceAll(t.Documentation, "\n", " ")
		fmt.Fprintf(fileWriter, doc, t.Name, dv)
		fmt.Fprintf(fileWriter, start, t.Name, typ)
	}
}

func ConvertType(s string) string {
	switch s {
	case "boolean":
		return "bool"
	case "uinteger":
		return "uint"
	case "integer":
		return "int"
	case "decimal":
		return "float64"
	case "LSPAny":
		return "interface{}"
	case "URI":
		// This is a base type but doesnt have a specific definition associated with it.
		// using string for now but consider unstable
		return "string"
	case "DocumentUri":
		// This is a base type but doesnt have a specific definition associated with it.
		// using string for now but consider unstable
		// TODO
		return "string"
	}

	return s
}

func ToTitleCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func usage() {
	// TODO: Replace os.Args[0] with binary name.
	fmt.Fprintf(os.Stderr, "Usage of lsp-gen:\n")
	// TODO: Add description here.
	fmt.Fprintf(os.Stderr, "Command Usage:\n")
	fmt.Fprintf(os.Stderr, "  generate      Generate LSP Grammar\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Flag Usage:\n")

	flag.PrintDefaults()
}
