# CLAUDE.md

## What is cfnpeek?

A read-only Go CLI for inspecting deployed AWS CloudFormation stacks. Displays resources, outputs, exports, and events in multiple formats (JSON, YAML, TOML, XML, INI, CSV, table).

**Core principles:**
- cfnpeek NEVER modifies stacks.
- cfnpeek must work with AWS read-only access (no write permissions required).

## CLI Structure

```
cfnpeek ls                        # List stacks
cfnpeek <stack>                   # All sections (resources, outputs, exports, events)
cfnpeek <stack> resources         # Resources only
cfnpeek <stack> outputs           # Outputs only
cfnpeek <stack> exports           # Exports only
cfnpeek <stack> events            # Stack events (deploy log)
cfnpeek <stack> resources,events  # Comma-separated
```

Key flags: `-f` format, `-r` region, `-p` profile, `-s` short, `--type` filter, `--grep` filter, `--limit` events.

## Build & Test

```bash
go build -o cfnpeek ./cmd/cfnpeek                    # Build
go test -race ./...                                   # Run tests
go vet ./...                                          # Vet
/Users/allixsenos/go/bin/golangci-lint run ./...      # Lint (v2)
```

## Project Structure

- `cmd/cfnpeek/` - Entrypoint with version ldflags
- `internal/cli/` - Cobra root command, ls subcommand, events handler
- `internal/aws/` - AWS SDK v2 client (interface-based for testing)
- `internal/model/` - Domain types (StackInfo, StackEvent, StackSummary, etc.)
- `internal/filter/` - Post-fetch filtering (--type, --grep)
- `internal/formatter/` - Output formatters (one file per format, ShortSetter interface)
- `.github/actions/install-cfnpeek/` - Composite GitHub Action
- `skill/` - Claude Code skill

## Conventions

- Conventional commits (`feat:`, `fix:`, `chore:`)
- release-please manages versioning and changelog
- GoReleaser handles cross-compilation (chained off release-please, not tag push)
- All formatters implement the `Formatter` interface in `internal/formatter/formatter.go`
- Table/CSV/INI formatters implement `ShortSetter` for `--short` mode
