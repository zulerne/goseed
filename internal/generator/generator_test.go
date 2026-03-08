package generator_test

import (
	"os"
	"path/filepath"
	"testing"

	goseed "github.com/zulerne/goseed"
	"github.com/zulerne/goseed/internal/config"
	"github.com/zulerne/goseed/internal/generator"
)

func TestGenerateLibrary(t *testing.T) {
	dir := t.TempDir()
	cfg := config.ProjectConfig{
		ProjectName:   "testlib",
		ModulePath:    "github.com/test/testlib",
		ProjectType:   "library",
		GoVersion:     "1.26",
		License:       "MIT",
		UseLinter:     true,
		BuildTool:     "taskfile",
		UseGoReleaser: false,
		UseDocker:     false,
		UseClaude:     true,
		UseClaudeCI:   false,
		UseDependabot:   false,
		GitHubOwner:   "test",
		Year:          2026,
	}

	err := generator.Generate(&cfg, goseed.Templates, dir)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	projectDir := filepath.Join(dir, "testlib")
	mustExist(t, projectDir, "testlib.go")
	mustExist(t, projectDir, "testlib_test.go")
	mustExist(t, projectDir, "go.mod")
	mustExist(t, projectDir, ".gitignore")
	mustExist(t, projectDir, ".golangci.yml")
	mustExist(t, projectDir, "Taskfile.yml")
	mustExist(t, projectDir, "LICENSE")
	mustExist(t, projectDir, "CLAUDE.md")
	mustExist(t, projectDir, ".claude/rules/go.md")
	mustNotExist(t, projectDir, "Dockerfile")
	mustNotExist(t, projectDir, "cmd")
}

func TestGenerateCLI(t *testing.T) {
	dir := t.TempDir()
	cfg := config.ProjectConfig{
		ProjectName:   "mycli",
		ModulePath:    "github.com/test/mycli",
		ProjectType:   "cli",
		GoVersion:     "1.26",
		License:       "none",
		UseLinter:     true,
		BuildTool:     "makefile",
		UseGoReleaser: true,
		UseDocker:     false,
		UseClaude:     false,
		UseClaudeCI:   false,
		UseDependabot:   true,
		GitHubOwner:   "test",
		Year:          2026,
	}

	err := generator.Generate(&cfg, goseed.Templates, dir)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	projectDir := filepath.Join(dir, "mycli")
	mustExist(t, projectDir, "cmd/mycli/main.go")
	mustExist(t, projectDir, "internal/cli/root.go")
	mustExist(t, projectDir, "internal/cli/version.go")
	mustExist(t, projectDir, "Makefile")
	mustExist(t, projectDir, ".goreleaser.yaml")
	mustExist(t, projectDir, ".github/dependabot.yml")
	mustNotExist(t, projectDir, "LICENSE")
	mustNotExist(t, projectDir, "Taskfile.yml")
	mustNotExist(t, projectDir, "Dockerfile")
	mustNotExist(t, projectDir, "CLAUDE.md")
}

func TestGenerateService(t *testing.T) {
	dir := t.TempDir()
	cfg := config.ProjectConfig{
		ProjectName:   "mysvc",
		ModulePath:    "github.com/test/mysvc",
		ProjectType:   "service",
		GoVersion:     "1.26",
		License:       "Apache-2.0",
		UseLinter:     true,
		BuildTool:     "taskfile",
		UseGoReleaser: true,
		UseDocker:     true,
		UseEnvExample: true,
		UseClaude:     true,
		UseClaudeCI:   true,
		UseDependabot:   true,
		GitHubOwner:   "test",
		HTTPFramework: "stdlib",
		Year:          2026,
	}

	err := generator.Generate(&cfg, goseed.Templates, dir)
	if err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	projectDir := filepath.Join(dir, "mysvc")
	mustExist(t, projectDir, "cmd/mysvc/main.go")
	mustExist(t, projectDir, "internal/handler/handler.go")
	mustExist(t, projectDir, "internal/config/config.go")
	mustExist(t, projectDir, "internal/server/server.go")
	mustExist(t, projectDir, "Dockerfile")
	mustNotExist(t, projectDir, "docker-compose.yml")
	mustExist(t, projectDir, ".env.example")
	mustExist(t, projectDir, ".goreleaser.yaml")
	mustExist(t, projectDir, ".github/workflows/claude-code-review.yml")
	mustExist(t, projectDir, ".github/workflows/claude.yml")
	mustExist(t, projectDir, ".github/workflows/release.yml")
	mustExist(t, projectDir, "LICENSE")
}

func TestManifestConditions(t *testing.T) {
	lib := &config.ProjectConfig{ProjectType: "library", UseLinter: true, BuildTool: "taskfile"}
	cli := &config.ProjectConfig{ProjectType: "cli", UseGoReleaser: true}
	svc := &config.ProjectConfig{ProjectType: "service", UseDocker: true}

	for _, fm := range generator.Manifest {
		// Library should not have Dockerfile or cmd/
		if fm.Source == "docker/Dockerfile.tmpl" && fm.Condition != nil && fm.Condition(lib) {
			t.Error("library should not include Dockerfile")
		}
		if fm.Source == "service/main.go.tmpl" && fm.Condition != nil && fm.Condition(lib) {
			t.Error("library should not include service main.go")
		}

		// CLI should have goreleaser when enabled
		if fm.Source == "tooling/goreleaser.yaml.tmpl" && fm.Condition != nil && !fm.Condition(cli) {
			t.Error("CLI with UseGoReleaser should include goreleaser")
		}

		// Service should have Dockerfile when enabled
		if fm.Source == "docker/Dockerfile.tmpl" && fm.Condition != nil && !fm.Condition(svc) {
			t.Error("service with UseDocker should include Dockerfile")
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
