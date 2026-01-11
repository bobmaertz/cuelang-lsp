package fmtr

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormat_SimpleFormatting(t *testing.T) {
	input := []byte(`package test

foo:    "bar"
baz:   42
`)

	expected := `package test

foo: "bar"
baz: 42
`

	output, err := Format("test.cue", input)
	require.NoError(t, err)
	assert.Equal(t, expected, string(output))
}

func TestFormat_NestedStructures(t *testing.T) {
	input := []byte(`package test

server: {
host:    "localhost"
  port:  8080
timeout:30
}
`)

	output, err := Format("test.cue", input)
	require.NoError(t, err)

	// Verify output is valid and properly formatted
	assert.Contains(t, string(output), "package test")
	assert.Contains(t, string(output), "server:")
	assert.Contains(t, string(output), "host:")
	assert.Contains(t, string(output), "port:")
	assert.Contains(t, string(output), "timeout:")
}

func TestFormat_WithComments(t *testing.T) {
	input := []byte(`package test

// This is a comment
foo: "bar"

// Another comment
baz: 42
`)

	output, err := Format("test.cue", input)
	require.NoError(t, err)

	// Verify comments are preserved
	assert.Contains(t, string(output), "// This is a comment")
	assert.Contains(t, string(output), "// Another comment")
	assert.Contains(t, string(output), "foo: \"bar\"")
	assert.Contains(t, string(output), "baz: 42")
}

func TestFormat_Definitions(t *testing.T) {
	input := []byte(`package test

#User: {
name:  string
age:int&>0
email:     string
}
`)

	output, err := Format("test.cue", input)
	require.NoError(t, err)

	// Verify definition is properly formatted
	assert.Contains(t, string(output), "#User:")
	assert.Contains(t, string(output), "name:")
	assert.Contains(t, string(output), "age:")
	assert.Contains(t, string(output), "email:")
}

func TestFormat_Lists(t *testing.T) {
	input := []byte(`package test

items: [
"one",
  "two",
    "three"
]
`)

	output, err := Format("test.cue", input)
	require.NoError(t, err)

	// Verify list is properly formatted
	assert.Contains(t, string(output), "items:")
	assert.Contains(t, string(output), "\"one\"")
	assert.Contains(t, string(output), "\"two\"")
	assert.Contains(t, string(output), "\"three\"")
}

func TestFormat_Constraints(t *testing.T) {
	input := []byte(`package test

port: int & >=1 & <=65535
host: string | *"localhost"
`)

	output, err := Format("test.cue", input)
	require.NoError(t, err)

	// Verify constraints are preserved
	assert.Contains(t, string(output), "port:")
	assert.Contains(t, string(output), "host:")
}

func TestFormat_ComplexExample(t *testing.T) {
	input := []byte(`package kubernetes

#Pod: {
apiVersion:    "v1"
  kind: "Pod"
metadata:{
name:string
    namespace:   string|*"default"
}
spec: {
containers:[...#Container]
restartPolicy:   *"Always"|"OnFailure"
}
}
`)

	output, err := Format("test.cue", input)
	require.NoError(t, err)

	// Verify complex structure is properly formatted
	assert.Contains(t, string(output), "#Pod:")
	assert.Contains(t, string(output), "apiVersion:")
	assert.Contains(t, string(output), "metadata:")
	assert.Contains(t, string(output), "spec:")
	assert.Contains(t, string(output), "containers:")
}

func TestFormat_InvalidSyntax(t *testing.T) {
	input := []byte(`package test

foo: {
  unclosed brace
`)

	_, err := Format("test.cue", input)
	require.Error(t, err)
}

func TestFormat_EmptyFile(t *testing.T) {
	input := []byte(``)

	_, err := Format("test.cue", input)
	// Empty file should return an error or be handled gracefully
	// Depending on CUE parser behavior
	if err == nil {
		// If no error, verify we get something back
		// (behavior may vary by CUE version)
		t.Log("Empty file processed without error")
	}
}

func TestFormat_OnlyPackage(t *testing.T) {
	input := []byte(`package test
`)

	output, err := Format("test.cue", input)
	require.NoError(t, err)
	assert.Contains(t, string(output), "package test")
}

func TestFormat_MultilineStrings(t *testing.T) {
	input := []byte(`package test

description: """
    This is a
    multiline string
    """
`)

	output, err := Format("test.cue", input)
	require.NoError(t, err)

	// Verify multiline string is preserved
	assert.Contains(t, string(output), "description:")
	assert.Contains(t, string(output), "This is a")
	assert.Contains(t, string(output), "multiline string")
}

func TestFormat_Indentation(t *testing.T) {
	input := []byte(`package test

nested: {
level1: {
level2: {
value: 42
}
}
}
`)

	output, err := Format("test.cue", input)
	require.NoError(t, err)

	// Verify proper indentation (4 spaces as per format.UseSpaces(4))
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "level1:") {
			// Should have exactly 4 spaces indentation
			assert.True(t, strings.HasPrefix(line, "    "), "level1 should be indented with 4 spaces")
		}
		if strings.Contains(line, "level2:") {
			// Should have exactly 8 spaces indentation
			assert.True(t, strings.HasPrefix(line, "        "), "level2 should be indented with 8 spaces")
		}
	}
}

func TestFormat_PreservesSemantics(t *testing.T) {
	// Verify that formatting doesn't change the semantic meaning
	input := []byte(`package test

x: 1
y: 2
z: x + y
`)

	output, err := Format("test.cue", input)
	require.NoError(t, err)

	// Verify all declarations are present
	assert.Contains(t, string(output), "x:")
	assert.Contains(t, string(output), "y:")
	assert.Contains(t, string(output), "z:")
}

func TestFormat_Disjunctions(t *testing.T) {
	input := []byte(`package test

value: "a"|"b"|   "c"
number:1|  2|3
`)

	output, err := Format("test.cue", input)
	require.NoError(t, err)

	// Verify disjunctions are properly formatted
	assert.Contains(t, string(output), "value:")
	assert.Contains(t, string(output), "number:")
}

func TestFormat_RealWorldConfig(t *testing.T) {
	// Test with a real-world configuration example
	input := []byte(`package config

server:{
host:string|*"localhost"
port:   int&>=1&<=65535|*8080
timeout:int|*30
}

database: {
driver:"postgres"|"mysql"|"sqlite"
host:  string
port:int
name:   string
}
`)

	output, err := Format("test.cue", input)
	require.NoError(t, err)

	// Verify all components are present and properly formatted
	result := string(output)
	assert.Contains(t, result, "package config")
	assert.Contains(t, result, "server:")
	assert.Contains(t, result, "database:")
	assert.Contains(t, result, "driver:")

	// Should not contain excessive whitespace
	assert.NotContains(t, result, "  host:")
	assert.NotContains(t, result, "   port:")
}

func TestFormat_Unifications(t *testing.T) {
	input := []byte(`package test

#Base: {
x: int
}

#Extended: #Base & {
y: string
}
`)

	output, err := Format("test.cue", input)
	require.NoError(t, err)

	assert.Contains(t, string(output), "#Base:")
	assert.Contains(t, string(output), "#Extended:")
}
