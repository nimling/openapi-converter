# OpenAPI Converter

CLI tool to convert OpenAPI specifications to various formats including Nginx configurations and VitePress documentation.

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
go build
```

## Usage

The converter processes OpenAPI 3.1.X YAML files and outputs Nginx configurations and/or VitePress documentation.

```bash
# Basic usage (requires at least one input file/pattern)
openapi-converter [flags] <input-files...>

# Flags:
-o, --output string          Output directory for nginx configuration files
-d, --docs string            Output directory for API docs
-i, --index string           Path to index.md to write VitePress features
    --file-prefix string     Prefix for generated filenames
    --common-prefix string   Path prefix for VitePress links
    --write-introduction     Generate API docs introduction file
    --merge-responses-inline Merges inline allOf definitions into a single inline object

# Examples:
# Convert single file to nginx config
openapi-converter api.yaml -o ./nginx/

# Convert multiple files
openapi-converter *.yaml -o ./nginx/

# Generate VitePress docs from directory
openapi-converter ./specs/ -d ./docs/

# Full example with all outputs
openapi-converter api.yaml \
  -o ./nginx/ \
  -d ./docs/api/ \
  -i ./docs/index.md \
  --file-prefix myapi- \
  --common-prefix /api/v1 \
  --write-introduction
```

The converter will:
- Process .yml and .yaml files from provided paths/patterns
- Generate nginx config files with .conf.template extension
- Create VitePress documentation structure
- Update VitePress index.md features if specified

## Common Issues

- If you get "module not found" errors, ensure you have configured git with:
```bash
git config --global url."git@github.com:nimling/".insteadOf "github.com/nimling/"
```

- If the binary isn't found after installation, ensure `$GOPATH/bin` is in your PATH

## License

MIT
