# CUE LSP Examples

This directory contains example CUE files for testing and demonstrating the LSP capabilities.

## Examples

### 1. Simple Configuration (`simple/`)
A basic configuration example demonstrating:
- Simple field definitions
- Default values
- Type constraints
- Disjunctions (OR types)

**Files:**
- `config.cue` - Server and database configuration with constraints

### 2. Schema Definitions (`schema/`)
A more complex example with Kubernetes-like resource definitions demonstrating:
- Definition syntax (`#Name`)
- Nested structures
- Optional fields
- List constraints
- Pattern matching with constraints

**Files:**
- `kubernetes.cue` - Kubernetes Pod resource schema with validation

### 3. Multi-file Project (`multi-file/`)
A multi-file CUE project demonstrating:
- Shared type definitions across files
- Package organization
- Type reuse
- Complex validation rules with regular expressions

**Files:**
- `types.cue` - Shared type definitions (#User, #Service, #Endpoint)
- `instances.cue` - Concrete instances using the shared types

## Using These Examples

### Validate CUE Files
```bash
cd examples/simple
cue vet config.cue

cd ../schema
cue vet kubernetes.cue

cd ../multi-file
cue vet *.cue
```

### Format CUE Files
```bash
cue fmt config.cue
```

### Export to JSON
```bash
cue export config.cue
cue export kubernetes.cue
cue export -e users multi-file/*.cue
```

## Testing with LSP

These examples are used in the LSP test suite to verify:
- Document formatting
- Syntax validation
- Type checking
- Cross-file references (multi-file example)
- Error diagnostics
