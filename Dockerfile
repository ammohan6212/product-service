# === Frontend Builder ===
FROM node:20-alpine AS frontend-builder

WORKDIR /app
COPY frontend/package*.json ./frontend/
RUN cd frontend && npm install

COPY frontend/ ./frontend
RUN cd frontend && npm run build

# === Backend Builder ===
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app
COPY src/main.go .
RUN go mod init src && \
    go get github.com/gin-gonic/gin && \
    go build -o server .

# === Final Stage ===
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app

# Copy Go binary
COPY --from=backend-builder /app/server .

# Copy React build output
COPY --from=frontend-builder /app/frontend/build ./frontend

# Serve React and API using Gin
EXPOSE 8000

CMD ["./server"]
