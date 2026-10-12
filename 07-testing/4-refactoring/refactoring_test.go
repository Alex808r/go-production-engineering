package refactoring

import (
	"testing"
)

// Test_rublesToUSD — интеграционный тест, зависящий от внешнего сетевого ресурса.
//
// ИСПОЛЬЗОВАНИЕ t.Skip():
// Если внешний API ЦБ РФ недоступен, тест не должен ломать сборку проекта в CI/CD (flaky test).
// Метод `t.Skipf(...)` помечает тест как пропущенный (статус SKIP), позволяя продолжить прогон остальных тестов.
func Test_rublesToUSD(t *testing.T) {
	const rubles = 1000
	res, err := rublesToUSD(rubles)
	if err != nil {
		t.Skipf("пропуск сетевого теста (внешний API курсов валют ЦБ недоступен: %v)", err)
	}
	t.Logf("%v рублей - это %.2f долларов\n", rubles, res)
}

// Test_calc — быстрый, изолированный модульный тест чистой функции.
// Не требует сетевого соединения, работает за микросекунды и гарантирует
// правильность математики независимо от доступности интернета.
func Test_calc(t *testing.T) {
	want := 10.0
	got := calc(1000, 100)
	if got != want {
		t.Fatalf("calc(1000, 100) = %v; want %v", got, want)
	}
}

