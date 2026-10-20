// Package main демонстрирует использование маршрутизатора Gorilla Mux (github.com/gorilla/mux).
//
// Gorilla Mux расширяет стандартный net/http следующими возможностями:
// 1. Извлечение параметров из пути URL (Path variables: /{name}, /{id:[0-9]+}).
// 2. Строгая фильтрация по HTTP-методам (.Methods(http.MethodGet, http.MethodPost)).
// 3. Подмаршрутизаторы (Subrouters) с общими префиксами путей (/api/v1/...).
// 4. Встроенная поддержка middleware через mux.Use().
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Создаем новый экземпляр роутера Gorilla Mux.
	// Он реализует стандартный интерфейс http.Handler (метод ServeHTTP(w, r)),
	// поэтому может передаваться напрямую в http.ListenAndServe или http.Server.
	router := mux.NewRouter()

	// Регистрация маршрута с динамической переменной пути {name}.
	// .Methods(...) гарантирует, что маршрут сработает только для GET-запросов.
	// Для остальных методов роутер автоматически вернет HTTP 405 Method Not Allowed.
	router.HandleFunc("/{name}", mainHandler).Methods(http.MethodGet)

	fmt.Println("Gorilla Mux сервер запущен на http://localhost:8080")
	fmt.Println("Попробуйте открыть: http://localhost:8080/Alex или http://localhost:8080/Gopher")

	// Передаем router в качестве HTTP-хэндлера.
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

// mainHandler — HTTP-обработчик с получением параметров пути.
func mainHandler(w http.ResponseWriter, r *http.Request) {
	// mux.Vars(r) извлекает разобранные переменные маршрута из контекста запроса r.
	// Возвращает map[string]string, где ключи совпадают с именами шаблонов в HandleFunc.
	vars := mux.Vars(r)
	name := vars["name"]

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "<html><body><h2>Привет, %s!</h2><p>Параметр извлечен через Gorilla Mux: vars[\"name\"]</p></body></html>", name)
}

