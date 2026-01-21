Create **`Dockerfile`** in the project root.

This uses **multi-stage builds** and works with SQLite.

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -o api ./cmd/music_app

# Runtime stage
FROM alpine:latest

RUN apk add --no-cache sqlite-libs

WORKDIR /app

COPY --from=builder /app/music_app .
COPY --from=builder /app/albums.db ./albums.db

EXPOSE 8080

CMD ["./music_app"]
