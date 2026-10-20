// Package bsearch демонстрирует алгоритмы поиска для сравнительного профилирования с помощью pprof.
package bsearch

// Binary реализует алгоритм бинарного поиска (двоичный поиск) O(log N) в отсортированном срезе.
//
// При профилировании CPU этот алгоритм выполняет минимальное число инструкций процессора
// (для N = 1_000_000 достаточно не более 20 итераций).
// Возвращает индекс найденного элемента или -1.
func Binary(data []int, item int) int {
	low, high := 0, len(data)-1
	for low <= high {
		mid := (low + high) / 2
		if data[mid] == item {
			return mid
		}
		if data[mid] < item {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1
}

// Simple реализует алгоритм линейного поиска O(N).
//
// При профилировании CPU этот алгоритм четко выделяется на FlameGraph и в топе функций (Hotspot),
// так как в среднем выполняет ~500_000 сравнений на каждый вызов.
// Возвращает индекс найденного элемента или -1.
func Simple(data []int, item int) int {
	for i := range data {
		if data[i] == item {
			return i
		}
	}
	return -1
}

