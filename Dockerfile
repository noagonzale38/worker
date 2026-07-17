# Build container
FROM golang:1.24 AS builder

RUN go version

RUN apt-get update && apt-get upgrade -y && apt-get install -y ca-certificates git zlib1g-dev

COPY . /go/src/github.com/TicketsBot-cloud/worker

# go.mod replaces github.com/TicketsBot-cloud/database with ../database;
# provide it via: docker buildx build --build-context database=../database
COPY --from=database / /go/src/github.com/TicketsBot-cloud/database

WORKDIR /go/src/github.com/TicketsBot-cloud/worker

RUN git submodule update --init --recursive --remote

RUN set -Eeux && \
    go mod download && \
    go mod verify

ARG TARGETOS
ARG TARGETARCH

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build \
    -tags=jsoniter \
    -trimpath \
    -o main cmd/worker/main.go

# Prod container
FROM ubuntu:latest

RUN apt-get update && apt-get upgrade -y && apt-get install -y ca-certificates curl

COPY --from=builder /go/src/github.com/TicketsBot-cloud/worker/locale /srv/worker/locale
COPY --from=builder /go/src/github.com/TicketsBot-cloud/worker/main /srv/worker/main

RUN chmod +x /srv/worker/main

RUN useradd -m container
USER container
WORKDIR /srv/worker

CMD ["/srv/worker/main"]
