package generator

import "github.com/zulerne/goseed/internal/config"

// FileMapping describes a single file to generate.
type FileMapping struct {
	Source     string                             // path inside embed.FS (under "templates/")
	Target     string                             // output path (may contain template expressions)
	Condition  func(c *config.ProjectConfig) bool // nil = always include
	IsTemplate bool                               // true = render with text/template
}

func wantsLinter(c *config.ProjectConfig) bool     { return c.UseLinter }
func wantsTaskfile(c *config.ProjectConfig) bool    { return c.BuildTool == "taskfile" }
func wantsMakefile(c *config.ProjectConfig) bool    { return c.BuildTool == "makefile" }
func wantsGoReleaser(c *config.ProjectConfig) bool  { return c.UseGoReleaser }
func wantsDocker(c *config.ProjectConfig) bool      { return c.UseDocker }
func wantsClaude(c *config.ProjectConfig) bool      { return c.UseClaude }
func wantsClaudeCI(c *config.ProjectConfig) bool    { return c.UseClaudeCI }
func wantsEnv(c *config.ProjectConfig) bool         { return c.UseEnvExample }
func wantsDependabot(c *config.ProjectConfig) bool  { return c.UseDependabot }
func wantsCI(c *config.ProjectConfig) bool          { return c.UseCI }
func wantsLicense(c *config.ProjectConfig) bool     { return c.License != "none" }

// Manifest is the single source of truth for all generated files.
var Manifest = []FileMapping{
	// ── Common (always) ──────────────────────────────────
	{"common/gitignore", ".gitignore", nil, false},
	{"common/editorconfig", ".editorconfig", nil, false},
	{"common/go.mod.tmpl", "go.mod", nil, true},
	{"common/README.md.tmpl", "README.md", nil, true},
	{"common/LICENSE.tmpl", "LICENSE", wantsLicense, true},
	{"common/main.go.tmpl", "cmd/{{.ProjectName}}/main.go", nil, true},
	{"common/gitkeep", "internal/.gitkeep", nil, false},

	// ── GitHub templates (always) ────────────────────────
	{"github/ISSUE_TEMPLATE/bug_report.yml.tmpl", ".github/ISSUE_TEMPLATE/bug_report.yml", nil, true},
	{"github/ISSUE_TEMPLATE/feature_request.yml", ".github/ISSUE_TEMPLATE/feature_request.yml", nil, false},
	{"github/ISSUE_TEMPLATE/config.yml", ".github/ISSUE_TEMPLATE/config.yml", nil, false},
	{"github/pull_request_template.md", ".github/pull_request_template.md", nil, false},

	// ── Tooling (conditional) ────────────────────────────
	{"tooling/golangci.yml.tmpl", ".golangci.yml", wantsLinter, true},
	{"tooling/Taskfile.yml.tmpl", "Taskfile.yml", wantsTaskfile, true},
	{"tooling/Makefile.tmpl", "Makefile", wantsMakefile, true},
	{"tooling/goreleaser.yaml.tmpl", ".goreleaser.yaml", wantsGoReleaser, true},
	{"tooling/env.example", ".env.example", wantsEnv, false},
	{"tooling/dependabot.yml", ".github/dependabot.yml", wantsDependabot, false},

	// ── Claude Code ──────────────────────────────────────
	{"claude/CLAUDE.md.tmpl", "CLAUDE.md", wantsClaude, true},
	{"claude/rules/go.md", ".claude/rules/go.md", wantsClaude, false},

	// ── CI workflows ─────────────────────────────────────
	{"ci/ci.yml.tmpl", ".github/workflows/ci.yml", wantsCI, true},
	{"ci/dependency-review.yml", ".github/workflows/dependency-review.yml", wantsCI, false},
	{"ci/claude-code-review.yml.tmpl", ".github/workflows/claude-code-review.yml", wantsClaudeCI, true},
	{"ci/claude.yml.tmpl", ".github/workflows/claude.yml", wantsClaudeCI, true},
	{"ci/release.yml.tmpl", ".github/workflows/release.yml", wantsGoReleaser, true},

	// ── Docker ───────────────────────────────────────────
	{"docker/Dockerfile.tmpl", "Dockerfile", wantsDocker, true},
	{"docker/dockerignore", ".dockerignore", wantsDocker, false},
}
