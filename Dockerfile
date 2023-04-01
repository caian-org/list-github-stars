FROM golang:1.20.1-alpine3.17@sha256:18da4399cedd9e383beb6b104d43aa1d48bd41167e312bb5306d72c51bd11548 AS image
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

FROM alpine:3.17.3@sha256:124c7d2707904eea7431fffe91522a01e5a861a624ee31d03372cc1d138a3126 AS runtime
COPY --from=build /go/app/list-github-stars /usr/local/bin
ENTRYPOINT ["list-github-stars"]
