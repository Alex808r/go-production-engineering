# syntax=docker/dockerfile:1

# Этап 1: Сборка статических бинарных файлов (Multi-Stage Build)
FROM golang:1.22-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git ca-certificates

# Кэширование зависимостей Go
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Компиляция микросервиса URL Shortener
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o /app/shortener ./homework-19/cmd/shortener

# Компиляция микросервиса Analytics
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o /app/analytics ./homework-19/cmd/analytics

# Этап 2: Финальный минимальный production-образ Alpine
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /app

# Непривилегированный пользователь
RUN adduser -D -u 1000 appuser
USER appuser

# Копирование скомпилированных бинарных файлов
COPY --from=builder /app/shortener /app/shortener
COPY --from=builder /app/analytics /app/analytics

EXPOSE 8080 8081

CMD ["/app/shortener"]
