// Поисковый сервис GoSearch (Итерация №7 — Тестирование, мокирование и защита от регрессий).
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"go-core-4/lessons/07-testing/homework-07/pkg/crawler"
	"go-core-4/lessons/07-testing/homework-07/pkg/crawler/membot"
	"go-core-4/lessons/07-testing/homework-07/pkg/crawler/spider"
	"go-core-4/lessons/07-testing/homework-07/pkg/engine"
	"go-core-4/lessons/07-testing/homework-07/pkg/index"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	searchFlag := flag.String("s", "", "слово для поиска через инвертированный индекс")
	fileFlag := flag.String("file", "docs.json", "путь к файлу для сохранения/загрузки данных поисковика")
	useMemBot := flag.Bool("mem", false, "использовать автономный membot вместо паука spider")
	forceScan := flag.Bool("rescan", false, "принудительно пересканировать сайты, даже если файл существует")
	flag.Parse()

	// 1. Инициализация зависимостей на верхнем уровне (main):
	var cr crawler.Interface
	if *useMemBot {
		fmt.Println("Инициализация автономного сканера membot...")
		cr = membot.New()
	} else {
		cr = spider.New()
	}
	idx := index.New()

	// 2. Внедрение зависимостей в объект поискового движка через конструктор New:
	eng := engine.New(cr, idx)

	// 3. Загрузка данных или сканирование сети:
	if !*forceScan && fileExists(*fileFlag) {
		fmt.Printf("Кэш-файл %q найден. Загрузка данных в объект движка...\n", *fileFlag)
		f, err := os.Open(*fileFlag)
		if err != nil {
			fmt.Printf("Ошибка открытия файла: %v\n", err)
			return
		}
		defer f.Close()

		if err := eng.Load(f); err != nil {
			fmt.Printf("Ошибка загрузки данных движком: %v\n", err)
			return
		}
		fmt.Printf("Движок успешно загрузил документов: %d\n", eng.Count())
	} else {
		fmt.Printf("Кэш-файл %q отсутствует или запрошен -rescan. Запуск сканирования...\n", *fileFlag)
		urls := []string{
			"https://go.dev",
			"https://golang.org",
		}

		if err := eng.Scan(urls, 1); err != nil || eng.Count() == 0 {
			fmt.Println("Сетевое сканирование не вернуло результатов. Переключаемся на membot...")
			bot := membot.New()
			engBot := engine.New(bot, idx)
			_ = engBot.Scan([]string{""}, 1)
			eng = engBot
		}

		fmt.Printf("Проиндексировано документов: %d. Сохранение в файл %q...\n", eng.Count(), *fileFlag)
		f, err := os.Create(*fileFlag)
		if err != nil {
			fmt.Printf("Ошибка создания файла %s: %v\n", *fileFlag, err)
			return
		}
		defer f.Close()

		if err := eng.Save(f); err != nil {
			fmt.Printf("Ошибка сохранения данных движком: %v\n", err)
			return
		}
		fmt.Println("Данные успешно сохранены в файл.")
	}

	// 4. Режим разового поиска через CLI-флаг -s:
	if *searchFlag != "" {
		executeSearch(eng, *searchFlag)
		return
	}

	// 5. Интерактивный режим поиска:
	runInteractive(eng)
}

func executeSearch(eng *engine.Engine, query string) {
	fmt.Printf("\n=== Результаты поиска для запроса: %q ===\n", query)
	results := eng.Search(query)
	if len(results) == 0 {
		fmt.Println("Ничего не найдено.")
		return
	}

	for i, doc := range results {
		fmt.Printf("[%d] ID: %d | %s\n    URL: %s\n", i+1, doc.ID, doc.Title, doc.URL)
	}
}

func runInteractive(eng *engine.Engine) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\n=== Интерактивный поиск GoSearch (ООП-движок) ===")
	fmt.Println("Введите поисковый запрос (или 'exit' для выхода):")

	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.EqualFold(line, "exit") || strings.EqualFold(line, "quit") {
			fmt.Println("Выход из GoSearch. До свидания!")
			break
		}

		executeSearch(eng, line)
	}
}
