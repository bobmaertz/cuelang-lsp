package main

import (
	"io"
	"log"
	"os"

	"github.com/bobmaertz/cuelang-lsp/pkg/protocol/serve"
	"github.com/bobmaertz/cuelang-lsp/pkg/version"
)

const (
	name = "cuelang-lsp"
)

func main() {
	args := os.Args[1:]

	var filepath string
	if len(args) > 0 {
		filepath = args[0]
	}

	l := getLogger(filepath)
	l.Printf("Version: %s", version.Version())
	serve.Serve(l)
}

func getLogger(filename string) *log.Logger {
	var file io.Writer
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		// If the file cannot be opened, continue to serve
		// TODO: need to refactor into a better option for this
		file = io.Discard
	}

	l := log.New(file, "["+name+"] ", log.Ldate|log.Ltime|log.Lshortfile)

	return l
}
