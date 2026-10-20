// Package api содержит цепочки middleware для авторизации, логирования и трейсинга запросов.
package api

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// sessionsMiddleware проверяет наличие валидного сессионного cookie.
// Если пользователь не аутентифицирован, возвращает 403 Forbidden.
func (api *API) sessionsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := api.store.Get(r, "session-cookie")
		if err != nil {
			http.Error(w, "ошибка сессии", http.StatusForbidden)
			return
		}

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "доступ запрещён: требуется авторизованная сессия", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// jwtMiddleware валидирует заголовок Authorization: Bearer <token>.
// Если токен передан, парсит и верифицирует его цифровую подпись.
func (api *API) jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			// Если заголовок отсутствует, пропускаем дальше (для публичных эндпоинтов /auth)
			next.ServeHTTP(w, r)
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 || splitted[0] != "Bearer" {
			http.Error(w, "неверный формат заголовка Authorization (ожидается: Bearer <token>)", http.StatusBadRequest)
			return
		}

		tokenString := splitted[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверяем метод подписи токена (HMAC)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
			}
			return JWTSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "недействительный или истекший JWT токен", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Сохраняем разобранные claims в контекст запроса для доступа в хэндлерах
			ctx := context.WithValue(r.Context(), "user_claims", claims)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// requestIDMiddleware генерирует уникальный идентификатор запроса (Request ID)
// и пробрасывает его в контекст запроса.
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := rand.Intn(1_000_000)
		ctx := context.WithValue(r.Context(), "request_id", reqID)
		newR := r.WithContext(ctx)

		next.ServeHTTP(w, newR)
	})
}

// logMiddleware логирует метод, адрес клиента, URL и Request ID входящего запроса.
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqID any = "none"
		if id := r.Context().Value("request_id"); id != nil {
			reqID = id
		}
		log.Printf("%s %s %s [req_id: %v]", r.Method, r.RemoteAddr, r.RequestURI, reqID)

		next.ServeHTTP(w, r)
	})
}

