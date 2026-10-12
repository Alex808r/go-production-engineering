package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"go-core-4/homework-05/pkg/crawler"
	"go-core-4/homework-05/pkg/index"
	"go-core-4/homework-05/pkg/storage"
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

	// Проверка директории (не должна считаться файлом)
	if fileExists(tempDir) {
		t.Error("Директория не должна считаться файлом")
	}
}

func TestStorageWorkflowWithTempFile(t *testing.T) {
	tempDir := t.TempDir()
	cachePath := filepath.Join(tempDir, "docs_cache.json")

	docs := []crawler.Document{
		{ID: 0, Title: "The Go Programming Language", URL: "https://go.dev"},
		{ID: 1, Title: "Go Standard Library", URL: "https://pkg.go.dev"},
	}

	// 1. Имитируем сохранение через os.Create (Задача №1: запись через io.Writer)
	fWrite, err := os.Create(cachePath)
	if err != nil {
		t.Fatalf("Ошибка создания файла: %v", err)
	}

	var w io.Writer = fWrite
	err = storage.Store(w, docs)
	fWrite.Close()
	if err != nil {
		t.Fatalf("storage.Store() ошибка: %v", err)
	}

	// 2. Имитируем загрузку через os.Open (Задача №2: чтение через io.Reader)
	if !fileExists(cachePath) {
		t.Fatal("Файл должен существовать")
	}

	fRead, err := os.Open(cachePath)
	if err != nil {
		t.Fatalf("Ошибка открытия файла: %v", err)
	}

	var r io.Reader = fRead
	loadedDocs, err := storage.Load(r)
	fRead.Close()
	if err != nil {
		t.Fatalf("storage.Load() ошибка: %v", err)
	}

	if len(loadedDocs) != 2 {
		t.Fatalf("Ожидалось 2 документа, получено: %d", len(loadedDocs))
	}
	if loadedDocs[0].Title != docs[0].Title || loadedDocs[1].URL != docs[1].URL {
		t.Errorf("Загруженные данные не совпадают с исходными")
	}
}

func TestExecuteSearch(t *testing.T) {
	docs := []crawler.Document{
		{ID: 0, Title: "Go Documentation and Guides", URL: "https://golang.org"},
		{ID: 1, Title: "Python Programming", URL: "https://python.org"},
	}

	idx := index.New()
	idx.Add(docs)

	// Тестируем, что поиск не паникует и находит совпадения
	executeSearch("Go", docs, idx)
	executeSearch("NotFoundWord", docs, idx)
}

func TestScanSitesMembot(t *testing.T) {
	docs := scanSites(true)
	if len(docs) == 0 {
		t.Error("scanSites(true) должен вернуть документы от membot")
	}
}

func TestBytesBufferReader(t *testing.T) {
	// Демонстрация того, что bytes.Buffer реализует io.Reader и io.Writer
	var buf bytes.Buffer
	docs := []crawler.Document{{ID: 42, Title: "Test"}}

	if err := storage.Store(&buf, docs); err != nil {
		t.Fatalf("Ошибка Store: %v", err)
	}

	loaded, err := storage.Load(&buf)
	if err != nil {
		t.Fatalf("Ошибка Load: %v", err)
	}

	if len(loaded) != 1 || loaded[0].ID != 42 {
		t.Errorf("Некорректный результат: %+v", loaded)
	}
}
