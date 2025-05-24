# Start from the official Golang image for building
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go app (static binary)
RUN go build -o product-service main.go

# Use a minimal base image for running the app
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/product-service .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./product-service"]