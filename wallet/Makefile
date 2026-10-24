.PHONY: help all build run test test-concurrency test-race docker-up docker-down docker-reset docker-logs docker-tools clean

# Команда по умолчанию выводит интерактивную справку
help: ## Показать список доступных команд
	@awk 'BEGIN {FS = ":.*##"; printf "\nИспользование:\n  make \033[36m<команда>\033[0m\n\nКоманды разработки:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

all: test build ## Запуск тестов и сборка приложения

build: ## Скомпилировать бинарный файл в bin/wallet-app
	go build -o bin/wallet-app ./cmd/app

run: ## Локальный запуск приложения через go run
	go run ./cmd/app

test: ## Запуск всех unit-тестов
	go test -v ./...

test-race: ## Запуск тестов с детектором гонок данных (Race Detector)
	go test -v -race ./...

test-concurrency: ## Запуск нагрузочных стресс-тестов параллелизма (1000 RPS)
	go test -v -run TestConcurrency ./tests

docker-up: ## Запустить весь стек в Docker (PostgreSQL + Wallet App)
	docker compose up --build -d

docker-down: ## Безопасно остановить контейнеры (БЕЗ удаления данных БД)
	docker compose down

docker-reset: ## ⚠️ ОПАСНО: Остановить контейнеры и УДАЛИТЬ все данные БД (тома)
	docker compose down -v

docker-logs: ## Просмотр логов всех сервисов в реальном времени
	docker compose logs -f

docker-tools: ## Запустить GUI Adminer для просмотра PostgreSQL (http://localhost:8081)
	docker compose --profile tools up -d adminer
	@echo "✓ Adminer запущен: http://localhost:8081"

clean: ## Очистить скомпилированные локальные бинарники
	rm -rf bin/

