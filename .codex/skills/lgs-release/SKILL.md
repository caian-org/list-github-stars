---
name: lgs-release
description: Release guidance for list-github-stars. Use when changing GoReleaser, Dockerfile, GitHub Actions release/CI/daily-run workflows, Devbox tooling, or CLI version metadata.
---

# LGS Release

Use this skill when changing GoReleaser, the Dockerfile, the GitHub Actions
workflows under `.github/workflows/`, Devbox packages, `.justfile`, or CLI
version metadata.

## Release Shape

Releases are tag-driven. Pushing a `v*` tag runs
`.github/workflows/release.yml`, which uses GoReleaser v2 to:

- Build binaries for `linux/darwin/windows × amd64/arm64`, packaged as
  `tar.gz` (or `zip` on Windows), with `README.md` and `LICENSE*` bundled.
- Publish a multi-stage Docker image to
  `ghcr.io/caian-org/list-github-stars`, tagged as both `{{ .Tag }}` and
  `latest`, built from the `goreleaser-runtime` stage of `Dockerfile`.
- Generate `checksums.txt` (SHA-256).
- Generate a GitHub Release with grouped changelog (Features / Bug Fixes
  / Others) filtered against `docs:` / `test:` / `chore:` / `ci:` and
  merge commits.

The image is pulled by `.github/workflows/daily-run.yml` on the cron
schedule, so any change to the published image name or tag scheme must be
mirrored there.

## Version Metadata

Release builds inject these variables in `cmd/lgs/cli`:

- `ProgramVersion`   ← `{{ .Version }}`
- `ProgramCommitSHA` ← `{{ .Commit }}`
- `ProgramBuildTime` ← `{{ .Date }}`

The local `just build lgs` path injects development metadata from the
current commit short SHA and a UTC timestamp.

## Dockerfile Stages

`Dockerfile` has three stages:

- `build`     — `golang:1.26.2-bookworm`, source-built static binary.
- `local-runtime` — distroless, copies the binary from `build`. Used when
  you run `docker build .` without `--target`.
- `goreleaser-runtime` — distroless, copies `linux/amd64/lgs` from the
  build context (goreleaser injects the pre-built binary there).
  `.goreleaser.yaml` selects this stage via
  `flags: [--target=goreleaser-runtime]`.

`CGO_ENABLED=0` is required everywhere — distroless static has no libc.

## Validation

Before considering release changes done:

```bash
just release-check         # goreleaser check
just release-snapshot      # builds dist/ locally, never publishes
docker build --target local-runtime -t lgs:local .
docker run --rm -e GITHUB_TOKEN=$GITHUB_TOKEN lgs:local
```

If GoReleaser is not installed globally, use:

```bash
devbox run release-check
devbox run release-snapshot
```

## Daily-Run Workflow

`.github/workflows/daily-run.yml` runs at 03:00 UTC and on
`workflow_dispatch`. It:

1. Logs into GHCR with the workflow's `GITHUB_TOKEN` (read scope on the
   package).
2. Runs `ghcr.io/caian-org/list-github-stars:latest` with `GITHUB_TOKEN`
   set to the long-lived `AUTH_TOKEN` secret (which carries the gist
   write scope).
3. Pipes the output into `my-github-stars.md`.
4. Calls `gh gist edit ada2117bd8c73a1e94e49580fd5c7cf7 --add ...` to
   update the public gist.

`AUTH_TOKEN` is a fine-grained PAT that the repository owner manages —
the workflow does not create it.

## Risk Checks

- Keep `go-version-file: go.mod` in the release and CI workflows so the
  toolchain stays in sync with `go.mod`.
- Keep `CGO_ENABLED=0` in all builds. The distroless `static` base has
  no glibc.
- Do not reintroduce Docker Hub publishing. The repo migrated to GHCR
  and the daily-run consumes the GHCR image.
- The CodeQL workflow was intentionally dropped during modernization;
  do not re-add it without a clear reason.
- Snapshot artifacts in `dist/` are never published. The `release.yml`
  workflow uses `args: release --clean` against the tagged commit only.
