package prompt

import (
	"errors"
	"fmt"
	"os/user"
	"regexp"
	"strings"

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

	form := huh.NewForm(
		// Group 1: Project basics
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
				Description("Go module path").
				Placeholder(fmt.Sprintf("github.com/%s/myapp", username)).
				Value(&cfg.ModulePath).
				Validate(func(s string) error {
					if !strings.Contains(s, "/") {
						return errors.New("must contain at least one /")
					}
					return nil
				}),

			huh.NewSelect[string]().
				Title("Project type").
				Options(
					huh.NewOption("Library", "library"),
					huh.NewOption("CLI application", "cli"),
					huh.NewOption("HTTP service", "service"),
				).
				Value(&projectType),

			huh.NewInput().
				Title("Description").
				Description("One-line project description (optional)").
				Value(&cfg.Description),

			huh.NewSelect[string]().
				Title("License").
				Options(
					huh.NewOption("MIT", "MIT"),
					huh.NewOption("Apache 2.0", "Apache-2.0"),
					huh.NewOption("None", "none"),
				).
				Value(&license),
		),

		// Group 2: Build tooling
		huh.NewGroup(
			huh.NewConfirm().
				Title("Include golangci-lint config?").
				Value(&cfg.UseLinter),

			huh.NewSelect[string]().
				Title("Build tool").
				Options(
					huh.NewOption("Taskfile", "taskfile"),
					huh.NewOption("Makefile", "makefile"),
					huh.NewOption("None", "none"),
				).
				Value(&buildTool),

			huh.NewConfirm().
				Title("Include GoReleaser?").
				Value(&cfg.UseGoReleaser),
		),

		// Group 3: Extras
		huh.NewGroup(
			huh.NewConfirm().
				Title("Include Dockerfile?").
				Value(&cfg.UseDocker),

			huh.NewConfirm().
				Title("Include .env.example?").
				Value(&cfg.UseEnvExample),

			huh.NewConfirm().
				Title("Include Renovate config?").
				Value(&cfg.UseRenovate),
		),

		// Group 4: Claude Code
		huh.NewGroup(
			huh.NewConfirm().
				Title("Include CLAUDE.md + .claude/rules?").
				Value(&cfg.UseClaude),

			huh.NewConfirm().
				Title("Include Claude CI workflows?").
				Value(&cfg.UseClaudeCI),
		),

		// Group 5: Service-specific
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
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil {
		return fmt.Errorf("prompt: %w", err)
	}

	cfg.ProjectType = projectType
	cfg.License = license
	cfg.BuildTool = buildTool
	cfg.HTTPFramework = httpFramework

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
	u, err := user.Current()
	if err != nil {
		return "user"
	}
	return u.Username
}
