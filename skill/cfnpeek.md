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

If not installed, download the latest release:

```bash
# macOS (Apple Silicon)
curl -sSL https://github.com/awesome-foundation/cfnpeek/releases/latest/download/cfnpeek_darwin_arm64.tar.gz | tar xz -C /usr/local/bin

# macOS (Intel)
curl -sSL https://github.com/awesome-foundation/cfnpeek/releases/latest/download/cfnpeek_darwin_amd64.tar.gz | tar xz -C /usr/local/bin

# Linux (amd64)
curl -sSL https://github.com/awesome-foundation/cfnpeek/releases/latest/download/cfnpeek_linux_amd64.tar.gz | tar xz -C /usr/local/bin

# Or with go install
go install github.com/awesome-foundation/cfnpeek/cmd/cfnpeek@latest
```

## Usage

```bash
# All info (resources, outputs, exports)
# Table format in TTY, JSON when piped
cfnpeek <stack-name-or-arn>

# Specific sections
cfnpeek <stack-name> --resources
cfnpeek <stack-name> --outputs
cfnpeek <stack-name> --exports

# Specific region/profile
cfnpeek <stack-name> -r us-west-2 -p production

# Output formats: json, yaml, toml, xml, ini, csv, table
cfnpeek <stack-name> -f json
cfnpeek <stack-name> -f yaml

# Pipe-friendly (auto-detects non-TTY, outputs JSON)
cfnpeek <stack-name> | jq '.outputs'
```

## Important

- cfnpeek is READ-ONLY. It never modifies stacks.
- Uses AWS default credential chain (env vars, profile, SSO, instance role).
- Output defaults to table in terminal, JSON when piped.
- Accepts both stack names and full ARNs.
