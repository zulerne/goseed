# Contributing to goseed

## Development

### Prerequisites

- Go 1.26+
- [Task](https://taskfile.dev) (optional, for build automation)
- [golangci-lint](https://golangci-lint.run) (for linting)

### Setup

```bash
git clone https://github.com/zulerne/goseed.git
cd goseed
go mod download
```

### Running

```bash
task run
# or
go run ./cmd/goseed
```

### Testing

```bash
task test
# or
go test -race ./...
```

### Linting

```bash
task lint
# or
golangci-lint run
```

### Pre-commit checklist

Run `task check` (lint + test) before opening a PR.

## Pull Requests

- Fork the repo and create your branch from `main`.
- Follow [Conventional Commits](https://www.conventionalcommits.org): `feat:`, `fix:`, `docs:`, etc.
- Ensure `task check` passes.
- Keep PRs focused — one feature or fix per PR.
