// Package engine реализует поисковый движок в объектно-ориентированном стиле.
// Согласно лекции №6 ("ООП в Go"), сущности, хранящие и изменяющие свое состояние
// между вызовами (состояние индекса, кэш документов), оформляются в виде структуры с методами.
package engine

import (
	"fmt"
	"io"
	"sort"

	"go-core-4/lessons/06-oop/homework-06/pkg/crawler"
	"go-core-4/lessons/06-oop/homework-06/pkg/index"
	"go-core-4/lessons/06-oop/homework-06/pkg/storage"
)

// Engine представляет поисковый движок как объект, инкапсулирующий состояние
// (проиндексированные документы) и зависимости (краулер, инвертированный индекс).
type Engine struct {
	crawler crawler.Interface
	index   *index.Service
	docs    []crawler.Document
}

// New — идиоматичный конструктор поискового движка (паттерн New() из лекции №6).
// Реализует внедрение зависимостей (Dependency Injection):
// конкретные реализации краулера и индекса создаются на верхнем уровне (в main)
// и передаются в конструктор подсистемы.
func New(cr crawler.Interface, idx *index.Service) *Engine {
	if idx == nil {
		idx = index.New()
	}
	return &Engine{
		crawler: cr,
		index:   idx,
		docs:    make([]crawler.Document, 0),
	}
}

// Scan выполняет обход списка сайтов с заданной глубиной,
// индексирует полученные документы и сохраняет их во внутреннем состоянии объекта.
func (e *Engine) Scan(urls []string, depth int) error {
	if e.crawler == nil {
		return fmt.Errorf("engine: поисковый робот crawler не инициализирован")
	}

	var allDocs []crawler.Document
	for _, u := range urls {
		scanned, err := e.crawler.Scan(u, depth)
		if err != nil {
			// Логируем предупреждение и продолжаем со следующими URL
			continue
		}
		allDocs = append(allDocs, scanned...)
	}

	// Назначаем уникальные последовательные ID
	for i := range allDocs {
		allDocs[i].ID = i
	}

	// Сортируем документы по ID для обеспечения бинарного поиска
	sort.Slice(allDocs, func(i, j int) bool {
		return allDocs[i].ID < allDocs[j].ID
	})

	e.docs = allDocs
	e.index.Add(e.docs)

	return nil
}

// Search выполняет поиск документов по поисковому запросу.
// Поиск использует обратный индекс для получения ID и бинарный поиск для извлечения документов.
func (e *Engine) Search(query string) []crawler.Document {
	ids := e.index.Search(query)
	if len(ids) == 0 {
		return nil
	}

	results := make([]crawler.Document, 0, len(ids))
	for _, id := range ids {
		pos := index.BinarySearch(e.docs, id)
		if pos != -1 {
			results = append(results, e.docs[pos])
		}
	}

	return results
}

// Documents возвращает копию среза документов поискового движка (инкапсуляция).
func (e *Engine) Documents() []crawler.Document {
	result := make([]crawler.Document, len(e.docs))
	copy(result, e.docs)
	return result
}

// SetDocuments вручную обновляет список документов движка и перестраивает поисковый индекс.
func (e *Engine) SetDocuments(docs []crawler.Document) {
	// Сортируем документы по ID
	sort.Slice(docs, func(i, j int) bool {
		return docs[i].ID < docs[j].ID
	})

	e.docs = make([]crawler.Document, len(docs))
	copy(e.docs, docs)

	// Переиндексация
	e.index = index.New()
	e.index.Add(e.docs)
}

// Save сохраняет состояние документов движка в io.Writer через пакет storage.
func (e *Engine) Save(w io.Writer) error {
	if w == nil {
		return fmt.Errorf("engine.Save: writer не может быть nil")
	}
	return storage.Store(w, e.docs)
}

// Load считывает документы из io.Reader, восстанавливает состояние движка и индексирует их.
func (e *Engine) Load(r io.Reader) error {
	if r == nil {
		return fmt.Errorf("engine.Load: reader не может быть nil")
	}

	docs, err := storage.Load(r)
	if err != nil {
		return fmt.Errorf("engine.Load: ошибка загрузки: %w", err)
	}

	e.SetDocuments(docs)
	return nil
}

// Count возвращает общее количество проиндексированных документов.
func (e *Engine) Count() int {
	return len(e.docs)
}
