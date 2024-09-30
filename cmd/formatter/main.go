package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bobmaertz/cuelang-lsp/pkg/fmtr"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "incorrect number of args\n")
		os.Exit(1)
	}

	filePath := args[0]
	fileInfo, err := os.Stat(filePath)
	if err != nil || fileInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "error locatingfile: %v\n", err)
		os.Exit(10)
	}

	fname := fileInfo.Name()
	b, rErr := os.ReadFile(filePath)
	if rErr != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v", rErr)
		os.Exit(11)
	}

	out, fErr := fmtr.Format(fname, b)
	if fErr != nil {
		fmt.Fprintf(os.Stderr, "error formatting: %v", fErr)
		os.Exit(12)
	}

	// TODO: Remove section in favor of Unix still stdout
	parts := strings.Split(fname, ".")
	wErr := os.WriteFile(fmt.Sprintf("%s_fmt.cue", parts[0]), out, os.ModePerm)
	if wErr != nil {
		fmt.Fprintf(os.Stderr, "error writing output to file: %v", wErr)
		os.Exit(13)
	}
}
