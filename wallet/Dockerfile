# Stage 1: Builder
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Download module dependencies (leveraging layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /bin/wallet-app ./cmd/app

# Stage 2: Lightweight runtime image
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary and configuration files
COPY --from=builder /bin/wallet-app /app/wallet-app
COPY config.env /app/config.env
COPY migrations /app/migrations

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=5s --timeout=3s --retries=5 --start-period=5s \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/wallet-app"]
