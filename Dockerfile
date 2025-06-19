# -------- Stage 1: Build Go Binary --------
FROM golang:1.21-alpine AS builder

# Build environment setup
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GO111MODULE=on

WORKDIR /go/src/go-app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire app source code
COPY . .

# Tidy and build the Go binary
RUN go mod tidy
RUN go build -o main .

# -------- Stage 2: Final Slim Container --------
FROM alpine:latest

# Add CA certificates (needed for HTTPS in Alpine)
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy built binary from builder stage
COPY --from=builder /go/src/go-app/main .

# Copy any other runtime files (if needed)
COPY --from=builder /go/src/go-app/data ./data

# Expose app port
EXPOSE 8080

# Set the expected path for Google credentials
# The actual file will be mounted at runtime
ENV GOOGLE_APPLICATION_CREDENTIALS=/app/credentials/service-account.json

# CMD to run the Go app
CMD ["./main"]
