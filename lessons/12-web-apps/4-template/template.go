// Package main демонстрирует использование промежуточного ПО (Middleware)
// и безопасную HTML-шаблонизацию с помощью html/template в Go.
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// =========================================================================
	// Конвейер Middleware (Цепочка обязанностей / Onion Architecture)
	// =========================================================================
	// router.Use добавляет middleware в глобальный стек маршрутизатора.
	// Они выполняются в строгом порядке добавления (слева направо / сверху вниз)
	// до вызова целевого обработчика, а затем в обратном порядке после возврата.
	router.Use(loggingMiddleware)
	router.Use(timingMiddleware)

	// Маршрут с несколькими переменными пути: /{first}/{last}
	router.HandleFunc("/{first}/{last}", mainHandler).Methods(http.MethodGet)

	fmt.Println("Сервер с Middleware и шаблонизацией запущен на http://localhost:8080")
	fmt.Println("Откройте в браузере: http://localhost:8080/Ivan/Petrov")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}

// loggingMiddleware логирует каждый входящий HTTP-запрос (метод, URL, RemoteAddr).
// Сигнатура стандартного middleware: func(http.Handler) http.Handler.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[LOG] Входящий запрос: %s %s от %s", r.Method, r.RequestURI, r.RemoteAddr)
		// Передаем управление следующему обработчику в цепочке:
		next.ServeHTTP(w, r)
		// Код ниже выполнится ПОСЛЕ завершения обработки запроса
		log.Printf("[LOG] Завершение запроса: %s %s", r.Method, r.RequestURI)
	})
}

// timingMiddleware измеряет время выполнения запроса.
func timingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("[TIMING] Запрос %s обработан за %v", r.RequestURI, duration)
	})
}

// Шаблон страницы, скомпилированный заранее в памяти для высокой производительности.
// Использование пакета html/template гарантирует контекстно-зависимое экранирование (auto-escaping).
// Если пользователь передаст имя вида `<script>alert('xss')</script>`, html/template
// автоматически преобразует его в `&lt;script&gt;alert('xss')&lt;/script&gt;`, предотвращая XSS.
var pageTemplate = template.Must(template.New("main").Parse(`<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Приветствие</title>
</head>
<body>
    <h2>Здравствуйте, {{ .FullName }}!</h2>
    <p>Имя: <strong>{{ .FirstName }}</strong></p>
    <p>Фамилия: <strong>{{ .LastName }}</strong></p>
</body>
</html>`))

// PageData — структура данных для передачи в шаблон.
type PageData struct {
	FirstName string
	LastName  string
	FullName  string
}

// mainHandler — обработчик, рендерящий HTML-страницу из шаблона.
func mainHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	first := vars["first"]
	last := vars["last"]

	data := PageData{
		FirstName: first,
		LastName:  last,
		FullName:  first + " " + last,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Execute записывает отрендеренный HTML напрямую в http.ResponseWriter.
	if err := pageTemplate.Execute(w, data); err != nil {
		http.Error(w, "ошибка при рендеринге шаблона", http.StatusInternalServerError)
		log.Printf("Ошибка рендеринга шаблона: %v", err)
	}
}

