// Package main демонстрирует архитектуру веб-приложения, оптимизированную
// для модульного тестирования без запуска реального сетевого сокета.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Создаем маршрутизатор и регистрируем маршруты
	r := newRouter()

	fmt.Println("Тестируемый сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}

// newRouter создает и настраивает маршрутизатор приложения.
// Вынесение настройки маршрутов в отдельную функцию позволяет переиспользовать
// её в тестах (testing_test.go) без дублирования кода.
func newRouter() *mux.Router {
	r := mux.NewRouter()
	endpoints(r)
	return r
}

// endpoints регистрирует маршруты веб-приложения.
func endpoints(r *mux.Router) {
	// Маршрут /{name} принимает только GET-запросы.
	r.HandleFunc("/{name}", mainHandler).Methods(http.MethodGet)
}

// mainHandler — HTTP-обработчик, приветствующий пользователя по имени из URL.
func mainHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "<html><body><h2>Hi, %v</h2></body></html>", name)
}

