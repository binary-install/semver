# Design Document

## Overview

The semver-resolver tool is a multi-interface Go application that resolves semantic version ranges against GitHub repository tags and releases. It provides three interfaces: a Go library, a CLI tool, and a GitHub Action. The tool leverages the Masterminds/semver library for semver operations and integrates with GitHub's REST API to fetch repository version data.

## Architecture

The application follows a layered architecture with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    Interface Layer                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   Library   │  │     CLI     │  │   GitHub Action     │ │
│  │   Package   │  │    Tool     │  │     Wrapper         │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                   Service Layer                             │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │              SemverResolver Service                     │ │
│  │  • Version resolution logic                             │ │
│  │  • GitHub API integration                               │ │
│  │  • Authentication management                            │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                  Infrastructure Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │   GitHub    │  │    Auth     │  │      Keyring        │ │
│  │   Client    │  │   Manager   │  │      Manager        │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Components and Interfaces

### Core Library Package (`pkg/semver`)

```go
// Resolver provides semver resolution functionality
type Resolver interface {
    MaxSatisfying(ctx context.Context, owner, repo, constraint string) (string, error)
    MaxSatisfyingWithClient(ctx context.Context, client GitHubClient, owner, repo, constraint string) (string, error)
}

// GitHubClient interface for GitHub API operations
type GitHubClient interface {
    ListTags(ctx context.Context, owner, repo string) ([]*Tag, error)
    ListReleases(ctx context.Context, owner, repo string) ([]*Release, error)
}

// Tag represents a Git tag
type Tag struct {
    Name string
    SHA  string
}

// Release represents a GitHub release
type Release struct {
    TagName string
    Draft   bool
    Prerelease bool
}
```

### Authentication Manager (`pkg/auth`)

```go
// AuthManager handles GitHub authentication
type AuthManager interface {
    GetToken(ctx context.Context) (string, error)
    SetToken(ctx context.Context, token string) error
    DeleteToken(ctx context.Context) error
}

// TokenSource defines token retrieval priority
type TokenSource int

const (
    TokenSourceEnv TokenSource = iota
    TokenSourceKeyring
    TokenSourceFlag
)
```

### GitHub Client (`pkg/github`)

```go
// Client wraps GitHub API operations
type Client struct {
    client *github.Client
    auth   AuthManager
}

// NewClient creates a new GitHub client with authentication
func NewClient(ctx context.Context, auth AuthManager) (*Client, error)

// NewClientWithToken creates a client with explicit token
func NewClientWithToken(token string) *Client
```

### CLI Package (`cmd/semver`)

The CLI uses cobra for command structure:

```
semver
├── resolve <owner>/<repo>@<constraint>
├── auth
│   ├── login
│   ├── logout
│   └── status
└── version
```

### GitHub Action (`action.yml`)

```yaml
name: 'Semver Resolver'
description: 'Resolve semantic version ranges against GitHub repositories'
inputs:
  resolve:
    description: 'Repository and constraint in format owner/repo@constraint'
    required: true
  token:
    description: 'GitHub token for API access'
    required: false
    default: ${{ github.token }}
outputs:
  version:
    description: 'Resolved semantic version'
runs:
  using: 'node20'
  main: 'dist/index.js'
```

## Data Models

### Version Resolution Flow

1. **Input Parsing**: Parse `owner/repo@constraint` format
2. **Authentication**: Retrieve GitHub token from environment, keyring, or input
3. **API Fetching**: Fetch tags and releases from GitHub API
4. **Version Filtering**: Filter out non-semver tags and draft/prerelease versions
5. **Constraint Matching**: Use Masterminds/semver to find highest matching version
6. **Result Return**: Return resolved version or appropriate error

### Error Types

```go
// Custom error types for better error handling
type ResolverError struct {
    Type    ErrorType
    Message string
    Cause   error
}

type ErrorType int

const (
    ErrorTypeInvalidInput ErrorType = iota
    ErrorTypeAuthentication
    ErrorTypeNetwork
    ErrorTypeNotFound
    ErrorTypeRateLimit
    ErrorTypeSemver
)
```

## Error Handling

### Rate Limiting Strategy

- Implement exponential backoff with jitter
- Respect GitHub's rate limit headers
- Use authenticated requests when possible (5000 req/hour vs 60 req/hour)
- Provide clear error messages with retry suggestions

### Network Error Handling

- Distinguish between DNS, timeout, and HTTP errors
- Implement context-aware timeouts
- Provide retry logic for transient failures

### Authentication Error Handling

- Clear error messages for missing/invalid tokens
- Guidance on token setup and keyring usage
- Fallback to unauthenticated requests when appropriate

## Testing Strategy

### Unit Testing

- Mock GitHub API responses for consistent testing
- Test semver constraint matching with various scenarios
- Test error handling for all error types
- Test authentication flows including keyring operations

### Integration Testing

- Test against real GitHub repositories with known version histories
- Test rate limiting behavior with actual API calls
- Test CLI commands end-to-end
- Test GitHub Action in actual workflow environment

### Test Data

- Use public repositories with well-defined version histories
- Create test scenarios for edge cases (no tags, invalid semver, etc.)
- Mock API responses for consistent unit testing

### Performance Testing

- Benchmark version resolution with large numbers of tags
- Test memory usage with large repository datasets
- Validate API call efficiency and caching strategies

## Dependencies

### Core Dependencies

- `github.com/Masterminds/semver/v3` - Semver operations
- `github.com/google/go-github/v57/github` - GitHub API client
- `github.com/zalando/go-keyring` - System keyring integration
- `github.com/spf13/cobra` - CLI framework
- `golang.org/x/oauth2` - OAuth2 authentication

### Development Dependencies

- `github.com/stretchr/testify` - Testing framework
- `github.com/golang/mock` - Mock generation
- `github.com/goreleaser/goreleaser` - Release automation

## Security Considerations

### Token Storage

- Use system keyring for secure token storage
- Never log or expose tokens in error messages
- Support token rotation and expiration

### API Security

- Use HTTPS for all GitHub API calls
- Validate all input parameters to prevent injection
- Implement proper timeout and context cancellation

### GitHub Action Security

- Use minimal required permissions
- Validate all inputs in the action wrapper
- Use compiled binary instead of source code for performance and security