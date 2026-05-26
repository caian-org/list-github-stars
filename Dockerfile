# syntax=docker/dockerfile:1
ARG GO_VERSION=1.26.2

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-bookworm AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags='-s -w' -o /out/lgs ./cmd/lgs

FROM gcr.io/distroless/static-debian12 AS local-runtime
COPY --from=build /out/lgs /usr/local/bin/lgs
ENTRYPOINT ["/usr/local/bin/lgs"]
