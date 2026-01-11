package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	lsp "github.com/bobmaertz/cuelang-lsp/pkg/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureStdout captures stdout during function execution
func captureStdout(f func()) string {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		panic(err) // Should never fail in tests, but handle it properly
	}
	return buf.String()
}

func TestHandleMessage_Initialize(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	request := lsp.InitializeRequest{
		Request: lsp.Request{
			RPC:    "2.0",
			ID:     1,
			Method: "initialize",
		},
		Params: lsp.InitializeParams{
			ClientInfo: lsp.ClientInfo{
				Name:    "test-client",
				Version: "1.0.0",
			},
		},
	}

	contents, err := json.Marshal(request)
	require.NoError(t, err)

	output := captureStdout(func() {
		HandleMessage(logger, nil, "initialize", contents)
	})

	// Verify response was written to stdout
	assert.Contains(t, output, "Content-Length:")
	assert.Contains(t, output, "cuelang-lsp")

	// Verify log output
	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "test-client")
	assert.Contains(t, logOutput, "1.0.0")

	// Verify response structure
	var response lsp.InitializeResponse
	// Extract JSON from Content-Length header format
	parts := strings.Split(output, "\r\n\r\n")
	if len(parts) >= 2 {
		err = json.Unmarshal([]byte(parts[1]), &response)
		require.NoError(t, err)
		assert.Equal(t, 1, response.ID)
		assert.Equal(t, "2.0", response.RPC)
		assert.True(t, response.Result.Capabilities.DocumentFormattingProvider)
		assert.Equal(t, "cuelang-lsp", response.Result.ServerInfo.Name)
	}
}

func TestHandleMessage_Initialize_InvalidJSON(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	invalidJSON := []byte(`{"invalid json}`)

	HandleMessage(logger, nil, "initialize", invalidJSON)

	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "unable to unmarshal initialize request")
}

func TestHandleMessage_DidOpen(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	notification := lsp.DidOpenNotification{
		Notification: lsp.Notification{
			RPC:    "2.0",
			Method: "textDocument/didOpen",
		},
		Params: lsp.DidOpenParams{
			TextDocument: lsp.TextDocumentItem{
				URI:        "file:///test.cue",
				LanguageID: "cue",
				Version:    1,
				Text:       "package test\n",
			},
		},
	}

	contents, err := json.Marshal(notification)
	require.NoError(t, err)

	HandleMessage(logger, nil, "textDocument/didOpen", contents)

	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "didOpen>")
	assert.Contains(t, logOutput, "file:///test.cue")
}

func TestHandleMessage_DidOpen_InvalidJSON(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	invalidJSON := []byte(`{"invalid json}`)

	HandleMessage(logger, nil, "textDocument/didOpen", invalidJSON)

	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "unable to unmarshal textDocument/didOpen notification")
}

func TestHandleMessage_DidChange(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	notification := lsp.DidChangeNotification{
		Notification: lsp.Notification{
			RPC:    "2.0",
			Method: "textDocument/didChange",
		},
		Params: lsp.DidChangeParams{
			TextDocument: lsp.VersionedTextDocumentIdentifier{
				TextDocumentIdentifier: lsp.TextDocumentIdentifier{
					URI: "file:///test.cue",
				},
				Version: 2,
			},
			ContentChanges: []lsp.TextDocumentContentChangeEvent{
				{
					Text: "package test\n\nfoo: \"bar\"\n",
				},
			},
		},
	}

	contents, err := json.Marshal(notification)
	require.NoError(t, err)

	HandleMessage(logger, nil, "textDocument/didChange", contents)

	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "didChange>")
	assert.Contains(t, logOutput, "file:///test.cue")
}

func TestHandleMessage_DidChange_InvalidJSON(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	invalidJSON := []byte(`{"invalid json}`)

	HandleMessage(logger, nil, "textDocument/didChange", invalidJSON)

	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "unable to unmarshal textDocument/didChange notification")
}

func TestHandleMessage_WillSave(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	contents := []byte(`{"textDocument":{"uri":"file:///test.cue"}}`)

	HandleMessage(logger, nil, "textDocument/willSave", contents)

	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "will Save:")
}

func TestHandleMessage_DidSave(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	contents := []byte(`{"textDocument":{"uri":"file:///test.cue"}}`)

	HandleMessage(logger, nil, "textDocument/didSave", contents)

	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "did Save:")
}

func TestHandleMessage_Formatting(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	// Create a temporary CUE file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.cue")
	unformattedCUE := `package test

foo:    "bar"
baz:   42
`
	err := os.WriteFile(testFile, []byte(unformattedCUE), 0644)
	require.NoError(t, err)

	request := lsp.TextFormatRequest{
		Request: lsp.Request{
			RPC:    "2.0",
			ID:     1,
			Method: "textDocument/formatting",
		},
		Params: lsp.DocumentFormattingParams{
			TextDocument: lsp.TextDocumentIdentifier{
				URI: "file://" + testFile,
			},
		},
	}

	contents, err := json.Marshal(request)
	require.NoError(t, err)

	output := captureStdout(func() {
		HandleMessage(logger, nil, "textDocument/formatting", contents)
	})

	// Verify response contains formatted text
	assert.Contains(t, output, "Content-Length:")
	assert.Contains(t, output, "newText")

	// Verify the response structure
	var response FormattingResponse
	parts := strings.Split(output, "\r\n\r\n")
	if len(parts) >= 2 {
		err = json.Unmarshal([]byte(parts[1]), &response)
		require.NoError(t, err)
		assert.Equal(t, 1, response.ID)
		assert.NotEmpty(t, response.Result)
		assert.Contains(t, response.Result[0].NewText, "package test")
	}
}

func TestHandleMessage_Formatting_InvalidJSON(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	invalidJSON := []byte(`{"invalid json}`)

	HandleMessage(logger, nil, "textDocument/formatting", invalidJSON)

	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "unable to unmarshal textDocument/formatting request")
}

func TestHandleMessage_Formatting_InvalidFile(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	request := lsp.TextFormatRequest{
		Request: lsp.Request{
			RPC:    "2.0",
			ID:     1,
			Method: "textDocument/formatting",
		},
		Params: lsp.DocumentFormattingParams{
			TextDocument: lsp.TextDocumentIdentifier{
				URI: "file:///nonexistent/file.cue",
			},
		},
	}

	contents, err := json.Marshal(request)
	require.NoError(t, err)

	HandleMessage(logger, nil, "textDocument/formatting", contents)

	logOutput := logBuf.String()
	// Should log an error when file doesn't exist
	assert.Contains(t, logOutput, "error")
}

func TestHandleMessage_Completion(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	request := lsp.TextCompletionRequest{
		Request: lsp.Request{
			RPC:    "2.0",
			ID:     1,
			Method: "textDocument/completion",
		},
		Params: lsp.CompletionParams{
			TextDocument: lsp.TextDocumentIdentifier{
				URI: "file:///test.cue",
			},
			Position: lsp.Position{
				Line:      5,
				Character: 10,
			},
		},
	}

	contents, err := json.Marshal(request)
	require.NoError(t, err)

	output := captureStdout(func() {
		HandleMessage(logger, nil, "textDocument/completion", contents)
	})

	// Verify response was written
	assert.Contains(t, output, "Content-Length:")

	// Verify response structure
	var response lsp.TextCompletionResponse
	parts := strings.Split(output, "\r\n\r\n")
	if len(parts) >= 2 {
		err = json.Unmarshal([]byte(parts[1]), &response)
		require.NoError(t, err)
		assert.Equal(t, 1, response.ID)
	}
}

func TestHandleMessage_Completion_InvalidJSON(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	invalidJSON := []byte(`{"invalid json}`)

	HandleMessage(logger, nil, "textDocument/completion", invalidJSON)

	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "unable to unmarshal textdocument/completion request")
}

func TestHandleMessage_UnknownMethod(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	contents := []byte(`{"some":"data"}`)

	HandleMessage(logger, nil, "unknown/method", contents)

	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "received method: unknown/method")
}
