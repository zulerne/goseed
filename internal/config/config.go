package config

import "time"

// ProjectConfig holds all user choices for project generation.
type ProjectConfig struct {
	// Core (always asked)
	ProjectName string // directory name and binary name
	ModulePath  string // e.g. "github.com/user/myapp"
	GoVersion   string // e.g. "1.26"
	License     string // "MIT" | "Apache-2.0" | "none"

	// Tooling
	UseLinter     bool   // .golangci.yml
	BuildTool     string // "taskfile" | "makefile" | "none"
	UseGoReleaser bool   // .goreleaser.yaml
	UseDocker     bool   // Dockerfile + .dockerignore
	UseEnvExample bool   // .env.example
	UseDependabot bool   // .github/dependabot.yml
	UseCI         bool   // ci.yml + dependency-review.yml

	// Claude Code
	UseClaude   bool // CLAUDE.md + .claude/rules/
	UseClaudeCI bool // claude-code-review.yml + claude.yml

	// Derived
	GitHubOwner string // inferred from ModulePath

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
		UseGoReleaser: false,
		UseDocker:     false,
		UseEnvExample: true,
		UseDependabot: false,
		UseCI:         true,
		UseClaude:     false,
		UseClaudeCI:   false,
		Year:          time.Now().Year(),
	}
}
