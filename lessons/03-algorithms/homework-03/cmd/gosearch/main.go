package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"go-core-4/lessons/03-algorithms/homework-03/pkg/crawler"
	"go-core-4/lessons/03-algorithms/homework-03/pkg/crawler/membot"
	"go-core-4/lessons/03-algorithms/homework-03/pkg/crawler/spider"
	"go-core-4/lessons/03-algorithms/homework-03/pkg/index"
)

func main() {
	// Задача №2: Обработка флага -s для быстрого поиска по индексу
	searchFlag := flag.String("s", "", "слово для поиска через инвертированный индекс")
	useMemBot := flag.Bool("mem", false, "использовать автономный membot вместо spider")
	flag.Parse()

	urls := []string{
		"https://go.dev",
		"https://golang.org",
	}

	var rawDocs []crawler.Document

	if *useMemBot {
		fmt.Println("Используется локальный эмулятор сканера (membot)...")
		bot := membot.New()
		docs, err := bot.Scan("", 1)
		if err != nil {
			fmt.Printf("Ошибка membot: %v\n", err)
			return
		}
		rawDocs = append(rawDocs, docs...)
	} else {
		sp := spider.New()
		fmt.Println("Запуск сканирования сайтов...")
		for _, u := range urls {
			fmt.Printf("Сканирование %s (глубина 1)...\n", u)
			docs, err := sp.Scan(u, 1)
			if err != nil {
				fmt.Printf("Предупреждение: ошибка при сканировании %s: %v\n", u, err)
				continue
			}
			rawDocs = append(rawDocs, docs...)
		}

		// Fallback при отсутствии сетевого доступа
		if len(rawDocs) == 0 {
			fmt.Println("Сетевой сканер не вернул данных. Используем membot...")
			bot := membot.New()
			docs, _ := bot.Scan("", 1)
			rawDocs = append(rawDocs, docs...)
		}
	}

	// Задача №1: Присваиваем каждому документу уникальный номер ID
	for i := range rawDocs {
		rawDocs[i].ID = i
	}

	// Задача №1: Сортируем документы по ID с помощью стандартной библиотеки
	sort.Slice(rawDocs, func(i, j int) bool {
		return rawDocs[i].ID < rawDocs[j].ID
	})

	// Задача №1: Создаем обратный поисковый индекс и индексируем документы
	idx := index.New()
	idx.Add(rawDocs)

	fmt.Printf("\nВсего проиндексировано документов: %d\n", len(rawDocs))

	query := strings.TrimSpace(*searchFlag)
	if query != "" {
		fmt.Printf("Поиск по индексу по слову %q:\n", query)
		results := SearchDocs(rawDocs, idx, query)
		if len(results) == 0 {
			fmt.Println("Ничего не найдено.")
			return
		}
		for _, doc := range results {
			fmt.Printf("  [ID: %d] %s — %s\n", doc.ID, doc.URL, doc.Title)
		}
		return
	}

	// Если слово для поиска не задано, выводим список документов
	fmt.Println("Список всех проиндексированных страниц (для поиска используйте флаг -s <слово>):")
	for _, doc := range rawDocs {
		fmt.Printf("  [ID: %d] %s — %s\n", doc.ID, doc.URL, doc.Title)
	}
}

// SearchDocs выполняет поиск номеров документов в индексе (Задача №2),
// а затем находит сами документы в отсортированном массиве бинарным поиском (Задача №3).
func SearchDocs(docs []crawler.Document, idx *index.Service, word string) []crawler.Document {
	docIDs := idx.Search(word)
	if len(docIDs) == 0 {
		return nil
	}

	var results []crawler.Document
	for _, id := range docIDs {
		pos := index.BinarySearch(docs, id)
		if pos != -1 {
			results = append(results, docs[pos])
		}
	}

	return results
}
