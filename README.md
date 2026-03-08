# go-template

Opinionated Go project template with CI, linting, and Claude Code integration — ready to clone and build.

## What's Included

```
.
├── .claude/rules/go.md          # Go coding rules for Claude Code
├── .editorconfig                # Editor formatting (tabs for Go, spaces for YAML)
├── .env.example                 # Environment variables skeleton
├── .github/workflows/
│   ├── ci.yml                   # Test + lint + govulncheck
│   ├── claude-code-review.yml   # Automated PR review via Claude Code
│   ├── claude.yml               # @claude mentions in issues/PRs
│   └── dependency-review.yml    # Block PRs with vulnerable deps
├── .gitignore
├── .golangci.yml                # golangci-lint v2 config
├── CLAUDE.md                    # Project instructions for Claude Code
├── Makefile                     # Build automation (traditional)
├── Taskfile.yml                 # Build automation (modern alternative)
├── cmd/app/main.go              # Entry point
└── internal/                    # Application packages
```

## Usage

1. **Clone or use as GitHub template**

2. **Find and replace** these placeholders:
   - `github.com/zulerne/go-template` → your module path (in `go.mod`)
   - `app` → your binary name (in `Makefile` / `Taskfile.yml`)
   - `./cmd/app` → your main package path
   - `zulerne` → your GitHub username (in `claude-code-review.yml`, `claude.yml`)

3. **Pick your build tool** — keep `Makefile` or `Taskfile.yml`, delete the other

4. **Enable project-specific linters** in `.golangci.yml`:
   - HTTP/API → uncomment `bodyclose`, `noctx`
   - Database → uncomment `sqlclosecheck`, `rowserrcheck`
   - TUI (Bubble Tea) → uncomment `hugeParam` disable
   - slog → uncomment `sloglint`

5. **Set up repository secrets:**
   - `CODECOV_TOKEN` — for coverage uploads
   - `CLAUDE_CODE_OAUTH_TOKEN` — for Claude Code workflows

## CI Pipeline

| Job | What it does |
|---|---|
| **test** | `go vet` → `go test -race -shuffle=on` → upload coverage to Codecov |
| **lint** | golangci-lint with 17+ linters |
| **vulncheck** | govulncheck — symbol-level vulnerability scanning |
| **dependency-review** | Blocks PRs introducing known-vulnerable dependencies |

## Linters

Enabled by default — the full list in `.golangci.yml`:

**Core:** errcheck, govet, ineffassign, staticcheck, unused

**Quality:** revive, gosec, gocritic, prealloc, unconvert, copyloopvar, intrange, modernize, nolintlint, perfsprint

**Error handling:** wrapcheck, errorlint

Commented out (enable per project): bodyclose, noctx, sqlclosecheck, rowserrcheck, sloglint

## Claude Code Rules

`.claude/rules/go.md` contains Go-specific rules covering:
- Error handling patterns
- Concurrency guidelines
- Modern Go idioms (organized by Go version, up to 1.26)
- Pre-commit checklist

> **Tip:** If you work on multiple Go projects, move this file to `~/.claude/rules/go.md` to apply it globally instead of duplicating per project.

## Build Commands

Both `Makefile` and `Taskfile.yml` provide the same targets:

| Command | Description |
|---|---|
| `run` | Run the application |
| `build` | Build binary |
| `test` | Run tests with race detector |
| `test-cover` / `test:cover` | Run tests with coverage report |
| `lint` | Run golangci-lint |
| `check` | Run lint + test |
| `fmt` | Format code |
| `clean` | Remove build artifacts |

## License

[MIT](LICENSE)
