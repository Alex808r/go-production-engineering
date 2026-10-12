package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"go-core-4/05-io/homework-05/pkg/crawler"
	"go-core-4/05-io/homework-05/pkg/crawler/membot"
	"go-core-4/05-io/homework-05/pkg/crawler/spider"
	"go-core-4/05-io/homework-05/pkg/index"
	"go-core-4/05-io/homework-05/pkg/storage"
)

func main() {
	searchFlag := flag.String("s", "", "слово для поиска через инвертированный индекс")
	fileFlag := flag.String("file", "docs.json", "путь к файлу для сохранения/загрузки данных сканера")
	useMemBot := flag.Bool("mem", false, "использовать автономный membot вместо spider")
	forceScan := flag.Bool("rescan", false, "принудительно пересканировать сайты, даже если файл существует")
	flag.Parse()

	var docs []crawler.Document

	// Задача №2: Если файл существует и не запрошен принудительный скан — загружаем данные из файла (io.Reader)
	if !*forceScan && fileExists(*fileFlag) {
		fmt.Printf("Файл %q найден. Загрузка данных (сканирование пропущено)...\n", *fileFlag)
		f, err := os.Open(*fileFlag)
		if err != nil {
			fmt.Printf("Ошибка открытия файла %s: %v\n", *fileFlag, err)
			return
		}
		defer f.Close()

		// Загрузка через io.Reader
		loaded, err := storage.Load(f)
		if err != nil {
			fmt.Printf("Ошибка чтения данных из файла: %v\n", err)
			return
		}
		docs = loaded
		fmt.Printf("Успешно загружено документов из файла: %d\n", len(docs))
	} else {
		// Если файла нет — запускаем сканирование
		fmt.Printf("Кэш-файл %q отсутствует. Запуск сканирования сети...\n", *fileFlag)
		docs = scanSites(*useMemBot)

		// Назначаем уникальные ID и сортируем по ID
		for i := range docs {
			docs[i].ID = i
		}
		sort.Slice(docs, func(i, j int) bool {
			return docs[i].ID < docs[j].ID
		})

		// Задача №1: Запись полученных от сканера данных в файл через io.Writer
		fmt.Printf("Сохранение %d документов в файл %q...\n", len(docs), *fileFlag)
		f, err := os.Create(*fileFlag)
		if err != nil {
			fmt.Printf("Ошибка создания файла %s: %v\n", *fileFlag, err)
			return
		}
		defer f.Close()

		if err := storage.Store(f, docs); err != nil {
			fmt.Printf("Ошибка записи в файл: %v\n", err)
			return
		}
		fmt.Println("Данные успешно сохранены.")
	}

	// Построение инвертированного поискового индекса
	idx := index.New()
	idx.Add(docs)

	// Если передан флаг -s — выполняем одиночный поиск
	if *searchFlag != "" {
		executeSearch(*searchFlag, docs, idx)
		return
	}

	// Иначе запускаем интерактивный поиск через bufio.NewScanner(os.Stdin)
	fmt.Println("\nРежим интерактивного поиска (для выхода введите 'exit' или нажмите Ctrl+C):")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nВведите поисковый запрос > ")
		if !scanner.Scan() {
			break
		}
		query := strings.TrimSpace(scanner.Text())
		if query == "" {
			continue
		}
		if query == "exit" || query == "quit" {
			fmt.Println("Завершение работы.")
			break
		}
		executeSearch(query, docs, idx)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка чтения ввода: %v\n", err)
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func scanSites(useMemBot bool) []crawler.Document {
	var rawDocs []crawler.Document
	urls := []string{
		"https://go.dev",
		"https://golang.org",
	}

	if useMemBot {
		bot := membot.New()
		docs, _ := bot.Scan("", 1)
		return docs
	}

	sp := spider.New()
	for _, u := range urls {
		docs, err := sp.Scan(u, 1)
		if err != nil {
			fmt.Printf("Предупреждение: ошибка сканирования %s: %v\n", u, err)
			continue
		}
		rawDocs = append(rawDocs, docs...)
	}

	// Fallback на membot при отсутствии интернета
	if len(rawDocs) == 0 {
		fmt.Println("Сетевой сканер не вернул данных. Используем встроенный membot...")
		bot := membot.New()
		docs, _ := bot.Scan("", 1)
		return docs
	}

	return rawDocs
}

func executeSearch(query string, docs []crawler.Document, idx *index.Service) {
	docIDs := idx.Search(query)
	if len(docIDs) == 0 {
		fmt.Printf("По запросу %q ничего не найдено.\n", query)
		return
	}

	fmt.Printf("Результаты поиска по слову %q (%d документов):\n", query, len(docIDs))
	for _, id := range docIDs {
		pos := index.BinarySearch(docs, id)
		if pos != -1 {
			d := docs[pos]
			fmt.Printf("  [%d] %s (%s)\n", d.ID, d.Title, d.URL)
		}
	}
}
