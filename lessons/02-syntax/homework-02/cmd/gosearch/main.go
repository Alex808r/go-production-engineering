package main

import (
	"flag"
	"fmt"
	"strings"

	"go-core-4/lessons/02-syntax/homework-02/pkg/crawler"
	"go-core-4/lessons/02-syntax/homework-02/pkg/crawler/membot"
	"go-core-4/lessons/02-syntax/homework-02/pkg/crawler/spider"
)

func main() {
	// Задача №3: Обработка флагов вызова с помощью пакета flag
	searchFlag := flag.String("s", "", "слово для поиска среди отсканированных документов")
	useMemBot := flag.Bool("mem", false, "использовать тестовый сканер membot вместо сетевого spider")
	flag.Parse()

	// Задача №2: Использование пакета crawler для сканирования сайтов go.dev и golang.org
	urls := []string{
		"https://go.dev",
		"https://golang.org",
	}

	var allDocs []crawler.Document

	if *useMemBot {
		fmt.Println("Используется локальный эмулятор сканера (membot)...")
		bot := membot.New()
		docs, err := bot.Scan("", 1)
		if err != nil {
			fmt.Printf("Ошибка сканирования membot: %v\n", err)
			return
		}
		allDocs = append(allDocs, docs...)
	} else {
		sp := spider.New()
		fmt.Println("Запуск сканирования сайтов...")
		for _, u := range urls {
			fmt.Printf("Сканирование сайта %s (глубина 1)...\n", u)
			docs, err := sp.Scan(u, 1)
			if err != nil {
				fmt.Printf("Предупреждение: ошибка при сканировании %s: %v\n", u, err)
				continue
			}
			// Объединение результатов сканирования сайтов
			allDocs = append(allDocs, docs...)
		}

		// Fallback при отсутствии сети
		if len(allDocs) == 0 {
			fmt.Println("Сетевое сканирование не вернуло данных. Переключаемся на membot...")
			bot := membot.New()
			docs, _ := bot.Scan("", 1)
			allDocs = append(allDocs, docs...)
		}
	}

	fmt.Printf("\nВсего найдено страниц: %d\n", len(allDocs))

	// Задача №3: Если пользователь ввел слово во флаге «s»,
	// приложение должно напечатать все ссылки, где это слово встречается.
	query := strings.TrimSpace(*searchFlag)
	if query != "" {
		fmt.Printf("Результаты поиска по запросу %q:\n", query)
		results := FilterDocs(allDocs, query)
		if len(results) == 0 {
			fmt.Println("Ничего не найдено.")
			return
		}
		for _, doc := range results {
			fmt.Printf("  • %s — %s\n", doc.URL, doc.Title)
		}
		return
	}

	// Если флаг -s не задан, выводим список всех найденных ссылок
	fmt.Println("Список всех просканированных ссылок (для фильтрации используйте флаг -s <слово>):")
	for _, doc := range allDocs {
		fmt.Printf("  • %s — %s\n", doc.URL, doc.Title)
	}
}

// FilterDocs выполняет поиск страниц, содержащих искомое слово в URL или Title (без учета регистра).
func FilterDocs(docs []crawler.Document, word string) []crawler.Document {
	word = strings.ToLower(word)
	var matched []crawler.Document

	for _, doc := range docs {
		titleLower := strings.ToLower(doc.Title)
		urlLower := strings.ToLower(doc.URL)

		if strings.Contains(titleLower, word) || strings.Contains(urlLower, word) {
			matched = append(matched, doc)
		}
	}

	return matched
}
