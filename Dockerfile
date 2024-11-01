# Ensure uniform casing for Dockerfile keywords
ARG GO_VERSION=1.23
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify
ENV SSL_CERT_DIR=/etc/ssl/certs

# Install CA certificates with curl
RUN apt-get update && \
    apt-get install -y ca-certificates curl && \
    curl -fsSL https://curl.se/ca/cacert.pem -o /etc/ssl/certs/ca-certificates.crt && \
    rm -rf /var/lib/apt/lists/*

# Copy all files from the build context into the container
COPY . .

# Specify the build path to the cmd/govote directory
RUN go build -v -o /run-app ./cmd/govote

FROM debian:bookworm

# Copy the compiled binary into the runtime image
COPY --from=builder /run-app /usr/local/bin/
CMD ["run-app"]