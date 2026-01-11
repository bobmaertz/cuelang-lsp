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
}

func main() {
	flag.Parse()
	flag.Usage = usage

	args := flag.Args()

	if len(args) < 1 {
		usage()
		return
	}

	// TODO: Read in from stdin ..
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

	switch args[0] {
	case "generate":
		generate(model)
	case "analyze":
		analyze(model)
	}
}

func analyze(model *MetaModel) {
	// structKeys := map[string]string{}
}

// generate generates the Go code from the provided MetaModel.
// It creates a file with the specified package name and writes the
// structures, enumerations, type aliases, and notifications to it.
func generate(model *MetaModel) {
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
		buf := GenerateStructure(s)
		if buf == nil {
			continue
		}
		if _, err := fileWriter.Write(buf.Bytes()); err != nil {
			fmt.Fprintf(os.Stderr, "error writing structure: %v\n", err)
			return
		}
		if err := fileWriter.Flush(); err != nil {
			fmt.Fprintf(os.Stderr, "error flushing: %v\n", err)
			return
		}
	}

	// For each enumeration..
	for _, e := range model.Enumerations {
		start := "type %s %s\n"
		fmt.Fprintf(fileWriter, start, e.Name, ConvertType(e.Type.Name))
	}

	// For each alias..
	for _, t := range model.TypeAliases {
		// TODO: Properly parse type.
		typ := "any"

		// TODO: Fix SelectionRange self reference - meaning add support for optional fields.
		doc := "// %s %s\n"
		start := "type %s %s\n"

		dv := strings.ReplaceAll(t.Documentation, "\n", " ")
		fmt.Fprintf(fileWriter, doc, t.Name, dv)
		fmt.Fprintf(fileWriter, start, t.Name, typ)
	}

	for _, n := range model.Notifications {
		buf := GenerateNotification(n)
		if buf == nil {
			continue
		}
		if _, err := fileWriter.Write(buf.Bytes()); err != nil {
			fmt.Fprintf(os.Stderr, "error writing notification: %v\n", err)
			return
		}
		if err := fileWriter.Flush(); err != nil {
			fmt.Fprintf(os.Stderr, "error flushing: %v\n", err)
			return
		}
	}
	// TODO Requests
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
		return "any"
	case "URI":
		// This is a base type but doesnt have a specific definition associated with it.
		// using string for now but consider unstable
		return "string"
	case "DocumentUri":
		// This is a base type but doesnt have a specific definition associated with it.
		// using string for now but consider unstable
		// TODO
		return "string"
	case "":
		return "any"
	}

	return s
}

func usage() {
	// TODO: Replace os.Args[0] with binary name.
	fmt.Fprintf(os.Stderr, "Usage of lsp-gen:\n")
	// TODO: Add description here.
	fmt.Fprintf(os.Stderr, "Command Usage:\n")
	fmt.Fprintf(os.Stderr, "  generate      Generate LSP Grammar\n")
	fmt.Fprintf(os.Stderr, "  analyze       Analyze metaModel\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Flag Usage:\n")
	fmt.Fprintf(os.Stderr, "  -f string\n")
	fmt.Fprintf(os.Stderr, "        location for metamodel file (default \"%s\")\n", defaultInputFile)
	fmt.Fprintf(os.Stderr, "  -o string\n")
	fmt.Fprintf(os.Stderr, "        location for output file (default \"%s\")\n", defaultOutputFile)
	fmt.Fprintf(os.Stderr, "  -p string\n")
	fmt.Fprintf(os.Stderr, "        package name for generated code (default \"%s\")\n", defaultPackageName)
	fmt.Fprintf(os.Stderr, "\n")

	flag.PrintDefaults()
}
