# CUE Language Server Protocol (LSP)

<>[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/bobmaertz/cuelang-lsp)](https://goreportcard.com/report/github.com/bobmaertz/cuelang-lsp)
[![Build Status](https://github.com/bobmaertz/cuelang-lsp/workflows/Build/badge.svg)](https://github.com/bobmaertz/cuelang-lsp/actions)

A Language Server Protocol implementation for the CUE language.

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

- Full cuelang file formatting in your IDE

## Installation

```bash
go get -u github.com/bobmaertz/cuelang-lsp
```

## Usage

Describe how to use your CUE LSP, including any command-line options or integration with popular editors.

```bash
lsp [options]
```

## Configuration

The lsp currently supports a single filepath argument for a debug log. 
```bash 
lsp ~/out.log
```

## Development

### Prerequisites

- Go 1.22 or higher
- CUE 0.10.0 or higher

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

There are two commands in the /cmd directory, one for the language server and the other for formatting. The LSP binay can be run using the binary located at /bin/lsp 

To test the formatting functionality, use the following command. 
```bash 
go run fmtr <path/to/cuefile>.cue 
```

### Register LSP with Nvim 
I've setup my local environment to use this LSP for *.cue files. 

```lua

local client = vim.lsp.start_client {
    name = "cuelang-lsp",
    cmd = { "/Users/bob/personal/cuelang-lsp/lsp" }
}

if not client then
    vim.notify "cuelang-lsp did not start"
    return
end

vim.api.nvim_create_autocmd("FileType", {
    pattern = "cue",
    callback = function()
        vim.lsp.buf_attach_client(0, client)
    end

})

```


### Register Filetype for Cuelang 
Only necessary if there the pattern used in the lsp configuration is not based on the standard treesitter filetype for cue.  

```
:setfiletype cue 
:LspRestart
```
## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.jjk
