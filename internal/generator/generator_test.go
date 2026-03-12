package generator_test

import (
	"os"
	"path/filepath"
	"testing"

	goseed "github.com/zulerne/goseed"
	"github.com/zulerne/goseed/internal/config"
	"github.com/zulerne/goseed/internal/generator"
)

func TestGenerateMinimal(t *testing.T) {
	dir := t.TempDir()
	cfg := config.ProjectConfig{
		ProjectName:   "myapp",
		ModulePath:    "github.com/test/myapp",
		GoVersion:     "1.26",
		License:       "MIT",
		UseLinter:     false,
		BuildTool:     "none",
		UseGoReleaser: false,
		UseDocker:     false,
		UseClaude:     false,
		UseClaudeCI:   false,
		UseDependabot: false,
		GitHubOwner:   "test",
		Year:          2026,
	}

	err := generator.Generate(&cfg, goseed.Templates, dir)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	projectDir := filepath.Join(dir, "myapp")
	mustExist(t, projectDir, "cmd/myapp/main.go")
	mustExist(t, projectDir, "internal/.gitkeep")
	mustExist(t, projectDir, "go.mod")
	mustExist(t, projectDir, ".gitignore")
	mustExist(t, projectDir, ".editorconfig")
	mustExist(t, projectDir, "README.md")
	mustExist(t, projectDir, "LICENSE")
	mustExist(t, projectDir, ".github/ISSUE_TEMPLATE/bug_report.yml")
	mustExist(t, projectDir, ".github/ISSUE_TEMPLATE/feature_request.yml")
	mustExist(t, projectDir, ".github/ISSUE_TEMPLATE/config.yml")
	mustExist(t, projectDir, ".github/pull_request_template.md")
	mustNotExist(t, projectDir, ".golangci.yml")
	mustNotExist(t, projectDir, "Taskfile.yml")
	mustNotExist(t, projectDir, "Makefile")
	mustNotExist(t, projectDir, "Dockerfile")
	mustNotExist(t, projectDir, ".dockerignore")
	mustNotExist(t, projectDir, "CLAUDE.md")
}

func TestGenerateFullFeatured(t *testing.T) {
	dir := t.TempDir()
	cfg := config.ProjectConfig{
		ProjectName:   "fullapp",
		ModulePath:    "github.com/test/fullapp",
		GoVersion:     "1.26",
		License:       "Apache-2.0",
		UseLinter:     true,
		BuildTool:     "taskfile",
		UseGoReleaser: true,
		UseDocker:     true,
		UseEnvExample: true,
		UseCI:         true,
		UseClaude:     true,
		UseClaudeCI:   true,
		UseDependabot: true,
		GitHubOwner:   "test",
		Year:          2026,
	}

	err := generator.Generate(&cfg, goseed.Templates, dir)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	projectDir := filepath.Join(dir, "fullapp")
	mustExist(t, projectDir, "cmd/fullapp/main.go")
	mustExist(t, projectDir, "internal/.gitkeep")
	mustExist(t, projectDir, "go.mod")
	mustExist(t, projectDir, ".golangci.yml")
	mustExist(t, projectDir, "Taskfile.yml")
	mustExist(t, projectDir, ".goreleaser.yaml")
	mustExist(t, projectDir, "Dockerfile")
	mustExist(t, projectDir, ".dockerignore")
	mustExist(t, projectDir, ".env.example")
	mustExist(t, projectDir, "CLAUDE.md")
	mustExist(t, projectDir, ".claude/rules/go.md")
	mustExist(t, projectDir, ".github/workflows/ci.yml")
	mustExist(t, projectDir, ".github/workflows/dependency-review.yml")
	mustExist(t, projectDir, ".github/workflows/release.yml")
	mustExist(t, projectDir, ".github/workflows/claude-code-review.yml")
	mustExist(t, projectDir, ".github/workflows/claude.yml")
	mustExist(t, projectDir, ".github/dependabot.yml")
	mustExist(t, projectDir, ".github/ISSUE_TEMPLATE/bug_report.yml")
	mustExist(t, projectDir, ".github/pull_request_template.md")
	mustExist(t, projectDir, "LICENSE")
}

func TestGenerateMakefile(t *testing.T) {
	dir := t.TempDir()
	cfg := config.ProjectConfig{
		ProjectName:   "mkapp",
		ModulePath:    "github.com/test/mkapp",
		GoVersion:     "1.26",
		License:       "none",
		UseLinter:     true,
		BuildTool:     "makefile",
		UseGoReleaser: false,
		UseDocker:     false,
		UseClaude:     false,
		UseClaudeCI:   false,
		UseDependabot: false,
		GitHubOwner:   "test",
		Year:          2026,
	}

	err := generator.Generate(&cfg, goseed.Templates, dir)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	projectDir := filepath.Join(dir, "mkapp")
	mustExist(t, projectDir, "Makefile")
	mustExist(t, projectDir, ".golangci.yml")
	mustNotExist(t, projectDir, "Taskfile.yml")
	mustNotExist(t, projectDir, "LICENSE")
	mustNotExist(t, projectDir, ".goreleaser.yaml")
}

func TestManifestConditions(t *testing.T) {
	minimal := &config.ProjectConfig{UseLinter: false, BuildTool: "none"}
	full := &config.ProjectConfig{UseLinter: true, BuildTool: "taskfile", UseGoReleaser: true, UseDocker: true}

	for _, fm := range generator.Manifest {
		// Minimal should not have linter config
		if fm.Source == "tooling/golangci.yml.tmpl" && fm.Condition != nil && fm.Condition(minimal) {
			t.Error("minimal config should not include golangci.yml")
		}

		// Full should have goreleaser when enabled
		if fm.Source == "tooling/goreleaser.yaml.tmpl" && fm.Condition != nil && !fm.Condition(full) {
			t.Error("full config with UseGoReleaser should include goreleaser")
		}

		// Full should have Dockerfile when enabled
		if fm.Source == "docker/Dockerfile.tmpl" && fm.Condition != nil && !fm.Condition(full) {
			t.Error("full config with UseDocker should include Dockerfile")
		}

		// Full should have dockerignore when enabled
		if fm.Source == "docker/dockerignore" && fm.Condition != nil && !fm.Condition(full) {
			t.Error("full config with UseDocker should include .dockerignore")
		}
	}
}

func mustExist(t *testing.T, dir, path string) {
	t.Helper()
	full := filepath.Join(dir, path)
	if _, err := os.Stat(full); err != nil {
		t.Errorf("expected file %s to exist", path)
	}
}

func mustNotExist(t *testing.T, dir, path string) {
	t.Helper()
	full := filepath.Join(dir, path)
	if _, err := os.Stat(full); err == nil {
		t.Errorf("expected file %s to NOT exist", path)
	}
}
