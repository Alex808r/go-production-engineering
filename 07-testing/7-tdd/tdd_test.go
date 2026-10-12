package tdd

import "testing"

// TestFact иллюстрирует цикл разработки по методологии TDD (Test-Driven Development):
//
// ФАЗЫ ЦИКЛА TDD (Red -> Green -> Refactor):
// 1. RED (Красный): Пишется минимальный тест для новой функциональности до написания самого кода.
//    Тест запускается и гарантированно падает (так как кода еще нет или он не полон).
// 2. GREEN (Зеленый): Пишется простейший рабочий код, достаточный ТОЛЬКО для прохождения теста
//    (без преждевременных оптимизаций и оверинжиниринга).
// 3. REFACTOR (Рефакторинг): Код улучшается, устраняются дублирования, повышается читаемость.
//    Тесты при этом остаются "зелеными", гарантируя сохранение требуемого поведения.
func TestFact(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Факториал 1! (базовый случай)",
			args: args{n: 1},
			want: 1,
		},
		{
			name: "Факториал 3! = 3 * 2 * 1 = 6",
			args: args{n: 3},
			want: 6,
		},
		{
			name: "Факториал 0! = 1 (математическое определение)",
			args: args{n: 0},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Fact(tt.args.n); got != tt.want {
				t.Errorf("Fact(%d) = %v, want %v", tt.args.n, got, tt.want)
			}
		})
	}
}

