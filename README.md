# cfnpeek

Read-only CLI for inspecting deployed AWS CloudFormation stacks. Lists resources, outputs, and exports in multiple formats.

## Install

Download the latest release for your platform from [Releases](https://github.com/awesome-foundation/cfnpeek/releases), or:

```bash
# macOS (Apple Silicon)
curl -sSL https://github.com/awesome-foundation/cfnpeek/releases/latest/download/cfnpeek_darwin_arm64.tar.gz | tar xz -C /usr/local/bin

# macOS (Intel)
curl -sSL https://github.com/awesome-foundation/cfnpeek/releases/latest/download/cfnpeek_darwin_amd64.tar.gz | tar xz -C /usr/local/bin

# Linux (amd64)
curl -sSL https://github.com/awesome-foundation/cfnpeek/releases/latest/download/cfnpeek_linux_amd64.tar.gz | tar xz -C /usr/local/bin

# From source
go install github.com/awesome-foundation/cfnpeek/cmd/cfnpeek@latest
```

## Usage

```bash
# Show everything (resources, outputs, exports)
cfnpeek my-stack

# Specific sections
cfnpeek my-stack --resources
cfnpeek my-stack --outputs
cfnpeek my-stack --exports

# Output formats: json, yaml, toml, xml, ini, csv, table
cfnpeek my-stack -f json
cfnpeek my-stack -f yaml

# Region and profile
cfnpeek my-stack -r eu-west-1 -p production

# Stack ARN works too
cfnpeek arn:aws:cloudformation:us-east-1:123456789:stack/my-stack/guid
```

Output defaults to **table** in a terminal, **JSON** when piped.

## Output Formats

| Format | Flag | Notes |
|--------|------|-------|
| Table | `-f table` | Default for TTY. Aligned columns. |
| JSON | `-f json` | Default when piped. Indented. |
| YAML | `-f yaml` | |
| TOML | `-f toml` | |
| XML | `-f xml` | With XML declaration. |
| INI | `-f ini` | Indexed keys for lists. |
| CSV | `-f csv` | Separate header rows per section. |

## Authentication

cfnpeek uses the standard AWS SDK credential chain:

1. Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION_TOKEN`)
2. Shared credentials / config file (`~/.aws/credentials`, `~/.aws/config`)
3. AWS SSO
4. EC2 instance role / ECS task role

Use `--profile` / `-p` to select a named profile, `--region` / `-r` to override the region.

## GitHub Action

```yaml
- uses: awesome-foundation/cfnpeek/.github/actions/install-cfnpeek@v1
  with:
    github-token: ${{ secrets.GITHUB_TOKEN }}
    version: latest  # or a specific tag like v1.0.0
```

## License

MIT
