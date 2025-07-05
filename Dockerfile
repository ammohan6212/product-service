# -------- Stage 1: Build Go Binary --------
FROM golang:1.21-alpine AS builder

# Set Go build environment
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GO111MODULE=on

# Set working directory inside builder container
WORKDIR /go/src/go-app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy application source code inside src directory
COPY src/ ./src/

# Change working directory to src where main.go is
WORKDIR /go/src/go-app/src

# Tidy modules and build binary
RUN go mod tidy && \
    go build -o main .

# -------- Stage 2: Final Runtime Image --------
FROM alpine:latest

# Install CA certificates (for HTTPS calls)
RUN apk --no-cache add ca-certificates

# Create application work directory
WORKDIR /app

# Copy compiled Go binary from builder stage
COPY --from=builder /go/src/go-app/src/main .

# Create credentials directory (empty, to be mounted at runtime)
RUN mkdir -p /app/credentials

# Set GCP credentials environment variable (path inside container)
ENV GOOGLE_APPLICATION_CREDENTIALS=/app/credentials/service-account.json

# Expose port used by the Go application
EXPOSE 8080

# Run the binary
CMD ["./main"]
