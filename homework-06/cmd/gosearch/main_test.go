package main

import (
	"os"
	"path/filepath"
	"testing"

	"go-core-4/homework-06/pkg/crawler/membot"
	"go-core-4/homework-06/pkg/engine"
	"go-core-4/homework-06/pkg/index"
)

func TestFileExists(t *testing.T) {
	tempDir := t.TempDir()
	existingFile := filepath.Join(tempDir, "test.json")

	if fileExists(existingFile) {
		t.Error("Файл еще не создан, но fileExists вернул true")
	}

	err := os.WriteFile(existingFile, []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("Ошибка создания файла: %v", err)
	}

	if !fileExists(existingFile) {
		t.Error("Файл создан, но fileExists вернул false")
	}

	if fileExists(tempDir) {
		t.Error("Директория не должна считаться файлом")
	}
}

func TestEngineIntegration(t *testing.T) {
	tempDir := t.TempDir()
	cachePath := filepath.Join(tempDir, "cache.json")

	bot := membot.New()
	idx := index.New()
	eng := engine.New(bot, idx)

	err := eng.Scan([]string{""}, 1)
	if err != nil {
		t.Fatalf("eng.Scan() failed: %v", err)
	}

	fWrite, err := os.Create(cachePath)
	if err != nil {
		t.Fatalf("os.Create() failed: %v", err)
	}
	if err := eng.Save(fWrite); err != nil {
		t.Fatalf("eng.Save() failed: %v", err)
	}
	fWrite.Close()

	// Загружаем в новый объект
	engLoaded := engine.New(nil, nil)
	fRead, err := os.Open(cachePath)
	if err != nil {
		t.Fatalf("os.Open() failed: %v", err)
	}
	defer fRead.Close()

	if err := engLoaded.Load(fRead); err != nil {
		t.Fatalf("engLoaded.Load() failed: %v", err)
	}

	if engLoaded.Count() != eng.Count() {
		t.Errorf("engLoaded.Count() = %d, want %d", engLoaded.Count(), eng.Count())
	}

	results := engLoaded.Search("google")
	if len(results) == 0 {
		t.Error("engLoaded.Search(\"google\") returned 0 results")
	}
}

func TestExecuteSearch(t *testing.T) {
	bot := membot.New()
	eng := engine.New(bot, nil)
	_ = eng.Scan([]string{""}, 1)

	executeSearch(eng, "golang")
	executeSearch(eng, "notfoundquery")
}
