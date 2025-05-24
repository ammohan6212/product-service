# Start from the official Golang image for building
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
RUN go mod download

# Copy the source code
COPY src/ ./src/

# Set working directory to src and build the Go app
WORKDIR /app/src
RUN go build -o /app/product-service main.go

# Use a minimal base image for running the app
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/product-service .

EXPOSE 8080

CMD ["./product-service"]