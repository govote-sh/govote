# Ensure uniform casing for Dockerfile keywords
ARG GO_VERSION=1.24.0
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy all files from the build context into the container
COPY . .

# Specify the build path to the cmd/govote directory
RUN go build -v -o /run-app ./cmd/govote

FROM debian:bookworm

# Copy the compiled binary into the runtime image
COPY --from=builder /run-app /usr/local/bin/
CMD ["run-app", "-keypath", "/data/govote"]
