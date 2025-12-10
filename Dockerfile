FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git for fetch dependencies
# RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
# COPY --from=builder /app/.env . # Optional: usually envs are passed via docker-compose

EXPOSE 8080

CMD ["./main"]
