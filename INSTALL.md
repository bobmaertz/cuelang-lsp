# Installation Guide

This guide covers installation and setup of the CUE Language Server for VS Code and Neovim.

## Prerequisites

- Go 1.25 or higher
- Git (for building from source)

## Installing the LSP Server

### Option 1: Install from Source (Recommended)

```bash
go install github.com/bobmaertz/cuelang-lsp/cmd/lsp@latest
```

The binary will be installed to `$GOPATH/bin/lsp` (typically `~/go/bin/lsp`).

### Option 2: Build from Source

```bash
git clone https://github.com/bobmaertz/cuelang-lsp.git
cd cuelang-lsp
make build
```

The binary will be available at `./bin/lsp`.

## Editor Setup

### VS Code

#### 1. Install the Generic LSP Extension

Install the [vscode-langservers-extracted](https://marketplace.visualstudio.com/items?itemName=ms-vscode.vscode-langservers-extracted) or use a generic LSP client like [vscode-glspc](https://marketplace.visualstudio.com/items?itemName=IWANABETHATGUY.vscode-glspc).

For manual configuration, add to your `settings.json`:

```json
{
  "cue.languageServer": {
    "enabled": true,
    "command": "/path/to/go/bin/lsp",
    "args": ["/tmp/cuelang-lsp.log"]
  }
}
```

#### 2. Configure File Association

Add to `settings.json`:

```json
{
  "files.associations": {
    "*.cue": "cue"
  }
}
```

#### 3. Debugging Configuration

Create `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug CUE LSP",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/lsp",
      "args": ["/tmp/cuelang-lsp-debug.log"],
      "showLog": true
    }
  ]
}
```

### Neovim

#### Using nvim-lspconfig (Recommended)

Add to your Neovim configuration:

```lua
-- init.lua or lua/lsp/cuelang.lua
local lspconfig = require('lspconfig')
local configs = require('lspconfig.configs')

-- Define cuelang-lsp if not already defined
if not configs.cuelang_lsp then
  configs.cuelang_lsp = {
    default_config = {
      cmd = { 'lsp', '/tmp/cuelang-lsp.log' },
      filetypes = { 'cue' },
      root_dir = lspconfig.util.root_pattern('.git', 'cue.mod'),
      settings = {},
    },
  }
end

-- Setup the LSP
lspconfig.cuelang_lsp.setup({
  on_attach = function(client, bufnr)
    -- Your keymaps and options here
    vim.keymap.set('n', 'gd', vim.lsp.buf.definition, { buffer = bufnr })
    vim.keymap.set('n', 'K', vim.lsp.buf.hover, { buffer = bufnr })
    vim.keymap.set('n', '<leader>f', vim.lsp.buf.format, { buffer = bufnr })
  end,
})
```

#### Manual Setup (Alternative)

```lua
local client = vim.lsp.start_client({
  name = 'cuelang-lsp',
  cmd = { 'lsp', '/tmp/cuelang-lsp.log' },
  root_dir = vim.fn.getcwd(),
})

if not client then
  vim.notify('cuelang-lsp failed to start', vim.log.levels.ERROR)
  return
end

vim.api.nvim_create_autocmd('FileType', {
  pattern = 'cue',
  callback = function()
    vim.lsp.buf_attach_client(0, client)
  end,
})
```

## Debugging the LSP

### Enable Debug Logging

Pass a log file path as the first argument when starting the LSP:

```bash
lsp /tmp/cuelang-lsp.log
```

The log file will contain detailed information about LSP operations, requests, and responses.

### Common Debug Steps

1. **Check if the LSP is running**:
   ```bash
   # In Neovim
   :LspInfo

   # Check log file
   tail -f /tmp/cuelang-lsp.log
   ```

2. **Verify the binary path**:
   ```bash
   which lsp
   # or
   ls -la ~/go/bin/lsp
   ```

3. **Test the LSP manually**:
   ```bash
   echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | lsp /tmp/test.log
   ```

4. **Check file type detection** (Neovim):
   ```vim
   :set filetype?
   ```

### VS Code Debugging

1. View LSP output: **View** → **Output** → Select "CUE Language Server"
2. Enable verbose logging in `settings.json`:
   ```json
   {
     "cue.trace.server": "verbose"
   }
   ```

### Neovim Debugging

```lua
-- Enable LSP logging
vim.lsp.set_log_level('debug')

-- View logs
:lua vim.cmd('e ' .. vim.lsp.get_log_path())
```

## Verifying Installation

Create a test file `test.cue`:

```cue
package test

message: "Hello, CUE!"
value: 42
```

Open the file in your editor and verify:
- Formatting works (`:lua vim.lsp.buf.format()` in Neovim or **Format Document** in VS Code)
- LSP server is connected (check status bar or run `:LspInfo` in Neovim)

## Troubleshooting

**LSP not starting**:
- Verify the binary exists and is executable: `ls -la $(which lsp)`
- Check the log file for errors
- Ensure Go bin directory is in your PATH: `echo $PATH | grep go/bin`

**No formatting**:
- Check that the CUE CLI tools are installed: `which cue`
- Verify file type is set to 'cue'

**Permission denied**:
```bash
chmod +x ~/go/bin/lsp
```

For additional help, see [TESTING.md](TESTING.md) or open an issue on GitHub.
