# -----------------------------
# Stage 1: Builder
# -----------------------------
FROM golang:1.20-alpine AS builder

# Install Git (required for some go modules)
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the image
COPY src/ ./src/

# Set working directory to the source folder and build the app
WORKDIR /app/src
RUN go build -o /app/product-service main.go

# -----------------------------
# Stage 2: Runtime
# -----------------------------
FROM alpine:latest

# Create work directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/product-service .

# Expose application port
EXPOSE 8080

# Run the application
CMD ["./product-service"]
