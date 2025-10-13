# syntax=docker/dockerfile:1.6

FROM golang:1.25.2-alpine AS builder
WORKDIR /app

RUN set -eux; \
    for i in 1 2 3 4 5; do \
        apk add --no-cache git ca-certificates && update-ca-certificates && break || \
        echo "Retrying apk add in 5s..." && sleep 5; \
    done

ENV GOPROXY=https://proxy.golang.org,direct \
    GOSUMDB=sum.golang.org \
    CGO_ENABLED=0 \
    GOOS=linux

# Copy go.mod first
COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked,id=gomodcache \
    set -eux; \
    for i in 1 2 3 4 5; do \
        go mod download -x && break || \
        echo "Retrying go mod download in 5s..." && sleep 5; \
    done

# Copy rest of the source
COPY --link . .

RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked,id=gomodcache \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked,id=gobuildcache \
    go build -trimpath -ldflags="-s -w" -o /server ./main.go

FROM alpine:3.20.3
WORKDIR /app

RUN apk add --no-cache ca-certificates && update-ca-certificates
RUN adduser -D appuser
USER appuser

COPY --from=builder /server /usr/local/bin/server
COPY ./config ./config

EXPOSE 8080
CMD ["server"]
