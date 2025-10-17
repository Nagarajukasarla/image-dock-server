# -----------------------------------------------------
# Stage 1: Build the Go binary
# -----------------------------------------------------
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o image-dock-server .

# -----------------------------------------------------
# Stage 2: Minimal runtime image (with certs)
# -----------------------------------------------------
FROM alpine:latest

# Install CA certificates (needed for HTTPS)
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/image-dock-server .

EXPOSE 8000
ENTRYPOINT ["./image-dock-server"]