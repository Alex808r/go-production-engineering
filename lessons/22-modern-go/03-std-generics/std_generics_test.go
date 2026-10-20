package stdgenerics

import (
	"maps"
	"slices"
	"sort"
	"testing"
)

func TestFindProductByTitle(t *testing.T) {
	products := []Product{
		{ID: 1, Title: "Apple iPhone", Price: 999.0},
		{ID: 2, Title: "Google Pixel", Price: 799.0},
		{ID: 3, Title: "Samsung Galaxy", Price: 899.0},
	}

	idx, found := FindProductByTitle(products, "Google Pixel")
	if !found || idx != 1 {
		t.Fatalf("expected found at idx 1, got found=%v, idx=%d", found, idx)
	}

	_, found = FindProductByTitle(products, "Nonexistent")
	if found {
		t.Fatalf("expected not found for nonexistent product")
	}
}

func TestSortProductsByPrice(t *testing.T) {
	products := []Product{
		{ID: 1, Title: "A", Price: 500.0},
		{ID: 2, Title: "B", Price: 100.0},
		{ID: 3, Title: "C", Price: 300.0},
	}

	SortProductsByPrice(products)

	if products[0].Price != 100.0 || products[1].Price != 300.0 || products[2].Price != 500.0 {
		t.Errorf("incorrect sort order: %+v", products)
	}
}

func TestDeduplicateAndCompact(t *testing.T) {
	nums := []int{5, 2, 8, 2, 5, 1, 8, 9}
	res := DeduplicateAndCompact(nums)

	expected := []int{1, 2, 5, 8, 9}
	if !slices.Equal(res, expected) {
		t.Errorf("expected %v, got %v", expected, res)
	}
}

func TestCloneAndFilterMap(t *testing.T) {
	orig := map[string]int{"a": 10, "b": 50, "c": 5, "d": 100}
	filtered := CloneAndFilterMap(orig, 50)

	expected := map[string]int{"b": 50, "d": 100}
	if !maps.Equal(filtered, expected) {
		t.Errorf("expected %v, got %v", expected, filtered)
	}

	// Убеждаемся, что исходная мапа не изменилась (иммутабельность)
	if len(orig) != 4 {
		t.Errorf("original map was mutated")
	}
}

func TestResolveServerPort(t *testing.T) {
	// 1. Все заполнены -> побеждает envPort
	if port := ResolveServerPort("8080", "9090", "3000"); port != "8080" {
		t.Errorf("expected 8080, got %s", port)
	}

	// 2. envPort пустой -> побеждает configPort
	if port := ResolveServerPort("", "9090", "3000"); port != "9090" {
		t.Errorf("expected 9090, got %s", port)
	}

	// 3. envPort и configPort пустые -> побеждает defaultPort
	if port := ResolveServerPort("", "", "3000"); port != "3000" {
		t.Errorf("expected 3000, got %s", port)
	}
}

func BenchmarkSlices_SortFunc(b *testing.B) {
	source := []int{9, 4, 1, 7, 3, 8, 2, 6, 5, 0}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := slices.Clone(source)
		slices.Sort(s)
	}
}

func BenchmarkClassic_SortSlice(b *testing.B) {
	source := []int{9, 4, 1, 7, 3, 8, 2, 6, 5, 0}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]int, len(source))
		copy(s, source)
		sort.Slice(s, func(i, j int) bool {
			return s[i] < s[j]
		})
	}
}
