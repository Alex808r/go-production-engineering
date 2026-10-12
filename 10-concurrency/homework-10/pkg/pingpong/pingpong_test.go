package pingpong

import (
	"bytes"
	"strings"
	"testing"
)

func TestPlayGame_CompletesAndWinnerValid(t *testing.T) {
	var buf bytes.Buffer
	targetScore := 5

	cfg := Config{
		Player1Name: "Тимо Болль",
		Player2Name: "Ма Лун",
		TargetScore: targetScore,
		LuckyChance: 0.25,
		Output:      &buf,
		Seed:        12345,
	}

	result := PlayGame(cfg)

	// Проверяем, что победитель набрал ровно targetScore
	if result.Scores[result.Winner] != targetScore {
		t.Errorf("Ожидался счёт победителя %d, получено %d", targetScore, result.Scores[result.Winner])
	}

	// Проверяем, что второй игрок набрал меньше targetScore
	otherPlayer := cfg.Player1Name
	if result.Winner == cfg.Player1Name {
		otherPlayer = cfg.Player2Name
	}

	if result.Scores[otherPlayer] >= targetScore {
		t.Errorf("Счёт проигравшего (%d) не должен превышать или равняться targetScore (%d)",
			result.Scores[otherPlayer], targetScore)
	}

	// Сумма очков должна равняться общему количеству подач
	totalPoints := result.Scores[cfg.Player1Name] + result.Scores[cfg.Player2Name]
	if totalPoints != result.Rounds {
		t.Errorf("Сумма очков (%d) должна быть равна количеству подач (%d)", totalPoints, result.Rounds)
	}
}

func TestPlayGame_OutputFormattingAndPhrases(t *testing.T) {
	var buf bytes.Buffer
	cfg := Config{
		Player1Name: "Игрок 1",
		Player2Name: "Игрок 2",
		TargetScore: 2,
		LuckyChance: 0.30,
		Output:      &buf,
		Seed:        999,
	}

	result := PlayGame(cfg)
	out := buf.String()

	// Проверяем наличие всех ключевых фраз по ТЗ
	requiredPhrases := []string{
		"Начало матча по пинг-понгу",
		"Игрок 1",
		"Игрок 2",
		"ping",
		"мощно загасил шарик в самый угол",
		"Партия окончена",
	}

	for _, phrase := range requiredPhrases {
		if !strings.Contains(out, phrase) {
			t.Errorf("Вывод не содержит обязательную фразу %q", phrase)
		}
	}

	if result.Winner == "" {
		t.Error("Победитель матча не определен")
	}
}

func TestPlayGame_DefaultConfig(t *testing.T) {
	// Проверка работы с пустым конфигом (дефолтные значения)
	var buf bytes.Buffer
	cfg := Config{
		TargetScore: 1, // 1 очко для быстрого завершения
		Output:      &buf,
		Seed:        42,
	}

	result := PlayGame(cfg)

	if result.Winner != "Игрок 1" && result.Winner != "Игрок 2" {
		t.Errorf("Некорректный победитель при дефолтных именах: %s", result.Winner)
	}

	if result.Rounds < 1 {
		t.Errorf("Должен быть сыгран как минимум 1 раунд, сыграно: %d", result.Rounds)
	}
}

func TestPlayGame_MultipleGamesUnderRaceDetector(t *testing.T) {
	// Стресс-тест: несколько игр подряд для проверки race detector
	for i := 0; i < 5; i++ {
		var buf bytes.Buffer
		res := PlayGame(Config{
			Player1Name: "Альфа",
			Player2Name: "Бета",
			TargetScore: 3,
			LuckyChance: 0.20,
			Output:      &buf,
			Seed:        int64(i * 100),
		})

		if res.Winner != "Альфа" && res.Winner != "Бета" {
			t.Fatalf("Некорректный победитель: %v", res.Winner)
		}
		if res.Scores[res.Winner] != 3 {
			t.Fatalf("Победитель должен набрать 3 очка, набрал %d", res.Scores[res.Winner])
		}
	}
}

func BenchmarkPlayGame(b *testing.B) {
	cfg := Config{
		Player1Name: "Игрок 1",
		Player2Name: "Игрок 2",
		TargetScore: 11,
		LuckyChance: 0.20,
		Output:      nil, // io.Discard
		Seed:        12345,
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = PlayGame(cfg)
	}
}

