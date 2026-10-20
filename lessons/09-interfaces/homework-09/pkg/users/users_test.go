package users

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func TestMaxAge(t *testing.T) {
	tests := []struct {
		name   string
		people []Ager
		want   int
	}{
		{
			name: "сотрудники и клиенты вместе",
			people: []Ager{
				NewEmployee("Алексей", 28),
				NewCustomer("Борис", 45),
				NewEmployee("Виктория", 32),
				NewCustomer("Григорий", 19),
			},
			want: 45,
		},
		{
			name: "только один человек",
			people: []Ager{
				NewCustomer("Ольга", 29),
			},
			want: 29,
		},
		{
			name:   "пустой список",
			people: []Ager{},
			want:   0,
		},
		{
			name: "список с nil-элементом",
			people: []Ager{
				NewEmployee("Иван", 25),
				nil,
				NewCustomer("Мария", 30),
			},
			want: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxAge(tt.people...)
			if got != tt.want {
				t.Errorf("MaxAge() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestOldestPerson(t *testing.T) {
	emp1 := EmployeeSimple{Name: "Иван", Age: 30}
	emp2 := EmployeeSimple{Name: "Дмитрий", Age: 50}
	cust1 := CustomerSimple{Name: "Анна", Age: 25}
	cust2 := CustomerSimple{Name: "Сергей", Age: 65}

	tests := []struct {
		name   string
		people []any
		want   any
	}{
		{
			name:   "поиск самого старшего среди разных типов",
			people: []any{emp1, cust1, emp2, cust2, "лишняя строка", 12345, true},
			want:   cust2,
		},
		{
			name:   "самый старший — сотрудник",
			people: []any{emp1, cust1, emp2},
			want:   emp2,
		},
		{
			name:   "пустой список аргументов",
			people: []any{},
			want:   nil,
		},
		{
			name:   "только посторонние типы данных",
			people: []any{"hello", 42, 3.14, false},
			want:   nil,
		},
		{
			name:   "указатели на структуры",
			people: []any{&emp1, &cust2},
			want:   &cust2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OldestPerson(tt.people...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OldestPerson() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

// errWriter — вспомогательный writer для тестирования сценариев ошибок
type errWriter struct{}

func (e *errWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("simulated write error")
}

func TestWriteOnlyStrings(t *testing.T) {
	t.Run("фильтрация и запись только строк", func(t *testing.T) {
		var buf bytes.Buffer

		ptrStr := " "
		items := []any{
			"Go",
			42,
			ptrStr,
			true,
			EmployeeSimple{Name: "Тест", Age: 20},
			"Interfaces",
			3.14,
			"!",
		}

		n, err := WriteOnlyStrings(&buf, items...)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expectedContent := "Go Interfaces!"
		if buf.String() != expectedContent {
			t.Errorf("buf.String() = %q, want %q", buf.String(), expectedContent)
		}
		if n != len(expectedContent) {
			t.Errorf("written bytes = %d, want %d", n, len(expectedContent))
		}
	})

	t.Run("нет строковых аргументов", func(t *testing.T) {
		var buf bytes.Buffer
		n, err := WriteOnlyStrings(&buf, 1, 2, 3, false, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n != 0 || buf.Len() != 0 {
			t.Errorf("expected empty buffer, got %d bytes: %q", n, buf.String())
		}
	})

	t.Run("nil writer возвращает ошибку", func(t *testing.T) {
		_, err := WriteOnlyStrings(nil, "test")
		if err == nil {
			t.Errorf("expected error for nil writer, got nil")
		}
	})

	t.Run("ошибка записи в writer", func(t *testing.T) {
		ew := &errWriter{}
		_, err := WriteOnlyStrings(ew, "error should occur")
		if err == nil {
			t.Errorf("expected write error, got nil")
		}
	})
}

// ============================================================================
// Бенчмарки производительности интерфейсов
// ============================================================================

func BenchmarkMaxAge(b *testing.B) {
	people := make([]Ager, 1000)
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			people[i] = NewEmployee("Employee", 20+(i%50))
		} else {
			people[i] = NewCustomer("Customer", 18+(i%60))
		}
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = MaxAge(people...)
	}
}

func BenchmarkOldestPerson(b *testing.B) {
	people := make([]any, 1000)
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			people[i] = EmployeeSimple{Name: "Employee", Age: 20 + (i % 50)}
		} else {
			people[i] = CustomerSimple{Name: "Customer", Age: 18 + (i % 60)}
		}
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = OldestPerson(people...)
	}
}

func BenchmarkWriteOnlyStrings(b *testing.B) {
	var buf bytes.Buffer
	items := []any{"Go", 100, "Interfaces", true, 3.14, "Benchmark", nil}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		_, _ = WriteOnlyStrings(&buf, items...)
	}
}

