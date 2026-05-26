# Claude Notes

Pointer doc for Claude Code agents working in `list-github-stars`.

**Read `AGENTS.md` first**. It is the canonical instruction layer. This
file covers only what is specific to Claude Code or to quick orientation.

## Bootstrap

1. `AGENTS.md` — canonical repository rules, structure, commands,
   testing, release, and agent-specific instructions.
2. `.codex/skills/lgs-dev-workflow/SKILL.md` — default local workflow
   guidance for this repo.
3. A subsystem skill when relevant:
   - `.codex/skills/lgs-github-fetch/SKILL.md`
   - `.codex/skills/lgs-markdown-render/SKILL.md`
   - `.codex/skills/lgs-release/SKILL.md`

`AGENTS.md` is authoritative. If `CLAUDE.md` disagrees with it, follow
`AGENTS.md` and reconcile this file in the same change set.

## Claude Code Specifics

- Use Claude Code's native subagent feature for delegated lanes when the
  session policy allows it. Specialists live under `.claude/agents/`,
  mirroring `.codex/agents/`.
- Skills are not duplicated for Claude. The canonical home is
  `.codex/skills/`; read skills from there.
- Claude Code settings live under `.claude/settings.json`. Hooks are
  currently empty.
- Use the local checkout as the source of truth. Do not use GitHub API
  reads as a substitute for inspecting files in this repo unless it is a
  narrow one-off check.
- Serialize shared-checkout mutations to one owner: edits, generated
  output, patch application, staging, committing, rebasing, branch
  switching, and pushing.
- After editing `AGENTS.md`, `.codex/skills/`, `.codex/agents/`, or
  `.claude/agents/`, update this file if the Claude-facing guidance
  changes.

## Quick Reference

Project purpose:

- Authenticate with the GitHub REST API via a personal access token.
- Paginate through the authenticated (or `--user`) account's starred
  repositories.
- Group them by language, sort alphabetically (case-insensitive), and
  render the result as a Markdown document.

Command surface:

- `devbox shell` — enter the pinned development environment.
- `just list` — list binaries under `cmd/`.
- `just build lgs` — build `bin/lgs` with dev metadata.
- `just run lgs` — build and run.
- `just test-race` — default validation command.
- `just release-check` — `goreleaser check`.
- `just release-snapshot` — local release dry run into `dist/`.

Runtime envs:

- Required: `GITHUB_TOKEN` (PAT with `read:user` scope minimum; `gist`
  scope only needed for the daily-run gist edit).
- Optional flags: `--token` / `-t` (overrides env), `--user` / `-u`
  (selects a non-authenticated user; the `/user` API lookup is skipped).

High-traffic paths:

- `cmd/lgs/cli/cli.go` — urfave/cli/v3 root command.
- `cmd/lgs/cli/meta.go` — `ProgramVersion` / `ProgramCommitSHA` /
  `ProgramBuildTime` symbols.
- `internal/github/client.go` — paginated star fetcher.
- `internal/render/render.go` — Markdown rendering.
- `internal/starlist/starlist.go` — orchestration glue.

Skills (canonical, read from `.codex/skills/`):

- `.codex/skills/lgs-dev-workflow/SKILL.md`
- `.codex/skills/lgs-github-fetch/SKILL.md`
- `.codex/skills/lgs-markdown-render/SKILL.md`
- `.codex/skills/lgs-release/SKILL.md`

Subagents (`.claude/agents/<name>.md` and `.codex/agents/<name>.toml`):

- `lgs-architect` — cross-package design questions across `cmd/lgs/` and
  `internal/{github,render,starlist}/`.
- `lgs-renderer-specialist` — Markdown layout, slug disambiguation,
  golden fixtures.
- `lgs-release-specialist` — `.goreleaser.yaml`, Dockerfile, GHCR,
  daily-run, version ldflags.
- `lgs-test-reviewer` — read-only review for tests and release dry-run
  coverage.
