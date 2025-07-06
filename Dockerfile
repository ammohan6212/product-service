# -------- Stage 1: Build Go Binary --------
FROM golang:1.21-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GO111MODULE=on

# Create base working directory
WORKDIR /go/src/go-app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./

# Download modules
RUN go mod download

# Copy full source code
COPY src/ ./src/

# Move go.mod into src/ to match code if needed
# (Optional: only needed if your code expects go.mod inside src)
RUN cp go.mod go.sum ./src/

# Change working directory to src
WORKDIR /go/src/go-app/src

# Tidy and build
RUN go mod tidy && \
    go build -o main .

# -------- Stage 2: Final Runtime Image --------
FROM alpine:latest

RUN apk --no-cache add ca-certificates && \
    apk --no-cache upgrade

WORKDIR /app

# Copy built binary from builder
COPY --from=builder /go/src/go-app/src/main .

# Prepare credentials directory (empty)
RUN mkdir -p /app/credentials

ENV GOOGLE_APPLICATION_CREDENTIALS=/app/credentials/service-account.json

EXPOSE 8080

CMD ["./main"]
