# Repository Guidelines

## Project Structure & Module Organization

`list-github-stars` is a Go CLI that lists a GitHub user's starred
repositories as a Markdown document grouped by language. The repository is
a single Go module (`github.com/caian-org/list-github-stars`); the binary
entrypoint is `cmd/lgs/main.go` and all domain logic lives under
`internal/`.

- `cmd/lgs/main.go` — minimal entrypoint that delegates to the CLI layer.
- `cmd/lgs/cli/` — `urfave/cli/v3` root command and the
  `ProgramVersion` / `ProgramCommitSHA` / `ProgramBuildTime` symbols
  injected by goreleaser ldflags.
- `internal/github/` — paginated starred-repos fetcher built on
  `go-github/v76`. Exposes `SetBaseURL` for `httptest.Server`-based tests.
- `internal/render/` — pure Markdown renderer with language grouping,
  slug disambiguation, and case-insensitive sorting.
- `internal/starlist/` — orchestration: fetch → render → write.
- `testdata/golden/` — golden Markdown fixtures regenerated with the
  `-update` test flag.

`bin/`, `dist/`, `coverage.out`, and `.devbox/` are generated artefacts —
never edit by hand.

## Build, Test, and Development Commands

Standard Go toolchain via devbox + just. Common targets:

- `just build lgs` — builds `bin/lgs` with dev-time ldflags (commit +
  timestamp).
- `just run lgs` — build, then run with the current environment.
- `just test` — runs `go test ./...`.
- `just test-race` — runs the suite with `-race`.
- `just cover` — coverage profile plus per-function totals.
- `just lint` — currently `go vet ./...`. CI additionally runs
  `golangci-lint` per `.golangci.yml`.
- `just tidy` — `go mod tidy`.
- `just update-goldens` — regenerate `testdata/golden/starred.md`.
- `just release-check` — `goreleaser check`.
- `just release-snapshot` — local goreleaser dry run into `dist/`.

## Coding Style & Naming Conventions

- Idiomatic Go ≥ 1.26 with `gofmt` and `goimports`, enforced via
  golangci-lint.
- Errors wrap with a package prefix in English — e.g.
  `github: failed to list starred (page %d): %w`,
  `starlist: missing token`. Keep the prefix when adding new error sites
  in an existing package.
- Code comments and package doc strings in English.
- The Markdown output layout is a public contract: the daily-run
  workflow rewrites a public gist with it. Any layout change requires
  regenerating the golden fixture and lands in the same commit.

## Testing Guidelines

`testify` is used for assertions. The GitHub fetcher is exercised through
`httptest.Server` via the `SetBaseURL` hook — tests must never call the
real GitHub API. The Markdown renderer is golden-tested against
`testdata/golden/starred.md`; pass `-update` to the renderer tests to
regenerate it.

Run `just test-race` before opening a PR. CI runs the same plus
golangci-lint and a `go mod tidy` check.

## Commit & Pull Request Guidelines

Conventional commits (`feat:`, `fix:`, `chore:`, `docs:`, `test:`,
`refactor:`, `build:`). CI (`.github/workflows/ci.yml`) runs on every PR
to `master`. Releases are cut by pushing a `v*` tag, which triggers
`.github/workflows/release.yml` and goreleaser (`.goreleaser.yaml`) to
publish multi-platform archives and a Docker image to
`ghcr.io/caian-org/list-github-stars`. The image is consumed by
`.github/workflows/daily-run.yml` on a 03:00 UTC cron.

## Agent-Specific Instructions

`AGENTS.md` is the canonical instruction layer for this repository.
`CLAUDE.md` is a Claude Code pointer and quick reference. Local Codex
skills live in `.codex/skills/`; Claude Code should read those skills
from `.codex/skills/` rather than duplicating them under
`.claude/skills/`. Specialist subagents are mirrored under
`.codex/agents/` and `.claude/agents/`.

Two rules every agent must honour in this repo:

1. **Never call the real GitHub API from tests.** The `github` package
   exposes `SetBaseURL` precisely so tests can point at `httptest.Server`.
   Calling production GitHub introduces flakiness and rate-limit
   failures.
2. **Renderer changes regenerate goldens in the same commit.** The
   daily-run workflow rewrites a public gist with the output, so layout
   drift is user-visible. `just update-goldens` followed by a careful
   review of the diff is the canonical workflow.

When changing repository guidance, keep `AGENTS.md`, `CLAUDE.md`, skills,
and subagents consistent in the same change set.
