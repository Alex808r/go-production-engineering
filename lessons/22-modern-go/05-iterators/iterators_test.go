package iterators

import (
	"slices"
	"testing"
)

func buildSampleTree() *TreeNode[int] {
	var root *TreeNode[int]
	values := []int{50, 30, 70, 20, 40, 60, 80}
	for _, v := range values {
		root = root.Insert(v)
	}
	return root
}

func TestTree_All_InOrder(t *testing.T) {
	root := buildSampleTree()

	var collected []int
	// Прямой обход через range (Go 1.23+ Range over func)!
	for val := range root.All() {
		collected = append(collected, val)
	}

	expected := []int{20, 30, 40, 50, 60, 70, 80}
	if !slices.Equal(collected, expected) {
		t.Fatalf("expected %v, got %v", expected, collected)
	}
}

func TestTree_All_EarlyBreak(t *testing.T) {
	root := buildSampleTree()

	var collected []int
	for val := range root.All() {
		collected = append(collected, val)
		if val == 40 {
			break // Проверяем досрочную остановку итерации
		}
	}

	expected := []int{20, 30, 40}
	if !slices.Equal(collected, expected) {
		t.Fatalf("expected %v, got %v", expected, collected)
	}
}

func TestFilter_Iterator(t *testing.T) {
	root := buildSampleTree()

	// Итератор, возвращающий только четные числа > 30
	evens := Filter(root.All(), func(x int) bool {
		return x%2 == 0 && x > 30
	})

	var collected []int
	for val := range evens {
		collected = append(collected, val)
	}

	expected := []int{40, 50, 60, 70, 80}
	if !slices.Equal(collected, expected) {
		t.Fatalf("expected %v, got %v", expected, collected)
	}
}

func TestPairs_Seq2(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	count := 0
	for k, v := range Pairs(m) {
		if m[k] != v {
			t.Errorf("mismatch for key %s", k)
		}
		count++
	}
	if count != 2 {
		t.Errorf("expected 2 pairs, got %d", count)
	}
}

func BenchmarkTree_IterSeq(b *testing.B) {
	root := buildSampleTree()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		for val := range root.All() {
			sum += val
		}
	}
}

func BenchmarkTree_FullSlice(b *testing.B) {
	root := buildSampleTree()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		s := root.Slice()
		for _, val := range s {
			sum += val
		}
	}
}
