# Start from the official Go image
FROM golang:1.24-alpine AS builder

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}

# Set the working directory
WORKDIR /app

# Install necessary tools
RUN apk add --no-cache bash git curl

RUN curl -sSL "https://github.com/sqlc-dev/sqlc/releases/download/v1.30.0/sqlc_1.30.0_${TARGETOS}_${TARGETARCH}.tar.gz" | tar -xz -C /usr/local/bin sqlc

# copy go mod and sum files and download deps
COPY go.mod go.sum ./
RUN go mod download

# copy the sqlc configuration and sql files
COPY sqlc.yaml query.sql schema.sql ./
RUN /usr/local/bin/sqlc generate

# copy the rest of the application source code
COPY . .

# Build the Go project
RUN go build -o app .

# the actual application image
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app .

RUN apt-get update && apt-get install -y ca-certificates

ENTRYPOINT ["./app"]