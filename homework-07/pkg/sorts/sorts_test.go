package sorts

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

// Задача №1 (Часть 1): Простой тест для функции sort.Ints().
func TestSortInts(t *testing.T) {
	data := []int{42, -5, 0, 17, -100, 8, 3}
	want := []int{-100, -5, 0, 3, 8, 17, 42}

	sort.Ints(data)

	if !sort.IntsAreSorted(data) {
		t.Errorf("sort.IntsAreSorted(data) = false, want true")
	}

	if !reflect.DeepEqual(data, want) {
		t.Errorf("sort.Ints() got %v, want %v", data, want)
	}
}

// Задача №1 (Часть 2): Табличный тест для функции sort.Strings().
func TestSortStrings(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "пустой срез",
			in:   []string{},
			want: []string{},
		},
		{
			name: "nil срез",
			in:   nil,
			want: nil,
		},
		{
			name: "один элемент",
			in:   []string{"golang"},
			want: []string{"golang"},
		},
		{
			name: "уже отсортированный срез",
			in:   []string{"apple", "banana", "cherry", "date"},
			want: []string{"apple", "banana", "cherry", "date"},
		},
		{
			name: "обратный порядок",
			in:   []string{"zebra", "tiger", "lion", "cat"},
			want: []string{"cat", "lion", "tiger", "zebra"},
		},
		{
			name: "дубликаты и повторяющиеся строки",
			in:   []string{"go", "rust", "go", "python", "rust", "c"},
			want: []string{"c", "go", "go", "python", "rust", "rust"},
		},
		{
			name: "регистрозависимость (ASCII порядок: заглавные перед строчными)",
			in:   []string{"banana", "Apple", "apple", "Banana"},
			want: []string{"Apple", "Banana", "apple", "banana"},
		},
		{
			name: "кириллица (Unicode UTF-8 побайтовый лексикографический порядок)",
			in:   []string{"яблоко", "арбуз", "банан", "груша"},
			want: []string{"арбуз", "банан", "груша", "яблоко"},
		},
		{
			name: "общие префиксы разной длины",
			in:   []string{"testing", "test", "tester", "testcase"},
			want: []string{"test", "testcase", "tester", "testing"},
		},
		{
			name: "строки со спецсимволами, цифрами и пробелами",
			in:   []string{" 10", "!important", "01", " 02", "@user"},
			want: []string{" 02", " 10", "!important", "01", "@user"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Делаем копию входного среза, чтобы не мутировать тестовую таблицу
			var input []string
			if tt.in != nil {
				input = make([]string, len(tt.in))
				copy(input, tt.in)
			}

			sort.Strings(input)

			if !sort.StringsAreSorted(input) {
				t.Errorf("StringsAreSorted() = false for input %v", input)
			}

			if !reflect.DeepEqual(input, tt.want) {
				t.Errorf("sort.Strings() = %v, want %v", input, tt.want)
			}
		})
	}
}

// Задача №2: Бенчмарки для sort.Ints() и sort.Float64s().
// Сравнение производительности на различных объемах данных: 10, 100, 1_000, 10_000 элементов.

func BenchmarkSortInts(b *testing.B) {
	sizes := []int{10, 100, 1_000, 10_000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("N=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			master := GenerateInts(size)
			buf := make([]int, size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				copy(buf, master)
				sort.Ints(buf)
			}
		})
	}
}

func BenchmarkSortFloat64s(b *testing.B) {
	sizes := []int{10, 100, 1_000, 10_000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("N=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			master := GenerateFloat64s(size)
			buf := make([]float64, size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				copy(buf, master)
				sort.Float64s(buf)
			}
		})
	}
}
