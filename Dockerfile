FROM golang:alpine AS builder
WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY cmd cmd
COPY crawler crawler
COPY tests tests

RUN go build -o /app/bin/crawler /app/cmd/crawler

# Binaries
FROM alpine

WORKDIR /app
COPY --from=builder /app/bin/crawler /app/crawler
