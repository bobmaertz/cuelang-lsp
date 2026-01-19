package analysis

import (
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/parser"
	"cuelang.org/go/cue/token"
)

// Position represents a position in a document (0-based line and character)
type Position struct {
	Line      int
	Character int
}

// Location represents a location in a source file
type Location struct {
	URI   string
	Range Range
}

// Range represents a range in a document
type Range struct {
	Start Position
	End   Position
}

// FindDefinition finds the definition of a symbol at the given position
func FindDefinition(uri string, content []byte, pos Position) ([]Location, error) {
	// Parse the CUE file
	file, err := parser.ParseFile(uri, content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Convert LSP position (0-based) to token position for searching
	targetLine := pos.Line + 1      // CUE uses 1-based line numbers
	targetChar := pos.Character + 1 // CUE uses 1-based columns

	var foundIdent *ast.Ident
	var foundPos token.Pos

	// Helper to convert token.Pos to line/column
	posToLineCol := func(p token.Pos) (int, int) {
		// Use the CUE token.Pos methods to get line and column
		if !p.IsValid() {
			return 0, 0
		}
		return p.Line(), p.Column()
	}

	// Walk the AST to find the identifier at the target position
	ast.Walk(file, func(node ast.Node) bool {
		if node == nil {
			return false
		}

		// Only look at identifiers
		ident, ok := node.(*ast.Ident)
		if !ok {
			return true
		}

		// Get position of this identifier
		line, col := posToLineCol(ident.Pos())

		// Check if cursor is on this identifier
		if line == targetLine && col <= targetChar && targetChar < col+len(ident.Name) {
			foundIdent = ident
			foundPos = ident.Pos()
			return false // Stop walking
		}

		return true
	}, nil)

	if foundIdent == nil {
		// No identifier found at this position
		return []Location{}, nil
	}

	// Now find the definition of this identifier
	identName := foundIdent.Name

	var defLocation *Location

	// Walk the AST again to find where this identifier is defined
	ast.Walk(file, func(node ast.Node) bool {
		if node == nil {
			return false
		}

		// Look for field definitions
		if field, ok := node.(*ast.Field); ok {
			if label, ok := field.Label.(*ast.Ident); ok {
				if label.Name == identName && label.Pos() != foundPos {
					// Found a definition
					line, col := posToLineCol(label.Pos())
					defLocation = &Location{
						URI: uri,
						Range: Range{
							Start: Position{
								Line:      line - 1, // Convert to 0-based
								Character: col - 1,  // Convert to 0-based
							},
							End: Position{
								Line:      line - 1,
								Character: col - 1 + len(label.Name),
							},
						},
					}
					return false // Stop walking
				}
			}
		}

		// Look for alias definitions
		if alias, ok := node.(*ast.Alias); ok {
			if alias.Ident.Name == identName && alias.Ident.Pos() != foundPos {
				line, col := posToLineCol(alias.Ident.Pos())
				defLocation = &Location{
					URI: uri,
					Range: Range{
						Start: Position{
							Line:      line - 1,
							Character: col - 1,
						},
						End: Position{
							Line:      line - 1,
							Character: col - 1 + len(alias.Ident.Name),
						},
					},
				}
				return false
			}
		}

		return true
	}, nil)

	if defLocation != nil {
		return []Location{*defLocation}, nil
	}

	// No definition found
	return []Location{}, nil
}
