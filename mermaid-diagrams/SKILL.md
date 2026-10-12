---
name: mermaid-diagrams
description: Создание и редактирование Mermaid-диаграмм в документации Go Course 4. Используй при упоминании диаграмм, схем, sequence diagram, flowchart, graph, erDiagram, нарисуй, сделай диаграмму, sequence, схема взаимодействия, диаграмма классов, диаграмма потока, er-диаграмма, mermaid, нарисуй схему, покажи поток, визуализируй архитектуру.
metadata:
  allowed_tools:
    - Read
    - Edit
    - Write
    - Grep
    - Glob
---

# Mermaid Diagrams Skill

Создание диаграмм в стиле проекта Go Course 4. Все паттерны извлечены из реальных диаграмм проекта.

## Общие правила

- Перенос строки в нодах: `<br/>` (не `\n`)
- Текст в нодах с разметкой: `["<b>Заголовок</b><br/>подпись"]`
- Все подписи на русском, технические идентификаторы — как есть
- Всегда использовать Mermaid — никаких других инструментов без явного запроса

---

## flowchart / graph

### Когда использовать

- Зависимости между компонентами системы
- Граф задач и их связей
- Архитектурные блок-схемы

### Палитра цветов (style-директивы)

| Роль                         | fill      | stroke    | color  |
| ---------------------------- | --------- | --------- | ------ |
| Клиент / точка входа         | `#2c3e50` | `#34495e` | `#fff` |
| Хардкод / проблема           | `#e74c3c` | `#c0392b` | `#fff` |
| Предупреждение / не критично | `#e67e22` | `#d35400` | `#fff` |
| Основной backend (smb2)      | `#27ae60` | `#229954` | `#fff` |
| Вторичный компонент          | `#3498db` | `#2980b9` | `#fff` |
| Партнёр / ЛК                 | `#f39c12` | `#d68910` | `#000` |
| Legacy (smb)                 | `#9b59b6` | `#8e44ad` | `#fff` |
| Вспомогательный (светлый)    | `#c8e6c9` | `#2e7d32` | `#000` |

### stroke-width

- Главный узел: `stroke-width:3px` или `4px`
- Обычный узел: `stroke-width:2px`

### Пример (из BASELINE.md §5)

```mermaid
flowchart TD
    Customer["<b>Клиент</b><br/>заказывает 199 дней"]

    subgraph smb2["smb2 (Rails 7.1)"]
        S2_MODEL["<b>OurSslCertificate#duration</b><br/>только целые годы<br/>Date.years2days(N)"]
        S2_API["GlobalSign API клиент<br/>найден: providers/global_sign/<br/>хардкод Months=duration/30"]
        S2_MODEL --> S2_PRICE
    end

    Customer -->|"заказ"| smb2

    style Customer fill:#2c3e50,stroke:#34495e,stroke-width:4px,color:#fff
    style S2_MODEL fill:#e74c3c,stroke:#c0392b,stroke-width:2px,color:#fff
    style S2_API   fill:#e67e22,stroke:#d35400,stroke-width:2px,color:#fff
```

### Рёбра

| Тип             | Синтаксис                       |
| --------------- | ------------------------------- | --------- | --- |
| Основной поток  | `==>` или `-->` с `\|"текст"\|` |
| Вторичный поток | `-.->` (пунктир)                |
| С подписью      | `-->                            | "подпись" | `   |

---

## sequenceDiagram

### Когда использовать

- Взаимодействие объектов во времени (Golang методы, HTTP, SOAP)
- Показ стека вызовов с activation boxes
- Разграничение слоёв абстракции

### Обязательный init-блок (тёмные рамки)

```
%%{init: {"themeVariables": {
  "activationBorderColor": "#1a252f",
  "labelBoxBorderColor": "#1a252f"
}}}%%
```

### Участники

- Объект — один participant. Self-call (`Svc->>Svc:`) — единственный способ показать внутренние вызовы.
- Не создавать двух участников для одного класса.
- Псевдонимы: `participant Svc as GlobalSign::Service`

### Activation boxes (границы методов)

- `activate` / `deactivate` обрамляют тело метода
- Вложенные активации поддерживаются — показывают стек вызовов
- Self-call стрелки с именем метода (`Svc->>Svc: order_request_parameter`) идут перед `activate` — это и есть подпись к activation box
- Возвраты: `-->>` (пунктир)

### Слои абстракции (rect)

Используй `rect rgba(r, g, b, 0.22)` + `Note over` для выделения слоёв:

| Слой                | Цвет                       |
| ------------------- | -------------------------- |
| Публичный интерфейс | `rgba(52, 73, 94, 0.22)`   |
| Построение запроса  | `rgba(41, 128, 185, 0.22)` |
| Доменная модель     | `rgba(39, 174, 96, 0.22)`  |
| Внешний API         | `rgba(230, 126, 34, 0.22)` |

**Важно**: `rect` в sequenceDiagram не имеет отдельного свойства для цвета рамки — `labelBoxBorderColor` действует только на `loop`/`alt`, не на `rect`. Граница слоя создаётся контрастом заливок.

### Ограничения Mermaid (sequenceDiagram)

- Activation boxes не поддерживают встроенные подписи — имя метода показывается стрелкой перед `activate`
- `rect` не имеет отдельного stroke — только rgba fill
- `rgba` поддерживается в Mermaid v10+

### Пример (из 02-Sequence-Diagrams.md §4.3)

```mermaid
%%{init: {"themeVariables": {
  "activationBorderColor": "#1a252f",
  "labelBoxBorderColor": "#1a252f"
}}}%%
sequenceDiagram
    participant Caller
    participant Svc as GlobalSign::Service
    participant OSC as OurSslCertificate
    participant GS as GlobalSign SOAP API

    rect rgba(52, 73, 94, 0.22)
        Note over Caller,GS: публичный интерфейс сервиса
        Caller->>Svc: dv_issue / ov_issue / ev_issue
        activate Svc

        rect rgba(41, 128, 185, 0.22)
            Note over Svc,OSC: построение SOAP-запроса
            Svc->>Svc: order_request_parameter
            activate Svc

            rect rgba(39, 174, 96, 0.22)
                Note over Svc,OSC: доменная модель
                Svc->>Svc: validity_months
                activate Svc
                Svc->>OSC: effective_validity_days
                activate OSC
                OSC-->>Svc: validity_days || duration
                deactivate OSC
                Svc-->>Svc: ceil(days / 30.0).to_s
                deactivate Svc
            end

            Svc-->>Svc: {"ValidityPeriod" => {"Months" => "7"}}
            deactivate Svc
        end

        rect rgba(230, 126, 34, 0.22)
            Note over Svc,GS: внешний API
            Svc->>GS: SSLOrder SOAP (Months: "7")
            activate GS
            Note over GS: cap to 199 days server-side
            GS-->>Svc: OrderResponse
            deactivate GS
        end

        Svc-->>Caller: order_id
        deactivate Svc
    end
```

---

## erDiagram

### Когда использовать

- Модель данных ActiveRecord
- Связи между таблицами

### Пример (из IMPLEMENTATION_PLAN.md)

```mermaid
erDiagram
    OurSslCertificate {
        integer id
        integer duration
        integer validity_days "nullable"
        boolean will_be_renewed "default: false"
        date valid_from
        date valid_to
        integer prolonged_id "FK nullable"
    }

    SslYearPrice {
        integer id
        integer year
        boolean buy
    }

    OurSslCertificate ||--o{ SslYearPrice : "через nomenclature"
    OurSslCertificate ||--o| OurSslCertificate : "prolonged_id"
```

---

## timeline

### Когда использовать

- Хронология событий, расписания изменений

### Пример (из BASELINE.md §1)

```mermaid
timeline
    title Жизненный цикл заказа GlobalSign (новая модель)
    Заказ клиента : Выпуск сертификата 1
                  : Срок действия 199 дней
    День 169 (-30) : Инициация renewal
    День 199       : Истечение сертификата 1
                   : Выпуск сертификата 2
    День 398       : Истечение сертификата 2
```

---

## Процесс создания диаграммы

1. Определи тип: поток → `flowchart`; взаимодействие объектов → `sequenceDiagram`; модель данных → `erDiagram`; хронология → `timeline`
2. Для sequenceDiagram: добавь init-блок с тёмными рамками
3. Для flowchart: применяй style-директивы из палитры по смысловой роли узла
4. Перенос строк в нодах: только `<br/>`
5. Русский текст в подписях, английские идентификаторы — как есть

## Источники диаграмм в проекте

- [doc/ssl/BASELINE.md](../../../doc/ssl/BASELINE.md) — flowchart, graph LR, timeline
- [doc/ssl/SMB2-455_TDR.md](../../../doc/ssl/SMB2-455_TDR.md) — sequenceDiagram с rect-слоями
- [doc/ssl/IMPLEMENTATION_PLAN.md](../../../doc/ssl/IMPLEMENTATION_PLAN.md) — erDiagram, flowchart зависимостей фаз
- [doc/architecture/ddd_framework/](../../../doc/architecture/ddd_framework/) — graph LR архитектуры, classDiagram, sequenceDiagram DDD-потока
- [doc/tech-debt/SMB2-http-streaming-proxy.md](../../../doc/tech-debt/SMB2-http-streaming-proxy.md) — sequenceDiagram с DDD-потоком
