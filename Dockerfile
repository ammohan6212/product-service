# Use official Go image as the build stage
FROM golang:1.21-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GO111MODULE=on

# Set the working directory inside the container
WORKDIR /go/src/go-postgres-app

# Copy go mod and sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go app binary
RUN go build -o main .

# ---- Final stage ----
FROM alpine:latest

# Set working directory in the final container
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /go/src/go-postgres-app/main .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the binary
CMD ["./main"]
