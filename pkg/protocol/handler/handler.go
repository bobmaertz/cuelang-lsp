package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bobmaertz/cuelang-lsp/pkg/fmtr"
	lsp "github.com/bobmaertz/cuelang-lsp/pkg/protocol"
	"github.com/bobmaertz/cuelang-lsp/pkg/protocol/rpc"
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
		var request lsp.TextFormatRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			l.Printf("unable to unmarshal textDocument/formatting request: %v\n", err)
			return
		}

		// Todo: fix this hacky implementation
		f := strings.TrimPrefix(request.Params.TextDocument.Uri, "file://")
		c, _ := os.ReadFile(f)
		l.Printf(f)
		update, err := fmtr.Format("", c)
		if err != nil {
			l.Printf("error %v", err)
			return
		}

		response := FormattingResponse{
			Response: lsp.Response{
				Id: request.Id,
			},
		}
		t := TextEdit{
			NewText: string(update),
		}

		response.Result = append(response.Result, t)
		out := rpc.EncodeMessage(response)
		fmt.Print(out)
		l.Printf("didFormat> %v\n", request)
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

type FormattingResponse struct {
	lsp.Response
	Result []TextEdit `json:"result,omitempty"`
}

// TODO: move me somewhere else
type TextEdit struct {
	/**
	 * The range of the text document to be manipulated. To insert
	 * text into a document create a range where start === end.
	 */
	Range Range `json:"range"`

	/**
	 * The string to be inserted. For delete operations use an
	 * empty string.
	 */
	NewText string `json:"newText"`
}

type Range struct {
	/**
	 * The range's start position.
	 */
	Start Position `json:"start"`

	/**
	 * The range's end position.
	 */
	End Position `json:"end"`
}

type Position struct {
	/**
	 * Line position in a document (zero-based).
	 */
	Line int `json:"line"`

	/**
	 * Character offset on a line in a document (zero-based). The meaning of this
	 * offset is determined by the negotiated `PositionEncodingKind`.
	 *
	 * If the character value is greater than the line length it defaults back
	 * to the line length.
	 */
	Character int `json:"character"`
}
