// Package index реализует обратный поисковый индекс для документов и бинарный поиск по ID.
package index

import (
	"sort"
	"strings"

	"go-core-4/lessons/07-testing/homework-07/pkg/crawler"
)

// Service представляет инвертированный поисковый индекс: слово -> срез ID документов.
type Service struct {
	data map[string][]int
}

// New инициализирует и возвращает указатель на новый экземпляр индекса.
func New() *Service {
	return &Service{
		data: make(map[string][]int),
	}
}

// Add индексирует переданный срез документов, извлекая слова из заголовка (Title) и URL.
// Номера документов внутри одного слова дедуплицируются.
func (s *Service) Add(docs []crawler.Document) {
	for _, doc := range docs {
		words := append(tokenize(doc.Title), tokenize(doc.URL)...)
		seen := make(map[string]bool)

		for _, word := range words {
			if seen[word] {
				continue
			}
			seen[word] = true
			s.data[word] = append(s.data[word], doc.ID)
		}
	}
}

// Search выполняет поиск слова по индексу и возвращает срез номеров (ID) документов.
func (s *Service) Search(word string) []int {
	word = strings.ToLower(word)
	word = cleanWord(word)
	if word == "" {
		return nil
	}
	return s.data[word]
}

// BinarySearch выполняет бинарный поиск документа по targetID в отсортированном срезе docs.
func BinarySearch(docs []crawler.Document, targetID int) int {
	idx := sort.Search(len(docs), func(i int) bool {
		return docs[i].ID >= targetID
	})
	if idx < len(docs) && docs[idx].ID == targetID {
		return idx
	}
	return -1
}

func tokenize(text string) []string {
	cleaned := strings.Map(func(r rune) rune {
		if r == '/' || r == '.' || r == ':' || r == '-' || r == '_' {
			return ' '
		}
		return r
	}, text)

	fields := strings.Fields(cleaned)
	var words []string
	for _, f := range fields {
		w := cleanWord(strings.ToLower(f))
		if w != "" {
			words = append(words, w)
		}
	}
	return words
}

func cleanWord(word string) string {
	return strings.Trim(word, ".,!?:;\"'()[]{}<>-–—«»/\\")
}
