# --------------------------------------------------------------------------------------------------
FROM golang:1.22.5-alpine3.20 AS build
LABEL maintainer="Caian Ertl <hi@caian.org>"
ENV GOCACHE=/root/.cache/go-build
WORKDIR /go/app

RUN apk update \
    && apk add --no-cache "make=~4.4" \
    && rm -rf /var/cache/apk/*

COPY vendor/ ./vendor/
COPY go.* .
COPY Makefile .
COPY main.go .
RUN --mount=type=cache,target="/root/.cache/go-build" make release

# --------------------------------------------------------------------------------------------------
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/app/lgs /
ENTRYPOINT ["/lgs"]
