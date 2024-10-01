package protocol

type TextFormatRequest struct {
	Request
	Params TextFormatParams `json:"params"`
}

type TextFormatParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument,omitempty"`
	Options      TextFormatOptions      `json:"options,omitempty"`
}

type TextFormatOptions struct {
	// Size of a tab in spaces.
	TabSize int `json:"tabSize"`
	// Prefer spaces over tabs.
	InsertSpaces bool `json:"insertSpaces"`
	// Trim trailing whitespace on a line.
	TrimTrailingWhitespace *bool `json:"trimTrailingWhitespace,omitempty"`
	// Insert a newline character at the end of the file if one does not exist.
	InsertFinalNewline *bool `json:"insertFinalNewline,omitempty"`
	// Trim all newlines after the final newline at the end of the file.
	TrimFinalNewlines *bool `json:"trimFinalNewlines,omitempty"`
}
