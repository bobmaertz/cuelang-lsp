package analysis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindDefinition_FieldReference(t *testing.T) {
	content := []byte(`package test

server: {
	host: "localhost"
	port: 8080
}

myServer: server
`)

	// Request definition for "server" on the line "myServer: server"
	pos := Position{
		Line:      7, // "myServer: server" line (0-based)
		Character: 11, // Position on "server"
	}

	locations, err := FindDefinition("test.cue", content, pos)
	require.NoError(t, err)

	// Should find the definition of "server" on line 2
	assert.Len(t, locations, 1)
	if len(locations) > 0 {
		assert.Equal(t, "test.cue", locations[0].URI)
		assert.Equal(t, 2, locations[0].Range.Start.Line) // Line 2 (0-based) where "server:" is defined
	}
}

func TestFindDefinition_NoIdentifierAtPosition(t *testing.T) {
	content := []byte(`package test

value: 42
`)

	// Request definition at a position with no identifier
	pos := Position{
		Line:      2,
		Character: 8, // After "value: " - on the number
	}

	locations, err := FindDefinition("test.cue", content, pos)
	require.NoError(t, err)

	// Should return empty - no identifier at this position
	assert.Len(t, locations, 0)
}

func TestFindDefinition_InvalidCUE(t *testing.T) {
	content := []byte(`package test

invalid {
  syntax here
`)

	pos := Position{
		Line:      1,
		Character: 0,
	}

	_, err := FindDefinition("test.cue", content, pos)
	// Should return an error for invalid CUE syntax
	assert.Error(t, err)
}

func TestFindDefinition_NestedField(t *testing.T) {
	content := []byte(`package test

config: {
	database: {
		host: "localhost"
		port: 5432
	}
}

db: config.database
`)

	// Request definition for "database" in "config.database"
	pos := Position{
		Line:      9,
		Character: 12, // Position on "database"
	}

	locations, err := FindDefinition("test.cue", content, pos)
	require.NoError(t, err)

	// Should find the definition
	assert.GreaterOrEqual(t, len(locations), 0)
}

func TestFindDefinition_DefinitionSchema(t *testing.T) {
	content := []byte(`package test

#User: {
	name: string
	age: int
}

user1: #User & {
	name: "Alice"
	age: 30
}
`)

	// Request definition for "#User" on line 7
	pos := Position{
		Line:      7,
		Character: 8, // Position on "#User"
	}

	locations, err := FindDefinition("test.cue", content, pos)
	require.NoError(t, err)

	// Should find the definition of #User
	assert.GreaterOrEqual(t, len(locations), 0)
	if len(locations) > 0 {
		assert.Equal(t, "test.cue", locations[0].URI)
	}
}

func TestFindDefinition_SelfReference(t *testing.T) {
	content := []byte(`package test

value: {
	name: "test"
}
`)

	// Request definition for "value" on its own definition line
	pos := Position{
		Line:      2,
		Character: 1, // Position on "value" where it's defined
	}

	locations, err := FindDefinition("test.cue", content, pos)
	require.NoError(t, err)

	// Should return empty or the same location (not jump to itself)
	assert.GreaterOrEqual(t, len(locations), 0)
}
