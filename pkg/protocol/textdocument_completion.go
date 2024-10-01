package protocol

import "github.com/bobmaertz/cuelang-lsp/pkg/protocol/rpc"

type TextCompletionRequest struct {
	Request
	Params TextCompletionParams `json:"params,omitempty"`
}

type TextCompletionParams struct {
	Context CompletionContext `json:"context"`
}

type CompletionContext struct {
	TriggerKind      int    `json:"triggerKind"`
	TriggerCharacter string `json:"triggerCharacter"`
}

type TextCompletionResponse struct {
	Response
	//TODO: finish
}

func NewTextCompletionResponse(id int) TextCompletionResponse {

	return TextCompletionResponse{
		Response: Response{
			Rpc: rpc.Version,
			Id:  id,
		},
	}

}
