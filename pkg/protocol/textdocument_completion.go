package protocol

import "github.com/bobmaertz/cuelang-lsp/pkg/protocol/rpc"

type TextCompletionRequest struct {
	Request
	Params TextCompletionParams `json:"params,omitempty"`
}

type TextCompletionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
	Context      CompletionContext      `json:"context,omitempty"`
}

type CompletionContext struct {
	TriggerKind      int    `json:"triggerKind"`
	TriggerCharacter string `json:"triggerCharacter"`
}

type TextCompletionResponse struct {
	Response
	// TODO: finish
}

func NewTextCompletionResponse(id int) TextCompletionResponse {
	return TextCompletionResponse{
		Response: Response{
			RPC: rpc.Version,
			ID:  id,
		},
	}
}
