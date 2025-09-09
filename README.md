# OpenAPI Converter

A powerful CLI tool and GitHub Action for converting OpenAPI 3.1.x specifications to various formats including Nginx configurations and VitePress documentation.

## Features

- **Multi-format conversion**: Generate Nginx configurations and VitePress documentation from OpenAPI specs
- **External reference resolution**: Automatically resolves `$ref` references to external files
- **Response merging**: Inline `allOf` definitions for cleaner output
- **Batch processing**: Process multiple specifications at once using glob patterns
- **Documentation sync**: Synchronize documentation files across repositories
- **GitHub Action support**: Available as a GitHub Action for CI/CD pipelines
- **Strict validation**: Ensures OpenAPI specifications meet documentation standards

## OpenAPI Compliance

For detailed information about OpenAPI specification requirements, validation rules, and best practices, see [OPENAPI.md](./OPENAPI.md).

## Prerequisites

- Go 1.23 or later
- Git installed
- GitHub authentication configured (SSH key or personal access token)

## Installation

### Using Go Install

```bash
go install github.com/nimling/openapi-converter@latest
```

For a specific version:

```bash
go install github.com/nimling/openapi-converter@v1.0.0
```

### From Source

```bash
git clone git@github.com:nimling/openapi-converter.git
cd openapi-converter
make build
```

### As GitHub Action

Add to your workflow:

```yaml
name: Convert OpenAPI Specs
on:
  push:
    paths:
      - 'api/**.yml'
      - 'api/**.yaml'

jobs:
  convert:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: nimling/openapi-converter@v1
        with:
          input-files: './api/*.yml'
          docs-output: './docs/api'
          nginx-output: './nginx'
          common-prefix: 'api'
```

## Usage

The tool provides two main commands: `convert` for processing OpenAPI specifications and `sync` for documentation synchronization.

### Convert Command

Processes OpenAPI 3.1.x YAML files and generates Nginx configurations and/or VitePress documentation.

```bash
openapi-converter convert [flags] <input-files...>
```

#### Flags

| Flag | Short | Description | Example |
|------|-------|-------------|---------|
| `--output` | `-o` | Output directory for Nginx configuration files | `-o ./nginx/` |
| `--docs` | `-d` | Output directory for VitePress API documentation | `-d ./docs/api/` |
| `--index` | `-i` | Path to generate/update VitePress index.md with features | `-i ./docs/index.md` |
| `--file-prefix` | | Prefix for generated file names | `--file-prefix api-` |
| `--common-prefix` | | URL path prefix for VitePress documentation links | `--common-prefix /api/v1` |
| `--write-introduction` | | Generate introduction page for API documentation | `--write-introduction` |
| `--merge-responses-inline` | | Merge allOf response definitions into single inline objects | `--merge-responses-inline` |

#### Examples

```bash
# Convert single file to Nginx config
openapi-converter convert api.yaml -o ./nginx/

# Convert multiple files with glob pattern
openapi-converter convert ./specs/*.yaml -o ./nginx/

# Generate VitePress documentation only
openapi-converter convert api.yaml -d ./docs/api/

# Generate both Nginx and VitePress with all options
openapi-converter convert api.yaml \
  -o ./nginx/ \
  -d ./docs/api/ \
  -i ./docs/index.md \
  --file-prefix myapi \
  --common-prefix /api/v1 \
  --write-introduction \
  --merge-responses-inline

# Process all YAML files in directory recursively
openapi-converter convert ./api/ -d ./documentation/
```

### Sync Command

Synchronize documentation files between directories using pattern-based mapping. Supports both individual file copying with renaming and full directory copying when target files exist.

```bash
openapi-converter sync -s <sync-map>
```

#### Flags

| Flag | Short | Description | Required |
|------|-------|-------------|----------|
| `--sync-map` | `-s` | JSON mapping file or inline JSON for sync command | Yes |

#### Mapping File Format

Create a JSON file or inline JSON with your synchronization rules. The destination must always specify a target filename:

```json
{
  "output/guides/myproject/index.md": [
    ".*docs/guide\\.md$",
    ".*docs/guide$"
  ],
  "output/api/myproject/index.md": [
    ".*docs/api\\.md$",
    ".*docs/api$"
  ],
  "output/tutorials/getting-started.md": [
    ".*tutorials/getting-started\\.md$",
    ".*docs/getting-started\\.md$"
  ]
}
```

**Behavior:**
- If source matches a **file**: Copies and renames to the destination filename
- If source matches a **directory**: Checks for a file with the destination filename inside, then copies the entire directory contents

#### Examples

```bash
# Sync documentation using mapping file
openapi-converter sync -s ./sync-config.json

# Typical CI/CD usage
openapi-converter sync --sync-map ./docs/mapping.json
```

### Output Formats

The converter generates:

#### Nginx Configuration (.conf.template)
- Location blocks with path patterns
- Method restrictions (GET, POST, PUT, DELETE)
- Upstream proxy configurations
- Security headers and CORS settings

#### VitePress Documentation
- Markdown files for each endpoint
- Interactive API documentation
- Request/response examples
- Schema definitions
- Navigation structure

## Validation

The converter enforces strict validation to ensure high-quality API documentation:

- **Required fields**: All paths must have `summary` and `description`
- **Operation IDs**: Every operation must have a unique `operationId`
- **Path format**: Paths must start with `/` and have valid segments
- **External references**: All `$ref` references must resolve successfully

For complete validation rules and OpenAPI compliance guidelines, see [OPENAPI.md](./OPENAPI.md).

## Common Issues

- If you get "module not found" errors, ensure you have configured git with:
```bash
git config --global url."git@github.com:nimling/".insteadOf "github.com/nimling/"
```

- If the binary isn't found after installation, ensure `$GOPATH/bin` is in your PATH

## License

MIT
