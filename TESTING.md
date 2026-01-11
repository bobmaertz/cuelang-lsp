# Testing Guide

This document describes the testing infrastructure and best practices for the cuelang-lsp project.

## Overview

The project uses Go's built-in testing framework along with the [testify](https://github.com/stretchr/testify) library for assertions. Our testing strategy includes:

- **Unit Tests** - Test individual components in isolation
- **Integration Tests** - Test complete LSP message flows
- **Coverage Reporting** - Track test coverage metrics

## Running Tests

### Basic Test Execution

Run all tests:
```bash
make test
```

Or using go directly:
```bash
go test -v ./...
```

### Coverage Reports

Generate coverage report:
```bash
make test-coverage
```

Generate and view HTML coverage report:
```bash
make test-coverage-html
```

View coverage by function:
```bash
make test-coverage-func
```

Generate and open HTML coverage in browser:
```bash
make test-coverage-view
```

## Test Structure

### Unit Tests

Unit tests are located alongside the code they test with the `_test.go` suffix:

- `pkg/fmtr/fmt_test.go` - Formatter tests
- `pkg/protocol/rpc/rpc_test.go` - RPC encoding/decoding tests
- `pkg/protocol/handler/handler_test.go` - Handler tests
- `cmd/generate/generate_test.go` - Code generation tests

### Integration Tests

Integration tests verify complete workflows:

- `pkg/protocol/integration_test.go` - End-to-end LSP message flow tests

## Test Examples

The `examples/` directory contains sample CUE files used for testing:

- **examples/simple/** - Basic configuration examples
- **examples/schema/** - Complex schema definitions (Kubernetes-like)
- **examples/multi-file/** - Multi-file projects with imports

These examples are used to:
- Test formatting functionality
- Verify LSP operations on real-world CUE code
- Provide documentation for users

## Writing Tests

### Unit Test Example

```go
func TestFormat_SimpleFormatting(t *testing.T) {
    input := []byte(`package test

foo:    "bar"
`)

    expected := `package test

foo: "bar"
`

    output, err := Format("test.cue", input)
    require.NoError(t, err)
    assert.Equal(t, expected, string(output))
}
```

### Integration Test Example

```go
func TestIntegration_InitializeWorkflow(t *testing.T) {
    var logBuf bytes.Buffer
    logger := log.New(&logBuf, "", 0)

    initRequest := protocol.InitializeRequest{
        Request: protocol.Request{
            RPC:    "2.0",
            ID:     1,
            Method: "initialize",
        },
        // ... params
    }

    // Test the handler
    handler.HandleMessage(logger, nil, "initialize", contents)

    // Verify response
    assert.Contains(t, output, "cuelang-lsp")
}
```

## Test Guidelines

### Best Practices

1. **Use Descriptive Names** - Test names should clearly describe what they test
   - Good: `TestFormat_NestedStructures`
   - Bad: `TestFormat1`

2. **Test One Thing** - Each test should verify a single behavior

3. **Use Table-Driven Tests** - For testing multiple scenarios:
   ```go
   tests := []struct {
       name     string
       input    []byte
       expected string
   }{
       {"case1", input1, expected1},
       {"case2", input2, expected2},
   }

   for _, tt := range tests {
       t.Run(tt.name, func(t *testing.T) {
           // test logic
       })
   }
   ```

4. **Use require vs assert**
   - `require.*` - Stops test execution on failure (use for prerequisites)
   - `assert.*` - Continues test execution on failure (use for checks)

5. **Clean Up Resources** - Use `t.TempDir()` for temporary files:
   ```go
   tmpDir := t.TempDir() // Automatically cleaned up
   ```

### What to Test

- ✅ Public functions and methods
- ✅ Error handling paths
- ✅ Edge cases (empty input, nil values, etc.)
- ✅ Integration between components
- ❌ Private implementation details
- ❌ External dependencies (mock them instead)

## Continuous Integration

Tests run automatically on every push and pull request via GitHub Actions:

- All tests must pass before merging
- Coverage reports are uploaded to Codecov
- HTML coverage reports are available as artifacts

See `.github/workflows/go.yml` for CI configuration.

## Coverage Goals

We aim for:
- **Overall coverage**: 70%+
- **Critical paths**: 90%+ (handlers, formatters)
- **New code**: 80%+ coverage required

Check current coverage:
```bash
make test-coverage-func
```

## Debugging Tests

### Run Specific Test

```bash
go test -v ./pkg/fmtr -run TestFormat_SimpleFormatting
```

### Run Tests with Race Detection

```bash
go test -race ./...
```

### Verbose Output

```bash
go test -v ./...
```

### Show Test Coverage While Running

```bash
go test -v -cover ./...
```

## Common Test Patterns

### Capturing stdout

```go
func captureStdout(f func()) string {
    old := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    f()

    w.Close()
    os.Stdout = old

    var buf bytes.Buffer
    io.Copy(&buf, r)
    return buf.String()
}
```

### Testing Logging

```go
var logBuf bytes.Buffer
logger := log.New(&logBuf, "", 0)

// ... use logger ...

logOutput := logBuf.String()
assert.Contains(t, logOutput, "expected log message")
```

### Creating Temporary Files

```go
tmpDir := t.TempDir()
testFile := filepath.Join(tmpDir, "test.cue")
err := os.WriteFile(testFile, []byte("content"), 0644)
require.NoError(t, err)
```

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Table Driven Tests](https://go.dev/wiki/TableDrivenTests)
- [Advanced Go Testing](https://about.sourcegraph.com/go/advanced-testing-in-go)
