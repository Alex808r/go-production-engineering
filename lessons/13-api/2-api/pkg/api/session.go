// Package api реализует аутентификацию с использованием HTTP-сессий (Cookie-based sessions).
package api

import (
	"encoding/json"
	"io"
	"net/http"
)

// authInfo содержит модель учетных данных для аутентификации.
type authInfo struct {
	Usr string `json:"usr"`
	Pwd string `json:"pwd"`
}

// authSession обрабатывает POST /api/v1/authSession:
// 1. Проверяет переданные логин и пароль.
// 2. Создает или обновляет сессионную cookie с помощью gorilla/sessions.
// 3. Устанавливает флаг session.Values["authenticated"] = true.
// 4. Отправляет заголовок Set-Cookie клиенту.
//
// Сравнение:
// - Cookie Sessions: состояние хранится на клиенте в зашифрованной/подписанной Cookie.
//   Сервер может легко инвалидировать сессию (logout).
// - JWT: полностью stateless, не требует сохранения состояния на сервере,
//   удобен для микросервисов и мобильных приложений.
func (api *API) authSession(w http.ResponseWriter, r *http.Request) {
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

	if auth.Usr == "Usr" && auth.Pwd == "Pwd" {
		session, _ := api.store.Get(r, "session-cookie")
		session.Values["Usr"] = auth.Usr
		session.Values["Pwd"] = auth.Pwd
		session.Values["authenticated"] = true

		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"authenticated"}`))
		return
	}

	http.Error(w, "неверный логин или пароль", http.StatusUnauthorized)
}

