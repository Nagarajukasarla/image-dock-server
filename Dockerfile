# -----------------------------------------------------
# Stage 1: Build the Go binary
# -----------------------------------------------------
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Disable CGO and build static binary (safer for Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o image-dock-server .

# -----------------------------------------------------
# Stage 2: Minimal runtime image
# -----------------------------------------------------
FROM scratch

# Copy binary
COPY --from=builder /app/image-dock-server /

# Expose app port
EXPOSE 8080

# Command
ENTRYPOINT ["/image-dock-server"]
