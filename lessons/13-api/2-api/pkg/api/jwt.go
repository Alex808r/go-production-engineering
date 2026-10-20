// Package api реализует выдачу и валидацию JWT токенов (JSON Web Token).
package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// JWTSecret — секретный ключ для симметричной подписи HMAC-SHA256.
// В production-системах должен храниться в защищенных переменных окружения.
var JWTSecret = []byte("secret-password")

// authJWT обрабатывает POST /api/v1/authJWT и генерирует подписанный JWT-токен
// в случае успешной проверки учетных данных (логин / пароль).
//
// Структура JWT:
// 1. Header (заголовок): алгоритм подписи (HS256) и тип токена (JWT).
// 2. Payload (полезная нагрузка): набор утверждений (Claims), например:
//    - "usr": имя пользователя
//    - "nbf" (Not Before): время, начиная с которого токен становится действительным.
//    - "exp" (Expiration): время истечения срока действия токена.
// 3. Signature (цифровая подпись): хэш HMAC-SHA256(Base64(Header) + "." + Base64(Payload), Secret).
func (api *API) authJWT(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var auth authInfo
	err = json.Unmarshal(body, &auth)
	if err != nil {
		http.Error(w, "некорректный JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Демонстрационная проверка учетных данных (Usr/Pwd)
	if auth.Usr == "Usr" && auth.Pwd == "Pwd" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"usr": auth.Usr,
			"nbf": time.Now().Unix(),
			"exp": time.Now().Add(24 * time.Hour).Unix(),
		})

		// Подписание токена секретным ключом
		tokenString, err := token.SignedString(JWTSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(tokenString))
		return
	}

	http.Error(w, "неверный логин или пароль", http.StatusUnauthorized)
}

