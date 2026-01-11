## Generate Command Documentation

The `generate` command in this repository is used to automate the creation of specific structures or configurations required for the cuelang-lsp project. 

The purpose is to generate a generic Go based SDK that can be used to build language servers. It's built and in many ways inspired by the gopls project - which performs meta-model generation based on the [metaModel.json ](https://github.com/microsoft/language-server-protocol/blob/gh-pages/_specifications/lsp/3.18/metaModel/metaModel.json) in Microsoft's Language Server Protocol repository. 


### Usage
1. Navigate to the `cmd/generate` directory.
2. Run the command using:
   ```bash
   go run main.go
   ```
3. Follow any prompts or review the generated output in the appropriate directories.

## TODO: 
- [x] Cleanup main.go and add tests for existing funcs
- Add suppport for Notifications
- Add support for Requests / Responses
-