# CLAUDE.md

## What is cfnpeek?

A read-only Go CLI for inspecting deployed AWS CloudFormation stacks. Displays resources, outputs, and exports in multiple formats (JSON, YAML, TOML, XML, INI, CSV, table).

**Core principles:**
- cfnpeek NEVER modifies stacks.
- cfnpeek must work with AWS read-only access (no write permissions required).

## Build & Test

```bash
go build -o cfnpeek ./cmd/cfnpeek   # Build
go test -race ./...                  # Run tests
go vet ./...                         # Vet
```

## Project Structure

- `cmd/cfnpeek/` - Entrypoint with version ldflags
- `internal/cli/` - Cobra command definitions
- `internal/aws/` - AWS SDK v2 client (interface-based for testing)
- `internal/model/` - Domain types (StackInfo, Resource, Output, Export)
- `internal/formatter/` - Output formatters (one file per format)
- `.github/actions/install-cfnpeek/` - Composite GitHub Action
- `skill/` - Claude Code skill

## Conventions

- Conventional commits (`feat:`, `fix:`, `chore:`)
- release-please manages versioning and changelog
- GoReleaser handles cross-compilation on tag push
- All formatters implement the `Formatter` interface in `internal/formatter/formatter.go`
