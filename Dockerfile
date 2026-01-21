# ===== Build stage =====
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -o music_api ./cmd/music_app

FROM alpine:latest

RUN apk add --no-cache sqlite-libs

WORKDIR /app

COPY --from=builder /app/music_api .
COPY --from=builder /app/albums.db ./albums.db

EXPOSE 8080

# Run the API
CMD ["./music_api"]
