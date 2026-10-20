package homework15

import (
	"os"
	"strings"
	"testing"
)

func TestFilesExistAndNonEmpty(t *testing.T) {
	files := []string{"schema.sql", "data.sql", "queries.sql"}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("failed to read file %s: %v", f, err)
		}
		if len(data) == 0 {
			t.Errorf("file %s is empty", f)
		}
	}
}

func TestSchemaSQL_ConstraintsAndTables(t *testing.T) {
	data, err := os.ReadFile("schema.sql")
	if err != nil {
		t.Fatalf("failed to read schema.sql: %v", err)
	}
	content := strings.ToLower(string(data))

	// Проверяем наличие всех требуемых таблиц
	requiredTables := []string{
		"create table studios",
		"create table actors",
		"create table directors",
		"create table movies",
		"create table movies_actors",
		"create table movies_directors",
	}

	for _, tbl := range requiredTables {
		if !strings.Contains(content, tbl) {
			t.Errorf("schema.sql is missing table: %s", tbl)
		}
	}

	// Проверяем обязательные ограничения (constraints)
	requiredConstraints := []string{
		"check (release_year >= 1800)",
		"check (rating in ('pg-10', 'pg-13', 'pg-18'))",
		"unique (title, release_year)",
		"references studios",
	}

	for _, c := range requiredConstraints {
		if !strings.Contains(content, c) {
			t.Errorf("schema.sql is missing required constraint: %s", c)
		}
	}

	// Проверяем наличие триггера и функции
	if !strings.Contains(content, "trigger") {
		t.Error("schema.sql is missing trigger definition")
	}
}

func TestQueriesSQL_AllTasksCovered(t *testing.T) {
	data, err := os.ReadFile("queries.sql")
	if err != nil {
		t.Fatalf("failed to read queries.sql: %v", err)
	}
	content := string(data)

	// Проверяем ключевые операторы Задачи №3
	checks := []struct {
		desc string
		term string
	}{
		{"JOIN со студиями", "JOIN studios"},
		{"Выборка для актера", "Leonardo"},
		{"Подсчет фильмов режиссера", "COUNT(md.movie_id)"},
		{"Подзапрос IN", "IN ("},
		{"Фильтрация HAVING", "HAVING"},
		{"Сборы > 1000", "box_office > 1000"},
		{"Уникальные режиссеры с кассой", "COUNT(DISTINCT"},
		{"Уникальные фамилии актеров", "SELECT DISTINCT"},
		{"Дубли по названию", "HAVING COUNT(*) > 1"},
	}

	for _, c := range checks {
		if !strings.Contains(content, c.term) {
			t.Errorf("queries.sql missing requirement [%s]: keyword %q not found", c.desc, c.term)
		}
	}
}

// BenchmarkSchemaValidation измеряет скорость синтаксического анализа и валидации файла схемы.
func BenchmarkSchemaValidation(b *testing.B) {
	data, _ := os.ReadFile("schema.sql")
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		content := strings.ToLower(string(data))
		_ = strings.Contains(content, "create table movies")
		_ = strings.Contains(content, "check (release_year >= 1800)")
	}
}

// BenchmarkQueriesValidation измеряет скорость валидации аналитических запросов.
func BenchmarkQueriesValidation(b *testing.B) {
	data, _ := os.ReadFile("queries.sql")
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		content := string(data)
		_ = strings.Contains(content, "JOIN studios")
		_ = strings.Contains(content, "HAVING")
	}
}

