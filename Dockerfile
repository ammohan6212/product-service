# ---- Build Stage ----
FROM golang:1.21-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GO111MODULE=on

WORKDIR /go/src/go-postgres-app

COPY go.mod go.sum ./
RUN go mod download

COPY . .              # Copies the source code AND data folder
RUN go mod tidy
RUN go build -o main .

# ---- Final Stage ----
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary
COPY --from=builder /go/src/go-postgres-app/main .

# Copy the data directory
COPY --from=builder /go/src/go-postgres-app/data ./data

EXPOSE 8080

CMD ["./main"]
