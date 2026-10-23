# syntax=docker/dockerfile:1
# ==============================================================================
# Этап 1: Сборка статических бинарных файлов (Multi-Stage Build)
# Используем официальный образ Golang 1.22 на базе Alpine Linux
# ==============================================================================
FROM golang:1.22-alpine AS builder

# Устанавливаем рабочую директорию сборщика
WORKDIR /build

# Установка git (для подгрузки зависимостей) и корневых SSL-сертификатов
RUN apk add --no-cache git ca-certificates

# Оптимизация кэширования слоев: копируем go.mod и go.sum
COPY go.mod go.sum ./
# Скачивание зависимостей (кэшируется, пока не изменились go.mod / go.sum)
RUN go mod download

# Копирование исходного кода проекта
COPY . .

# Компиляция микросервиса URL Shortener:
# CGO_ENABLED=0            — автономный статический бинарник без CGO-зависимостей
# GOOS=linux GOARCH=amd64  — сборка под Linux 64-bit
# -ldflags="-s -w"         — стриппинг таблицы символов и DWARF-отладки для минимизации размера
# -extldflags '-static'    — статическая компоновка
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o /app/shortener ./homework-19/cmd/shortener

# Компиляция микросервиса Analytics (обработчик событий из Kafka):
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o /app/analytics ./homework-19/cmd/analytics

# ==============================================================================
# Этап 2: Финальный минимальный production-образ Alpine
# Чистая операционная система без компилятора (~15 МБ)
# ==============================================================================
FROM alpine:3.19

# Установка сертификатов HTTPS, данных часовых поясов и утилиты wget для healthcheck
RUN apk --no-cache add ca-certificates tzdata wget

# Рабочая папка для запуска бинарников
WORKDIR /app

# Создаем системного непривилегированного пользователя appuser (UID 1000)
RUN adduser -D -u 1000 appuser
# Переключаемся на непривилегированного пользователя (Security Best Practice)
USER appuser

# Копирование обоих скомпилированных бинарных файлов с этапа сборщика builder
COPY --from=builder /app/shortener /app/shortener
COPY --from=builder /app/analytics /app/analytics

# Декларируем порты: 8080 (веб-сервис shortener) и 8081 (метрики/эндпоинты analytics)
EXPOSE 8080 8081

# Команда по умолчанию: запуск сервиса shortener (может быть переопределена в docker-compose)
CMD ["/app/shortener"]
