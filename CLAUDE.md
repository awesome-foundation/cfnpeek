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

## cfntop

Live TUI monitor for CloudFormation stacks. Bubbletea + lipgloss.

```
cfntop -r <region>        # Monitor stacks
cfntop -n 10              # Custom poll interval
cfntop --absolute-time    # Disable humanized times
```

- `internal/tui/app.go` - Bubbletea model, view, update loop
- `internal/tui/poller.go` - Stack/resource fetching, deploy boundary detection, deleted resource tracking
- `internal/tui/ecs.go` - ECS service deployment + failed task details
- `internal/tui/format.go` - Resource type shortening, humanized times
- `internal/tui/styles.go` - Lipgloss styles, status color coding

Active stacks poll every cycle. Inactive expanded stacks poll once per minute.

## Build & Test

```bash
go build -o cfnpeek ./cmd/cfnpeek                    # Build cfnpeek
go build -o cfntop ./cmd/cfntop                       # Build cfntop
go test -race ./...                                   # Run tests
go vet ./...                                          # Vet
/Users/allixsenos/go/bin/golangci-lint run ./...      # Lint (v2)
./watchexec.sh -r <region>                            # Dev mode with hot reload
```

## Project Structure

- `cmd/cfnpeek/` - cfnpeek entrypoint with version ldflags
- `cmd/cfntop/` - cfntop entrypoint (TUI)
- `internal/cli/` - Cobra root command, ls subcommand, events handler
- `internal/aws/` - AWS SDK v2 CloudFormation client (interface-based for testing)
- `internal/model/` - Domain types (StackInfo, StackEvent, StackSummary, etc.)
- `internal/filter/` - Post-fetch filtering (--type, --grep)
- `internal/formatter/` - Output formatters (one file per format, ShortSetter interface)
- `internal/tui/` - Bubbletea TUI for cfntop (poller, ECS, styles, formatting)
- `.github/actions/install-cfnpeek/` - Composite GitHub Action
- `skill/` - Claude Code skill

## Conventions

- Conventional commits (`feat:`, `fix:`, `chore:`)
- release-please manages versioning and changelog
- GoReleaser handles cross-compilation (chained off release-please, not tag push)
- All formatters implement the `Formatter` interface in `internal/formatter/formatter.go`
- Table/CSV/INI formatters implement `ShortSetter` for `--short` mode
