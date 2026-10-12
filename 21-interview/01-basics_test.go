package interview

import (
	"testing"
)

func TestTypicalErrorsNoSync(t *testing.T) {
	TypicalErrorsNoSync()
}

func TestTypicalErrorsIndex(t *testing.T) {
	TypicalErrorsIndex()
}

func TestSlices(t *testing.T) {
	Slices()
}

func TestStrings(t *testing.T) {
	Strings()
}

func TestNilInterface(t *testing.T) {
	// Демонстрация ловушки: интерфейс не nil, несмотря на nil внутри!
	err := FailReturningCustomError()
	if err == nil {
		t.Fatalf("ожидалось, что типизированный nil-интерфейс НЕ равен nil")
	}

	// Корректная обработка: возвращен честный нетипизированный nil
	errCorrect := CorrectFail()
	if errCorrect != nil {
		t.Fatalf("ожидалось nil для CorrectFail, получено: %v", errCorrect)
	}
}

func TestDeferWithNamedReturn(t *testing.T) {
	got := DeferWithNamedReturn()
	if got != 20 {
		t.Fatalf("ожидалось 20, получено: %d", got)
	}
}

func BenchmarkSliceAppend_NoPrealloc(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var s []int
		for j := 0; j < 1000; j++ {
			s = append(s, j)
		}
	}
}

func BenchmarkSliceAppend_WithPrealloc(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]int, 0, 1000)
		for j := 0; j < 1000; j++ {
			s = append(s, j)
		}
	}
}
