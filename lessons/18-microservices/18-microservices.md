# 18-microservices — Микросервисная архитектура, контейнеризация Docker и Kubernetes

В данном модуле изучаются принципы проектирования микросервисов в соответствии с методологией «12 Factor App», контейнеризация Go-приложений в Docker с использованием многоэтапной сборки (Multi-stage builds), обеспечение безопасности контейнеров (non-root пользователи), а также декларативное развертывание в Kubernetes (манифесты Deployment, Service, Liveness и Readiness Probes).

---

## 📋 Структура каталога

```text
18-microservices/
├── microservice/                # Демонстрационный микросервис для Kubernetes
│   ├── cmd/app/main.go          # Точка входа HTTP-сервера
│   ├── pkg/api/                 # Обработчики API (api.go, api_test.go), возврат hostname пода
│   ├── go.mod, go.sum           # Изолированный Go-модуль (microapp)
│   └── build/
│       ├── Dockerfile           # Docker-образ микросервиса
│       ├── deployment.yaml      # Kubernetes Deployment: 2 реплики, порты
│       └── service.yaml         # Kubernetes Service: маршрутизация TCP трафика (port 8080)
└── homework-18/                 # Домашнее задание №18: Сервис сокращения ссылок (URL Shortener)
    ├── Dockerfile               # Production Multi-stage образ (<25 МБ, non-root appuser)
    ├── docker-compose.yml       # Декларативный запуск с Healthcheck
    ├── README.md                # Спецификация API, архитектура и примеры curl
    ├── cmd/shortener/main.go    # Точка входа с Graceful Shutdown
    └── pkg/
        ├── api/                 # REST API (CRUD ссылок, redirect, Prometheus /metrics, /healthz)
        ├── config/              # Чтение конфигурации из ENV (PORT, BASE_URL)
        └── storage/             # Хранилище ссылок на базе среза в памяти ([]Link)
```

---

## 🔍 Ключевые концепции

1. **Многоэтапная сборка Docker (Multi-stage build):**
   - **Этап сборки (`builder`):** компиляция статически слинкованного бинарника (`CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w"`) в окружении `golang:alpine`.
   - **Этап выполнения (`runtime`):** копирование скомпилированного бинарника в минимальный базовый образ `alpine` или `scratch`.
   - Размер итогового контейнера снижается с 1 ГБ до ~20 МБ, исключая уязвимости исходного кода и инструментария компилятора в проде.

2. **Безопасность контейнеров:**
   - Никогда не запускайте контейнеры под `root`. Создавайте непривилегированного пользователя:
     ```dockerfile
     RUN addgroup -S appgroup && adduser -S appuser -G appgroup
     USER appuser
     ```

3. **Методология 12 Factor App:**
   - **III. Конфигурация:** хранение настроек в переменных окружения (`os.Getenv`), а не в коде.
   - **VI. Процессы:** stateless-приложение (все сохраняемое состояние выносится в кэш или БД).
   - **IX. Утилизируемость:** быстрая остановка через Graceful Shutdown (`SIGTERM`/`SIGINT`) с закрытием сетевых слушателей и транзакций.
   - **XI. Логи:** логирование в `stdout` / `stderr` в структурированном виде (JSON).

4. **Kubernetes Healthchecks и оркестрация:**
   - **Liveness Probe (`/healthz`):** проверяет, завис ли процесс; перезапускает под при отсутствии ответа.
   - **Readiness Probe (`/ready`):** проверяет, готов ли под принимать клиентский трафик (подключение к БД установлено).
   - **Service Load Balancer:** распределяет входящие запросы между подами Deployment (в демонстрационном `microservice` запрос к `/` возвращает имя конкретного пода).

---

## 🚀 Запуск тестов и бенчмарков

```bash
# Тесты домашнего задания (URL Shortener)
cd go_course_4
go test -v -race ./lessons/18-microservices/homework-18/...

# Запуск с бенчмарками
go test -v -race -bench=. ./lessons/18-microservices/homework-18/...

# Тесты демонстрационного микросервиса K8s
cd 18-microservices/microservice
go test -v -race -bench=. ./pkg/api/...
```

---

## 🔗 Файлы и связанные материалы
- [homework-18/README.md](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/18-microservices/homework-18/README.md) — Домашнее задание №18: URL Shortener микросервис
- [homework-18/Dockerfile](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/18-microservices/homework-18/Dockerfile) — Multi-stage production Dockerfile
- [homework-18/docker-compose.yml](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/18-microservices/homework-18/docker-compose.yml) — Docker Compose оркестрация
- [microservice/build/deployment.yaml](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/18-microservices/microservice/build/deployment.yaml) — Манифест Kubernetes Deployment
- [microservice/build/service.yaml](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/18-microservices/microservice/build/service.yaml) — Манифест Kubernetes Service
- **Теория:** [Том 5: Тестирование и архитектура (Раздел 23: Микросервисы, Docker, K8s)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/05-go-core-testing-arch.md)
