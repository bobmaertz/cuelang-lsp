package protocol

/* LSP defines the interation and types of the LSP protocol.*/
type Request struct {
	Rpc    string `json:"jsonrpc"`
	Id     int    `json:"id"`
	Method string `json:"method"`
}

type Response struct {
	Rpc   string `json:"jsonrpc"`
	Id    int    `json:"id,omitempty"`
	Error *Error `json:"error,omitempty"`
}

type Error struct {
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Notification struct {
	Rpc    string `json:"jsonrpc"`
	Method string `json:"method"`
}
