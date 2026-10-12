// Package stdgenerics демонстрирует возможности стандартных generic-пакетов
// slices, maps и cmp, появившихся в Go 1.21+.
package stdgenerics

import (
	"cmp"
	"maps"
	"slices"
	"strings"
)

// Product представляет товар в каталоге.
type Product struct {
	ID    int
	Title string
	Price float64
}

// FindProductByTitle находит индекс товара по точному названию в отсортированном срезе
// с использованием generic-функции slices.BinarySearchFunc.
func FindProductByTitle(products []Product, targetTitle string) (int, bool) {
	return slices.BinarySearchFunc(products, targetTitle, func(p Product, target string) int {
		return strings.Compare(p.Title, target)
	})
}

// SortProductsByPrice сортирует срез товаров по возрастанию цены с помощью slices.SortFunc.
func SortProductsByPrice(products []Product) {
	slices.SortFunc(products, func(a, b Product) int {
		return cmp.Compare(a.Price, b.Price)
	})
}

// DeduplicateAndCompact удаляет повторяющиеся подряд элементы с помощью slices.Compact.
func DeduplicateAndCompact(numbers []int) []int {
	sorted := slices.Clone(numbers)
	slices.Sort(sorted)
	return slices.Compact(sorted)
}

// CloneAndFilterMap клонирует мапу и удаляет элементы, не удовлетворяющие предикату.
func CloneAndFilterMap(original map[string]int, minVal int) map[string]int {
	cloned := maps.Clone(original)
	maps.DeleteFunc(cloned, func(k string, v int) bool {
		return v < minVal
	})
	return cloned
}

// ResolveServerPort выбирает первый ненулевой порт из переданных вариантов с помощью cmp.Or (Go 1.22+).
// Заменяет цепочки `if envPort != "" { ... } else if configPort != "" { ... }`.
func ResolveServerPort(envPort, configPort, defaultPort string) string {
	return cmp.Or(envPort, configPort, defaultPort)
}
