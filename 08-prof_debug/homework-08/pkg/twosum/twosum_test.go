package twosum

import (
	"reflect"
	"testing"
)

func TestTwoSumBuggy_FailsOnDuplicateElement(t *testing.T) {
	// Первый тестовый кейс из урока - срабатывает случайно успешно
	res1 := TwoSumBuggy([]int{2, 7, 11, 15}, 13)
	if !reflect.DeepEqual(res1, [2]int{0, 2}) {
		t.Errorf("TwoSumBuggy() got %v, want [0, 2]", res1)
	}

	// Второй кейс из урока: nums = [2, 10, 6, 12, 14], target = 20
	// Правильный ответ: 6 + 14 = 20 (индексы [2, 4]).
	// Баг: алгоритм находит nums[1] + nums[1] == 10 + 10 == 20 и возвращает [1, 1]!
	res2 := TwoSumBuggy([]int{2, 10, 6, 12, 14}, 20)
	if reflect.DeepEqual(res2, [2]int{1, 1}) {
		t.Logf("Ожидаемо воспроизведен баг: TwoSumBuggy() вернул %v (использовал один элемент дважды)", res2)
	}
}

func TestTwoSumSolutions(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   [2]int
	}{
		{
			name:   "урок пример 1 (2+11=13)",
			nums:   []int{2, 7, 11, 15},
			target: 13,
			want:   [2]int{0, 2},
		},
		{
			name:   "урок пример 2 (6+14=20, обход ложного 10+10)",
			nums:   []int{2, 10, 6, 12, 14},
			target: 20,
			want:   [2]int{2, 4},
		},
		{
			name:   "leetcode пример: 3+2+4 target 6",
			nums:   []int{3, 2, 4},
			target: 6,
			want:   [2]int{1, 2},
		},
		{
			name:   "два одинаковых числа на разных индексах",
			nums:   []int{3, 3},
			target: 6,
			want:   [2]int{0, 1},
		},
		{
			name:   "отрицательные числа (-3 + 5 = 2)",
			nums:   []int{-3, 4, 5, 90},
			target: 2,
			want:   [2]int{0, 2},
		},
		{
			name:   "решение отсутствует",
			nums:   []int{1, 2, 3},
			target: 10,
			want:   [2]int{-1, -1},
		},
	}

	for _, tt := range tests {
		t.Run("BruteForce_"+tt.name, func(t *testing.T) {
			got := TwoSumBruteForce(tt.nums, tt.target)
			if got != tt.want {
				t.Errorf("TwoSumBruteForce() = %v, want %v", got, tt.want)
			}
		})

		t.Run("Optimal_"+tt.name, func(t *testing.T) {
			got := TwoSumOptimal(tt.nums, tt.target)
			if got != tt.want {
				t.Errorf("TwoSumOptimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func generateData(size int) ([]int, int) {
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = i * 2
	}
	// Target в самом конце, чтобы протестировать худший случай O(N^2)
	target := data[size-2] + data[size-1]
	return data, target
}

func BenchmarkTwoSumBruteForce_1000(b *testing.B) {
	nums, target := generateData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = TwoSumBruteForce(nums, target)
	}
}

func BenchmarkTwoSumOptimal_1000(b *testing.B) {
	nums, target := generateData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = TwoSumOptimal(nums, target)
	}
}
