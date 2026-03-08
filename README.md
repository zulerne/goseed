# gostart

Interactive CLI tool that scaffolds Go projects with best practices — CI, linting, Claude Code integration, and more.

## Install

```bash
brew install zulerne/tap/gostart
```

Or build from source:

```bash
go install github.com/zulerne/goseed/cmd/gostart@latest
```

## Usage

### Interactive mode

```bash
gostart
```

Walks you through 4 groups of questions (project basics, tooling, Claude Code, service-specific) and generates a ready-to-build project.

### Non-interactive mode

```bash
gostart --name myapp --module github.com/user/myapp --type service --no-interactive
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `--name` | | Project name |
| `--module` | | Go module path |
| `--type` | | `library`, `cli`, or `service` |
| `--description` | | One-line project description |
| `--go-version` | `1.26` | Go version |
| `--license` | `MIT` | `MIT`, `Apache-2.0`, or `none` |
| `--build-tool` | `taskfile` | `taskfile`, `makefile`, or `none` |
| `--http-framework` | `stdlib` | `stdlib` or `chi` (service only) |
| `--linter` | `true` | Include golangci-lint config |
| `--goreleaser` | `true` | Include GoReleaser |
| `--docker` | `false` | Include Dockerfile |
| `--env-example` | `false` | Include .env.example |
| `--renovate` | `true` | Include Renovate config |
| `--claude` | `true` | Include Claude Code files |
| `--claude-ci` | `true` | Include Claude CI workflows |
| `--no-interactive` | `false` | Skip TUI, use flags + defaults |
| `--output-dir` | `.` | Output directory |

## Project Types

### Library

Generates a Go package with exported functions and table-driven tests.

### CLI

Generates a cobra-based CLI application with version subcommand and GoReleaser config.

### Service

Generates an HTTP service with graceful shutdown, health endpoint, config from environment, and optional Docker support.

## What's Generated

Every project includes:
- `.gitignore`, `.editorconfig`
- `go.mod`, `README.md`
- CI workflow (test + lint + govulncheck)
- Dependency review workflow

Optional (based on choices):
- `.golangci.yml` — 17+ linters
- `Taskfile.yml` or `Makefile`
- `.goreleaser.yaml` + release workflow
- `Dockerfile` + `docker-compose.yml`
- `CLAUDE.md` + `.claude/rules/go.md`
- Claude Code CI workflows
- `renovate.json`
- `LICENSE` (MIT or Apache 2.0)

## License

[MIT](LICENSE)
