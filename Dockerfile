FROM golang:1.21-bullseye AS builder

# Set working directory for the build
WORKDIR /go/src/git.vdb.to/cerc-io/laconicd

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Add source files
COPY . .

# Make the binary
RUN make build

# Final image
FROM ubuntu:22.04

# Install ca-certificates, jq, curl, bash, and other necessary packages
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    jq curl netcat bash \
    && rm -rf /var/lib/apt/lists/*

# Copy over binary from the builder
COPY --from=builder /go/src/git.vdb.to/cerc-io/laconicd/build/laconicd /usr/bin/laconicd

WORKDIR /

# Run laconicd by default
CMD ["laconicd"]
