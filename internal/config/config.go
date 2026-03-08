package config

import "time"

// ProjectConfig holds all user choices for project generation.
type ProjectConfig struct {
	// Core (always asked)
	ProjectName string // directory name and binary name
	ModulePath  string // e.g. "github.com/user/myapp"
	ProjectType string // "library" | "cli" | "service"
	Description string // one-liner for README
	GoVersion   string // e.g. "1.26"
	License     string // "MIT" | "Apache-2.0" | "none"

	// Tooling
	UseLinter     bool   // .golangci.yml
	BuildTool     string // "taskfile" | "makefile" | "none"
	UseGoReleaser bool   // .goreleaser.yaml
	UseDocker     bool   // Dockerfile + docker-compose
	UseEnvExample bool   // .env.example
	UseRenovate   bool   // renovate.json

	// Claude Code
	UseClaude   bool // CLAUDE.md + .claude/rules/
	UseClaudeCI bool // claude-code-review.yml + claude.yml

	// Derived
	GitHubOwner string // inferred from ModulePath

	// Service-specific
	HTTPFramework string // "stdlib" | "chi"

	// Computed at generation time
	Year int
}

// Defaults returns a ProjectConfig with sensible default values.
func Defaults() ProjectConfig {
	return ProjectConfig{
		GoVersion:     "1.26",
		License:       "MIT",
		UseLinter:     true,
		BuildTool:     "taskfile",
		UseGoReleaser: true,
		UseDocker:     false,
		UseEnvExample: false,
		UseRenovate:   true,
		UseClaude:     true,
		UseClaudeCI:   true,
		HTTPFramework: "stdlib",
		Year:          time.Now().Year(),
	}
}
