ARG GO_VERSION=1.25.2
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /usr/src/app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy all files from the build context into the container
COPY . .

# Specify the build path to the cmd/govote directory
RUN CGO_ENABLED=0 go build -v -o /run-app ./cmd/govote

FROM gcr.io/distroless/static-debian12

# Copy the compiled binary into the runtime image
COPY --from=builder /run-app /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/run-app"]
CMD ["-keypath", "/data/govote"]
