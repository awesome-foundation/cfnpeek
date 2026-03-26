---
name: cfnpeek
description: Download and use cfnpeek to inspect AWS CloudFormation stacks (resources, outputs, exports, events)
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

# Or with go install (installs to ~/go/bin)
go install github.com/awesome-foundation/cfnpeek/cmd/cfnpeek@latest
```

## Usage

```bash
cfnpeek ls                        # List all stacks in region
cfnpeek <stack>                   # All sections (resources, outputs, exports, events)
cfnpeek <stack> resources         # Resources only
cfnpeek <stack> outputs           # Outputs only
cfnpeek <stack> exports           # Exports only
cfnpeek <stack> events            # Stack events (deploy log)
cfnpeek <stack> resources,events  # Combine with commas
```

### Flags

```bash
-r, --region     AWS region
-p, --profile    AWS profile
-f, --format     auto, table, json, yaml, toml, xml, ini, csv (default: auto)
-s, --short      Compact output (fewer columns in table/csv/ini)
    --type       Filter resources by type (fuzzy, e.g. --type ec2)
    --grep       Filter outputs/exports by key or value (e.g. --grep vpc)
    --limit      Max events to show (default: 20, 0 for all)
```

### Examples

```bash
cfnpeek my-stack -r eu-west-1 -p production
cfnpeek my-stack resources --type lambda
cfnpeek my-stack outputs --grep endpoint
cfnpeek my-stack events --limit 10
cfnpeek my-stack -f json | jq '.resources[] | select(.type == "AWS::Lambda::Function")'
cfnpeek my-stack -s                          # compact table output
```

## Important

- cfnpeek is READ-ONLY. It never modifies stacks.
- Uses AWS default credential chain (env vars, profile, SSO, instance role).
- Output defaults to table in terminal, JSON when piped.
- Accepts both stack names and full ARNs.
