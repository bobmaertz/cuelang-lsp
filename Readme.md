# CUE Language Server Protocol (LSP)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/bobmaertz/cuelang-lsp)](https://goreportcard.com/report/github.com/bobmaertz/cuelang-lsp)
[![Build Status](https://github.com/bobmaertz/cuelang-lsp/actions/workflows/go.yml/badge.svg)](https://github.com/bobmaertz/cuelang-lsp/actions)

A Language Server Protocol implementation for the CUE language.

🚨 🏗️ **This project is a work in progress and may not be completely functional yet** 🏗️ 🚨

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Development](#development)
  - [Prerequisites](#prerequisites)
  - [Building](#building)
  - [Testing](#testing)
- [License](#license)

## Features

- **Document Formatting** - Format CUE files with proper indentation and spacing
- **Go to Definition** - Jump to field and type definitions within your CUE files
- **Document Synchronization** - Real-time updates as you edit

## Installation

### Quick Start

```bash
go install github.com/bobmaertz/cuelang-lsp/cmd/lsp@latest
```

For detailed installation instructions, editor setup (VS Code, Neovim), and debugging, see [INSTALL.md](INSTALL.md).

## Usage

Start the LSP server with optional debug logging:

```bash
lsp [log-file-path]
```

Example:
```bash
lsp /tmp/cuelang-lsp.log
```

The server communicates via JSON-RPC over stdin/stdout and is designed to be used with LSP clients in editors like VS Code and Neovim.

## Development

### Prerequisites

- Go 1.25 or higher
- CUE 0.15.3 or higher

### Building

```bash
git clone https://github.com/bobmaertz/cuelang-lsp.git
cd cuelang-lsp
make build
```

### Testing

```bash
make test
```

### Running

The LSP server is located at `cmd/lsp`. After building:

```bash
./bin/lsp /tmp/debug.log
```

For editor integration examples, see [INSTALL.md](INSTALL.md).
## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
