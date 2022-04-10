FROM golang:alpine AS builder
WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# TODO: Should probably replace with COPY . . and .dockerignore file
COPY cmd cmd
COPY client client
COPY crawler crawler
COPY server server
COPY utils utils
COPY tests tests

RUN go build -o /app/bin/history /app/cmd/history
RUN go build -o /app/bin/server /app/cmd/server

# Binaries
FROM alpine

WORKDIR /app
COPY --from=builder /app/bin/history /app/crawler-history
COPY --from=builder /app/bin/server /app/crawler-server
