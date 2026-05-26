set shell := ["/bin/bash", "-c"]

_help:
    @just --list

# list all available binaries under cmd/
list:
    @find cmd -mindepth 1 -maxdepth 1 -type d -exec basename {} \;

# build cmd/<program> into bin/<program> with dev-time ldflags
build program *extra_flags:
    #!/usr/bin/env bash
    set -euo pipefail

    CLIPKG="github.com/caian-org/list-github-stars/cmd/{{ program }}/cli"
    COMMIT_HASH="$(git rev-parse --short HEAD)"
    BUILD_TS="$(date -u '+%Y-%m-%dT%H:%M:%S')"

    LDFLAGS=(
      "-X '${CLIPKG}.ProgramVersion=0.0.0-dev'"
      "-X '${CLIPKG}.ProgramCommitSHA=${COMMIT_HASH}'"
      "-X '${CLIPKG}.ProgramBuildTime=${BUILD_TS}'"
    )

    mkdir -p bin
    cd cmd/{{ program }} \
      && go build \
        -trimpath \
        -ldflags="${LDFLAGS[*]}" \
        {{ extra_flags }} \
        -o ../../bin/{{ program }}

# build then run bin/<program> with the given args
run program *args:
    @just build {{ program }}
    @./bin/{{ program }} {{ args }}

# run go test ./...
test:
    @go test ./...

# run the test suite with the race detector
test-race:
    @go test ./... -race

# coverage profile + per-function totals
cover:
    @go test -coverprofile=coverage.out ./...
    @go tool cover -func=coverage.out | tail -20

# go vet (CI also runs golangci-lint)
lint:
    @go vet ./...

# go mod tidy
tidy:
    @go mod tidy

# remove build outputs
clean:
    @rm -rf bin coverage.out dist

# regenerate golden Markdown fixtures
update-goldens:
    @go test ./internal/render/... -update

# validate the goreleaser configuration
release-check:
    @goreleaser check

# build release artifacts locally without publishing
release-snapshot:
    @goreleaser release --snapshot --clean
