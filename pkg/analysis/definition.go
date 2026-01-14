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

	// Convert LSP position (0-based) to token.Pos
	// Note: CUE uses 1-based line numbers
	targetLine := pos.Line + 1
	targetChar := pos.Character

	var foundIdent *ast.Ident
	var foundPos token.Pos

	// Walk the AST to find the identifier at the target position
	ast.Walk(file, func(node ast.Node) bool {
		if node == nil {
			return false
		}

		nodePos := file.Pos(node.Pos(), token.NoRelPos)
		nodeEnd := file.Pos(node.End(), token.NoRelPos)

		// Check if this node contains our target position
		if nodePos.Line() == targetLine {
			// Check if the character position is within this node
			if ident, ok := node.(*ast.Ident); ok {
				// Check if cursor is on this identifier
				startCol := nodePos.Column() - 1 // Convert to 0-based
				endCol := startCol + len(ident.Name)

				if targetChar >= startCol && targetChar < endCol {
					foundIdent = ident
					foundPos = node.Pos()
				}
			}
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
					labelPos := file.Pos(label.Pos(), token.NoRelPos)
					defLocation = &Location{
						URI: uri,
						Range: Range{
							Start: Position{
								Line:      labelPos.Line() - 1, // Convert to 0-based
								Character: labelPos.Column() - 1,
							},
							End: Position{
								Line:      labelPos.Line() - 1,
								Character: labelPos.Column() - 1 + len(label.Name),
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
				aliasPos := file.Pos(alias.Ident.Pos(), token.NoRelPos)
				defLocation = &Location{
					URI: uri,
					Range: Range{
						Start: Position{
							Line:      aliasPos.Line() - 1,
							Character: aliasPos.Column() - 1,
						},
						End: Position{
							Line:      aliasPos.Line() - 1,
							Character: aliasPos.Column() - 1 + len(alias.Ident.Name),
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
