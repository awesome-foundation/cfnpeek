# cfnpeek

Read-only CLI for inspecting deployed AWS CloudFormation stacks. Lists resources, outputs, and exports in multiple formats.

## Install

Download the latest release for your platform from [Releases](https://github.com/awesome-foundation/cfnpeek/releases), or:

```bash
# macOS (Apple Silicon)
curl -sSL https://github.com/awesome-foundation/cfnpeek/releases/latest/download/cfnpeek_0.3.1_darwin_arm64.tar.gz | tar xz -C /usr/local/bin # x-release-please-version

# macOS (Intel)
curl -sSL https://github.com/awesome-foundation/cfnpeek/releases/latest/download/cfnpeek_0.3.1_darwin_amd64.tar.gz | tar xz -C /usr/local/bin # x-release-please-version

# Linux (amd64)
curl -sSL https://github.com/awesome-foundation/cfnpeek/releases/latest/download/cfnpeek_0.3.1_linux_amd64.tar.gz | tar xz -C /usr/local/bin # x-release-please-version

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
| INI | `-f ini` | Named keys and sections. |
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

## Example Output

### Table (default in terminal)

```
Stack: awesome-vpc-prod
Status: UPDATE_COMPLETE

Resources (4)
LOGICAL ID          PHYSICAL ID                          TYPE                            STATUS
InternetGateway     igw-0a1b2c3d4e5f6g7h8                AWS::EC2::InternetGateway       CREATE_COMPLETE
PublicRouteTable    rtb-0f1e2d3c4b5a6978                 AWS::EC2::RouteTable            CREATE_COMPLETE
VPC                 vpc-0a1b2c3d4e5f6g7h8                AWS::EC2::VPC                   CREATE_COMPLETE
VPCGatewayAttach    aweso-VPCGa-1A2B3C4D5E6F             AWS::EC2::VPCGatewayAttachment  CREATE_COMPLETE

Outputs (2)
KEY        VALUE                      EXPORT NAME
VpcId      vpc-0a1b2c3d4e5f6g7h8      awesome-vpc-prod-VpcId
VpcCidr    10.0.0.0/16                awesome-vpc-prod-VpcCidr

Exports (2)
NAME                        VALUE
awesome-vpc-prod-VpcCidr    10.0.0.0/16
awesome-vpc-prod-VpcId      vpc-0a1b2c3d4e5f6g7h8
```

### JSON (default when piped)

```json
{
  "stack_name": "awesome-vpc-prod",
  "stack_id": "arn:aws:cloudformation:eu-west-1:123456789:stack/awesome-vpc-prod/abc123",
  "status": "UPDATE_COMPLETE",
  "resources": [
    {
      "logical_id": "VPC",
      "physical_id": "vpc-0a1b2c3d4e5f6g7h8",
      "type": "AWS::EC2::VPC",
      "status": "CREATE_COMPLETE",
      "last_updated": "2026-01-15T10:30:00Z"
    }
  ],
  "outputs": [
    {
      "key": "VpcId",
      "value": "vpc-0a1b2c3d4e5f6g7h8",
      "export_name": "awesome-vpc-prod-VpcId"
    }
  ],
  "exports": [
    {
      "name": "awesome-vpc-prod-VpcId",
      "value": "vpc-0a1b2c3d4e5f6g7h8"
    }
  ]
}
```

### YAML

```yaml
stack_name: awesome-vpc-prod
stack_id: arn:aws:cloudformation:eu-west-1:123456789:stack/awesome-vpc-prod/abc123
status: UPDATE_COMPLETE
resources:
  - logical_id: VPC
    physical_id: vpc-0a1b2c3d4e5f6g7h8
    type: AWS::EC2::VPC
    status: CREATE_COMPLETE
    last_updated: "2026-01-15T10:30:00Z"
outputs:
  - key: VpcId
    value: vpc-0a1b2c3d4e5f6g7h8
    export_name: awesome-vpc-prod-VpcId
exports:
  - name: awesome-vpc-prod-VpcId
    value: vpc-0a1b2c3d4e5f6g7h8
```

### TOML

```toml
[stack]
name = "awesome-vpc-prod"
id = "arn:aws:cloudformation:eu-west-1:123456789:stack/awesome-vpc-prod/abc123"
status = "UPDATE_COMPLETE"

[resources.VPC]
physical_id = "vpc-0a1b2c3d4e5f6g7h8"
type = "AWS::EC2::VPC"
status = "CREATE_COMPLETE"
last_updated = "2026-01-15T10:30:00Z"

[outputs]
VpcId = "vpc-0a1b2c3d4e5f6g7h8"

[exports]
awesome-vpc-prod-VpcId = "vpc-0a1b2c3d4e5f6g7h8"
```

### INI

```ini
[stack]
name = awesome-vpc-prod
id = arn:aws:cloudformation:eu-west-1:123456789:stack/awesome-vpc-prod/abc123
status = UPDATE_COMPLETE

[resource.VPC]
physical_id = vpc-0a1b2c3d4e5f6g7h8
type = AWS::EC2::VPC
status = CREATE_COMPLETE
last_updated = 2026-01-15T10:30:00Z

[outputs]
VpcId = vpc-0a1b2c3d4e5f6g7h8

[exports]
awesome-vpc-prod-VpcId = vpc-0a1b2c3d4e5f6g7h8
```

### `cfnpeek ls`

```
Stacks (3)
NAME                STATUS            UPDATED                DESCRIPTION
awesome-bastion     CREATE_COMPLETE   2026-02-10T14:22:00Z   SSH bastion host
awesome-vpc-prod    UPDATE_COMPLETE   2026-03-20T09:15:00Z   Production VPC
awesome-web-prod    UPDATE_COMPLETE   2026-03-25T16:45:00Z   ECS cluster and ALB
```

## License

MIT
