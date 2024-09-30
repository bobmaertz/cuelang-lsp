package lsp

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
			Rpc: "2.0",
			Id:  id,
		},
	}

}
