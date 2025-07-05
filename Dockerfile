# -------- Stage 1: Build Go Binary --------
FROM golang:1.21-alpine AS builder

# Set Go build environment
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GO111MODULE=on

WORKDIR /go/src/go-app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy all files inside src directory (main.go + subfolders)
COPY src/ ./src/

# Change working directory to src
WORKDIR /go/src/go-app/src

# Tidy and build the Go binary
RUN go mod tidy && \
    go build -o main .

# -------- Stage 2: Final Runtime Image --------
FROM alpine:latest

# Install CA certificates (required for HTTPS)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy built binary from builder stage
COPY --from=builder /go/src/go-app/src/main .

# Ensure credentials directory exists (even if mounted at runtime)
RUN mkdir -p /app/credentials

# Do NOT copy service-account.json into the image
# COPY service-account.json /app/credentials/service-account.json

# Set env var for service account key to be mounted later
ENV GOOGLE_APPLICATION_CREDENTIALS=/app/credentials/service-account.json

EXPOSE 8080

CMD ["./main"]
