package resolver_test

import (
	"context"
	"testing"

	"github.com/binary-install/semver/pkg/github"
	"github.com/binary-install/semver/pkg/resolver"
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

func TestResolver_MaxSatisfying(t *testing.T) {
	tests := []struct {
		name       string
		tags       []string
		releases   []github.Release
		constraint string
		options    resolver.Options
		want       string
		wantErr    bool
	}{
		{
			name: "exact version",
			tags: []string{"v1.20.0", "v1.21.0", "v1.22.0"},
			constraint: "1.21.0",
			want:       "v1.21.0",
			wantErr:    false,
		},
		{
			name: "range constraint",
			tags: []string{"v1.20.0", "v1.21.0", "v1.21.5", "v1.22.0"},
			constraint: ">=1.20.0 <1.22.0",
			want:       "v1.21.5",
			wantErr:    false,
		},
		{
			name: "tilde range",
			tags: []string{"v1.21.0", "v1.21.5", "v1.22.0"},
			constraint: "~1.21.0",
			want:       "v1.21.5",
			wantErr:    false,
		},
		{
			name: "caret range",
			tags: []string{"v1.21.0", "v1.21.5", "v2.0.0"},
			constraint: "^1.21.0",
			want:       "v1.21.5",
			wantErr:    false,
		},
		{
			name: "no matching version",
			tags: []string{"v1.0.0", "v2.0.0"},
			constraint: ">=99.0.0",
			want:       "",
			wantErr:    true,
		},
		{
			name: "invalid constraint",
			tags: []string{"v1.0.0"},
			constraint: "invalid",
			want:       "",
			wantErr:    true,
		},
		{
			name: "empty owner",
			constraint: "^1.0.0",
			want:       "",
			wantErr:    true,
		},
		{
			name: "exclude prereleases by default",
			tags: []string{"v1.0.0", "v1.1.0-beta", "v1.2.0"},
			constraint: "^1.0.0",
			want:       "v1.2.0",
			wantErr:    false,
		},
		{
			name: "include prereleases when option set",
			tags: []string{"v1.0.0", "v1.1.0-beta", "v1.2.0"},
			constraint: "^1.0.0",
			options:    resolver.Options{IncludePrerelease: true},
			want:       "v1.2.0",
			wantErr:    false,
		},
		{
			name: "exclude draft releases by default",
			releases: []github.Release{
				{Tag: "v1.0.0", Draft: false},
				{Tag: "v1.1.0", Draft: true},
				{Tag: "v1.2.0", Draft: false},
			},
			constraint: "^1.0.0",
			want:       "v1.2.0",
			wantErr:    false,
		},
		{
			name: "include draft releases when option set",
			releases: []github.Release{
				{Tag: "v1.0.0", Draft: false},
				{Tag: "v1.1.0", Draft: true},
				{Tag: "v1.2.0", Draft: false},
			},
			constraint: "^1.0.0",
			options:    resolver.Options{IncludeDraft: true},
			want:       "v1.2.0",
			wantErr:    false,
		},
		{
			name: "handle go prefix",
			tags: []string{"go1.20.0", "go1.21.0", "go1.22.0"},
			constraint: "1.21.0",
			want:       "go1.21.0",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockGitHubClient{
				tags:     tt.tags,
				releases: tt.releases,
			}

			r := resolver.New(mockClient, tt.options)

			// Use default owner/repo for tests unless testing empty values
			owner := "test"
			repo := "repo"
			if tt.name == "empty owner" {
				owner = ""
			}

			got, err := r.MaxSatisfying(context.Background(), owner, repo, tt.constraint)

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
