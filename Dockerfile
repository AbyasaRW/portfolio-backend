# Stage 1: Build binary
FROM golang:alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o portfolio-backend .

# Stage 2: Minimal runtime image
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/portfolio-backend .

# Expose port
EXPOSE 8080

# Run the app
CMD ["./portfolio-backend"]
