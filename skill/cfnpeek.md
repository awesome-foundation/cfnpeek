---
name: cfnpeek
description: Download and use cfnpeek to inspect AWS CloudFormation stacks (resources, outputs, exports)
---

# cfnpeek - CloudFormation Stack Inspector

cfnpeek is a read-only CLI tool for inspecting deployed CloudFormation stacks.

## Installation

Check if cfnpeek is already installed:

```bash
cfnpeek --version
```

If not installed, download the latest release to `~/.local/bin`:

```bash
# macOS (Apple Silicon)
curl -sSL https://github.com/awesome-foundation/cfnpeek/releases/latest/download/cfnpeek_darwin_arm64.tar.gz | tar xz -C ~/.local/bin

# macOS (Intel)
curl -sSL https://github.com/awesome-foundation/cfnpeek/releases/latest/download/cfnpeek_darwin_amd64.tar.gz | tar xz -C ~/.local/bin

# Linux (amd64)
curl -sSL https://github.com/awesome-foundation/cfnpeek/releases/latest/download/cfnpeek_linux_amd64.tar.gz | tar xz -C ~/.local/bin

# Or with go install (installs to ~/go/bin)
go install github.com/awesome-foundation/cfnpeek/cmd/cfnpeek@latest
```

## Usage

### List stacks

```bash
# List all active stacks in current region
cfnpeek ls

# List stacks in specific region
cfnpeek ls -r us-west-2

# List stacks with specific AWS profile
cfnpeek ls -p production
```

### Inspect stacks

```bash
# All info (resources, outputs, exports)
# Displays as table in TTY, JSON when piped
cfnpeek <stack-name-or-arn>

# Show only resources
cfnpeek <stack-name> --resources

# Show only outputs
cfnpeek <stack-name> --outputs

# Show only exports
cfnpeek <stack-name> --exports

# Specific region/profile
cfnpeek <stack-name> -r us-west-2 -p production

# Output formats: auto (default), table, json, yaml, toml, xml, ini, csv
cfnpeek <stack-name> -f json
cfnpeek <stack-name> -f yaml
cfnpeek <stack-name> -f table

# Pipe-friendly (auto-detects non-TTY, outputs JSON)
cfnpeek <stack-name> | jq '.outputs'
cfnpeek api-stack --resources | jq '.resources[] | select(.type == "AWS::Lambda::Function")'
```

### Flags

- `-r, --region` — AWS region (default: from AWS config)
- `-p, --profile` — AWS profile name (default: default)
- `-f, --format` — Output format: auto, table, json, yaml, toml, xml, ini, csv (default: auto)
- `--resources` — Show resources only
- `--outputs` — Show outputs only
- `--exports` — Show exports only

## Important

- cfnpeek is READ-ONLY. It never modifies stacks.
- Uses AWS default credential chain (env vars, profile, SSO, instance role).
- Output defaults to table in terminal, JSON when piped.
- Accepts both stack names and full ARNs.
