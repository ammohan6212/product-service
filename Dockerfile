# Stage 1: Build Go backend
FROM golang:1.21 AS builder

WORKDIR /app

COPY backend/main.go .

RUN go mod init backend && \
    go get github.com/gin-gonic/gin && \
    go build -o server .

# Stage 2: Minimal image with frontend
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy Go binary
COPY --from=builder /app/server .

# Copy React build files
COPY frontend/ ./frontend/

EXPOSE 8080

CMD ["./server"]
