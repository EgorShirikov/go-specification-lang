# AGENTS.md

This file provides guidance to agents when working with code in this repository.

## Build Commands

- **Build framework**: `cd framework && make build` (runs `go generate` for ANTLR parser, then compiles to `build/framework`)
- **Generate ANTLR parser only**: `cd framework && make gosl_parser` or `go generate ./...` (requires Java)
- **Run tests**: `cd framework && go test -v ./...`
- **Run single test**: `cd framework && go test -v ./plugin/dynamic_checks/... -run TestName`
- **Update golden files**: `cd framework && go test ./plugin/dynamic_checks/... -update`

## Critical Non-Obvious Requirements

- **ANTLR parser generation required before building**: The framework uses ANTLR-generated parsers for ACSL/GOSL specifications. Always run `make gosl_parser` or `go generate ./...` in the `framework/` directory before building or testing. Without this, imports to `gpma/framework/gosl/front/generated/grammar/*` will fail.
- **Plugin registration via init()**: All plugins MUST register themselves in their `init()` function using `plugin.Register("name", factory)`. The CLI imports plugins as blank imports (`_ "gpma/framework/plugin/..."`) to trigger registration.
- **Test data location**: Golden test files are in `framework/plugin/dynamic_checks/testdata/golden/` with structure matching `samples/projects/`. Use `-update` flag to regenerate golden files.
- **Go workspace**: This project uses Go workspaces (`go.work`). The framework module is at `./framework` with module path `gpma/framework`.
- **AST modification pattern**: Plugins receive `*data.ASTData` containing Go AST with specification annotations bound via `NodeInfo map[ast.Node]*AdditionalInfo`. Modifications should update the AST directly; the code generator handles formatting.
- **Specification language**: GPMA (Go Program Modelling Annotations) is an ACSL dialect for Go. Annotations are in `/*@ ... */` comments and parsed by the ANTLR-generated parser in `framework/gosl/front/`.
