FROM golang:alpine AS build-env

# Install dependencies
RUN apk add --update git build-base linux-headers

# Set working directory for the build
WORKDIR /go/src/git.vdb.to/cerc-io/laconic2d

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Add source files
COPY . .

# Make the binary
RUN make build

# Final image
FROM alpine:3.17.0

# Install ca-certificates
RUN apk add --update ca-certificates jq curl

# Copy over binaries from the build-env
COPY --from=build-env /go/src/git.vdb.to/cerc-io/laconic2d/build/laconic2d /usr/bin/laconic2d

WORKDIR /

# Run laconic2d by default
CMD ["laconic2d"]
