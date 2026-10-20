package main

import (
	"net/http"

	"go-core-4/lessons/01-intro/demoapp/pkg/stringutils"
)

func main() {
	http.ListenAndServe(
		":8080",
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				s := r.URL.Query()["s"]
				if len(s) == 0 {
					http.Error(w, "bad query", http.StatusBadRequest)
				}

				rev := stringutils.Rev(s[0])

				w.Write([]byte(rev))
			},
		),
	)
}

/*
Пояснения к работе кода:

1. package main:
   Точка входа в исполняемое Go-приложение.

2. Импорты:
   - "net/http" — стандартная библиотека Go для работы с HTTP (сервер и клиент).
   - "go-core-4/lessons/01-intro/demoapp/pkg/stringutils" — локальный пакет проекта, содержащий функцию Rev.

3. http.ListenAndServe(":8080", ...):
   - Запускает HTTP-сервер, который слушает TCP-порт 8080 на всех сетевых интерфейсах.
   - Первым аргументом принимает адрес и порт (":8080").
   - Вторым аргументом принимает корневой обработчик запросов (http.Handler).
   - Блокирует выполнение горутины main, пока сервер не завершит работу с ошибкой.

4. http.HandlerFunc(...):
   - Это адаптер (type HandlerFunc func(ResponseWriter, *Request)), который реализует
     интерфейс http.Handler (метод ServeHTTP).
   - Позволяет передавать обычную функцию в качестве обработчика HTTP-запросов.

5. r.URL.Query()["s"]:
   - r.URL.Query() парсит параметры строки запроса (query string) и возвращает url.Values
     (тип map[string][]string).
   - ["s"] извлекает срез значений для ключа "s" (например, при запросе ?s=hello вернёт []string{"hello"}).

6. if len(s) == 0:
   - Проверяет, передан ли параметр "s" в запросе.
   - Если параметр отсутствует, http.Error отправляет клиенту статус 400 Bad Request с текстом ошибки.
   - return прерывает дальнейшее выполнение обработчика, чтобы не произошло паники при обращении к s[0].

7. stringutils.Rev(s[0]):
   - Берёт первое значение переданного параметра s[0].
   - Переворачивает строку задом наперёд с корректной обработкой многобайтовых UTF-8 символов (через руны []rune).

8. w.Write([]byte(rev)):
   - Преобразует перевёрнутую строку в срез байт []byte и записывает её в тело ответа клиенту.
*/
