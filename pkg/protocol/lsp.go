package protocol

/* LSP defines the interation and types of the LSP protocol.*/
type Request struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`
}

type Response struct {
	RPC   string `json:"jsonrpc"`
	ID    int    `json:"id,omitempty"`
	Error *Error `json:"error,omitempty"`
}

type Error struct {
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Notification struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
}
