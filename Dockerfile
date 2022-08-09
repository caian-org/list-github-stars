FROM golang:1.18.1-alpine3.15 AS image

FROM image AS base
WORKDIR /go/app
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go list -m all \
    | tail -n +2 \
    | cut -f 1 -d " " \
    | awk 'NF{print $0 "/..."}' \
    | CGO_ENABLED=0 GOOS=linux xargs -n1 \
    go build -v -installsuffix cgo -i; echo "done"

FROM base AS build
WORKDIR /go/app
COPY . .
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -v -installsuffix cgo -o list-github-stars -ldflags "-s -w"

FROM alpine:3.16.2 as runtime
COPY --from=build /go/app/list-github-stars /usr/local/bin
ENTRYPOINT ["list-github-stars"]
