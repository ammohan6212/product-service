# Stage 1: Build React frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app
COPY frontend/package*.json ./frontend/
RUN cd frontend && npm install
COPY frontend ./frontend
RUN cd frontend && npm run build

# Stage 2: Build Go backend
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app
COPY src/main.go .
RUN go mod init src \
    && go get github.com/gin-gonic/gin \
    && go build -o server main.go

# Stage 3: Final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=backend-builder /app/server .
COPY --from=frontend-builder /app/frontend/build ./frontend/build

EXPOSE 8000
CMD ["./server"]
