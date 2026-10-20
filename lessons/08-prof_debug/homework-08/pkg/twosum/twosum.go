package twosum

// TwoSumBuggy демонстрирует исходный баг из материалов урока №8:
// вложенный цикл начинает j с 0 и не проверяет условие i != j,
// из-за чего при nums[i] * 2 == target функция ошибочно использует один и тот же элемент дважды.
func TwoSumBuggy(nums []int, target int) [2]int {
	for i := 0; i < len(nums); i++ {
		for j := 0; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return [2]int{i, j}
			}
		}
	}
	return [2]int{-1, -1}
}

// TwoSumBruteForce исправляет баг: j начинается с i + 1, исключая использование
// одного и того же элемента дважды и повторные проверки пар.
// Временная сложность: O(N^2), Память: O(1).
func TwoSumBruteForce(nums []int, target int) [2]int {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return [2]int{i, j}
			}
		}
	}
	return [2]int{-1, -1}
}

// TwoSumOptimal решает задачу за один линейный проход O(N) с помощью хеш-таблицы (map).
// Временная сложность: O(N), Память: O(N).
func TwoSumOptimal(nums []int, target int) [2]int {
	seen := make(map[int]int, len(nums)) // complement -> index
	for i, num := range nums {
		diff := target - num
		if idx, ok := seen[diff]; ok {
			return [2]int{idx, i}
		}
		seen[num] = i
	}
	return [2]int{-1, -1}
}
