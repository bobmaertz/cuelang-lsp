# cuelang-lsp Development Guide

## Project Overview

This is a Language Server Protocol (LSP) implementation for the CUE language, providing IDE integration for CUE files. The project is built in Go and follows LSP specification for text document synchronization and formatting.

**Status:** Work in progress - formatting works, other features (completion, hover, definitions) are stubbed but not implemented.

**Core Capability:** Document formatting using CUE's native formatter

## Architecture

```
cmd/lsp/           - Main LSP server entry point
pkg/protocol/      - LSP protocol types and handlers
  ├── handler/     - Message routing and request handling
  ├── rpc/         - JSON-RPC message encoding/decoding
  └── serve/       - Server initialization and message loop
pkg/fmtr/          - CUE formatting wrapper
pkg/version/       - Build version information
examples/          - Example CUE files for testing
```

## Key Protocol Types

**CRITICAL:** The protocol types use specific naming conventions that differ from standard LSP:
- ✅ `DidOpenTextDocumentParams` (NOT `DidOpenParams`)
- ✅ `DidChangeTextDocumentParams` (NOT `DidChangeParams`)
- ✅ `TextFormatParams` (NOT `DocumentFormattingParams`)
- ✅ `TextCompletionParams` (includes `TextDocument`, `Position`, and `Context`)

**Always check existing types before adding tests or code that references them!**

## Development Workflow

### Building
```bash
make build              # Builds to ./bin/lsp
make test               # Runs all tests
make test-coverage      # Runs tests with coverage report
make test-coverage-html # Generate HTML coverage report
make lint               # Runs golangci-lint
```

### Testing
- **Unit tests:** Test individual handlers and formatters in isolation
- **Integration tests:** Test complete LSP message flows (init → didOpen → format)
- **Example files:** Use files in `examples/` for real-world testing

**Test Conventions:**
- Use `require.NoError(t, err)` for errors that should stop the test
- Use `assert.*` for checks that should continue
- Always handle errors from `os.Pipe()` and `io.Copy()` in tests
- Accept both tabs and spaces for indentation tests (CUE formatter behavior varies)

### Error Handling

**NEVER ignore errors silently!** Common patterns:

```go
// In production code
c, err := os.ReadFile(f)
if err != nil {
    l.Printf("error reading file %s: %v", f, err)
    return
}

// In tests
r, w, err := os.Pipe()
require.NoError(t, err)  // Use require for setup

_, err = io.Copy(&buf, r)
require.NoError(t, err)  // Use require for operations
```

## Common Pitfalls

### 1. Wrong Protocol Type Names
❌ **Don't:**
```go
Params: protocol.DidOpenParams{  // This type doesn't exist!
```

✅ **Do:**
```go
Params: protocol.DidOpenTextDocumentParams{
```

### 2. Ignored Errors
❌ **Don't:**
```go
c, _ := os.ReadFile(f)  // Silently ignores errors!
```

✅ **Do:**
```go
c, err := os.ReadFile(f)
if err != nil {
    // Handle the error
}
```

### 3. Hardcoded Indentation Expectations
❌ **Don't:**
```go
assert.True(t, strings.HasPrefix(line, "    "))  // Only spaces!
```

✅ **Do:**
```go
assert.True(t, strings.HasPrefix(line, "\t") || strings.HasPrefix(line, "    "))
```

### 4. Missing imports in aliases
The project uses `lsp` as an alias for the protocol package:
```go
import lsp "github.com/bobmaertz/cuelang-lsp/pkg/protocol"
```

## Adding New LSP Features

To add a new LSP capability (hover, definition, etc.):

1. **Define protocol types** in `pkg/protocol/textdocument_*.go`
   - Request/Response structs
   - Parameter types
   - Follow existing naming conventions

2. **Add handler** in `pkg/protocol/handler/handler.go`
   - Add case to switch statement
   - Unmarshal request
   - Implement logic
   - Encode and write response

3. **Update capabilities** in `pkg/protocol/initialize.go`
   - Set capability flag to `true`

4. **Add tests**
   - Unit tests in `handler/handler_test.go`
   - Integration tests in `integration_test.go`

5. **Document** in README.md

## Testing Strategy

### Unit Tests
Test each handler independently:
- Valid requests
- Invalid JSON
- Error cases (file not found, parse errors)
- Edge cases (empty files, etc.)

### Integration Tests
Test complete workflows:
- Client initialization
- Document lifecycle (open → change → format)
- Multi-client scenarios
- Error propagation

### Example-Based Tests
Use files in `examples/` to test with real CUE:
- Simple configs
- Complex schemas
- Multi-file projects

## Code Quality Standards

- **Linting:** All code must pass `golangci-lint run`
- **Coverage:** Aim for 70%+ overall, 90%+ for critical paths
- **Error handling:** No ignored errors without justification
- **Tests:** All new features need tests
- **Documentation:** Update README for user-facing changes

## Git Workflow

- **Branch naming:** `claude/*` for AI-generated changes
- **Commits:** Clear, descriptive messages explaining "why"
- **PRs:** Wait for CI to pass before merging
- **Never push to main directly**

## CI/CD

GitHub Actions runs on every push/PR:
- Build (`go build -v ./...`)
- Test with coverage and race detection
- Upload coverage to Codecov
- Generate HTML coverage reports

**Important:** Tests run with Go 1.25+ (required version)

## Current Limitations & TODOs

**Working:**
- ✅ LSP initialization
- ✅ Document open/change notifications
- ✅ Document formatting

**Not Working (Stubbed):**
- ❌ Code completion (empty response)
- ❌ Hover provider (declared but not implemented)
- ❌ Go to definition (declared but not implemented)
- ❌ Diagnostics (no error checking)

**Known Issues:**
- Formatting reads from filesystem (hacky - should use document state)
- No document state management (each request is stateless)
- No workspace awareness (can't handle multi-file projects properly)

## Debugging

**Enable debug logging:**
```bash
lsp ~/debug.log
```

**Check what's being logged:**
```bash
tail -f ~/debug.log
```

**Common debug scenarios:**
- Message unmarshaling errors → check JSON structure
- No response → check if handler wrote to stdout
- Formatting errors → check CUE syntax validity

## Resources

- [LSP Specification](https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/)
- [CUE Documentation](https://cuelang.org/)
- [Go LSP Example (gopls)](https://github.com/golang/tools/tree/master/gopls)

## Quick Reference

**Add a new test:**
```go
func TestHandleMessage_NewFeature(t *testing.T) {
    var logBuf bytes.Buffer
    logger := log.New(&logBuf, "", 0)

    request := lsp.NewFeatureRequest{
        Request: lsp.Request{
            RPC: "2.0",
            ID: 1,
            Method: "textDocument/newFeature",
        },
        Params: lsp.NewFeatureParams{
            // params here
        },
    }

    contents, err := json.Marshal(request)
    require.NoError(t, err)

    handler.HandleMessage(logger, nil, "textDocument/newFeature", contents)

    // Assert expected behavior
}
```

**Run specific test:**
```bash
go test -v ./pkg/protocol/handler -run TestHandleMessage_NewFeature
```

**Common Make targets:**
```bash
make build              # Build binary
make test               # Run all tests
make test-coverage      # Tests with coverage
make test-coverage-html # Generate HTML report
make lint               # Run linter
make clean              # Clean build artifacts
```

## Protocol Type Reference

Quick reference for commonly used types and their correct names:

### Requests
- `InitializeRequest` with `InitializeParams`
- `TextFormatRequest` with `TextFormatParams`
- `TextCompletionRequest` with `TextCompletionParams`

### Notifications
- `DidOpenNotification` with `DidOpenTextDocumentParams`
- `DidChangeNotification` with `DidChangeTextDocumentParams`

### Responses
- `InitializeResponse` with `InitializeResult`
- `TextCompletionResponse`
- Custom `FormattingResponse` with `[]TextEdit`

### Common Structures
- `TextDocumentIdentifier` (has `URI`)
- `VersionedTextDocumentIdentifier` (has `URI` and `Version`)
- `TextDocumentItem` (has `URI`, `LanguageID`, `Version`, `Text`)
- `Position` (has `Line`, `Character`)
- `Range` (has `Start`, `End` positions)

## Test File Patterns

### Handler Tests (`handler/handler_test.go`)
Tests individual message handlers in isolation. Uses `captureStdout()` helper to capture LSP responses.

### Integration Tests (`integration_test.go`)
Tests complete message flows from request encoding through handler to response decoding.

### Formatter Tests (`fmtr/fmt_test.go`)
Tests CUE formatting logic. **Important:** Tests must accept both tabs and spaces for indentation.
