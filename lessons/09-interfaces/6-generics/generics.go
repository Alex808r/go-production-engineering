// Package main демонстрирует основы параметрического полиморфизма (Generics) в Go (1.18+):
// обобщенные функции, ограничения типов (type constraints), оператор тильды `~` и вывод типов (type inference).
package main

import (
	"cmp"
	"fmt"
	"slices"
)

// CustomInt — пользовательский тип на базе int.
type CustomInt int

// ordered — кастомный интерфейс-ограничение (type constraint).
// Знак тильды `~` (approximation element) означает: "любой тип, базовым типом которого является указанный".
// Без `~` тип CustomInt не удовлетворял бы ограничению int!
type ordered interface {
	~int | ~int64 | ~float64 | ~float32 | ~byte | ~string
}

// min — обобщенная функция для поиска минимального элемента в срезе.
// Синтаксис:
//   - [T ordered] — список параметров типа (type parameters) с ограничением ordered.
//   - in []T — срез элементов типа T.
//   - (T, bool) — возвращает минимальный элемент типа T и флаг успешности.
func min[T ordered](in []T) (T, bool) {
	var res T
	if len(in) == 0 {
		return res, false
	}
	res = in[0]
	for _, v := range in {
		// Благодаря ограничению ordered компилятор знает, что для T допустима операция сравнения `<`
		if v < res {
			res = v
		}
	}
	return res, true
}

// maxGeneric — пример использования стандартного ограничения cmp.Ordered (Go 1.21+).
func maxGeneric[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func main() {
	// 1. Обобщенная функция со стандартными типами
	ints := []int{-2, 0, 2, 3}
	floats := []float32{2.3, 3.4, 4.5}
	strings := []string{"def", "Абв", "ABC"}

	// Явное указание параметра типа: min[int](ints)
	if val, ok := min[int](ints); ok {
		fmt.Printf("Min int: %v\n", val)
	}

	// Неявный вывод типа компилятором (Type Inference): min(floats)
	if val, ok := min(floats); ok {
		fmt.Printf("Min float32: %v\n", val)
	}

	if val, ok := min(strings); ok {
		fmt.Printf("Min string: %v\n", val)
	}

	// 2. Использование кастомного типа с оператором тильды `~`
	customSlice := []CustomInt{100, 42, 500}
	if val, ok := min(customSlice); ok {
		fmt.Printf("Min CustomInt (благодаря ~int): %v\n", val)
	}

	// 3. Стандартная библиотека slices (Go 1.21+)
	var slice = []int{1, 2, 3}
	fmt.Printf("slices.Contains(slice, 2): %v\n", slices.Contains(slice, 2))
	fmt.Printf("slices.Contains(slice, 20): %v\n", slices.Contains(slice, 20))

	// 4. Обобщенная функция cmp.Ordered
	fmt.Printf("maxGeneric(10, 25): %v\n", maxGeneric(10, 25))
	fmt.Printf("maxGeneric(\"apple\", \"banana\"): %v\n", maxGeneric("apple", "banana"))
}

