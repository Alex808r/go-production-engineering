package fuzzing

import (
	"testing"
)

func TestParseQuery_Valid(t *testing.T) {
	input := "name=alice&role=admin&city=moscow"
	res, err := ParseQuery(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res["name"] != "alice" || res["role"] != "admin" || res["city"] != "moscow" {
		t.Errorf("mismatch: %+v", res)
	}
}

func TestParseQuery_Invalid(t *testing.T) {
	cases := []string{
		"",
		"invalid_pair",
		"=val",
		"key=",
		"\xff\xfe\xfd", // invalid utf-8
	}

	for _, c := range cases {
		_, err := ParseQuery(c)
		if err == nil {
			t.Errorf("expected error for input %q, got nil", c)
		}
	}
}

// FuzzParseQuery — фаззинг-тест, использующий встроенный фаззер Go (testing.F).
// Фаззер мутирует сид-корпус (Seed Corpus) и проверяет инвариант:
// парсер НИКОГДА не должен паниковать на произвольных входных байтах.
func FuzzParseQuery(f *testing.F) {
	// Добавляем сид-корпус (начальные эталонные значения)
	f.Add("name=alice&age=30")
	f.Add("key=val")
	f.Add("k1=v1&foo=bar")
	f.Add("&&&===&&&")
	f.Add("\x00\x01\x02")

	f.Fuzz(func(t *testing.T, input string) {
		// Главный инвариант надежности: вызов ParseQuery не должен приводить к панике (panic)
		res, err := ParseQuery(input)
		if err == nil {
			// Если парсинг успешен, результат должен содержать как минимум один ключ
			if len(res) == 0 {
				t.Errorf("успешный парсинг вернул пустой результат для %q", input)
			}
		}
	})
}
