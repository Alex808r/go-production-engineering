package storage

import (
	"errors"
	"sort"
	"strings"
	"sync"
)

// ErrNotFound возвращается, когда запрашиваемый документ не найден.
var ErrNotFound = errors.New("document not found")

// Document представляет индексируемый документ поисковика GoSearch.
type Document struct {
	ID    int    `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Interface определяет контракт хранилища документов и поискового индекса.
type Interface interface {
	Search(query string) []Document
	Documents() []Document
	Document(id int) (Document, error)
	Add(doc Document) Document
	Update(doc Document) error
	Delete(id int) error
}

// DB реализует потокобезопасное хранилище документов и инвертированный индекс.
type DB struct {
	mu    sync.RWMutex
	docs  map[int]Document
	idx   map[string][]int
	maxID int
}

// New создает новый экземпляр DB.
func New() *DB {
	return &DB{
		docs: make(map[int]Document),
		idx:  make(map[string][]int),
	}
}

// Search выполняет поиск документов по текстовому запросу с использованием инвертированного индекса.
func (db *DB) Search(query string) []Document {
	db.mu.RLock()
	defer db.mu.RUnlock()

	words := tokenize(query)
	if len(words) == 0 {
		return nil
	}

	seenIDs := make(map[int]bool)
	var matchedDocs []Document

	for _, w := range words {
		ids, ok := db.idx[w]
		if !ok {
			continue
		}
		for _, id := range ids {
			if seenIDs[id] {
				continue
			}
			seenIDs[id] = true
			if doc, exists := db.docs[id]; exists {
				matchedDocs = append(matchedDocs, doc)
			}
		}
	}

	sort.Slice(matchedDocs, func(i, j int) bool {
		return matchedDocs[i].ID < matchedDocs[j].ID
	})

	return matchedDocs
}

// Documents возвращает отсортированный по ID список всех документов в хранилище.
func (db *DB) Documents() []Document {
	db.mu.RLock()
	defer db.mu.RUnlock()

	res := make([]Document, 0, len(db.docs))
	for _, doc := range db.docs {
		res = append(res, doc)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ID < res[j].ID
	})

	return res
}

// Document возвращает документ по его ID или ErrNotFound.
func (db *DB) Document(id int) (Document, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	doc, ok := db.docs[id]
	if !ok {
		return Document{}, ErrNotFound
	}
	return doc, nil
}

// Add добавляет документ в хранилище и индексирует его.
// Если ID <= 0 или уже занят, ему автоматически присваивается инкрементальный ID.
func (db *DB) Add(doc Document) Document {
	db.mu.Lock()
	defer db.mu.Unlock()

	if doc.ID <= 0 || db.docs[doc.ID].ID != 0 {
		db.maxID++
		doc.ID = db.maxID
	} else if doc.ID > db.maxID {
		db.maxID = doc.ID
	}

	db.docs[doc.ID] = doc

	// Индексируем документ
	words := tokenize(doc.Title + " " + doc.URL + " " + doc.Body)
	seen := make(map[string]bool)
	for _, w := range words {
		if seen[w] {
			continue
		}
		seen[w] = true
		db.idx[w] = append(db.idx[w], doc.ID)
	}

	return doc
}

// Update обновляет документ по его ID.
// По условию домашнего задания перестраивать индекс при изменении документов не обязательно.
func (db *DB) Update(doc Document) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.docs[doc.ID]; !ok {
		return ErrNotFound
	}

	db.docs[doc.ID] = doc
	return nil
}

// Delete удаляет документ из хранилища по его ID.
// По условию домашнего задания перестраивать индекс при удалении документов не обязательно.
func (db *DB) Delete(id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.docs[id]; !ok {
		return ErrNotFound
	}

	delete(db.docs, id)
	return nil
}

// tokenize разбивает строку на нормализованные токены без спецсимволов.
func tokenize(text string) []string {
	cleaned := strings.Map(func(r rune) rune {
		if r == '/' || r == '.' || r == ':' || r == '-' || r == '_' || r == ',' || r == '!' || r == '?' {
			return ' '
		}
		return r
	}, text)

	fields := strings.Fields(cleaned)
	var words []string
	for _, f := range fields {
		w := strings.Trim(strings.ToLower(f), ".,!?:;\"'()[]{}<>-–—«»/\\")
		if len(w) > 1 {
			words = append(words, w)
		}
	}
	return words
}
