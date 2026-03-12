package config

import (
	"testing"
	"time"
)

func TestDefaults(t *testing.T) {
	cfg := Defaults()

	checks := []struct {
		name string
		got  any
		want any
	}{
		{"GoVersion", cfg.GoVersion, "1.26"},
		{"License", cfg.License, "MIT"},
		{"BuildTool", cfg.BuildTool, "taskfile"},
		{"UseLinter", cfg.UseLinter, true},
		{"UseGoReleaser", cfg.UseGoReleaser, false},
		{"UseEnvExample", cfg.UseEnvExample, true},
		{"UseCI", cfg.UseCI, true},
		{"UseDocker", cfg.UseDocker, false},
		{"UseDependabot", cfg.UseDependabot, false},
		{"UseClaude", cfg.UseClaude, false},
		{"UseClaudeCI", cfg.UseClaudeCI, false},
		{"Year", cfg.Year, time.Now().Year()},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %v, want %v", c.name, c.got, c.want)
		}
	}

	// Zero-value fields
	if cfg.ProjectName != "" {
		t.Errorf("ProjectName = %q, want empty", cfg.ProjectName)
	}
	if cfg.ModulePath != "" {
		t.Errorf("ModulePath = %q, want empty", cfg.ModulePath)
	}
	if cfg.GitHubOwner != "" {
		t.Errorf("GitHubOwner = %q, want empty", cfg.GitHubOwner)
	}
}
