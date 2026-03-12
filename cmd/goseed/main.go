package main

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"

	goseed "github.com/zulerne/goseed"
	"github.com/zulerne/goseed/internal/config"
	"github.com/zulerne/goseed/internal/generator"
	"github.com/zulerne/goseed/internal/prompt"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	cfg := config.Defaults()
	var noInteractive bool
	var outputDir string

	cmd := &cobra.Command{
		Use:     "goseed",
		Short:   "Scaffold Go projects with a clean foundation",
		Long: `Interactive CLI tool that scaffolds Go projects with CI, linting, and Claude Code integration.

Run without flags to start the interactive TUI.
Use --no-interactive with --name and --module to scaffold non-interactively.`,
		Version: version + " (commit: " + commit + ")",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !noInteractive {
				if err := prompt.Run(&cfg); err != nil {
					return err
				}
			} else {
				if cfg.ProjectName == "" {
					return errors.New("--name is required in non-interactive mode")
				}
				if cfg.ModulePath == "" {
					return errors.New("--module is required in non-interactive mode")
				}
				cfg.GitHubOwner = inferOwner(cfg.ModulePath)
			}

			return generator.Generate(&cfg, goseed.Templates, outputDir)
		},
	}

	f := cmd.Flags()
	f.StringVar(&cfg.ProjectName, "name", "", "Project name")
	f.StringVar(&cfg.ModulePath, "module", "", "Go module path")
	f.StringVar(&cfg.GoVersion, "go-version", cfg.GoVersion, "Go version")
	f.StringVar(&cfg.License, "license", cfg.License, "License: MIT, Apache-2.0, none")
	f.StringVar(&cfg.BuildTool, "build-tool", cfg.BuildTool, "Build tool: taskfile, makefile, none")
	f.BoolVar(&cfg.UseLinter, "linter", cfg.UseLinter, "Include golangci-lint config")
	f.BoolVar(&cfg.UseGoReleaser, "goreleaser", cfg.UseGoReleaser, "Include GoReleaser")
	f.BoolVar(&cfg.UseCI, "ci", cfg.UseCI, "Include CI workflows")
	f.BoolVar(&cfg.UseDocker, "docker", cfg.UseDocker, "Include Dockerfile")
	f.BoolVar(&cfg.UseEnvExample, "env-example", cfg.UseEnvExample, "Include .env.example")
	f.BoolVar(&cfg.UseDependabot, "dependabot", cfg.UseDependabot, "Include Dependabot config")
	f.BoolVar(&cfg.UseClaude, "claude", cfg.UseClaude, "Include Claude Code files")
	f.BoolVar(&cfg.UseClaudeCI, "claude-ci", cfg.UseClaudeCI, "Include Claude CI workflows")
	f.BoolVar(&noInteractive, "no-interactive", false, "Use flags and defaults only")
	f.StringVar(&outputDir, "output-dir", ".", "Output directory")

	return cmd
}

func inferOwner(modulePath string) string {
	parts := strings.Split(modulePath, "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}
