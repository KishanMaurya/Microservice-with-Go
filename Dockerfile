# ---------- Build Stage ----------
FROM golang:1.26-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0

RUN apk add --no-cache git

# Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Source
COPY . .

# Build
RUN go build -ldflags="-s -w" -o main ./cmd

# ---------- Run Stage ----------
FROM alpine:3.20

WORKDIR /app

# Security: non-root user
RUN adduser -D appuser
USER appuser

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]