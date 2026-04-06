# ---------- Build Stage ----------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Fix CGO issues
ENV CGO_ENABLED=0

# Install git (important for dependencies)
RUN apk add --no-cache git

# Copy dependency files first
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

# Copy full project
COPY . .

# ✅ Build from cmd folder
RUN go build -o main ./cmd

# ---------- Run Stage ----------
FROM alpine:latest

WORKDIR /app

# Copy binary
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]