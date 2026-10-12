package index

import (
	"sort"
	"strings"

	"go-core-4/GoSearch/pkg/crawler"
)

// Service представляет обратный поисковый индекс.
type Service struct {
	data map[string][]int
}

// New возвращает новый экземпляр индекса.
func New() *Service {
	return &Service{
		data: make(map[string][]int),
	}
}

// Add индексирует срез документов, извлекая слова из Title каждого документа.
func (s *Service) Add(docs []crawler.Document) {
	for _, doc := range docs {
		words := tokenize(doc.Title)
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

// BinarySearch находит документ по его ID в отсортированном по ID срезе документов
// с помощью бинарного поиска (sort.Search из стандартной библиотеки).
// Возвращает индекс найденного документа в docs или -1, если документ не найден.
func BinarySearch(docs []crawler.Document, targetID int) int {
	idx := sort.Search(len(docs), func(i int) bool {
		return docs[i].ID >= targetID
	})
	if idx < len(docs) && docs[idx].ID == targetID {
		return idx
	}
	return -1
}

// tokenize разбивает строку на отдельные нормализованные слова.
func tokenize(text string) []string {
	fields := strings.Fields(text)
	var words []string
	for _, f := range fields {
		w := cleanWord(strings.ToLower(f))
		if w != "" {
			words = append(words, w)
		}
	}
	return words
}

// cleanWord очищает слово от знаков препинания по краям.
func cleanWord(word string) string {
	return strings.Trim(word, ".,!?:;\"'()[]{}<>-–—«»/\\")
}
