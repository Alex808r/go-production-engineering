# Mermaid Diagrams — Стандарт архитектурной визуализации Go Course 4

В данном каталоге собраны стандарты, шаблоны и скилл ([`SKILL.md`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/go-documentation/mermaid-diagrams/SKILL.md)) для построения качественных архитектурных диаграмм в формате **Mermaid** в рамках курса `Go Course 4`.

---

## 📋 Структура каталога

```text
mermaid-diagrams/
├── SKILL.md     # Скилл агента с полным описанием синтаксиса, правил стилей и палитры
└── README.md    # Данный справочник с примерами и каталогом диаграмм курса
```

---

## 🎨 Стандарты оформления и палитра

### 1. Цветовая палитра для компонентов (Flowchart)

| Роль узла | Назначение | fill | stroke | color |
| :--- | :--- | :--- | :--- | :--- |
| **Client / Entrypoint** | Браузер, curl, клиентский шлюз | `#2c3e50` | `#34495e` | `#ffffff` |
| **Core Service** | Основной backend-сервис (Go) | `#27ae60` | `#229954` | `#ffffff` |
| **Secondary / Cache** | Кэш (Redis), вспомогательные микросервисы | `#3498db` | `#2980b9` | `#ffffff` |
| **Database / Storage** | Персистентные СУБД (PostgreSQL) | `#8e44ad` | `#7d3c98` | `#ffffff` |
| **Queue / Broker** | Брокеры сообщений (Apache Kafka, RabbitMQ) | `#e67e22` | `#d35400` | `#ffffff` |
| **Warning / Trap** | Проблемные зоны, антипаттерны, дедлоки | `#e74c3c` | `#c0392b` | `#ffffff` |

### 2. Типографика и форматирование узлов
- Переносы строк внутри блоков выполняются строго через `<br/>` (не `\n`).
- Заголовки оформляются жирным шрифтом с подписью: `["<b>Имя сервиса</b><br/>порт :8080"]`.
- Подписи узлов и действий — на русском языке, технические идентификаторы — в оригинале.

---

## 🔍 Каталог примеров из курса Go Course 4

### 1. Архитектура Cache-Aside ([`20-NoSQL`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/20-NoSQL/20-NoSQL.md) и [`lynks`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lynks/README.md))

```mermaid
sequenceDiagram
    autonumber
    actor Client as Клиент / Браузер
    participant S as Shortener (:8080)
    participant M as Memcache (:8082)
    participant PG as PostgreSQL (:5432)

    Client->>S: GET /{short} (редирект)
    S->>M: GET /api/v1/cache/{short}
    alt Cache Hit (в кэше найдено)
        M-->>S: 200 OK + destination_url
        S-->>Client: 302 Found (редирект)
    else Cache Miss (в кэше нет или TTL истёк)
        M-->>S: 404 Not Found
        S->>PG: SELECT destination FROM links WHERE short = $1
        PG-->>S: destination_url
        S->>M: POST /api/v1/cache (short, destination, TTL=24h)
        S-->>Client: 302 Found (редирект)
    end
```

---

### 2. Топология микросервисов и баз данных ([`18-microservices`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/18-microservices/18-microservices.md))

```mermaid
flowchart TD
    Client["<b>Внешний клиент</b><br/>HTTP REST / JSON"]

    subgraph AppNetwork["Docker Bridge Network (lynks)"]
        S["<b>Lynks Shortener</b><br/>Go microservice (:8080)"]
        M["<b>Lynks Memcache</b><br/>Go cache service (:8082)"]
        PG[("<b>PostgreSQL 15</b><br/>Таблица links (:5432)")]
        R[("<b>Redis 7</b><br/>In-memory кэш (:6379)")]
    end

    Client ==>|"POST / (сократить)<br/>GET /{short} (перейти)"| S
    S -.->|"Cache-Aside"| M
    M ==>|"GET / SET (TTL 24h)"| R
    S ==>|"Fallback SELECT / INSERT"| PG

    style Client fill:#2c3e50,stroke:#34495e,stroke-width:3px,color:#fff
    style S fill:#27ae60,stroke:#229954,stroke-width:2px,color:#fff
    style M fill:#3498db,stroke:#2980b9,stroke-width:2px,color:#fff
    style PG fill:#8e44ad,stroke:#7d3c98,stroke-width:2px,color:#fff
    style R fill:#e67e22,stroke:#d35400,stroke-width:2px,color:#fff
```

---

### 3. Модель данных сущностей реляционной БД ([`15-sql`](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/15-sql/15-sql.md))

```mermaid
erDiagram
    DIRECTORS ||--o{ MOVIES : "снимает"
    MOVIES ||--o{ REVIEWS : "получает"

    DIRECTORS {
        bigserial id PK
        varchar(100) name
        date birth_date
        varchar(100) country
    }

    MOVIES {
        bigserial id PK
        bigint director_id FK
        varchar(255) title
        integer release_year
        numeric rating
    }

    REVIEWS {
        bigserial id PK
        bigint movie_id FK
        varchar(100) author
        text content
        integer score
    }
```

---

## 🔗 Модули с диаграммами в проекте
- [15-sql — Реляционные базы данных и ER-диаграммы](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/15-sql/15-sql.md)
- [17-system-design — Чистая архитектура и слои Clean Architecture](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/17-system-design/17-system-design.md)
- [18-microservices — Контейнеризация и топология микросервисов](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/18-microservices/18-microservices.md)
- [19-queue — Потоки сообщений Apache Kafka](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/19-queue/19-queue.md)
- [20-NoSQL — Архитектура кэширования Cache-Aside](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/20-NoSQL/20-NoSQL.md)
- [lynks — Архитектура итогового проекта](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lynks/README.md)
