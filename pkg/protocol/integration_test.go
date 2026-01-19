package protocol_test

import (
	"bytes"
	"encoding/json"
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

	require.NoError(t, w.Close())
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
	err := os.WriteFile(testFile, []byte(initialContent), 0o600)
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

	err = os.WriteFile(testFile, []byte(changedContent), 0o600)
	require.NoError(t, err)

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

	handler.HandleMessage(logger, nil, "textDocument/didChange", didChangeJSON)

	logOutput = logBuf.String()
	assert.Contains(t, logOutput, "didChange>")

	// Step 3: formatting request
	logBuf.Reset()
	formatRequest := protocol.TextFormatRequest{
		Request: protocol.Request{
			RPC:    "2.0",
			ID:     2,
			Method: "textDocument/formatting",
		},
		Params: protocol.TextFormatParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: "file://" + testFile,
			},
			Options: protocol.TextFormatOptions{

				TabSize:      4,
				InsertSpaces: true,
			},
		},
	}

	formatJSON, err := json.Marshal(formatRequest)
	require.NoError(t, err)

	// Capture stdout for formatting response
	old := os.Stdout
	r, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout = w

	handler.HandleMessage(logger, nil, "textDocument/formatting", formatJSON)

	require.NoError(t, w.Close())
	os.Stdout = old

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)
	output := buf.String()

	// Verify formatting response
	assert.Contains(t, output, "Content-Length:")
	assert.Contains(t, output, "result")
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
		example := example
		t.Run(example.name, func(t *testing.T) {
			// Read example file
			absPath, err := filepath.Abs(example.path)
			require.NoError(t, err)

			content, err := os.ReadFile(absPath) //nolint:gosec // test reads known example files
			if os.IsNotExist(err) {
				t.Skipf("Example file not found: %s", absPath)
			}
			require.NoError(t, err)

			// Create formatting request
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
					Options: protocol.TextFormatOptions{

						TabSize:      4,
						InsertSpaces: true,
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

			logger := log.New(io.Discard, "", 0)
			handler.HandleMessage(logger, nil, "textDocument/formatting", formatJSON)

			require.NoError(t, w.Close())
			os.Stdout = old

			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			require.NoError(t, err)
			output := buf.String()

			// Verify response contains edits
			assert.Contains(t, output, "Content-Length:")
			assert.Contains(t, output, "newText")

			// Verify the formatted output differs from original for poorly formatted files
			// For well-formatted files, output may be same
			parts := strings.Split(output, "\r\n\r\n")
			if len(parts) >= 2 {
				var response struct {
					Result []struct {
						NewText string `json:"newText"`
					} `json:"result"`
				}
				err = json.Unmarshal([]byte(parts[1]), &response)
				require.NoError(t, err)
				assert.NotEmpty(t, response.Result)
				// Verify formatted text is valid CUE (should at least contain package declaration)
				assert.Contains(t, response.Result[0].NewText, "package")
			}

			// For multi-file example, just ensure it formats without error
			_ = content
		})
	}
}

// TestIntegration_MultiClientScenario tests handling multiple clients
func TestIntegration_MultiClientScenario(t *testing.T) {
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)

	clients := []struct {
		name    string
		version string
	}{
		{"client1", "1.0"},
		{"client2", "2.0"},
	}

	for i, client := range clients {
		client := client
		t.Run(client.name, func(t *testing.T) {
			// Reset log buffer
			logBuf.Reset()

			initRequest := protocol.InitializeRequest{
				Request: protocol.Request{
					RPC:    "2.0",
					ID:     i + 1,
					Method: "initialize",
				},
				Params: protocol.InitializeParams{
					ClientInfo: protocol.ClientInfo{
						Name:    client.name,
						Version: client.version,
					},
				},
			}

			contents, err := json.Marshal(initRequest)
			require.NoError(t, err)

			// Capture stdout
			old := os.Stdout
			r, w, pipeErr := os.Pipe()
			require.NoError(t, pipeErr)
			os.Stdout = w

			handler.HandleMessage(logger, nil, "initialize", contents)

			require.NoError(t, w.Close())
			os.Stdout = old

			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			require.NoError(t, err)
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
