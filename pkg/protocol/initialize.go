package protocol

import (
	"github.com/bobmaertz/cuelang-lsp/pkg/protocol/rpc"
	"github.com/bobmaertz/cuelang-lsp/pkg/version"
)

type InitializeRequest struct {
	Request
	Params InitializeParams `json:"params,omitempty"`
}

type InitializeParams struct {
	ClientInfo ClientInfo `json:"clientInfo,omitempty"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo,omitempty"`
}

type ServerCapabilities struct {
	TextDocumentSync           int            `json:"textDocumentSync"`
	HoverProvider              bool           `json:"hoverProvider"`
	DefinitionProvider         bool           `json:"definitionProvider"`
	CodeActionProvider         bool           `json:"codeActionProvider"`
	CompletionProvider         map[string]any `json:"completionProvider"`
	DocumentFormattingProvider bool           `json:"documentFormattingProvider"`
}

type TextDocumentSyncOptions struct {
	OpenClose bool `json:"openClose,omitempty"`
	Change    int  `json:"change,omitempty"`
}

type CompletionProvider struct {
	WorkDoneProgress bool `json:"workDoneProgress,omitempty"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			Rpc: rpc.Version,
			Id:  id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync:           1,
				DocumentFormattingProvider: true,
			},
			ServerInfo: ServerInfo{
				Name:    "cuelang-lsp",
				Version: version.Version(),
			},
		},
	}
}
