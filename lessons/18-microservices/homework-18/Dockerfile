# syntax=docker/dockerfile:1

# Этап 1: Сборка статического бинарного файла
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Установка необходимых утилит
RUN apk add --no-cache git ca-certificates

# Копируем go.mod и go.sum из корня репозитория курса
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем легковесный статический бинарник без CGO
WORKDIR /build/homework-18/cmd/shortener
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o /app/shortener .

# Этап 2: Финальный минимальный production-образ
FROM alpine:3.19

# Сертификаты для HTTPS-запросов и часовой пояс
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Создаем непривилегированного пользователя для безопасности (Security Best Practice)
RUN adduser -D -u 1000 appuser
USER appuser

# Копируем скомпилированный бинарный файл
COPY --from=builder /app/shortener /app/shortener

# Переменные окружения по умолчанию
ENV PORT=8080
ENV BASE_URL=http://localhost:8080

EXPOSE 8080

ENTRYPOINT ["/app/shortener"]
