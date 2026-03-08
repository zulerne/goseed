# gostart

Interactive CLI tool that scaffolds Go projects (library, CLI, or service) with CI, linting, and Claude Code integration.

## Project structure

- `cmd/gostart/main.go` — CLI entry point (cobra)
- `internal/config/` — ProjectConfig struct and defaults
- `internal/prompt/` — TUI form (charmbracelet/huh)
- `internal/generator/` — file generation from templates
- `templates/` — embedded template files
- `embed.go` — go:embed directive

## Code Quality

Always run `task lint` before committing. Run `task check` (lint + test) for full validation.
