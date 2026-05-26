---
name: lgs-dev-workflow
description: Local development workflow for list-github-stars. Use when orienting in the workspace, choosing just/go/Docker/GoReleaser commands, or validating general changes.
---

# LGS Dev Workflow

Use this skill for repo orientation, command selection, and everyday
implementation hygiene across the workspace.

## Repo Shape

`list-github-stars` is a single-module Go project (`github.com/caian-org/list-github-stars`).
The main binary is `cmd/lgs`, supported by repo-local packages in `internal/`.

- `cmd/lgs/main.go` — minimal entrypoint that delegates to the CLI layer.
- `cmd/lgs/cli/` — `urfave/cli/v3` root command and the `ProgramVersion` /
  `ProgramCommitSHA` / `ProgramBuildTime` symbols injected by ldflags.
- `internal/github/` — paginated GitHub starred-repos fetcher built on
  `go-github/v76`. Exposes `SetBaseURL` for `httptest.Server`-based tests.
- `internal/render/` — pure Markdown renderer with language grouping,
  slug disambiguation, and case-insensitive sorting.
- `internal/starlist/` — orchestration: fetch → render → write.
- `testdata/golden/` — golden Markdown fixtures regenerated with the
  `-update` test flag.

## Commands

Prefer `devbox shell` when available.

```bash
just list                  # list binaries under cmd/
just build lgs             # build bin/lgs with dev ldflags (commit + timestamp)
just run lgs               # build, then run bin/lgs (reads GITHUB_TOKEN)
just test                  # go test ./...
just test-race             # go test ./... -race
just cover                 # coverage profile + per-function totals
just lint                  # go vet ./... (CI also runs golangci-lint)
just tidy                  # go mod tidy
just update-goldens        # regenerate testdata/golden/starred.md
just release-check         # goreleaser check
just release-snapshot      # local goreleaser dry-run into dist/
```

## Validation Lanes

Run the smallest useful check first.

- After a small code change: `just test`.
- Before opening a PR: `just test-race` plus `just lint` (CI also runs
  golangci-lint and a tidy check).
- After touching `internal/render/*` or any layout helper: regenerate
  goldens (`just update-goldens`), inspect the diff, then `just test-race`.
- After touching `.goreleaser.yaml`, `Dockerfile`, or the release workflow:
  `just release-check` plus `just release-snapshot`.

## Runtime Env

The CLI needs a GitHub personal access token:

```bash
export GITHUB_TOKEN=ghp_xxx
./bin/lgs > my-stars.md
```

`--token` (alias `-t`) overrides the env var. `--user` (alias `-u`) selects
a user other than the authenticated one (no `/user` lookup is performed in
that case).

## Implementation Rules

- Idiomatic Go ≥ 1.26 with `gofmt` / `goimports`, enforced by golangci-lint.
- Keep `internal/github` HTTP-pure: tests go through `SetBaseURL` against
  `httptest.Server`. Never make real GitHub API calls from tests.
- Keep `internal/render` I/O-free. Anything that writes to a `Writer`
  belongs in `internal/starlist`.
- Version metadata symbols (`ProgramVersion`, `ProgramCommitSHA`,
  `ProgramBuildTime`) are injected by ldflags. Renaming requires touching
  `.justfile` and `.goreleaser.yaml` in the same change.

## Commit Hygiene

Conventional commits (`feat:`, `fix:`, `chore:`, `docs:`, `test:`,
`refactor:`, `build:`). Keep commits focused. Renderer changes that
regenerate goldens land in the same commit as the change that caused
them, never separately.
