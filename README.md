# semver-resolver

A Go library and CLI tool for resolving semantic version constraints against GitHub repository tags and releases.

## Features

- 🔍 Resolve semantic version constraints against GitHub repositories
- 📚 Available as both a Go library and CLI tool
- 🔐 GitHub token authentication support
- 🏷️ Support for GitHub releases (tags optional via flag)
- 🚀 Prerelease and draft release filtering
- ⚡ Efficient pagination and rate limit handling
- 🛡️ Comprehensive error handling

## Installation

### CLI Tool

```bash
go install github.com/binary-install/semver/cmd/semver@latest
```

### Library

```bash
go get github.com/binary-install/semver
```

## Usage

### CLI

#### Version Resolution

```bash
# Find exact version
semver resolve golang/go@1.21.0

# Find latest patch version
semver resolve golang/go@~1.21.0

# Find latest minor version
semver resolve golang/go@^1.21.0

# Find version within range
semver resolve golang/go@">=1.20.0 <1.22.0"

# Include prereleases
semver resolve --prerelease golang/go@^1.21.0

# Include git tags (by default only releases are used)
semver resolve --tags golang/go@^1.21.0

# With custom token
semver resolve --token ghp_xxxxx golang/go@^1.21.0
```

#### Token Management

```bash
# Set token interactively (hidden input)
semver token set

# Set token via argument
semver token set ghp_your_token_here

# Set token via stdin
echo "ghp_your_token_here" | semver token set

# Get current token (masked for security)
semver token get

# Delete stored token
semver token delete
```

### Library

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/binary-install/semver"
)

func main() {
    ctx := context.Background()

    // Basic usage
    version, err := semver.MaxSatisfying(ctx, "golang", "go", "~1.21.0", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Latest version: %s\n", version)

    // With options
    opts := &semver.Options{
        Token:             "ghp_your_token",
        IncludePrerelease: true,
        IncludeDraft:      false,
        IncludeTags:       false, // Only use releases by default
    }
    version, err = semver.MaxSatisfying(ctx, "owner", "repo", "^1.0.0", opts)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Latest version: %s\n", version)
}
```

## Authentication

The tool supports GitHub authentication through:

1. `GITHUB_TOKEN` environment variable
2. `GH_TOKEN` environment variable
3. System keyring (via `zalando/go-keyring`)
4. `--token` flag (CLI) or `Token` field (library)

Authentication increases API rate limits from 60 to 5,000 requests per hour.

## Semantic Version Constraints

The tool uses [Masterminds/semver](https://github.com/Masterminds/semver) for constraint parsing:

- `1.2.3` - Exact version
- `~1.2.3` - Patch releases: `>=1.2.3 <1.3.0`
- `^1.2.3` - Minor releases: `>=1.2.3 <2.0.0`
- `>=1.2.3` - Minimum version
- `>1.2.3 <2.0.0` - Version range
- `1.2.x` - Wildcard

## Development

### Running Tests

```bash
# Unit tests
go test ./...

# Integration tests (requires GitHub token)
go test -tags=integration ./test/integration
```

### Building

```bash
# Build CLI
go build -o semver ./cmd/semver

# Build with version
go build -ldflags "-X github.com/binary-install/semver/cmd/semver/cli.Version=v1.0.0" -o semver ./cmd/semver
```

## License

MIT License - see LICENSE file for details.
