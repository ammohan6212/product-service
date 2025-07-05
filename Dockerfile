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

# Copy the rest of the application source code
COPY . .

# Tidy and build the Go binary
RUN go mod tidy && \
    go build -o main .

# -------- Stage 2: Final Runtime Image --------
FROM alpine:latest

# Install CA certificates (required for HTTPS)
RUN apk --no-cache add ca-certificates

# Set working directory inside the container
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /go/src/go-app/main .
# Ensure credentials directory exists (even if credentials are mounted later)
RUN mkdir -p /app/credentials

#COPY service-account.json /app/credentials/service-account.json

ENV GOOGLE_APPLICATION_CREDENTIALS=/app/credentials/service-account.json



# Expose the application port
EXPOSE 8080

# Start the Go app
CMD ["./main"]
