# syntax=docker/dockerfile:1
ARG GO_VERSION=1.26.2

# Build stage shared by the local-runtime target. Goreleaser bypasses this and
# uses the pre-built binary from its own build context (see the
# goreleaser-runtime stage below).
FROM golang:${GO_VERSION}-bookworm AS build
WORKDIR /src
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags='-s -w' -o /out/lgs ./cmd/lgs


# `docker build .` (no --target) builds from source. Handy for development.
FROM gcr.io/distroless/static-debian12 AS local-runtime
COPY --from=build /out/lgs /usr/local/bin/lgs
ENTRYPOINT ["/usr/local/bin/lgs"]


# Used by goreleaser's dockers_v2 step at release time. The build context
# already contains linux/amd64/lgs (the binary goreleaser cross-built).
FROM gcr.io/distroless/static-debian12 AS goreleaser-runtime
COPY linux/amd64/lgs /usr/local/bin/lgs
ENTRYPOINT ["/usr/local/bin/lgs"]
