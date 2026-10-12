package storage

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"go-core-4/05-io/homework-05/pkg/crawler"
)

func sampleDocs() []crawler.Document {
	return []crawler.Document{
		{ID: 0, URL: "https://go.dev", Title: "The Go Programming Language", Body: "Go is an open source programming language."},
		{ID: 1, URL: "https://golang.org", Title: "Go Documentation", Body: "Documentation and packages for Go."},
		{ID: 2, URL: "https://pkg.go.dev", Title: "Go Packages Search", Body: "Find packages and tools."},
	}
}

func TestStoreAndLoadJSON(t *testing.T) {
	orig := sampleDocs()

	var buf bytes.Buffer
	// Тестируем Store (io.Writer)
	err := Store(&buf, orig)
	if err != nil {
		t.Fatalf("Store() ошибка: %v", err)
	}

	// Тестируем Load (io.Reader)
	loaded, err := Load(&buf)
	if err != nil {
		t.Fatalf("Load() ошибка: %v", err)
	}

	if !reflect.DeepEqual(orig, loaded) {
		t.Errorf("Загруженные данные не совпадают с исходными:\nОжидалось: %+v\nПолучено:   %+v", orig, loaded)
	}
}

func TestStoreAndLoadGOB(t *testing.T) {
	orig := sampleDocs()

	var buf bytes.Buffer
	// Тестируем StoreGOB (io.Writer)
	err := StoreGOB(&buf, orig)
	if err != nil {
		t.Fatalf("StoreGOB() ошибка: %v", err)
	}

	// Тестируем LoadGOB (io.Reader)
	loaded, err := LoadGOB(&buf)
	if err != nil {
		t.Fatalf("LoadGOB() ошибка: %v", err)
	}

	if !reflect.DeepEqual(orig, loaded) {
		t.Errorf("Загруженные GOB данные не совпадают:\nОжидалось: %+v\nПолучено:   %+v", orig, loaded)
	}
}

func TestStoreAndLoadEmpty(t *testing.T) {
	var empty []crawler.Document
	var buf bytes.Buffer

	err := Store(&buf, empty)
	if err != nil {
		t.Fatalf("Store() пустого среза ошибка: %v", err)
	}

	loaded, err := Load(&buf)
	if err != nil {
		t.Fatalf("Load() пустого среза ошибка: %v", err)
	}

	if len(loaded) != 0 {
		t.Errorf("Ожидался пустой срез, получено: %d элементов", len(loaded))
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	corrupted := strings.NewReader("not a valid json string {")
	_, err := Load(corrupted)
	if err == nil {
		t.Error("Load() должен вернуть ошибку для поврежденного JSON, получено nil")
	}
}

func TestNilParams(t *testing.T) {
	if err := Store(nil, sampleDocs()); err == nil {
		t.Error("Store(nil, ...) должен вернуть ошибку")
	}
	if _, err := Load(nil); err == nil {
		t.Error("Load(nil) должен вернуть ошибку")
	}
	if err := StoreGOB(nil, sampleDocs()); err == nil {
		t.Error("StoreGOB(nil, ...) должен вернуть ошибку")
	}
	if _, err := LoadGOB(nil); err == nil {
		t.Error("LoadGOB(nil) должен вернуть ошибку")
	}
}
