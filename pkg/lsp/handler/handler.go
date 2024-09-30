package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bobmaertz/cuelang-lsp/pkg/lsp"
	"github.com/bobmaertz/cuelang-lsp/pkg/lsp/rpc"
)

func HandleMessage(l *log.Logger, _ interface{}, method string, contents []byte) {
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			l.Printf("unable to unmarshal initialize request: %v\n", err)
			return
		}
		l.Println(request)
		l.Printf("connected to client %s - %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
		response := lsp.NewInitializeResponse(1)
		out := rpc.EncodeMessage(response)
		fmt.Print(out)
		l.Println(out)
	case "textDocument/didOpen":
		var notification lsp.DidOpenNotification
		if err := json.Unmarshal(contents, &notification); err != nil {
			l.Printf("unable to unmarshal textDocument/didOpen notification: %v\n", err)
			return
		}
		l.Printf("didOpen> %v\n", notification.Params.TextDocument.Uri)
		// state.OpenDocument(notification.Params.TextDocument.Uri, notification.Params.TextDocument.Text)
	case "textDocument/didChange":
		var notification lsp.DidChangeNotification
		if err := json.Unmarshal(contents, &notification); err != nil {
			l.Printf("unable to unmarshal textDocument/didChange notification: %v\n", err)
			return
		}

		l.Printf("didChange> %v\n", notification.Params.TextDocument.Uri)
		//for _, change := range notification.Params.ContentChanges {
		//	state.UpdateDocument(notification.Params.TextDocument.Uri, change.Text)
		//}
	case "textDocument/willSave":
		l.Printf("will Save: %v", string(contents))
	case "textDocument/didSave":
		l.Printf("did Save: %v", string(contents))
	case "textDocument/formatting":
		l.Printf("formatting: %v", "<>")
	case "textDocument/completion":
		var request lsp.TextCompletionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			l.Printf("unable to unmarshal textdocument/completion request: %v\n", err)
			return
		}
		response := lsp.NewTextCompletionResponse(request.Id)
		out := rpc.EncodeMessage(response)
		fmt.Print(out)
	default:
		l.Printf("received method: %s, message: %s\n", method, contents)
		// l.Printf("state: %v", state)
	}
}
