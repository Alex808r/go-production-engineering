// Package api предоставляет продвинутый программный интерфейс REST API
// с поддержкой аутентификации через сессии и токены JWT, контекстов запросов
// и структурированного логирования.
package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
)

// API инкапсулирует состояние и зависимости REST-сервиса:
// маршрутизатор, хранилище сессионных кук и логгер.
type API struct {
	router *mux.Router
	store  *sessions.CookieStore
	log    *zerolog.Logger
}

// New создаёт и конфигурирует новый экземпляр API.
func New() *API {
	api := API{
		router: mux.NewRouter(),
		// CookieStore шифрует и подписывает cookie ключом для предотвращения подделки сессий клиентом
		store: sessions.NewCookieStore([]byte("secret_password")),
	}

	api.endpoints()

	return &api
}

// Router возвращает настроенный HTTP-маршрутизатор (реализует http.Handler).
func (api *API) Router() *mux.Router {
	return api.router
}

// endpoints регистрирует цепочки middleware и конечные точки REST API.
func (api *API) endpoints() {
	// Глобальные middleware: выполняются для каждого запроса
	api.router.Use(requestIDMiddleware)
	api.router.Use(logMiddleware)
	api.router.Use(api.jwtMiddleware)
	// api.router.Use(api.sessionsMiddleware) // опциональная проверка сессий

	// Маршруты аутентификации
	api.router.HandleFunc("/api/v1/authSession", api.authSession).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/api/v1/authJWT", api.authJWT).Methods(http.MethodPost, http.MethodOptions)

	// Бизнес-эндпоинты работы с книгами
	api.router.HandleFunc("/api/v1/books", api.books).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/api/v1/newBook", api.newBook).Methods(http.MethodPost, http.MethodOptions)
}

