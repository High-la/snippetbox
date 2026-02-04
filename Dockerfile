# ===============================
# 1. BUILD stage
# ===============================
FROM golang:1.25-alpine3.23 AS build-stage

# Install build dependencies
# Install git and certs (needed for modules and external HTTPS requests)
# Install git (needed for go mod)
RUN apk add --no-cache git ca-certificates

# Set working directory inside container
WORKDIR /app


# 1. Cache dependencies (standard optimization)
# Copy go mod files first (better caching)
# COPY go.mod go.sum -> Enables Docker layer caching
# Dependencies won’t re-download on every build
COPY go.mod go.sum ./
RUN go mod download


# 2. Copy source and build

# Copy the rest of the source code 
COPY . /app

# Build the binary
# CGO_ENABLED=0 -> Produces static binary
# No libc dependency
# Works cleanly in Alpine
# -ldflags="-w -s" reduces binary size by removing debug info
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o snippetbox ./cmd/web


# ===============================
# 2 Runtime stage
# ===============================
# Using Debian 12 Distroless for maximum security (no shell/vulnerabilities)
FROM gcr.io/distroless/static-debian12


# Set working directory
WORKDIR /app 

# Copy the statically linked binary from build-stage
COPY --from=build-stage /app/snippetbox /app/snippetbox

# Copy SSR assets (templates and static files)
COPY --from=build-stage /app/ui /app/ui

# Optional: Copy migrations if you run them within the app container
COPY --from=build-stage /app/internal/models/testdata /app/migrations


# Expose app port (internal only)
EXPOSE 52100

# Distroless automatically runs as 'nonroot' for safety
USER nonroot:nonroot

# Start the app
CMD [ "/app/snippetbox" ]