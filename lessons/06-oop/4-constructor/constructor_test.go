package guitar

import (
	"testing"
)

// TestNew проверяет корректную инициализацию объекта через конструктор New
// и форматирование информации методом Info().
func TestNew(t *testing.T) {
	g := New("Fender", "Stratocaster")

	got := g.Info()
	want := "Гитара Fender Stratocaster"
	if got != want {
		t.Errorf("получили %s, ожидалось %s", got, want)
	}

	// Проверяем, что мапа Tabs инициализирована конструктором
	// и запись в нее не приводит к панике
	if g.Tabs == nil {
		t.Fatal("конструктор New() должен инициализировать поле Tabs")
	}
	g.Tabs["Deep Purple"] = "Smoke on the Water"
}

// TestGuitar демонстрирует ручную инициализацию структуры без конструктора.
// Если объявить `var g Guitar`, то `g.Tabs == nil`.
// Без вызова `make(map[string]string)` последующая запись `g.Tabs["..."] = "..."`
// завершилась бы с `panic: assignment to entry in nil map`.
func TestGuitar(t *testing.T) {
	var g Guitar

	// Без этой строки возникнет паника:
	g.Tabs = make(map[string]string)
	g.Tabs["Metallica"] = "One"

	if g.Tabs["Metallica"] != "One" {
		t.Errorf("ожидалась табулатура 'One', получили %s", g.Tabs["Metallica"])
	}
}

