package search

// Binary выполняет логарифмический бинарный поиск O(log N) в отсортированном срезе.
// Возвращает индекс элемента или -1, если элемент не найден.
func Binary(data []int, item int) int {
	low, high := 0, len(data)-1
	for low <= high {
		mid := int(uint(low+high) >> 1) // Защита от переполнения int
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

// Linear выполняет линейный последовательный поиск O(N).
// Возвращает индекс первого найденного элемента или -1.
func Linear(data []int, item int) int {
	for i, v := range data {
		if v == item {
			return i
		}
	}
	return -1
}
