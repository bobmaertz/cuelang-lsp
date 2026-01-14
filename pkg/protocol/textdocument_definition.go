package protocol

type DefinitionRequest struct {
	Request
	Params DefinitionParams `json:"params"`
}

type DefinitionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type DefinitionResponse struct {
	Response
	Result []Location `json:"result,omitempty"`
}

type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

func NewDefinitionResponse(id int, locations []Location) DefinitionResponse {
	return DefinitionResponse{
		Response: Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: locations,
	}
}
