# Implementation Plan

- [ ] 1. Set up project structure and core interfaces
  - Create Go module with proper directory structure
  - Define core interfaces for Resolver, GitHubClient, and TokenManager
  - Set up basic error types and constants
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_

- [ ] 2. Implement token management functionality
  - [ ] 2.1 Create TokenManager interface and implementation
    - Implement environment variable token retrieval
    - Implement keyring token retrieval with fallback
    - Add proper error handling for missing tokens
    - _Requirements: 1.5, 2.5_

  - [ ] 2.2 Write unit tests for token management
    - Test environment variable retrieval
    - Test keyring retrieval and fallback behavior
    - Test error cases for missing/invalid tokens
    - _Requirements: 5.1, 5.2, 5.3_

- [ ] 3. Implement GitHub API client
  - [ ] 3.1 Create GitHub client wrapper
    - Implement GitHubClient interface using go-github
    - Add authentication support using TokenManager
    - Implement ListTags and ListReleases methods
    - _Requirements: 1.1, 1.3, 4.2_

  - [ ] 3.2 Add rate limiting and retry logic
    - Implement exponential backoff with jitter
    - Handle GitHub rate limit headers properly
    - Add context-aware timeouts
    - _Requirements: 4.1, 4.3_

  - [ ] 3.3 Write unit tests for GitHub client
    - Test API calls with mocked responses
    - Test rate limiting behavior
    - Test authentication flows
    - _Requirements: 5.1, 5.2, 5.3_

- [ ] 4. Implement core semver resolution logic
  - [ ] 4.1 Create semver resolver service
    - Parse owner/repo@constraint input format
    - Fetch tags and releases from GitHub API
    - Filter non-semver tags and draft/prerelease versions
    - Use Masterminds/semver to find highest matching version
    - _Requirements: 1.1, 1.2_

  - [ ] 4.2 Add comprehensive error handling
    - Handle invalid input formats
    - Handle network and API errors
    - Handle semver parsing errors
    - Provide descriptive error messages with context
    - _Requirements: 1.2, 1.3, 1.4, 5.1, 5.2, 5.3, 5.4_

  - [ ] 4.3 Write unit tests for resolver logic
    - Test version resolution with various constraints
    - Test error handling for all error types
    - Test edge cases (no tags, invalid semver, etc.)
    - _Requirements: 5.1, 5.2, 5.3, 5.4_

- [ ] 5. Implement CLI interface
  - [ ] 5.1 Create CLI structure using Cobra
    - Implement resolve command with owner/repo@constraint argument
    - Implement version command
    - Add proper flag handling for token input
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

  - [ ] 5.2 Integrate CLI with core resolver
    - Wire CLI commands to use semver resolver service
    - Handle CLI-specific error formatting and exit codes
    - Add proper output formatting for resolved versions
    - _Requirements: 2.1, 2.2, 2.3_

  - [ ] 5.3 Write CLI integration tests
    - Test resolve command with various inputs
    - Test error handling and exit codes
    - Test token handling via environment variables
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

- [ ] 6. Create library package interface
  - [ ] 6.1 Implement public library API
    - Create MaxSatisfying function for simple usage
    - Create MaxSatisfyingWithClient for advanced usage
    - Ensure proper context handling throughout
    - _Requirements: 1.1, 1.5_

  - [ ] 6.2 Write library integration tests
    - Test library functions with real GitHub repositories
    - Test authentication flows with library interface
    - Test error propagation from library functions
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_

- [ ] 7. Add comprehensive integration testing
  - [ ] 7.1 Create integration test suite
    - Test against real GitHub repositories with known version histories
    - Test rate limiting behavior with actual API calls
    - Test both authenticated and unauthenticated scenarios
    - _Requirements: 4.1, 4.2, 4.3_

  - [ ] 7.2 Add performance and edge case testing
    - Test with repositories having large numbers of tags
    - Test memory usage with large datasets
    - Test network error scenarios and recovery
    - _Requirements: 4.1, 5.1, 5.2, 5.3_