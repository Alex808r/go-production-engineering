# 🗄️ Задачник по Live SQL (Топ-33 практических задач с собеседований: Junior, Middle, Senior)

Практическое руководство для подготовки к написанию SQL-запросов в реальном времени (Live SQL Coding) на позиции Junior, Middle и Senior Backend Go Developer.

> [!NOTE]
> Все примеры написаны для **PostgreSQL** (диалектные особенности отмечены). Чтобы потренироваться, поднимите локальную БД по инструкции из [Тома 4, раздел 12.9 (Docker Compose)](04-go-core-backend-web.md#129-docker-compose-для-postgresql-в-локальной-разработке), создайте таблицы из условия задач и наполните их тестовыми данными.
> 
> 💡 **Рекомендация по уровню подготовки:**
> - Для **Junior и Middle разработчиков** (а также при переходе на Go/Postgres с других стеков) рекомендуется начать с **[Раздела 7 (Базовый SQL и фундамент реляционных БД)](#7-базовый-sql-и-фундамент-реляционных-бд-junior--middle-практикум)**, **[Раздела 8 (Скрининг Ozon)](#8-задачи-со-скрининга-и-проектирования-бд-ozon-technical-screening)** и **[Раздела 5 (Все виды JOIN)](#5-все-виды-join-и-каверзные-ловушки)** — они охватывают группировку, ловушки `NULL`, пагинацию в API, блокировки и `UPSERT`.
> - Для позиций **Middle+ и Senior** ключевыми являются **[Раздел 1 (Оконные функции)](#1-оконные-функции-и-ранжирование)**, **[Раздел 3 (Временные ряды)](#3-аналитика-временных-рядов-и-накопительные-итоги)** и **[Раздел 6 (EXPLAIN ANALYZE и оптимизация)](#6-разбор-explain-analyze-и-live-оптимизация)**.

> [!TIP]
> 📚 **Связанные материалы и практический код:**  
> - 📖 **[README.md (Главное оглавление базы знаний по Go)](README.md)** — язык Go, память и рантайм; краткая справка — [go-core-cheatsheet.md](go-core-cheatsheet.md).  
> - 🗄️ **[go-database-interview-guide.md (Гайд по Базам Данных)](go-database-interview-guide.md)** — PostgreSQL, `pgxpool`, ACID, MVCC, индексы, EXPLAIN ANALYZE и 30 вопросов.  
> - 💻 **[go-livecoding-guide.md (Live-Coding на Go)](go-livecoding-guide.md)** — 52 задачи по многопоточности, структурам данных и алгоритмам в Go (Junior, Middle, Senior).  
> - 🐧 **[go-os-and-linux-interview-guide.md (Гайд по ОС, Ядру Linux и Сетям)](go-os-and-linux-interview-guide.md)** — виртуальная память, сисколы, futex, файловые дескрипторы, OOM, процессы vs потоки, HTTP, балансировка и прокси.  
> - 🏛️ **[go-system-design-guide.md (System Design для Go)](go-system-design-guide.md)** — 4-шаговый фреймворк, расчеты и 4 кейса архитектуры.  
> - 🎯 **[homework-tasks.md (Практический задачник Goflex)](homework-tasks.md)** — постановка и разбор всех 32 домашних заданий.  
> - 🗺️ **[course-codebase-guide.md (Путеводитель по кодовой базе)](course-codebase-guide.md)** — сквозная карта курса, лекционные примеры (00–22), проекты GoSearch, Lynks и Wallet (1000+ RPS), ДЗ (homework-02–20).  
> - 💻 **Код лекций:** [15-sql](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/15-sql) (схемы DDL и аналитические запросы SQL), [16-db-apps](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/16-db-apps) (интерфейсы БД и репозитории).  
> - 📁 **Решения ДЗ и проекты:** [homework-15](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/15-sql/homework-15) (схема БД кинотеатра), [homework-16](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/lessons/16-db-apps/homework-16) (PostgreSQL репозиторий фильмов), [wallet](file:///Users/user/workspace/golang_projects/go-core-4/go-production-engineering/wallet) (High-Load сервис кошелька: pgxpool, транзакции, атомарный UPDATE).  

---

## 📑 Оглавление

1. [Оконные функции и Ранжирование](#1-оконные-функции-и-ранжирование)
   - [Задача 1: N-я максимальная зарплата в компании](#задача-1-n-я-максимальная-зарплата-в-компании)
   - [Задача 2: Вторая по величине зарплата в каждом отделе](#задача-2-вторая-по-величине-зарплата-в-каждом-отделе)
   - [Задача 3: Топ-3 самых продаваемых товаров в каждой категории](#задача-3-топ-3-самых-продаваемых-товаров-в-каждой-категории)
2. [Дедупликация и Фильтрация](#2-дедупликация-и-фильтрация)
   - [Задача 4: Поиск дубликатов записей в таблице](#задача-4-поиск-дубликатов-записей-в-таблице)
   - [Задача 5: Удаление дубликатов с сохранением минимального ID](#задача-5-удаление-дубликатов-с-сохранением-минимального-id)
   - [Задача 6: Пользователи без заказов (LEFT JOIN vs NOT EXISTS vs NOT IN)](#задача-6-пользователи-без-заказов-left-join-vs-not-exists-vs-not-in)
3. [Аналитика Временных Рядов и Накопительные Итоги](#3-аналитика-временных-рядов-и-накопительные-итоги)
   - [Задача 7: Накопительный итог (Running Total / Cumulative Sum)](#задача-7-накопительный-итог-running-total--cumulative-sum)
   - [Задача 8: Скользящее среднее за 7 дней (Moving Average)](#задача-8-скользящее-среднее-за-7-дней-moving-average)
   - [Задача 9: Разница с предыдущим днем (Delta через LAG)](#задача-9-разница-с-предыдущим-днем-delta-через-lag)
   - [Задача 10: Непрерывные цепочки входа (Gaps and Islands Problem)](#задача-10-непрерывные-цепочки-входа-gaps-and-islands-problem)
4. [Иерархии, Самообъединения и Продвинутые запросы](#4-иерархии-самообъединения-и-продвинутые-запросы)
   - [Задача 11: Сотрудники, получающие больше своего руководителя (Self Join)](#задача-11-сотрудники-получающие-больше-своего-руководителя-self-join)
   - [Задача 12: Рекурсивное дерево категорий (WITH RECURSIVE)](#задача-12-рекурсивное-дерево-категорий-with-recursive)
   - [Задача 13: Клиенты, купившие товары во ВСЕХ категориях](#задача-13-клиенты-купившие-товары-во-всех-категориях)
   - [Задача 14: Расчет Retention 1-го дня (Когортный анализ)](#задача-14-расчет-retention-1-го-дня-когортный-анализ)
   - [Задача 15: Агрегация и поиск по JSONB-документам](#задача-15-агрегация-и-поиск-по-jsonb-документам)
5. [Все виды JOIN и каверзные ловушки](#5-все-виды-join-и-каверзные-ловушки)
   - [Задача 16: Ловушка фильтрации при LEFT JOIN (ON vs WHERE)](#задача-16-ловушка-фильтрации-при-left-join-on-vs-where)
   - [Задача 17: Декартово произведение и правильный подсчет агрегатов при связях 1:N](#задача-17-декартово-произведение-и-правильный-подсчет-агрегатов-при-связях-1n)
   - [Задача 18: FULL OUTER JOIN для сверки и поиска расхождений между реестрами](#задача-18-full-outer-join-для-сверки-и-поиска-расхождений-между-реестрами)
6. [Разбор EXPLAIN ANALYZE и Live-Оптимизация](#6-разбор-explain-analyze-и-live-оптимизация)
   - [Задача 19: Диагностика реального плана EXPLAIN (ANALYZE, BUFFERS) и устранение узкого места](#задача-19-диагностика-реального-плана-explain-analyze-buffers-и-устранение-узкого-места)
   - [Задача 20: Устранение Sort Method: external merge Disk через индексы и настройку памяти сортировки](#задача-20-устранение-sort-method-external-merge-disk-через-индексы-и-настройку-памяти-сортировки)
7. [Базовый SQL и фундамент реляционных БД (Junior / Middle Практикум)](#7-базовый-sql-и-фундамент-реляционных-бд-junior--middle-практикум)
   - [Задача 21: Фильтрация агрегатов: WHERE vs HAVING и группировка (GROUP BY)](#задача-21-фильтрация-агрегатов-where-vs-having-и-группировка-group-by)
   - [Задача 22: Трёхзначная логика, ловушки NULL и оператор IS DISTINCT FROM](#задача-22-трёхзначная-логика-ловушки-null-и-оператор-is-distinct-from)
   - [Задача 23: Условная агрегация (Pivot через CASE WHEN и PostgreSQL FILTER)](#задача-23-условная-агрегация-pivot-через-case-when-и-postgresql-filter)
   - [Задача 24: Пагинация в высоконагруженных API: OFFSET vs Keyset Pagination (Cursor)](#задача-24-пагинация-в-высоконагруженных-api-offset-vs-keyset-pagination-cursor)
   - [Задача 25: Конкурентная обработка и блокировки строк: SELECT FOR UPDATE и SKIP LOCKED](#задача-25-конкурентная-обработка-и-блокировки-строк-select-for-update-и-skip-locked)
   - [Задача 26: Атомарный UPSERT (INSERT ... ON CONFLICT DO UPDATE)](#задача-26-атомарный-upsert-insert--on-conflict-do-update)
   - [Задача 27: Заполнение пропусков во временных рядах через GENERATE_SERIES](#задача-27-заполнение-пропусков-во-временных-рядах-через-generate_series)
   - [Задача 28: Индексируемые условия дат (SARGable запросы vs DATE_TRUNC)](#задача-28-индексируемые-условия-дат-sargable-запросы-vs-date_trunc)
   - [Задача 29: Поиск смежных свободных мест (Consecutive Available Seats / Кинотеатр)](#задача-29-поиск-смежных-свободных-мест-consecutive-available-seats--кинотеатр)
8. [Задачи со Скрининга и Проектирования БД (Ozon Technical Screening)](#8-задачи-со-скрининга-и-проектирования-бд-ozon-technical-screening)
   - [Задача 30: Группировка и фильтрация заказов (Ozon Screening: Orders)](#задача-30-группировка-и-фильтрация-заказов-ozon-screening-orders)
   - [Задача 31: Покупки до блокировки пользователя и фильтрация сумм (Ozon Screening: Purchases & Ban List)](#задача-31-покупки-до-блокировки-пользователя-и-фильтрация-сумм-ozon-screening-purchases--ban-list)
   - [Задача 32: Моделирование архитектуры чатов и сообщений (Ozon RDBMS: Chats & Messages)](#задача-32-моделирование-архитектуры-чатов-и-сообщений-ozon-rdbms-chats--messages)
   - [Задача 33: Моделирование каталога библиотеки (Ozon RDBMS: Library Catalog)](#задача-33-моделирование-каталога-библиотеки-ozon-rdbms-library-catalog)

---

## 1. Оконные функции и Ранжирование

### Задача 1: N-я максимальная зарплата в компании

#### 🎯 Условие:
Дана таблица `employees (id INT, name VARCHAR, salary INT)`. Напишите запрос, возвращающий вторую (или $N$-ю) по величине уникальную зарплату. Если второй зарплаты нет, запрос должен вернуть `NULL`.

#### ❌ Наивное решение с ловушкой:
```sql
-- ❌ ОШИБКА: Если у двух людей одинаковая максимальная зарплата (например, 100k и 100k), 
-- OFFSET 1 вернет те же 100k вместо реального второго уровня зарплат!
SELECT salary FROM employees ORDER BY salary DESC LIMIT 1 OFFSET 1;
```

#### ✅ Решение 1: Самое простое и быстрое (LeetCode Standard, $O(1)$ по индексу)
```sql
SELECT (
    SELECT DISTINCT salary 
    FROM employees 
    ORDER BY salary DESC 
    LIMIT 1 OFFSET 1 -- Для N-й максимальной зарплаты укажите: OFFSET (N - 1)
) AS second_highest_salary;
```

* **Почему это самое оптимальное решение на live-coding:**
  1. **Всего 5 строк кода:** пишется за 20 секунд без громоздких CTE.
  2. **Мгновенная работа ($O(1)$):** при наличии индекса по `salary` планировщик выполняет `Index Scan` + `Limit` и читает ровно 2 кортежа индекса, не сканируя всю таблицу.
  3. **Автоматический `NULL`:** внешний скалярный подзапрос `SELECT (...)` по стандарту SQL возвращает `NULL`, если второй зарплаты нет в базе (не нужны `COALESCE` или костыли).
  4. **Легко масштабируется:** для 3-й зарплаты ставим `OFFSET 2`, для $N$-й — `OFFSET (N - 1)`.

#### 💡 Решение 2: Через оконную функцию `DENSE_RANK()` (Для вывода строк сотрудников)
```sql
WITH RankedSalaries AS (
    SELECT 
        salary,
        DENSE_RANK() OVER (ORDER BY salary DESC) AS rank_num
    FROM employees
)
SELECT (
    SELECT salary FROM RankedSalaries WHERE rank_num = 2 LIMIT 1
) AS second_highest_salary;
```

* **Когда предпочесть `DENSE_RANK()`:** В отличие от `ROW_NUMBER()` (который просто нумерует строки) и `RANK()` (который оставляет пропуски `1, 1, 3`), `DENSE_RANK()` ранжирует без пропусков (`1, 1, 2`). Это решение незаменимо, если нужно вывести не только саму сумму, но и имя сотрудника или сгруппировать по отделам (как в Задаче 2).

---

### Задача 2: Вторая по величине зарплата в каждом отделе

#### 🎯 Условие:
Дана таблица `employees (id, name, department_id, salary)`. Найдите сотрудника со второй по величине зарплатой в каждом отделе.

#### ✅ Эталонное решение:
```sql
WITH RankedEmployees AS (
    SELECT 
        id,
        name,
        department_id,
        salary,
        DENSE_RANK() OVER (
            PARTITION BY department_id 
            ORDER BY salary DESC
        ) AS salary_rank
    FROM employees
)
SELECT department_id, id, name, salary
FROM RankedEmployees
WHERE salary_rank = 2;
```

* **План и Индекс:** Для ускорения необходим составной индекс:  
  `CREATE INDEX idx_emp_dept_salary ON employees (department_id, salary DESC);`

---

### Задача 3: Топ-3 самых продаваемых товаров в каждой категории

#### 🎯 Условие:
Таблицы: `products (id, category_id, name)` и `order_items (id, product_id, quantity)`. Найти топ-3 товара с наибольшим суммарным объемом продаж в каждой категории.

#### ✅ Эталонное решение:
```sql
WITH ProductSales AS (
    SELECT 
        p.category_id,
        p.id AS product_id,
        p.name AS product_name,
        COALESCE(SUM(oi.quantity), 0) AS total_sold,
        ROW_NUMBER() OVER (
            PARTITION BY p.category_id 
            ORDER BY COALESCE(SUM(oi.quantity), 0) DESC
        ) AS sales_rank
    FROM products p
    LEFT JOIN order_items oi ON p.id = oi.product_id
    GROUP BY p.category_id, p.id, p.name
)
SELECT category_id, product_id, product_name, total_sold
FROM ProductSales
WHERE sales_rank <= 3
ORDER BY category_id, sales_rank;
```

> ⚠️ `ROW_NUMBER()` при равных продажах выберет «случайного» победителя. Если нужно, чтобы все товары с одинаковыми продажами попадали в топ вместе — используйте `RANK()`/`DENSE_RANK()`; если важен детерминизм — добавьте второй ключ сортировки (`, p.id`).

> [!TIP]
> 💡 **Супер-лайфхак для Топ-1 в каждой категории (PostgreSQL `DISTINCT ON`):**  
> Если на собеседовании просят вывести **не Топ-3, а именно Топ-1** (самый продаваемый товар, последний заказ клиента, последнее сообщение в каждом чате), используйте идиоматичный для PostgreSQL синтаксис `DISTINCT ON` без CTE и без оконных функций:
> ```sql
> SELECT DISTINCT ON (p.category_id)
>     p.category_id,
>     p.id AS product_id,
>     p.name AS product_name,
>     COALESCE(SUM(oi.quantity), 0) AS total_sold
> FROM products p
> LEFT JOIN order_items oi ON p.id = oi.product_id
> GROUP BY p.category_id, p.id, p.name
> ORDER BY p.category_id, total_sold DESC;
> ```
> *Интервьюеры очень ценят знание `DISTINCT ON`: запрос пишется в несколько строк и выполняется оптимизатором за один единственный проход.*

---

## 2. Дедупликация и Фильтрация

### Задача 4: Поиск дубликатов записей в таблице

#### 🎯 Условие:
Таблица `users (id, email, created_at)`. Найти все email-адреса, которые встречаются более одного раза, и количество их повторений.

#### ✅ Эталонное решение:
```sql
SELECT 
    email, 
    COUNT(*) AS duplicates_count
FROM users
WHERE email IS NOT NULL -- Защита: незаполненные email не должны считаться «дубликатами»
GROUP BY email
HAVING COUNT(*) > 1
ORDER BY duplicates_count DESC;
```

> [!NOTE]
> 💡 **Подводный камень с `NULL`:** По стандарту SQL в `GROUP BY` все значения `NULL` объединяются в **одну общую группу**. Если в таблице 10 пользователей без email (`email IS NULL`), запрос без фильтра `WHERE email IS NOT NULL` ошибочно выдаст строку `(NULL, 10)` как «дубликат». Всегда фильтруйте `IS NOT NULL`, если колонка допускает `NULL`.

---

### Задача 5: Удаление дубликатов с сохранением минимального ID

#### 🎯 Условие:
В таблице `users (id, email)` есть дублирующиеся строки с одинаковым email. Удалите дубликаты, сохранив для каждого email только запись с наименьшим `id`.

#### ✅ Решение 1: Самое простое и быстрое в PostgreSQL (`DELETE ... USING`)
```sql
DELETE FROM users a 
USING users b 
WHERE a.email = b.email 
  AND a.id > b.id;
```

* **Почему это эталон в PostgreSQL:**
  1. **Всего 4 строки:** без вложенных подзапросов, без тяжелых CTE и без оконных функций.
  2. **Максимальная производительность:** движок выполняет быстрое соединение по индексу `(email)` и удаляет все строки с `id`, превышающим минимальный, за один проход.
  3. **Если в таблице нет уникального столбца `id`:** для различения дубликатов используют скрытый физический адрес строки `ctid`:
     ```sql
     DELETE FROM users a USING users b WHERE a.email = b.email AND a.ctid > b.ctid;
     ```

#### 💡 Решение 2: Универсальный ANSI SQL (через CTE и `ROW_NUMBER`)
```sql
WITH RankedUsers AS (
    SELECT 
        id,
        ROW_NUMBER() OVER (PARTITION BY email ORDER BY id ASC) AS rn
    FROM users
)
DELETE FROM users
WHERE id IN (
    SELECT id FROM RankedUsers WHERE rn > 1
);
```

* **Когда предпочесть:** В реляционных СУБД без поддержки конструкции `USING` (например, Oracle или стандартном ANSI SQL).

---

### Задача 6: Пользователи без заказов (LEFT JOIN vs NOT EXISTS vs NOT IN)

#### 🎯 Условие:
Таблицы: `users (id, name)` и `orders (id, user_id, amount)`. Найти всех пользователей, которые не сделали ни одного заказа.

#### 💡 Ловушка `NOT IN` со значениями `NULL`:
Если колонка `orders.user_id` содержит хотя бы один `NULL`, запрос `WHERE id NOT IN (SELECT user_id FROM orders)` вернет **0 строк**, так как сравнение с `NULL` возвращает `UNKNOWN`!

#### ✅ Эталонное решение (Оптимально по производительности: `NOT EXISTS`):
```sql
-- Вариант 1: NOT EXISTS (Идеально для PostgreSQL, оптимизатор использует Anti-Join)
SELECT u.id, u.name
FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM orders o WHERE o.user_id = u.id
);

-- Вариант 2: LEFT JOIN ... WHERE IS NULL
SELECT u.id, u.name
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE o.id IS NULL;
```

---

## 3. Аналитика Временных Рядов и Накопительные Итоги

### Задача 7: Накопительный итог (Running Total / Cumulative Sum)

#### 🎯 Условие:
Таблица `payments (id, amount, paid_at DATE)`. Вывести дату, ежедневную выручку и накопительный итог выручки с начала периода по текущий день.

#### ✅ Решение 1: Самое простое и быстрое (в один запрос без CTE)
```sql
SELECT 
    paid_at,
    SUM(amount) AS daily_amount,
    SUM(SUM(amount)) OVER (ORDER BY paid_at) AS running_total
FROM payments
GROUP BY paid_at
ORDER BY paid_at;
```

* **Почему это супер-решение на собеседовании:**
  1. **Всего 7 строк:** пишется мгновенно без создания громоздких CTE.
  2. **Вложенная оконная агрегация `SUM(SUM(...)) OVER (...)`:** стандарт SQL и PostgreSQL позволяют вызывать оконную функцию над результатом обычной групповой агрегации. СУБД сначала группирует строки по `paid_at`, а затем сразу считает кумулятивную сумму по полученным агрегатам за один проход.
  3. **Рамка окна по умолчанию:** при указании `ORDER BY paid_at` окно по стандарту неявно принимает рамку `RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW`, корректно накапливая сумму от начала ряда до текущего дня.

#### 💡 Решение 2: Пошагово через CTE (Классический подход)
```sql
WITH DailyRevenue AS (
    SELECT 
        paid_at,
        SUM(amount) AS daily_amount
    FROM payments
    GROUP BY paid_at
)
SELECT 
    paid_at,
    daily_amount,
    SUM(daily_amount) OVER (
        ORDER BY paid_at 
        ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
    ) AS running_total
FROM DailyRevenue
ORDER BY paid_at;
```
* **Когда предпочесть:** Когда требуется явно разбить расчет на промежуточные шаги или переиспользовать `daily_amount` в других подзапросах.

---

### Задача 8: Скользящее среднее за 7 дней (Moving Average)

#### 🎯 Условие:
Рассчитать 7-дневное скользящее среднее дневной выручки (включая текущий день и 6 предшествующих).

#### ✅ Решение 1: Самое компактное и быстрое (в один запрос без CTE)
```sql
SELECT 
    paid_at,
    SUM(amount) AS daily_amount,
    ROUND(
        AVG(SUM(amount)) OVER (
            ORDER BY paid_at
            ROWS BETWEEN 6 PRECEDING AND CURRENT ROW
        ), 2
    ) AS moving_avg_7d
FROM payments
GROUP BY paid_at
ORDER BY paid_at;
```

* **Почему так:** агрегат `SUM(amount)` рассчитывается на шаге `GROUP BY`, а внешняя оконная функция `AVG(...) OVER (...)` вычисляет скользящее среднее по уже сгруппированным дням в рамках единого плана запроса.

#### 💡 Решение 2: Пошагово через CTE
```sql
WITH DailyRevenue AS (
    SELECT paid_at, SUM(amount) AS daily_amount
    FROM payments
    GROUP BY paid_at
)
SELECT 
    paid_at,
    daily_amount,
    ROUND(AVG(daily_amount) OVER (
        ORDER BY paid_at
        ROWS BETWEEN 6 PRECEDING AND CURRENT ROW
    ), 2) AS moving_avg_7d
FROM DailyRevenue
ORDER BY paid_at;
```

> ⚠️ `ROWS BETWEEN 6 PRECEDING` берёт 6 **предыдущих строк**, а не 6 предыдущих **дней**: если в данных есть дни без платежей (пропуски дат), окно «растянется» на больше чем 7 календарных дней. Для настоящего календарного окна используйте `RANGE BETWEEN INTERVAL '6 days' PRECEDING AND CURRENT ROW` (PostgreSQL 11+) либо сначала дополните ряд недостающими датами через `generate_series`. Также учтите, что в первые 6 дней среднее считается по неполному окну.

---

### Задача 9: Разница с предыдущим днем (Delta через LAG)

#### 🎯 Условие:
Для каждого дня вывести сумму выручки и её изменение по сравнению со вчерашним днем (в процентах).

#### ✅ Решение 1: Самое чистое и читаемое (через CTE)
```sql
WITH DailySales AS (
    SELECT paid_at, SUM(amount) AS current_day_sales
    FROM payments
    GROUP BY paid_at
)
SELECT 
    paid_at,
    current_day_sales,
    LAG(current_day_sales) OVER (ORDER BY paid_at) AS prev_day_sales,
    ROUND(
        (current_day_sales - LAG(current_day_sales) OVER (ORDER BY paid_at)) * 100.0 / 
        NULLIF(LAG(current_day_sales) OVER (ORDER BY paid_at), 0), 
        2
    ) AS growth_percentage
FROM DailySales
ORDER BY paid_at;
```
* **Почему через CTE:** вычисление процента прироста требует повторного обращения к вчерашней выручке. Использование CTE избавляет от дублирования громоздких выражений `SUM(amount)` и `LAG(SUM(amount))`.

#### 💡 Решение 2: В один запрос без CTE (Компактный однопроходный SQL)
```sql
SELECT 
    paid_at,
    SUM(amount) AS current_day_sales,
    LAG(SUM(amount)) OVER (ORDER BY paid_at) AS prev_day_sales,
    ROUND(
        (SUM(amount) - LAG(SUM(amount)) OVER (ORDER BY paid_at)) * 100.0 / 
        NULLIF(LAG(SUM(amount)) OVER (ORDER BY paid_at), 0), 
        2
    ) AS growth_percentage
FROM payments
GROUP BY paid_at
ORDER BY paid_at;
```

---

### Задача 10: Непрерывные цепочки входа (Gaps and Islands Problem)

#### 🎯 Условие:
Таблица `user_activity (user_id, login_date DATE)`. Найти максимальную длину непрерывной серии дней (стрик), в которые пользователь заходил в приложение каждый день подряд.

#### ✅ Эталонное решение (Паттерн «Острова и Пропуски»):
```sql
WITH UniqueLogins AS (
    -- Избавляемся от множественных логинов за один день
    SELECT DISTINCT user_id, login_date 
    FROM user_activity
),
NumberedLogins AS (
    -- Вычисляем разницу между датой и порядковым номером
    -- Для непрерывных дат разность login_date - ROW_NUMBER() будет КОНСТАНТОЙ!
    SELECT 
        user_id,
        login_date,
        login_date - (ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY login_date))::int AS island_group
    FROM UniqueLogins
),
StreakLengths AS (
    SELECT 
        user_id,
        island_group,
        COUNT(*) AS streak_days,
        MIN(login_date) AS streak_start,
        MAX(login_date) AS streak_end
    FROM NumberedLogins
    GROUP BY user_id, island_group
)
SELECT user_id, MAX(streak_days) AS max_streak
FROM StreakLengths
GROUP BY user_id;
```

> [!TIP]
> 💡 **Что если поле времени — `TIMESTAMPTZ` / `TIMESTAMP`?**  
> Если в таблице хранятся точные метки времени (например, `created_at TIMESTAMPTZ`), перед вычислением островков их обязательно нужно привести к чистой календарной дате в CTE `UniqueLogins`:  
> `SELECT DISTINCT user_id, (login_time AT TIME ZONE 'UTC')::date AS login_date`. Иначе вычитание целого числа из `TIMESTAMPTZ` вычтет целые секунды или завершится ошибкой типов.

---

## 4. Иерархии, Самообъединения и Продвинутые запросы

### Задача 11: Сотрудники, получающие больше своего руководителя (Self Join)

#### 🎯 Условие:
Таблица `employees (id, name, salary, manager_id)`. Найти всех сотрудников, зарплата которых выше зарплаты их непосредственного начальника.

#### ✅ Эталонное решение:
```sql
SELECT 
    e.name AS employee_name,
    e.salary AS employee_salary,
    m.name AS manager_name,
    m.salary AS manager_salary
FROM employees e
JOIN employees m ON e.manager_id = m.id
WHERE e.salary > m.salary;
```

---

### Задача 12: Рекурсивное дерево категорий (WITH RECURSIVE)

#### 🎯 Условие:
Таблица `categories (id, name, parent_id)`. Построить полное дерево категорий, вычислив для каждой глубину вложенности `level` и полный хлебный путь `path` (например: `"Электроника -> Смартфоны -> Apple"`).

#### ✅ Эталонное решение:
```sql
WITH RECURSIVE CategoryTree AS (
    -- Якорный запрос: корневые категории (parent_id IS NULL)
    SELECT 
        id,
        name,
        parent_id,
        1 AS level,
        name::TEXT AS path
    FROM categories
    WHERE parent_id IS NULL

    UNION ALL

    -- Рекурсивный шаг: присоединяем дочерние категории
    SELECT 
        c.id,
        c.name,
        c.parent_id,
        ct.level + 1,
        ct.path || ' -> ' || c.name
    FROM categories c
    JOIN CategoryTree ct ON c.parent_id = ct.id
)
SELECT id, name, level, path
FROM CategoryTree
ORDER BY path;
```

> ⚠️ Если в данных есть цикл (`A → B → A`), рекурсия не остановится. Защита: ограничение глубины (`WHERE ct.level < 20`), накопление массива посещённых `id` или, в PostgreSQL 14+, предложение `CYCLE id SET is_cycle USING path_ids`.

---

### Задача 13: Клиенты, купившие товары во ВСЕХ категориях

#### 🎯 Условие:
Таблицы: `customers (id, name)`, `orders (id, customer_id)`, `order_items (order_id, product_id)`, `products (id, category_id)`, `categories (id)`. Найти покупателей, которые купили товары из абсолютно всех существующих категорий.

#### ✅ Эталонное решение (Реляционное деление):
```sql
SELECT 
    c.id, 
    c.name
FROM customers c
JOIN orders o ON c.id = o.customer_id
JOIN order_items oi ON o.id = oi.order_id
JOIN products p ON oi.product_id = p.id
GROUP BY c.id, c.name
HAVING COUNT(DISTINCT p.category_id) = (SELECT COUNT(*) FROM categories);
```

> Нюанс: если в справочнике есть категория, в которой вообще нет товаров, то никто не сможет «купить во всех категориях» — при необходимости считайте только категории, где есть хотя бы один товар (`SELECT COUNT(DISTINCT category_id) FROM products`).

---

### Задача 14: Расчет Retention 1-го дня (Когортный анализ)

#### 🎯 Условие:
Таблица `logins (user_id, login_date DATE)`. Вычислить Retention первого дня (долю пользователей, которые зашли в сервис на следующий день после своего первого визита).

#### ✅ Эталонное решение:
```sql
WITH UserInstall AS (
    -- Дата первой активности (инсталл/регистрация)
    SELECT user_id, MIN(login_date) AS install_date
    FROM logins
    GROUP BY user_id
)
SELECT 
    ui.install_date,
    COUNT(DISTINCT ui.user_id) AS cohort_size,
    COUNT(DISTINCT l.user_id) AS returned_day_1,
    ROUND(
        COUNT(DISTINCT l.user_id) * 100.0 / COUNT(DISTINCT ui.user_id), 
        2
    ) AS retention_d1_pct
FROM UserInstall ui
LEFT JOIN logins l 
    ON ui.user_id = l.user_id 
    AND l.login_date = ui.install_date + 1  -- для типа DATE «+ 1» означает «плюс один день»
GROUP BY ui.install_date
ORDER BY ui.install_date;
```

---

### Задача 15: Агрегация и поиск по JSONB-документам

#### 🎯 Условие:
Таблица `user_profiles (id, metadata JSONB)`. Поле `metadata` имеет структуру:  
`{"roles": ["admin", "editor"], "settings": {"theme": "dark", "notify": true}}`.  
Написать запрос, выбирающий всех пользователей с ролью `"admin"` и включенной темой `"dark"`.

#### ✅ Эталонное решение:
```sql
-- 1. Элегантный поиск через оператор включения @> (поддерживает GIN индекс!):
SELECT id, metadata->'settings'->>'theme' AS current_theme
FROM user_profiles
WHERE metadata @> '{"roles": ["admin"], "settings": {"theme": "dark"}}';

-- 2. Необходимый индекс для моментального поиска в миллионных таблицах:
CREATE INDEX idx_user_profiles_meta ON user_profiles USING GIN (metadata jsonb_path_ops);
```

---

## 5. Все виды JOIN и каверзные ловушки

### Задача 16: Ловушка фильтрации при LEFT JOIN (ON vs WHERE)

#### 🎯 Условие:
Даны таблицы:
- `users (id, name)`
- `orders (id, user_id, status, created_at, amount)`

Напишите запрос, возвращающий список **всех** пользователей и сумму их оплаченных заказов (`status = 'PAID'`) за **2026 год**.  
Если у пользователя не было таких заказов, он **всё равно должен остаться в итоговой выборке** с суммой `0`.

#### ❌ Наивная ошибка (превращение в `INNER JOIN`):
```sql
-- ❌ ОШИБКА: Секция WHERE отфильтрует всех пользователей без заказов!
SELECT 
    u.id, 
    u.name, 
    COALESCE(SUM(o.amount), 0) AS total_paid
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE o.status = 'PAID' 
  AND o.created_at >= '2026-01-01' 
  AND o.created_at < '2027-01-01'
GROUP BY u.id, u.name;
```
* **Почему это баг:** Для пользователей без заказов `LEFT JOIN` сгенерирует поля `orders.*` равными `NULL`. Затем условие `NULL = 'PAID'` вернет `UNKNOWN`, и блок `WHERE` молча выкинет этих пользователей из результата!

#### ✅ Эталонное решение (условие в `ON`):
```sql
SELECT 
    u.id, 
    u.name, 
    COALESCE(SUM(o.amount), 0) AS total_paid
FROM users u
LEFT JOIN orders o 
    ON u.id = o.user_id
    AND o.status = 'PAID'
    AND o.created_at >= '2026-01-01'
    AND o.created_at < '2027-01-01'
GROUP BY u.id, u.name
ORDER BY total_paid DESC;
```
* **Правило:** Условия на правую таблицу в `LEFT JOIN` обязаны находиться в секции `ON`, если мы хотим сохранить строки левой таблицы.

---

### Задача 17: Декартово произведение и правильный подсчет агрегатов при связях 1:N

#### 🎯 Условие:
Даны таблицы:
- `users (id, name)`
- `orders (id, user_id, amount)` (у одного пользователя может быть много заказов)
- `reviews (id, user_id, rating)` (у одного пользователя может быть много отзывов)

Для каждого пользователя необходимо посчитать:
1. Общую сумму всех его заказов (`total_spent`).
2. Общее количество оставленных им отзывов (`total_reviews`).

#### ❌ Наивная ошибка (Fan-out / Раздувание строк):
```sql
-- ❌ ОШИБКА: Декартово произведение между orders и reviews умножит суммы!
SELECT 
    u.id,
    u.name,
    COALESCE(SUM(o.amount), 0) AS total_spent,
    COUNT(r.id) AS total_reviews
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
LEFT JOIN reviews r ON u.id = r.user_id
GROUP BY u.id, u.name;
```
* **Почему это баг:** Если у пользователя 3 заказа и 4 отзыва, объединение таблиц породит $3 \times 4 = 12$ строк! Сумма заказов `SUM(o.amount)` посчитается 4 раза, а `COUNT(r.id)` вернет 12 вместо 4!

#### ✅ Эталонное решение (Предварительная агрегация в CTE):
```sql
WITH UserOrders AS (
    SELECT user_id, SUM(amount) AS total_spent
    FROM orders
    GROUP BY user_id
),
UserReviews AS (
    SELECT user_id, COUNT(*) AS total_reviews
    FROM reviews
    GROUP BY user_id
)
SELECT 
    u.id,
    u.name,
    COALESCE(uo.total_spent, 0) AS total_spent,
    COALESCE(ur.total_reviews, 0) AS total_reviews
FROM users u
LEFT JOIN UserOrders uo ON u.id = uo.user_id
LEFT JOIN UserReviews ur ON u.id = ur.user_id;
```

---

### Задача 18: FULL OUTER JOIN для сверки и поиска расхождений между реестрами

#### 🎯 Условие:
При интеграции с платёжным шлюзом необходимо провести сверку (реконсиляцию) реестра транзакций в нашей базе `app_payments (tx_id, amount)` и выписки из банка `bank_statements (tx_id, amount)`.

Напишите запрос, возвращающий:
1. Транзакции, которые есть у нас, но банк их не видит.
2. Транзакции, которые есть в банке, но их нет у нас.
3. Транзакции, где сумма `amount` в нашей базе расходится с суммой банка.

#### ✅ Эталонное решение:
```sql
SELECT 
    COALESCE(ap.tx_id, bs.tx_id) AS tx_id,
    ap.amount AS app_amount,
    bs.amount AS bank_amount,
    CASE 
        WHEN bs.tx_id IS NULL THEN 'MISSING_IN_BANK'
        WHEN ap.tx_id IS NULL THEN 'MISSING_IN_APP'
        WHEN ap.amount IS DISTINCT FROM bs.amount THEN 'AMOUNT_MISMATCH'
    END AS mismatch_reason
FROM app_payments ap
FULL OUTER JOIN bank_statements bs ON ap.tx_id = bs.tx_id
WHERE ap.tx_id IS NULL 
   OR bs.tx_id IS NULL 
   OR ap.amount IS DISTINCT FROM bs.amount;
```

> [!NOTE]
> 💡 **Почему `IS DISTINCT FROM` вместо `<>`:** Если транзакция присутствует в обоих реестрах, но в одном из них сумма оказалась `NULL` (сбой валидации/незавершенный процессинг), стандартное сравнение `ap.amount <> bs.amount` вернет `UNKNOWN`, и расхождение **молча выпадет из отчета**! Оператор `IS DISTINCT FROM` трактует `NULL` как различимое значение и гарантирует 100% точность сверки.

---

## 6. Разбор EXPLAIN ANALYZE и Live-Оптимизация

### Задача 19: Диагностика реального плана EXPLAIN (ANALYZE, BUFFERS) и устранение узкого места

#### 🎯 Условие:
На собеседовании вам показывают запрос и реальный вывод команды `EXPLAIN (ANALYZE, BUFFERS)`:

```sql
SELECT id, message, created_at 
FROM notifications 
WHERE user_id = 105 AND is_read = false 
ORDER BY created_at DESC 
LIMIT 20;
```

#### Листинг плана:
```text
Limit  (cost=45120.00..45120.05 rows=20 width=48) (actual time=785.120..785.125 rows=12 loops=1)
  Buffers: shared hit=410 read=21800
  ->  Sort  (cost=45120.00..45120.30 rows=120 width=48) (actual time=785.118..785.120 rows=12 loops=1)
        Sort Key: created_at DESC
        Sort Method: quicksort  Memory: 26kB
        Buffers: shared hit=410 read=21800
        ->  Seq Scan on notifications  (cost=0.00..45100.00 rows=120 width=48) (actual time=0.065..784.850 rows=12 loops=1)
              Filter: ((NOT is_read) AND (user_id = 105))
              Rows Removed by Filter: 1999988
              Buffers: shared hit=410 read=21800
Planning Time: 0.185 ms
Execution Time: 785.160 ms
```

#### ❓ Вопросы интервьюера:
1. Что является главным источником медленной работы запроса?
2. Почему показатель `shared read` так велик и к чему это ведет?
3. Какую DDL команду нужно выполнить, чтобы оптимизировать данный запрос до <1 мс?

#### 💡 Эталонный ответ кандидата:
1. **Главная проблема:** Узел `Seq Scan on notifications`. База данных выполняет полное последовательное сканирование таблицы из 2 миллионов строк. Параметр `Rows Removed by Filter: 1999988` показывает, что движок прочитал почти 2 000 000 строк и отбросил их, чтобы найти всего 12 нужных записей.
2. **Дисковый I/O:** `Buffers: shared read=21800` означает, что 21800 страниц не оказались в `shared_buffers` и были запрошены у операционной системы (с диска или из кэша ОС): $21800 \times 8\text{ КБ} \approx 170\text{ МБ}$ данных. Это перегружает подсистему ввода-вывода и занимает основную часть из 785 мс.
3. **Идеальный индекс:** Поскольку непрочитанных уведомлений обычно мало (`is_read = false`), а фильтр идет по `user_id` с сортировкой по `created_at DESC`, наилучшим решением является **частичный составной индекс (Partial Index)**:

```sql
CREATE INDEX idx_notifications_unread_user 
ON notifications (user_id, created_at DESC) 
WHERE is_read = false;
```
* **Результат после добавления индекса:** Метод доступа сменится на `Index Scan` (или `Index Only Scan`), `shared read` упадет до 2–3 буферов, сортировка `Sort` исчезнет (данные уже упорядочены в B-Tree), а время выполнения упадет до **0.05–0.1 мс**.

---

### Задача 20: Устранение Sort Method: external merge Disk через индексы и настройку памяти сортировки

#### 🎯 Условие:
Запрос формирования аналитического отчета сортирует большой массив данных:
```sql
SELECT order_id, customer_id, total_amount, created_at
FROM order_history
WHERE created_at >= '2026-01-01'
ORDER BY total_amount DESC;
```

#### В плане `EXPLAIN (ANALYZE, BUFFERS)` обнаружено:
```text
Sort  (cost=125400.00..128900.00 rows=350000 width=32) (actual time=2450.120..2810.350 rows=350000 loops=1)
  Sort Key: total_amount DESC
  Sort Method: external merge  Disk: 28540kB
  Buffers: shared hit=18500, temp read=3568 written=3570
```

#### ❓ Вопросы интервьюера:
1. Что означает `Sort Method: external merge Disk` и `temp read/written`?
2. Назовите 2 независимых способа оптимизировать эту сортировку.

#### 💡 Эталонный ответ кандидата:
1. **Диагноз:** Памяти, выделенной процессу на сортировку (параметр `work_mem`, по умолчанию в PostgreSQL всего 4 МБ), оказалось недостаточно для удержания 350 000 строк в RAM. Движок был вынужден сбросить промежуточные куски сортировки во временные файлы на диск (`temp read/written`, `Disk: 28540kB`), что привело к сильной деградации времени до 2.8 секунды.
2. **Способ 1 (Тюнинг памяти сессии):**  
   Для тяжелых фоновых аналитических запросов или воркеров можно локально поднять лимит памяти в рамках текущей транзакции:
   ```sql
   SET LOCAL work_mem = '64MB';
   ```
   Сортировка переключится в `quicksort Memory` и выполнится целиком в RAM без обращения к временным файлам на диске.
3. **Способ 2 (Уменьшить объём сортировки):**
   - Если отчёту нужны не все строки, а «топ-N», добавьте `LIMIT N`: PostgreSQL применит алгоритм *top-N heapsort* и удержит в памяти только N лучших строк, не сортируя всё множество (`Sort Method: top-N heapsort Memory: ...`).
   - Выбирайте только необходимые колонки (чем уже строка, тем меньше памяти нужно на сортировку).
   - Для запросов вида «топ-N по `total_amount`» помогает индекс по самому сортируемому столбцу: `CREATE INDEX idx_order_history_amount ON order_history (total_amount DESC);` — планировщик читает строки уже в нужном порядке и останавливается после N подходящих (`Index Scan` + `Limit`, без узла `Sort`).
4. **Ловушка «очевидного» ответа:** составной индекс `(created_at, total_amount DESC)` **не устранит** сортировку в этом запросе. Условие `created_at >= '2026-01-01'` — это **диапазон**, а не равенство, поэтому строки в индексе упорядочены по `created_at`, а внутри диапазона общий порядок по `total_amount` не сохраняется. Индекс убирает `Sort` только тогда, когда ведущие столбцы фильтруются по равенству (`WHERE customer_id = 42 ORDER BY total_amount DESC` → индекс `(customer_id, total_amount DESC)`).

---

## 7. Базовый SQL и фундамент реляционных БД (Junior / Middle Практикум)

Этот раздел содержит практические задачи, составляющие костяк технических собеседований на позиции **Junior / Junior+ / Middle Backend Developer** в компании любого масштаба (от стартапов до Ozon, Т-Банка, Авито и Яндекса). Задачи охватывают фундаментальные принципы реляционной модели: правильное разделение `WHERE` и `HAVING`, трёхзначную логику и ловушки `NULL`, пагинацию в высоконагруженных API, транзакционные блокировки и атомарный `UPSERT`.

---

### Задача 21: Фильтрация агрегатов: WHERE vs HAVING и группировка (GROUP BY)

#### 🎯 Условие:
Дана таблица заказов интернет-магазина `orders (id INT, customer_id INT, order_date DATE, amount NUMERIC, status VARCHAR)`.  
Найти всех клиентов (`customer_id`), которые за последние 30 дней совершили более 3 успешных заказов (`status = 'completed'`) со средней суммой заказа выше 1500 рублей.  
Вывести: `customer_id`, общее количество заказов и среднюю сумму (округлить до 2 знаков), отсортировав по убыванию средней суммы.

#### ❌ Наивное решение с синтаксической ошибкой:
```sql
-- ❌ ОШИБКА: Использование агрегатных функций внутри WHERE запрещено стандартом SQL!
-- ERROR: aggregate functions are not allowed in WHERE
SELECT customer_id, COUNT(*), AVG(amount)
FROM orders
WHERE status = 'completed' 
  AND COUNT(*) > 3 
  AND AVG(amount) > 1500
GROUP BY customer_id;
```

#### 💡 Ловушки и замечания с собеседований:
1. **Фазы выполнения SQL-запроса:**  
   Порядок выполнения запроса движком: `FROM` $\to$ `WHERE` $\to$ `GROUP BY` $\to$ `HAVING` $\to$ `SELECT` $\to$ `ORDER BY` $\to$ `LIMIT`.
   - `WHERE` фильтрует отдельные **строки ДО группировки** (до агрегации), активно используя B-Tree индексы таблицы.
   - `HAVING` фильтрует уже **сформированные группы ПОСЛЕ агрегации**.
2. **Ловушка производительности:**  
   Если поместить условие `status = 'completed'` в `HAVING` (`HAVING status = 'completed' ...`), запрос синтаксически может выполниться (если столбец добавлен в `GROUP BY`), но база данных сначала прочитает и сгруппирует миллионы отмененных и ошибочных заказов, и лишь в самом конце отбросит их! **Строчные предикаты обязаны находиться строго в `WHERE`**.

#### ✅ Эталонное решение:
```sql
SELECT 
    customer_id,
    COUNT(*) AS total_orders,
    ROUND(AVG(amount), 2) AS avg_amount
FROM orders
WHERE status = 'completed' 
  AND order_date >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY customer_id
HAVING COUNT(*) > 3 
   AND AVG(amount) > 1500
ORDER BY avg_amount DESC;
```

---

### Задача 22: Трёхзначная логика, ловушки NULL и оператор IS DISTINCT FROM

#### 🎯 Условие:
Дана таблица пользователей `users (id INT, name VARCHAR, discount INT)`. Поле `discount` хранит процент персональной скидки и может содержать `NULL` (скидка не назначена/неизвестна).  
Требуется выбрать всех пользователей, у которых скидка **не равна 10%** (включая пользователей, у которых скидка равна `NULL`).

#### ❌ Наивная ошибка (потеря данных):
```sql
-- ❌ КАТАСТРОФА: Этот запрос НЕ вернет пользователей с discount IS NULL!
SELECT id, name, discount 
FROM users 
WHERE discount != 10;
```

#### 💡 Разбор трёхзначной логики (Three-Valued Logic):
В реляционных базах данных логические выражения вычисляются не в булеву систему (`TRUE` / `FALSE`), а в трёхзначную: `TRUE`, `FALSE` и `UNKNOWN`.
- Любое сравнение с `NULL` через стандартные операторы (`=`, `!=`, `<>`, `>`, `<`) возвращает `UNKNOWN`.
- Выражение `NULL != 10` дает `UNKNOWN`.
- Предложение `WHERE` оставляет строку **только если предикат строго равен `TRUE`**. Значения `FALSE` и `UNKNOWN` отбрасываются! В итоге все пользователи с `NULL` скидкой молча исчезают из выборки.

#### ✅ Эталонные решения:

**1. Идиоматичный PostgreSQL (оператор `IS DISTINCT FROM`):**
```sql
-- Рекомендуется на собеседованиях по PostgreSQL:
SELECT id, name, discount
FROM users
WHERE discount IS DISTINCT FROM 10;
```
*Почему:* `IS DISTINCT FROM` трактует `NULL` как обычное различимое значение: `NULL IS DISTINCT FROM 10` вернет `TRUE`, а `NULL IS DISTINCT FROM NULL` вернет `FALSE`.

**2. Универсальный переносимый ANSI SQL:**
```sql
SELECT id, name, discount
FROM users
WHERE discount != 10 OR discount IS NULL;
```

**3. Через `COALESCE` (с оговоркой по индексам):**
```sql
SELECT id, name, discount
FROM users
WHERE COALESCE(discount, -1) != 10;
```
> ⚠️ **Предостережение для Senior:** Оборачивание поля в `COALESCE(discount, -1)` отключает стандартный B-Tree индекс по столбцу `discount` (требуется построение функционального индекса по выражению).

---

### Задача 23: Условная агрегация (Pivot через CASE WHEN и PostgreSQL FILTER)

#### 🎯 Условие:
Дана таблица заказов `orders (id INT, seller_id INT, status VARCHAR, amount NUMERIC)`.  
Возможные статусы: `'new'`, `'delivered'`, `'cancelled'`.  
Построить сводную строку аналитики по каждому продавцу (`seller_id`):
1. `total_orders` — общее число заказов;
2. `delivered_count` — число доставленных заказов;
3. `cancelled_count` — число отмененных заказов;
4. `delivered_revenue` — суммарная выручка только по доставленным заказам (если доставленных нет, вернуть `0`).

#### ❌ Наивные решения и ловушки:
- **Антипаттерн со множественными JOIN:** попытка сделать 3 отдельных подзапроса на каждый статус и соединить их по `seller_id` приводит к многократному чтению одной и той же таблицы ($3 \times \text{I/O}$).
- **Ловушка с `COUNT`:** выражение `COUNT(CASE WHEN status = 'delivered' THEN 0 END)` ошибочно посчитает все строки! Функция `COUNT(expr)` считает любые строки, где `expr` не равен `NULL` (а `0` — это не `NULL`). Нужно писать либо `THEN 1 ELSE NULL`, либо использовать `SUM`.

#### ✅ Эталонные решения:

**Вариант 1: Идиоматичный PostgreSQL через конструкцию `FILTER (WHERE ...)` (Золотой стандарт):**
```sql
SELECT 
    seller_id,
    COUNT(*) AS total_orders,
    COUNT(*) FILTER (WHERE status = 'delivered') AS delivered_count,
    COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled_count,
    COALESCE(SUM(amount) FILTER (WHERE status = 'delivered'), 0) AS delivered_revenue
FROM orders
GROUP BY seller_id
ORDER BY delivered_revenue DESC;
```
*Преимущество:* Читаемость на порядок выше, а оптимизатор PostgreSQL вычисляет все фильтры за один единственный проход по таблице.

**Вариант 2: Стандартный SQL через `CASE WHEN`:**
```sql
SELECT 
    seller_id,
    COUNT(*) AS total_orders,
    SUM(CASE WHEN status = 'delivered' THEN 1 ELSE 0 END) AS delivered_count,
    SUM(CASE WHEN status = 'cancelled' THEN 1 ELSE 0 END) AS cancelled_count,
    COALESCE(SUM(CASE WHEN status = 'delivered' THEN amount END), 0) AS delivered_revenue
FROM orders
GROUP BY seller_id
ORDER BY delivered_revenue DESC;
```

---

### Задача 24: Пагинация в высоконагруженных API: OFFSET vs Keyset Pagination (Cursor)

#### 🎯 Условие:
Таблица постов социальной сети `posts (id BIGINT, user_id INT, title TEXT, created_at TIMESTAMPTZ)` насчитывает 20 000 000 строк.  
Бэкенд на Go отдает мобильному приложению бесконечную ленту постов порциями по 20 штук от самых свежих к старым (`ORDER BY created_at DESC, id DESC`).  
Спроектируйте SQL-запрос пагинации, который будет работать стабильно быстро ($<1$ мс) как на первой, так и на 100 000-й странице.

#### ❌ Проблема стандартного `OFFSET`:
```sql
-- Запрос 5000-й страницы в наивном API:
SELECT id, title, created_at
FROM posts
ORDER BY created_at DESC, id DESC
LIMIT 20 OFFSET 100000;
```
*Почему это не работает в проде:*  
Движок PostgreSQL обязан вычислить и отсортировать **100 020 строк**, затем выбросить первые 100 000 и вернуть только 20! Временная сложность такой пагинации — $O(N)$. На глубоких страницах запрос уходит в секунды. Кроме того, если во время чтения пользователь добавит новый пост, смещение сдвинется, и клиент увидит дубликат на следующей странице.

#### ✅ Эталонное решение (Keyset / Cursor Pagination):
Бэкенд передает клиенту в ответе курсор — составное значение последней записи на странице: `(last_created_at, last_id)`. Следующий запрос запрашивает данные строго после курсора:

```sql
-- 1. Первая страница (без курсора):
SELECT id, title, created_at
FROM posts
ORDER BY created_at DESC, id DESC
LIMIT 20;

-- 2. Запрос любой последующей страницы (передаем cursor):
SELECT id, title, created_at
FROM posts
WHERE (created_at, id) < (:last_created_at, :last_id)
ORDER BY created_at DESC, id DESC
LIMIT 20;
```

#### 🚀 Необходимый композитный B-Tree индекс:
```sql
CREATE INDEX idx_posts_keyset ON posts (created_at DESC, id DESC);
```

* **Результат:** Запрос выполняет `Index Scan` сразу по нужной ветке дерева за константное время $O(\log N)$ и читает ровно **20 строк**, независимо от того, листает пользователь 1-ю страницу или миллионную.

> [!IMPORTANT]
> ⚠️ **Критическое требование к полям курсора:** Все столбцы, входящие в курсор (`created_at`, `id`), **обязаны иметь ограничение `NOT NULL`**. Если хотя бы в одной строке `created_at` окажется `NULL`, кортежное сравнение `(created_at, id) < (...)` вернёт `UNKNOWN`, и пагинация сломается.

---

### Задача 25: Конкурентная обработка и блокировки строк: SELECT FOR UPDATE и SKIP LOCKED

#### 🎯 Условие:
1. **Транзакционный перевод:** Реализовать SQL-логику перевода 500 рублей со счета клиента `A (id = 10)` на счет клиента `B (id = 20)` с защитой от гонки баланса (Lost Update) и дедлоков (Deadlock).
2. **Очередь задач (Task Queue):** Несколько экземпляров микросервиса на Go параллельно опрашивают таблицу задач `task_queue (id BIGINT, status VARCHAR, payload JSONB)`. Каждый воркер должен забрать одну задачу в статусе `'pending'`, перевести её в статус `'processing'` и вернуть себе payload, не блокируя остальных воркеров.

#### 💡 Ловушки многопоточности:
1. **Дедлок взаимных переводов:** Если клиент 1 переводит $10 \to 20$, а клиент 2 одновременно переводит $20 \to 10$, они заблокируют строки в противоположном порядке и транзакции завершатся фатальной ошибкой `deadlock detected`. Решение: **сортировка захвата строк по первичному ключу**.
2. **Блокировка всей очереди:** Если воркеры вызовут простой `SELECT ... WHERE status = 'pending' LIMIT 1 FOR UPDATE`, первый воркер заблокирует строку, а остальные 9 воркеров **встанут в ожидание**, превращая параллельную обработку в последовательную!

#### ✅ Эталонное решение:

**1. Перевод денег с сортировкой блокировок (Deadlock-Free Transfer):**
```sql
BEGIN;

-- Блокируем ОБА счета одновременно строго в порядке возрастания ID:
SELECT id, balance 
FROM accounts 
WHERE id IN (10, 20) 
ORDER BY id 
FOR UPDATE;

-- ⚠️ На бэкенде в Go ОБЯЗАТЕЛЬНО проверить: len(rows) == 2!
-- Если один из счетов не существует (неверный ID или удален), SELECT вернет 1 строку.
-- Без этой проверки произойдет ошибочное списание без зачисления.

-- Выполняем списание и зачисление (проверка balance >= 500 выполняется в Go):
UPDATE accounts SET balance = balance - 500 WHERE id = 10;
UPDATE accounts SET balance = balance + 500 WHERE id = 20;

COMMIT;
```

**2. Захват задачи из очереди через `SKIP LOCKED`:**
```sql
-- Атомарный выбор и захват задачи воркером:
WITH candidate_task AS (
    SELECT id 
    FROM task_queue
    WHERE status = 'pending'
    ORDER BY id ASC
    LIMIT 1
    FOR UPDATE SKIP LOCKED -- Пропускает задачи, уже заблокированные другими горутинами
)
UPDATE task_queue t
SET status = 'processing',
    updated_at = NOW()
FROM candidate_task c
WHERE t.id = c.id
RETURNING t.id, t.payload;
```
* **Почему `SKIP LOCKED`:** Воркер не ждет освобождения заблокированной строки, а мгновенно переходит к следующей свободной. Десятки воркеров работают параллельно с максимальной пропускной способностью.

---

### Задача 26: Атомарный UPSERT (INSERT ... ON CONFLICT DO UPDATE)

#### 🎯 Условие:
В микросервисе сбора статистики ведется учет просмотров статей блога:  
Таблица `article_stats (article_id INT PRIMARY KEY, views_count BIGINT NOT NULL, last_viewed_at TIMESTAMPTZ NOT NULL)`.  
При каждом визите пользователя необходимо атомарно:
- если статьи еще нет в таблице — создать запись с `views_count = 1`;
- если статья уже есть — увеличить `views_count` на 1 и обновить метку `last_viewed_at`.

#### ❌ Антипаттерн "SELECT then INSERT/UPDATE":
```go
// ❌ ГРУБАЯ ОШИБКА:
// Горутина 1 делает SELECT -> записи нет
// Горутина 2 делает SELECT -> записи нет
// Горутина 1 делает INSERT -> успешно
// Горутина 2 делает INSERT -> ERROR: duplicate key value violates unique constraint!
```

#### ✅ Эталонное решение (PostgreSQL UPSERT):
```sql
INSERT INTO article_stats (article_id, views_count, last_viewed_at)
VALUES (42, 1, NOW())
ON CONFLICT (article_id) 
DO UPDATE SET 
    views_count = article_stats.views_count + 1,
    last_viewed_at = EXCLUDED.last_viewed_at
RETURNING views_count;
```

* **Как работает `EXCLUDED`:** Псевдотаблица `EXCLUDED` представляет значения, которые мы пытались вставить через `VALUES(...)`.
* **Условный UPSERT:** Можно добавлять условия обновления, например:  
  `WHERE EXCLUDED.last_viewed_at > article_stats.last_viewed_at`.

---

### Задача 27: Заполнение пропусков во временных рядах через GENERATE_SERIES

#### 🎯 Условие:
Таблица продаж `orders (id INT, created_at DATE, total_amount NUMERIC)`.  
Сформировать аналитический отчет по дням за период с `'2026-09-01'` по `'2026-09-10'`.  
Если в какой-то из дней не было ни одного заказа (например, в выходные), эта дата **обязательно должна присутствовать в отчете** со значением выручки `0` и количеством заказов `0`.

#### ❌ Наивная ошибка:
```sql
-- ❌ Дни без продаж просто выпадут из выборки:
SELECT created_at, COUNT(*), SUM(total_amount)
FROM orders
WHERE created_at BETWEEN '2026-09-01' AND '2026-09-10'
GROUP BY created_at;
```

#### ✅ Эталонное решение (Календарная сетка через `generate_series`):
```sql
WITH calendar AS (
    -- Генерируем полный непрерывный ряд дат:
    SELECT generate_series(
        '2026-09-01'::date, 
        '2026-09-10'::date, 
        '1 day'::interval
    )::date AS report_date
)
SELECT 
    c.report_date,
    COUNT(o.id) AS orders_count,
    COALESCE(SUM(o.total_amount), 0) AS total_revenue
FROM calendar c
LEFT JOIN orders o 
    ON o.created_at = c.report_date
GROUP BY c.report_date
ORDER BY c.report_date ASC;
```

* **Важные детали:**
  1. `LEFT JOIN` сохраняет все дни из виртуальной таблицы `calendar`.
  2. В `COUNT(o.id)` передается конкретный столбец из правой таблицы, а не `COUNT(*)`, иначе дни с `NULL` вернули бы ошибочную единицу.
  3. `COALESCE(SUM(...), 0)` заменяет результирующий `NULL` на `0`.

---

### Задача 28: Индексируемые условия дат (SARGable запросы vs DATE_TRUNC)

#### 🎯 Условие:
В таблице биллинга `transactions (id BIGINT, account_id INT, amount NUMERIC, created_at TIMESTAMPTZ)` имеется индекс:
```sql
CREATE INDEX idx_transactions_created_at ON transactions (created_at);
```
Напишите запрос для выборки всех транзакций за конкретный день `'2026-09-29'` в часовом поясе `+03:00` так, чтобы планировщик PostgreSQL **гарантированно использовал B-Tree индекс**.

#### ❌ Антипаттерн (Non-SARGable предикат):
```sql
-- ❌ ХУДШЕЕ РЕШЕНИЕ: Оборачивание индексированного столбца в функцию:
SELECT * FROM transactions WHERE DATE(created_at) = '2026-09-29';

-- Или:
SELECT * FROM transactions WHERE date_trunc('day', created_at) = '2026-09-29 00:00:00+03';
```
*Почему это убивает производительность:*  
Предикат **Non-SARGable** (Search Argument Non-Able). Индекс `idx_transactions_created_at` упорядочивает чистые значения времени `TIMESTAMPTZ`. Чтобы проверить условие `DATE(created_at)`, PostgreSQL вынужден выполнить функцию `DATE()` для каждой из миллионов строк таблицы через медленный `Seq Scan` (Full Table Scan), полностью игнорируя индекс!

#### ✅ Эталонное SARGable-решение (Диапазон по полуинтервалу):
```sql
-- Столбец остается чистым, а константа разворачивается в полуинтервал [StartOfDay, EndOfDay):
SELECT id, account_id, amount, created_at
FROM transactions
WHERE created_at >= '2026-09-29 00:00:00+03'
  AND created_at <  '2026-09-30 00:00:00+03';
```

* **Результат в плане EXPLAIN:** Метод доступа меняется на `Index Scan` по индексу `idx_transactions_created_at` со сложностью $O(\log N + K)$, где $K$ — число транзакций за выбранный день. Время падает с сотен миллисекунд до десятых долей миллисекунды.

---

### Задача 29: Поиск смежных свободных мест (Consecutive Available Seats / Кинотеатр)

#### 🎯 Условие:
Таблица бронирования мест в кинотеатре `cinema (seat_id INT PRIMARY KEY, is_free INT)`.  
Значение `is_free = 1` означает, что место свободно, `0` — занято. Поле `seat_id` непрерывно возрастает с шагом 1.

Напишите SQL-запрос, который возвращает все **последовательно доступные места** (минимальная длина непрерывной цепочки свободных мест — 2 места). Одиночные свободные места без свободных соседей слева или справа выводиться **не должны**.  
Результат отсортируйте по `seat_id` по возрастанию.

**Пример данных:**
| seat_id | is_free | Попадает в результат? | Причина |
| :--- | :--- | :--- | :--- |
| 1 | 1 | ❌ Нет | Одиночное место (сосед 2 занят) |
| 2 | 0 | ❌ Нет | Занято |
| 3 | 1 | ✅ Да | Имеет свободного соседа 4 |
| 4 | 1 | ✅ Да | Имеет свободных соседей 3 и 5 |
| 5 | 1 | ✅ Да | Имеет свободного соседа 4 |

**Ожидаемый вывод:** `[3, 4, 5]`.

---

#### ✅ Решение 1: Самое короткое (Классика LeetCode 603, всего 7 строк)
Соединяем таблицу `cinema` саму с собой по условию соседства (`ABS(a.seat_id - b.seat_id) = 1`):

```sql
SELECT DISTINCT a.seat_id
FROM cinema a
JOIN cinema b 
    ON ABS(a.seat_id - b.seat_id) = 1
WHERE a.is_free = 1 
  AND b.is_free = 1
ORDER BY a.seat_id ASC;
```

* **Почему это супер-решение на доске:** Пишется за 30 секунд в любой реляционной СУБД (PostgreSQL, MySQL, SQLite, Oracle).

> [!WARNING]
> **Ловушка дубликатов при `SELF JOIN` (Любимый вопрос интервьюера!):**  
> Если подряд идут 3 свободных места (например, 3, 4, 5), то у центрального места 4 **два свободных соседа** (3 и 5). При объединении строка 4 сдублируется дважды (`4-3` и `4-5`).  
> Поэтому использование **`SELECT DISTINCT` здесь строго обязательно!**

---

#### 💡 Решение 2: Однопроходное через `LAG()` и `LEAD()` ($O(N)$ без соединений)
Для каждого свободного места смотрим статус соседа слева (`LAG`) и справа (`LEAD`):

```sql
SELECT seat_id
FROM (
    SELECT 
        seat_id, 
        is_free,
        LAG(is_free) OVER (ORDER BY seat_id) AS prev_free,
        LEAD(is_free) OVER (ORDER BY seat_id) AS next_free
    FROM cinema
) t
WHERE is_free = 1 
  AND (prev_free = 1 OR next_free = 1)
ORDER BY seat_id ASC;
```

* **Плюсы:** Работает за один проход по индексу $O(N)$ без джойнов и без необходимости в `DISTINCT`.
* **Защита от пропусков в нумерации мест:** если между рядами есть пропуски в `seat_id` (например, места 10 и 21 разделены проходом), то дополнительно проверяют равенство соседних ID:  
  `AND ((prev_free = 1 AND prev_id = seat_id - 1) OR (next_free = 1 AND next_id = seat_id + 1))`.

---

#### 💡 Решение 3: Обобщение на произвольную длину $K \ge N$ (Паттерн Gaps and Islands)
Если на собеседовании усложнят задачу: *«Найти места, входящие в непрерывную свободную цепочку длиной от 3 (или 5) мест подряд»*:

```sql
WITH NumberedFreeSeats AS (
    -- Фильтруем свободные места и нумеруем их:
    SELECT 
        seat_id,
        seat_id - ROW_NUMBER() OVER (ORDER BY seat_id) AS island_group
    FROM cinema
    WHERE is_free = 1
),
IslandSizes AS (
    -- Считаем размер каждого непрерывного островка:
    SELECT 
        seat_id,
        COUNT(*) OVER (PARTITION BY island_group) AS streak_len
    FROM NumberedFreeSeats
)
SELECT seat_id
FROM IslandSizes
WHERE streak_len >= 3 -- Параметр минимальной длины цепочки
ORDER BY seat_id ASC;
```

---

#### 🎁 Бонус с собеседования: Замена `NULL` на дефолт (`COALESCE` vs `CASE WHEN`)
Частый вопрос интервьюера: *«Как вывести имя человека и город (`LEFT JOIN address`), а если города нет — вывести вместо NULL строку 'Unknown'?»*

```sql
-- ✅ Идиоматично и лаконично (функция возвращает первый аргумент, не равный NULL):
SELECT 
    p.name, 
    COALESCE(a.city, 'Unknown') AS city
FROM person p
LEFT JOIN address a ON p.id = a.person_id;

-- Эквивалент через CASE WHEN (более громоздкий):
-- CASE WHEN a.city IS NOT NULL THEN a.city ELSE 'Unknown' END AS city
```

---

## 8. Задачи со Скрининга и Проектирования БД (Ozon Technical Screening)

### Задача 30: Группировка и фильтрация заказов (Ozon Screening: Orders)

#### 🎯 Условие:
Дана таблица `orders`:
```sql
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    price_total NUMERIC NOT NULL
);
```
Напишите SQL-запрос, который вернёт количество заказов по каждому пользователю с `price_total >= 1000` в виде `(user_id, order_count)`, отсортированному по количеству заказов в обратном порядке (по убыванию).

#### ✅ Эталонное решение:
```sql
SELECT 
    user_id, 
    COUNT(*) AS order_count
FROM orders 
WHERE price_total >= 1000 
GROUP BY user_id 
ORDER BY order_count DESC;
```

#### 💡 Ловушка и вопросы интервьюера:
1. **Почему `WHERE`, а не `HAVING`?**  
   Условие `price_total >= 1000` фильтрует **исходные строки** до агрегации. Написание этого условия в `HAVING` либо вызовет синтаксическую ошибку (поле `price_total` не входит в `GROUP BY`), либо заставит СУБД агрегировать ненужные строки. `WHERE` отсекает строки на раннем этапе (Index Scan / Seq Scan), экономя память и CPU.
2. **`COUNT(*)` vs `COUNT(id)`:**  
   В современных версиях PostgreSQL и MySQL `COUNT(*)` оптимизирован и работает идентично `COUNT(1)` / `COUNT(id)`. Но `COUNT(col)` игнорирует строки с `NULL`, а `COUNT(*)` считает все физические строки группы.

---

### Задача 31: Покупки до блокировки пользователя и фильтрация сумм (Ozon Screening: Purchases & Ban List)

#### 🎯 Условие:
Дана схема интернет-магазина:
```sql
-- Пользователи:
CREATE TABLE users (
    id INT PRIMARY KEY,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    birth DATE NOT NULL
);

-- Покупки:
CREATE TABLE purchase (
    sku INT NOT NULL,
    price NUMERIC NOT NULL,
    user_id INT NOT NULL REFERENCES users(id),
    date DATE NOT NULL
);

-- Список заблокированных пользователей:
CREATE TABLE ban_list (
    user_id INT NOT NULL REFERENCES users(id),
    date_from DATE NOT NULL
);
```

#### Задание 1: Покупки до блокировки
Вывести уникальные комбинации пользователя и SKU товара для всех покупок, совершенных пользователями **до того, как их забанили** (или если пользователь никогда не был забанен). Отсортировать по фамилии, имени, ID пользователя и SKU.

##### ✅ Решение 1:
```sql
SELECT DISTINCT 
    u.id, 
    u.firstname, 
    u.lastname, 
    p.sku
FROM users u
JOIN purchase p ON u.id = p.user_id
LEFT JOIN ban_list bl ON u.id = bl.user_id
WHERE bl.user_id IS NULL       -- 1. Пользователь никогда не был заблокирован
   OR bl.date_from > p.date    -- 2. Либо покупка совершена ДО даты блокировки
ORDER BY u.lastname, u.firstname, u.id, p.sku;
```

> [!WARNING]
> 🚨 **Ловушка с `INNER JOIN ban_list`:**  
> Если соединить `ban_list` через обычный `JOIN`, в результирующую выборку попадут **ТОЛЬКО пользователи, которые есть в таблице банов**! Все честные незаблокированные пользователи будут безвозвратно отброшены. Обязателен **`LEFT JOIN`** с проверкой `bl.user_id IS NULL`.

---

#### Задание 2: VIP-покупатели с суммой $> 5000$
Найти пользователей, совершивших покупок на общую сумму больше 5000 рублей. Вывести: `id`, `firstname`, `lastname`, `total_amount`.

##### ✅ Решение 2 (Идиоматичный `HAVING`):
```sql
SELECT 
    u.id, 
    u.firstname, 
    u.lastname, 
    SUM(p.price) AS total_amount
FROM users u 
JOIN purchase p ON u.id = p.user_id
GROUP BY u.id, u.firstname, u.lastname
HAVING SUM(p.price) > 5000
ORDER BY total_amount DESC;
```

---

### Задача 32: Моделирование архитектуры чатов и сообщений (Ozon RDBMS: Chats & Messages)

#### 🎯 Условие:
Спроектировать схему реляционной базы данных для мессенджера (Ozon Chat):
1. Сущности: **Пользователь**, **Чат**, **Сообщение**.
2. Пользователь может состоять в нескольких чатах одновременно.
3. В чате может быть много пользователей (Связь N:M).
4. Сообщение принадлежит строго одному чату и строго одному автору (Связь N:1).

#### ✅ DDL-схема таблиц (PostgreSQL):
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Таблица связки многие-ко-многим (N:M):
CREATE TABLE user_chats (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    chat_id INT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, chat_id)
);

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    chat_id INT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    author_id INT NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индекс для мгновенной выборки переписки чата в хронологическом порядке:
CREATE INDEX idx_messages_chat_created ON messages (chat_id, created_at DESC);
```

#### Запрос: Выбрать все чаты пользователя 'Вася' в формате (chat_id, chat_name):
```sql
SELECT 
    uc.chat_id, 
    c.name AS chat_name
FROM users u
JOIN user_chats uc ON uc.user_id = u.id
JOIN chats c ON c.id = uc.chat_id
WHERE u.name = 'Вася';
```

#### 🛡️ Вопросы собеседования сеньор-грейда (18–19 грейд Ozon):
1. **Что будет, если вместо `JOIN` применить `LEFT JOIN`?**  
   Ничего не изменится для существующих пользователей с чатами. Но если Вася есть в базе, а чатов у него нет: обычный `JOIN` вернет 0 строк, а `LEFT JOIN` вернет 1 строку с `chat_id = NULL`.
2. **Сколько и какие индексы созданы по этой схеме?**  
   Автоматически создано **4 уникальных B-Tree индекса** на Primary Key каждой таблицы (`users_pkey`, `chats_pkey`, `user_chats_pkey`, `messages_pkey`).
3. **Запрос поиска чатов Васи стал медленным при росте базы. Почему?**  
   Условие `WHERE u.name = 'Вася'` делает полный скан таблицы `users` (**`Seq Scan`**), так как по полю `name` нет индекса!  
   *Решение:* добавить `CREATE INDEX idx_users_name ON users(name);`.
4. **Что показывает команда `EXPLAIN ANALYZE` и какой узел указывает на нехватку индекса?**  
   Показывает реальный план выполнения, буферы и время в миллисекундах. Узел **`Seq Scan`** с большим `Rows Removed by Filter` однозначно сигнализирует об отсутствии индекса.
5. **Какой первичный ключ лучше: `BIGSERIAL` или `UUID`?**  
   - `BIGSERIAL (8 байт)`: компактный, монотонно возрастает, идеален для B-Tree (нет фрагментации страниц), быстрый `Index Scan`. Но выдает объем бизнеса и требует централизованной генерации СУБД.
   - `UUID (16 байт)`: распределенная генерация в Go без обращения к базе. Но случайный `UUIDv4` фрагментирует B-Tree страницы и увеличивает кэш в 2 раза. В современном проде используют **`UUIDv7`** (сортируемый по времени).
6. **Шардирование:** чаты растут, диск переполнен. Что делать?  
   Шардировать таблицу **`messages`** (она занимает 95% диска из-за текста). Ключ шардирования — **`chat_id`**. Все сообщения одного чата будут лежать на одной ноде, позволяя читать переписку без межшардовых распределенных запросов.

---

### Задача 33: Моделирование каталога библиотеки (Ozon RDBMS: Library Catalog)

#### 🎯 Условие:
Описать модель библиотеки. Три сущности: **«Автор»**, **«Книга»**, **«Читатель»**.  
Физически книга (экземпляр) только одна и может быть на руках только у одного читателя одновременно. У книги может быть несколько авторов.

#### ✅ Схема DDL:
```sql
CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE readers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    card_number TEXT NOT NULL UNIQUE
);

CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    -- Физически экземпляр книги один. Если null — книга в библиотеке (на полке):
    reader_id INT REFERENCES readers(id) ON DELETE SET NULL,
    borrowed_at TIMESTAMP
);

-- Связка книги и авторов (многие-ко-многим):
CREATE TABLE book_authors (
    book_id INT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    author_id INT NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, author_id)
);
```

#### Запросы собеседования:

##### 1. Выбрать названия всех книг, которые сейчас на руках у читателей:
```sql
SELECT b.title, r.name AS reader_name, b.borrowed_at
FROM books b
JOIN readers r ON b.reader_id = r.id;
```

##### 2. Выбрать названия книг, у которых более 3 соавторов:
```sql
SELECT b.title, COUNT(ba.author_id) AS authors_count
FROM books b
JOIN book_authors ba ON b.id = ba.book_id
GROUP BY b.id, b.title
HAVING COUNT(ba.author_id) > 3;
```

##### 3. Топ-3 самых читаемых авторов на данный момент:
```sql
SELECT 
    a.id, 
    a.name, 
    COUNT(*) AS books_on_hands
FROM authors a
JOIN book_authors ba ON a.id = ba.author_id
JOIN books b ON b.id = ba.book_id
WHERE b.reader_id IS NOT NULL -- книга физически на руках
GROUP BY a.id, a.name
ORDER BY books_on_hands DESC
LIMIT 3;
```


