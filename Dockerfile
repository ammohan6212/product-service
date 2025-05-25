# Stage 1: Build React frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm install

COPY frontend/ ./
RUN npm run build


# Stage 2: Build Go backend
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

COPY backend/main.go .

# Add Gin package
RUN go mod init backend \
    && go get github.com/gin-gonic/gin \
    && go build -o server main.go


# Stage 3: Final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy compiled Go server
COPY --from=backend-builder /app/server .

# Copy React build files
COPY --from=frontend-builder /app/frontend/build ./frontend/build

# Expose the port used by the server
EXPOSE 8000

# Run the Go server
CMD ["./server"]
