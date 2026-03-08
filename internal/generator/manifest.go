package generator

import "github.com/zulerne/goseed/internal/config"

// FileMapping describes a single file to generate.
type FileMapping struct {
	Source     string                               // path inside embed.FS (under "templates/")
	Target     string                               // output path (may contain template expressions)
	Condition  func(c *config.ProjectConfig) bool   // nil = always include
	IsTemplate bool                                 // true = render with text/template
}

func isLibrary(c *config.ProjectConfig) bool      { return c.ProjectType == "library" }
func isCLI(c *config.ProjectConfig) bool          { return c.ProjectType == "cli" }
func isService(c *config.ProjectConfig) bool      { return c.ProjectType == "service" }
func wantsLinter(c *config.ProjectConfig) bool    { return c.UseLinter }
func wantsTaskfile(c *config.ProjectConfig) bool  { return c.BuildTool == "taskfile" }
func wantsMakefile(c *config.ProjectConfig) bool  { return c.BuildTool == "makefile" }
func wantsGoReleaser(c *config.ProjectConfig) bool { return c.UseGoReleaser }
func wantsDocker(c *config.ProjectConfig) bool    { return c.UseDocker }
func wantsClaude(c *config.ProjectConfig) bool    { return c.UseClaude }
func wantsClaudeCI(c *config.ProjectConfig) bool  { return c.UseClaudeCI }
func wantsEnv(c *config.ProjectConfig) bool       { return c.UseEnvExample }
func wantsDependabot(c *config.ProjectConfig) bool { return c.UseDependabot }
func wantsLicense(c *config.ProjectConfig) bool   { return c.License != "none" }

// Manifest is the single source of truth for all generated files.
var Manifest = []FileMapping{
	// ── Common (always) ──────────────────────────────────
	{"common/gitignore", ".gitignore", nil, false},
	{"common/editorconfig", ".editorconfig", nil, false},
	{"common/go.mod.tmpl", "go.mod", nil, true},
	{"common/README.md.tmpl", "README.md", nil, true},
	{"common/LICENSE.tmpl", "LICENSE", wantsLicense, true},

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
	{"ci/ci.yml.tmpl", ".github/workflows/ci.yml", nil, true},
	{"ci/dependency-review.yml", ".github/workflows/dependency-review.yml", nil, false},
	{"ci/claude-code-review.yml.tmpl", ".github/workflows/claude-code-review.yml", wantsClaudeCI, true},
	{"ci/claude.yml.tmpl", ".github/workflows/claude.yml", wantsClaudeCI, true},
	{"ci/release.yml.tmpl", ".github/workflows/release.yml", wantsGoReleaser, true},

	// ── Docker ───────────────────────────────────────────
	{"docker/Dockerfile.tmpl", "Dockerfile", wantsDocker, true},

	// ── Library type ─────────────────────────────────────
	{"library/lib.go.tmpl", "{{.ProjectName}}.go", isLibrary, true},
	{"library/lib_test.go.tmpl", "{{.ProjectName}}_test.go", isLibrary, true},

	// ── CLI type ─────────────────────────────────────────
	{"cli/main.go.tmpl", "cmd/{{.ProjectName}}/main.go", isCLI, true},
	{"cli/cmd/root.go.tmpl", "internal/cli/root.go", isCLI, true},
	{"cli/cmd/version.go.tmpl", "internal/cli/version.go", isCLI, true},

	// ── Service type ─────────────────────────────────────
	{"service/main.go.tmpl", "cmd/{{.ProjectName}}/main.go", isService, true},
	{"service/internal/handler/handler.go.tmpl", "internal/handler/handler.go", isService, true},
	{"service/internal/config/config.go.tmpl", "internal/config/config.go", isService, true},
	{"service/internal/server/server.go.tmpl", "internal/server/server.go", isService, true},
}
