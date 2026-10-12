package webapp

import (
	"strings"
	"sync"
)

// SearchDB представляет потокобезопасное хранилище документов и инвертированный поисковый индекс.
// Реализует интерфейс DataProvider.
type SearchDB struct {
	mu   sync.RWMutex
	docs []Document
	idx  map[string][]int
}

// NewSearchDB создает новое хранилище документов и поискового индекса.
func NewSearchDB() *SearchDB {
	return &SearchDB{
		docs: make([]Document, 0),
		idx:  make(map[string][]int),
	}
}

// Add индексирует и добавляет документы в хранилище.
func (db *SearchDB) Add(docs []Document) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, doc := range docs {
		db.docs = append(db.docs, doc)

		words := tokenize(doc.Title + " " + doc.URL + " " + doc.Body)
		seen := make(map[string]bool)

		for _, w := range words {
			if seen[w] {
				continue
			}
			seen[w] = true
			db.idx[w] = append(db.idx[w], doc.ID)
		}
	}
}

// Documents возвращает копию списка всех документов в хранилище.
func (db *SearchDB) Documents() []Document {
	db.mu.RLock()
	defer db.mu.RUnlock()

	res := make([]Document, len(db.docs))
	copy(res, db.docs)
	return res
}

// Index возвращает копию инвертированного поискового индекса.
func (db *SearchDB) Index() map[string][]int {
	db.mu.RLock()
	defer db.mu.RUnlock()

	res := make(map[string][]int, len(db.idx))
	for k, v := range db.idx {
		docIDs := make([]int, len(v))
		copy(docIDs, v)
		res[k] = docIDs
	}
	return res
}

// Search находит ID документов по заданному слову в индексе.
func (db *SearchDB) Search(word string) []int {
	db.mu.RLock()
	defer db.mu.RUnlock()

	cleaned := cleanWord(strings.ToLower(word))
	if cleaned == "" {
		return nil
	}
	return db.idx[cleaned]
}

// tokenize разбивает строку на нормализованные лексемы.
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
		w := cleanWord(strings.ToLower(f))
		if len(w) > 1 {
			words = append(words, w)
		}
	}
	return words
}

// cleanWord очищает токен от спецсимволов.
func cleanWord(word string) string {
	return strings.Trim(word, ".,!?:;\"'()[]{}<>-–—«»/\\")
}
