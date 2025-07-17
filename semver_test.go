package semver_test

import (
	"context"
	"testing"

	"github.com/binary-install/semver"
	"github.com/binary-install/semver/pkg/github"
)

// mockGitHubClient is a mock implementation of github.Client for testing.
type mockGitHubClient struct {
	tags     []string
	releases []github.Release
	err      error
}

func (m *mockGitHubClient) ListTags(ctx context.Context, owner, repo string) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.tags, nil
}

func (m *mockGitHubClient) ListReleases(ctx context.Context, owner, repo string) ([]github.Release, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.releases, nil
}

func TestMaxSatisfying(t *testing.T) {
	tests := []struct {
		name       string
		owner      string
		repo       string
		constraint string
		tags       []string
		releases   []github.Release
		opts       *semver.Options
		want       string
		wantErr    bool
	}{
		{
			name:       "basic version resolution",
			owner:      "test",
			repo:       "repo",
			constraint: "^1.0.0",
			releases: []github.Release{
				{Tag: "v1.0.0", Draft: false},
				{Tag: "v1.1.0", Draft: false},
				{Tag: "v1.2.0", Draft: false},
				{Tag: "v2.0.0", Draft: false},
			},
			want:    "v1.2.0",
			wantErr: false,
		},
		{
			name:       "with custom GitHub client",
			owner:      "test",
			repo:       "repo",
			constraint: "~1.1.0",
			tags:       []string{"v1.0.0", "v1.1.0", "v1.1.5", "v1.2.0"},
			opts: &semver.Options{
				IncludeTags: true,
				GitHubClient: &mockGitHubClient{
					tags: []string{"v1.0.0", "v1.1.0", "v1.1.5", "v1.2.0"},
				},
			},
			want:    "v1.1.5",
			wantErr: false,
		},
		{
			name:       "include prereleases",
			owner:      "test",
			repo:       "repo",
			constraint: "^1.0.0",
			releases: []github.Release{
				{Tag: "v1.0.0", Prerelease: false},
				{Tag: "v1.1.0-beta", Prerelease: true},
				{Tag: "v1.2.0", Prerelease: false},
			},
			opts: &semver.Options{
				IncludePrerelease: true,
				GitHubClient: &mockGitHubClient{
					releases: []github.Release{
						{Tag: "v1.0.0", Prerelease: false},
						{Tag: "v1.1.0-beta", Prerelease: true},
						{Tag: "v1.2.0", Prerelease: false},
					},
				},
			},
			want:    "v1.2.0",
			wantErr: false,
		},
		{
			name:       "include draft releases",
			owner:      "test",
			repo:       "repo",
			constraint: "^1.0.0",
			releases: []github.Release{
				{Tag: "v1.0.0", Draft: false},
				{Tag: "v1.1.0", Draft: true},
				{Tag: "v1.2.0", Draft: false},
			},
			opts: &semver.Options{
				IncludeDraft: true,
				GitHubClient: &mockGitHubClient{
					releases: []github.Release{
						{Tag: "v1.0.0", Draft: false},
						{Tag: "v1.1.0", Draft: true},
						{Tag: "v1.2.0", Draft: false},
					},
				},
			},
			want:    "v1.2.0",
			wantErr: false,
		},
		{
			name:       "invalid constraint",
			owner:      "test",
			repo:       "repo",
			constraint: "invalid",
			releases: []github.Release{
				{Tag: "v1.0.0", Draft: false},
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.opts
			if opts == nil {
				opts = &semver.Options{}
			}
			if opts.GitHubClient == nil {
				opts.GitHubClient = &mockGitHubClient{
					tags:     tt.tags,
					releases: tt.releases,
				}
			}

			got, err := semver.MaxSatisfying(context.Background(), tt.owner, tt.repo, tt.constraint, opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("MaxSatisfying() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("MaxSatisfying() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxSatisfyingWithToken(t *testing.T) {
	mockClient := &mockGitHubClient{
		releases: []github.Release{
			{Tag: "v1.0.0", Draft: false},
			{Tag: "v1.1.0", Draft: false},
			{Tag: "v2.0.0", Draft: false},
		},
	}

	// We can't easily test the token functionality without a real GitHub API,
	// but we can verify the function works with our mock
	ctx := context.Background()
	opts := &semver.Options{
		Token:        "fake-token",
		GitHubClient: mockClient,
	}

	got, err := semver.MaxSatisfying(ctx, "test", "repo", "^1.0.0", opts)
	if err != nil {
		t.Fatalf("MaxSatisfyingWithToken() unexpected error: %v", err)
	}

	want := "v1.1.0"
	if got != want {
		t.Errorf("MaxSatisfyingWithToken() = %v, want %v", got, want)
	}
}
