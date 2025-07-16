package github_test

import (
	"context"
	"testing"

	"github.com/binary-install/semver/pkg/github"
)

func TestClient_ListTags(t *testing.T) {
	tests := []struct {
		name    string
		owner   string
		repo    string
		wantErr bool
	}{
		{
			name:    "empty owner",
			owner:   "",
			repo:    "repo",
			wantErr: true,
		},
		{
			name:    "empty repo",
			owner:   "owner",
			repo:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := github.NewClient("")
			_, err := client.ListTags(context.Background(), tt.owner, tt.repo)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ListReleases(t *testing.T) {
	tests := []struct {
		name    string
		owner   string
		repo    string
		wantErr bool
	}{
		{
			name:    "empty owner",
			owner:   "",
			repo:    "repo",
			wantErr: true,
		},
		{
			name:    "empty repo",
			owner:   "owner",
			repo:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := github.NewClient("")
			_, err := client.ListReleases(context.Background(), tt.owner, tt.repo)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListReleases() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
