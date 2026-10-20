# 08-prof_debug — Профилирование, отладка и трассировка в Go

В данной директории изучаются инструменты анализа производительности приложений на Go: профайлер `pprof` (анализ процессорного времени и потребления памяти), отладка и сборка трасс исполнения (`go tool trace`).

---

## 📋 Структура каталога

```text
08-prof_debug/
├── 1-bench_profile/  # Профилирование тестов (-cpuprofile cpu.out, -memprofile mem.out)
├── 2-app_profile/    # Профилирование работающего приложения через net/http/pprof
├── 3-debug/          # Инструменты отладки (стектрейсы, сигналы, Delve)
├── 4-trace/          # Трассировка выполнения приложения (runtime/trace)
│   ├── trace_cpu/    # Анализ планировщика горутин и загрузки ядер CPU
│   └── trace_mem/    # Анализ аллокаций памяти и работы Garbage Collector
└── homework-08/      # Практическое ДЗ: Оптимизация Two Sum, профилирование и трассировка
```

---

## 🔍 Ключевые концепции

1. **CPU и Memory профилирование (`pprof`):**
   - **CPU profile:** сэмплирует состояние программы каждые 10 мс, определяя функции, потребляющие наибольшее процессорное время ("hot spots").
   - **Memory profile:** отслеживает аллокации памяти в куче (heap) и помогает выявлять утечки памяти (`inuse_space`, `alloc_objects`).

2. **Интерактивный веб-интерфейс pprof:**
   ```bash
   go tool pprof -http=:8080 cpu.out
   ```
   Позволяет строить графы вызовов (Callgraph), диаграммы пламени (Flame Graph) и просматривать ассемблерный код с подсветкой затрат.

3. **Трассировка исполнения (`go tool trace`):**
   - Дает детальную покадровую картину: переключение контекста горутин, задержки сетевого ввода-вывода (Network I/O Wait), паузы сборщика мусора (STW / GC Sweep) и блокировки мьютексов.

---

## 🚀 Запуск профилирования

```bash
# Сборка профилей из бенчмарков
cd ./lessons/08-prof_debug/1-bench_profile
go test -bench=. -cpuprofile=cpu.out -memprofile=mem.out

# Анализ в веб-браузере
go tool pprof -http=:8081 cpu.out
```

---

## 🔗 Связанные материалы
- **Теория:** [Том 5: Тестирование и архитектура (Раздел 20: Профилирование и отладка)](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/go-documentation/05-go-core-testing-arch.md)
- **Практическое ДЗ:** [homework-08](file:///Users/user/workspace/golang_projects/go-core-4/go_course_4/lessons/08-prof_debug/homework-08) (Оптимизация алгоритма Two Sum, профилирование GoSearch с pprof и trace)
