package fmtr

import (
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/parser"
)

func Format(name string, input []byte) ([]byte, error) {
	parseOpts := []parser.Option{
		parser.AllErrors,
		parser.ParseComments,
	}
	ast, err := parser.ParseFile("testFile", input, parseOpts...)
	if err != nil {
		return nil, err
	}

	options := []format.Option{
		//	format.Simplify(),
		format.UseSpaces(4),
	}
	final, err := format.Node(ast, options...)
	return final, err
}
