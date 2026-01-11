package protocol_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bobmaertz/cuelang-lsp/pkg/protocol"
	"github.com/bobmaertz/cuelang-lsp/pkg/protocol/handler"
	"github.com/bobmaertz/cuelang-lsp/pkg/protocol/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegration_InitializeWorkflow tests the complete initialization workflow
func TestIntegration_InitializeWorkflow(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	// Step 1: Send initialize request
	initRequest := protocol.InitializeRequest{
		Request: protocol.Request{
			RPC:    "2.0",
			ID:     1,
			Method: "initialize",
		},
		Params: protocol.InitializeParams{
			ClientInfo: protocol.ClientInfo{
				Name:    "test-editor",
				Version: "1.0.0",
			},
		},
	}

	initMessage := rpc.EncodeMessage(initRequest)

	// Parse the message back
	method, contents, err := rpc.DecodeMessage([]byte(initMessage))
	require.NoError(t, err)
	assert.Equal(t, "initialize", method)

	// Capture stdout for response
	old := os.Stdout
	r, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout = w

	handler.HandleMessage(logger, nil, method, contents)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)
	output := buf.String()

	// Verify response
	assert.Contains(t, output, "Content-Length:")
	assert.Contains(t, output, "cuelang-lsp")

	// Parse response
	parts := strings.Split(output, "\r\n\r\n")
	if len(parts) >= 2 {
		var response protocol.InitializeResponse
		err = json.Unmarshal([]byte(parts[1]), &response)
		require.NoError(t, err)
		assert.Equal(t, "2.0", response.RPC)
		assert.Equal(t, 1, response.ID)
		assert.True(t, response.Result.Capabilities.DocumentFormattingProvider)
	}
}

// TestIntegration_DocumentLifecycle tests document open, change, and formatting
func TestIntegration_DocumentLifecycle(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	// Create a temporary CUE file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.cue")
	initialContent := `package test

foo:    "bar"
`
	err := os.WriteFile(testFile, []byte(initialContent), 0644)
	require.NoError(t, err)

	// Step 1: didOpen notification
	didOpenNotif := protocol.DidOpenNotification{
		Notification: protocol.Notification{
			RPC:    "2.0",
			Method: "textDocument/didOpen",
		},
		Params: protocol.DidOpenTextDocumentParams{
			TextDocument: protocol.TextDocumentItem{
				URI:        "file://" + testFile,
				LanguageID: "cue",
				Version:    1,
				Text:       initialContent,
			},
		},
	}

	didOpenJSON, err := json.Marshal(didOpenNotif)
	require.NoError(t, err)

	handler.HandleMessage(logger, nil, "textDocument/didOpen", didOpenJSON)

	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "didOpen>")

	// Step 2: didChange notification
	logBuf.Reset()
	changedContent := `package test

foo:    "bar"
baz:   42
`
	didChangeNotif := protocol.DidChangeNotification{
		Notification: protocol.Notification{
			RPC:    "2.0",
			Method: "textDocument/didChange",
		},
		Params: protocol.DidChangeTextDocumentParams{
			TextDocument: protocol.VersionedTextDocumentIdentifier{
				TextDocumentIdentifier: protocol.TextDocumentIdentifier{
					URI: "file://" + testFile,
				},
				Version: 2,
			},
			ContentChanges: []protocol.TextDocumentContentChangeEvent{
				{Text: changedContent},
			},
		},
	}

	didChangeJSON, err := json.Marshal(didChangeNotif)
	require.NoError(t, err)

	// Update file for formatting test
	err = os.WriteFile(testFile, []byte(changedContent), 0644)
	require.NoError(t, err)

	handler.HandleMessage(logger, nil, "textDocument/didChange", didChangeJSON)

	logOutput = logBuf.String()
	assert.Contains(t, logOutput, "didChange>")

	// Step 3: formatting request
	formatRequest := protocol.TextFormatRequest{
		Request: protocol.Request{
			RPC:    "2.0",
			ID:     1,
			Method: "textDocument/formatting",
		},
		Params: protocol.TextFormatParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: "file://" + testFile,
			},
		},
	}

	formatJSON, err := json.Marshal(formatRequest)
	require.NoError(t, err)

	// Capture stdout
	old := os.Stdout
	r, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout = w

	handler.HandleMessage(logger, nil, "textDocument/formatting", formatJSON)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)
	output := buf.String()

	// Verify formatting response
	assert.Contains(t, output, "Content-Length:")
	assert.Contains(t, output, "newText")
}

// TestIntegration_MessageEncoding tests encoding and decoding of various message types
func TestIntegration_MessageEncoding(t *testing.T) {
	tests := []struct {
		name    string
		message interface{}
		method  string
	}{
		{
			name: "Initialize Request",
			message: protocol.InitializeRequest{
				Request: protocol.Request{
					RPC:    "2.0",
					ID:     1,
					Method: "initialize",
				},
				Params: protocol.InitializeParams{
					ClientInfo: protocol.ClientInfo{
						Name:    "test",
						Version: "1.0",
					},
				},
			},
			method: "initialize",
		},
		{
			name: "DidOpen Notification",
			message: protocol.DidOpenNotification{
				Notification: protocol.Notification{
					RPC:    "2.0",
					Method: "textDocument/didOpen",
				},
				Params: protocol.DidOpenTextDocumentParams{
					TextDocument: protocol.TextDocumentItem{
						URI:        "file:///test.cue",
						LanguageID: "cue",
						Version:    1,
						Text:       "package test",
					},
				},
			},
			method: "textDocument/didOpen",
		},
		{
			name: "Completion Request",
			message: protocol.TextCompletionRequest{
				Request: protocol.Request{
					RPC:    "2.0",
					ID:     2,
					Method: "textDocument/completion",
				},
				Params: protocol.TextCompletionParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "file:///test.cue",
					},
					Position: protocol.Position{
						Line:      5,
						Character: 10,
					},
				},
			},
			method: "textDocument/completion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode message
			encoded := rpc.EncodeMessage(tt.message)

			// Verify Content-Length header
			assert.Contains(t, encoded, "Content-Length:")

			// Decode message
			method, contents, err := rpc.DecodeMessage([]byte(encoded))
			require.NoError(t, err)
			assert.Equal(t, tt.method, method)

			// Verify contents can be unmarshaled back
			var result map[string]interface{}
			err = json.Unmarshal(contents, &result)
			require.NoError(t, err)

			// Verify method is present
			if methodField, ok := result["method"]; ok {
				assert.Equal(t, tt.method, methodField)
			}
		})
	}
}

// TestIntegration_ErrorHandling tests error scenarios
func TestIntegration_ErrorHandling(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	tests := []struct {
		name          string
		method        string
		contents      []byte
		expectedError string
	}{
		{
			name:          "Invalid JSON for initialize",
			method:        "initialize",
			contents:      []byte(`{invalid json}`),
			expectedError: "unable to unmarshal initialize request",
		},
		{
			name:          "Invalid JSON for didOpen",
			method:        "textDocument/didOpen",
			contents:      []byte(`{invalid json}`),
			expectedError: "unable to unmarshal textDocument/didOpen notification",
		},
		{
			name:          "Invalid JSON for formatting",
			method:        "textDocument/formatting",
			contents:      []byte(`{invalid json}`),
			expectedError: "unable to unmarshal textDocument/formatting request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logBuf.Reset()

			handler.HandleMessage(logger, nil, tt.method, tt.contents)

			logOutput := logBuf.String()
			assert.Contains(t, logOutput, tt.expectedError)
		})
	}
}

// TestIntegration_MultipleClients simulates multiple clients connecting
func TestIntegration_MultipleClients(t *testing.T) {
	clients := []struct {
		name    string
		version string
	}{
		{"vscode", "1.70.0"},
		{"neovim", "0.9.0"},
		{"emacs", "29.1"},
	}

	for _, client := range clients {
		t.Run(client.name, func(t *testing.T) {
			var logBuf bytes.Buffer
			logger := log.New(&logBuf, "", 0)

			initRequest := protocol.InitializeRequest{
				Request: protocol.Request{
					RPC:    "2.0",
					ID:     1,
					Method: "initialize",
				},
				Params: protocol.InitializeParams{
					ClientInfo: protocol.ClientInfo{
						Name:    client.name,
						Version: client.version,
					},
				},
			}

			initJSON, err := json.Marshal(initRequest)
			require.NoError(t, err)

			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			handler.HandleMessage(logger, nil, "initialize", initJSON)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Verify each client gets proper response
			assert.Contains(t, output, "Content-Length:")

			// Verify log contains client info
			logOutput := logBuf.String()
			assert.Contains(t, logOutput, client.name)
			assert.Contains(t, logOutput, client.version)
		})
	}
}

// TestIntegration_FormattingWithExamples tests formatting with the example files
func TestIntegration_FormattingWithExamples(t *testing.T) {
	// Get the examples directory
	examplesDir := filepath.Join("..", "..", "examples")

	examples := []struct {
		name string
		path string
	}{
		{"Simple Config", filepath.Join(examplesDir, "simple", "config.cue")},
		{"Kubernetes Schema", filepath.Join(examplesDir, "schema", "kubernetes.cue")},
		{"Multi-file Types", filepath.Join(examplesDir, "multi-file", "types.cue")},
	}

	for _, example := range examples {
		t.Run(example.name, func(t *testing.T) {
			// Check if example file exists
			absPath, err := filepath.Abs(example.path)
			if err != nil {
				t.Skip("Example file path resolution failed")
				return
			}

			if _, err := os.Stat(absPath); os.IsNotExist(err) {
				t.Skip(fmt.Sprintf("Example file not found: %s", absPath))
				return
			}

			var logBuf bytes.Buffer
			logger := log.New(&logBuf, "", 0)

			formatRequest := protocol.TextFormatRequest{
				Request: protocol.Request{
					RPC:    "2.0",
					ID:     1,
					Method: "textDocument/formatting",
				},
				Params: protocol.TextFormatParams{
					TextDocument: protocol.TextDocumentIdentifier{
						URI: "file://" + absPath,
					},
				},
			}

			formatJSON, err := json.Marshal(formatRequest)
			require.NoError(t, err)

			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			handler.HandleMessage(logger, nil, "textDocument/formatting", formatJSON)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Verify formatting response
			assert.Contains(t, output, "Content-Length:")
			assert.Contains(t, output, "newText")
		})
	}
}
