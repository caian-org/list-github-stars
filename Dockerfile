FROM golang:1.20.1-alpine3.17@sha256:48f336ef8366b9d6246293e3047259d0f614ee167db1869bdbc343d6e09aed8a AS image
LABEL maintainer="Caian Ertl <hi@caian.org>"

FROM image AS base
WORKDIR /go/app
COPY go.mod .
COPY go.sum .
RUN go mod download \
    && go list -m all \
        | tail -n +2 \
        | cut -f 1 -d " " \
        | awk 'NF{print $0 "/..."}' \
        | CGO_ENABLED=0 GOOS=linux xargs -n1 \
            go build -v -trimpath -installsuffix cgo -i; echo "done"

FROM base AS build
WORKDIR /go/app
COPY . .
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -v -trimpath -installsuffix cgo -o list-github-stars -ldflags "-s -w"

FROM alpine:3.17.2@sha256:e2e16842c9b54d985bf1ef9242a313f36b856181f188de21313820e177002501 AS runtime
COPY --from=build /go/app/list-github-stars /usr/local/bin
ENTRYPOINT ["list-github-stars"]
