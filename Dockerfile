# ===============================
# 1. Builder stage
# ===============================
FROM golang:1.25-alpine3.23 AS builder

# Install build dependencies
RUN apk add --no-cache ca-certificates tzdata


# Set working directory inside container
WORKDIR /app

# Copy go mod files first (better caching)
# COPY go.mod go.sum -> Enables Docker layer caching
# Dependencies won’t re-download on every build
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code 
COPY . .

# Build the binary
# CGO_ENABLED=0 -> Produces static binary
# No libc dependency
# Works cleanly in Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o snippetbox ./cmd/web


# ===============================
# 2 Runtime stage
# ===============================
FROM alpine:3.23

# Install runtime dependencies 
RUN apk add --no-cache ca-certificates tzdata

# Create a non-root user
RUN addgroup -S snippetbox && adduser -S snippetbox -G snippetbox

# Set working directory
WORKDIR /app 

# Copy binary from builder 
COPY --from=builder /app/snippetbox .

# Copy template & static files
COPY --from=builder /app/ui ./ui
COPY --from=builder /app/tls ./tls

# Change ownership 
RUN chown -R snippetbox:snippetbox /app

# Switch to non-root user
USER snippetbox

# Expose app port (documentation only)
EXPOSE 4000

# Start the app
CMD [ "./snippetbox" ]