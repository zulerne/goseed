package prompt

import (
	"errors"
	"fmt"
	"os/exec"
	"os/user"
	"regexp"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"

	"github.com/zulerne/goseed/internal/config"
)

var validName = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// Run displays the interactive TUI form and returns the filled config.
func Run(cfg *config.ProjectConfig) error {
	username := guessGitHubUser()

	var projectType string
	var license string
	var buildTool string
	var httpFramework string
	var features []string
	var ciFeatures []string
	var claudeFeatures []string

	km := huh.NewDefaultKeyMap()
	km.MultiSelect.Toggle = key.NewBinding(key.WithKeys(" ", "x"), key.WithHelp("space", "toggle"))

	form := huh.NewForm(
		// Group 1: Project identity
		huh.NewGroup(
			huh.NewInput().
				Title("Project name").
				Description("Directory name and binary name").
				Value(&cfg.ProjectName).
				Validate(func(s string) error {
					if !validName.MatchString(s) {
						return errors.New("must match [a-z][a-z0-9-]*")
					}
					return nil
				}),

			huh.NewInput().
				Title("Module path").
				Description("Go module path (leave empty for default)").
				PlaceholderFunc(func() string {
					name := cfg.ProjectName
					if name == "" {
						name = "myapp"
					}
					return fmt.Sprintf("github.com/%s/%s", username, name)
				}, &cfg.ProjectName).
				Value(&cfg.ModulePath).
				Validate(func(s string) error {
					if s == "" {
						return nil
					}
					if !strings.Contains(s, "/") {
						return errors.New("must contain at least one /")
					}
					return nil
				}),
		),

		// Group 2: Project type and license
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Project type").
				Options(
					huh.NewOption("Library", "library"),
					huh.NewOption("CLI application", "cli"),
					huh.NewOption("HTTP service", "service"),
				).
				Value(&projectType),

			huh.NewSelect[string]().
				Title("License").
				Options(
					huh.NewOption("MIT", "MIT"),
					huh.NewOption("Apache 2.0", "Apache-2.0"),
					huh.NewOption("None", "none"),
				).
				Value(&license),
		),

		// Group 3: Build tool
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Build tool").
				Options(
					huh.NewOption("Taskfile", "taskfile"),
					huh.NewOption("Makefile", "makefile"),
					huh.NewOption("None", "none"),
				).
				Value(&buildTool),
		),

		// Group 4: Optional features
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Optional features").
				Options(
					huh.NewOption("golangci-lint", "linter").Selected(cfg.UseLinter),
					huh.NewOption("GoReleaser", "goreleaser").Selected(cfg.UseGoReleaser),
					huh.NewOption(".env.example", "env").Selected(cfg.UseEnvExample),
					huh.NewOption("Dockerfile", "docker").Selected(cfg.UseDocker),
				).
				Value(&features).
				Filterable(false),
		),

		// Group 5: CI/CD
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("CI/CD").
				Options(
					huh.NewOption("CI workflows (lint, test, vulncheck)", "ci").Selected(cfg.UseCI),
					huh.NewOption("Dependabot (gomod + actions)", "dependabot").Selected(cfg.UseDependabot),
				).
				Value(&ciFeatures).
				Filterable(false),
		),

		// Group 6: Claude Code
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Claude Code integration").
				Options(
					huh.NewOption("CLAUDE.md + .claude/rules", "claude").Selected(cfg.UseClaude),
					huh.NewOption("CI workflows (review + agent)", "claude-ci").Selected(cfg.UseClaudeCI),
				).
				Value(&claudeFeatures).
				Filterable(false),
		),

		// Group 7: Service-specific
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("HTTP framework").
				Options(
					huh.NewOption("stdlib (net/http)", "stdlib"),
					huh.NewOption("chi", "chi"),
				).
				Value(&httpFramework),
		).WithHideFunc(func() bool {
			return projectType != "service"
		}),
	).WithKeyMap(km).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil {
		return fmt.Errorf("prompt: %w", err)
	}

	cfg.ProjectType = projectType
	cfg.License = license
	cfg.BuildTool = buildTool
	cfg.HTTPFramework = httpFramework
	cfg.UseLinter = slices.Contains(features, "linter")
	cfg.UseGoReleaser = slices.Contains(features, "goreleaser")
	cfg.UseDocker = slices.Contains(features, "docker")
	cfg.UseEnvExample = slices.Contains(features, "env")
	cfg.UseCI = slices.Contains(ciFeatures, "ci")
	cfg.UseDependabot = slices.Contains(ciFeatures, "dependabot")
	cfg.UseClaude = slices.Contains(claudeFeatures, "claude")
	cfg.UseClaudeCI = slices.Contains(claudeFeatures, "claude-ci")

	// Infer GitHubOwner from ModulePath
	cfg.GitHubOwner = inferOwner(cfg.ModulePath)

	// Set defaults based on project type
	if cfg.ModulePath == "" {
		cfg.ModulePath = fmt.Sprintf("github.com/%s/%s", username, cfg.ProjectName)
		cfg.GitHubOwner = username
	}

	return nil
}

func inferOwner(modulePath string) string {
	parts := strings.Split(modulePath, "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

func guessGitHubUser() string {
	// 1. git config github.user (fast, local)
	if out, err := exec.Command("git", "config", "github.user").Output(); err == nil {
		if name := strings.TrimSpace(string(out)); name != "" {
			return name
		}
	}

	// 2. gh CLI (requires auth, may do network call)
	if out, err := exec.Command("gh", "api", "user", "-q", ".login").Output(); err == nil {
		if name := strings.TrimSpace(string(out)); name != "" {
			return name
		}
	}

	// 3. System username as fallback
	u, err := user.Current()
	if err != nil {
		return "user"
	}
	return u.Username
}
