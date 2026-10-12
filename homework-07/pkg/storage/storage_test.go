package storage

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"go-core-4/homework-07/pkg/crawler"
)

func TestJSONStorageRoundtrip(t *testing.T) {
	docs := []crawler.Document{
		{ID: 1, URL: "https://go.dev", Title: "Go Dev"},
		{ID: 2, URL: "https://golang.org", Title: "Golang Org"},
	}

	buf := new(bytes.Buffer)

	if err := Store(buf, docs); err != nil {
		t.Fatalf("Store() failed: %v", err)
	}

	loaded, err := Load(buf)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if !reflect.DeepEqual(docs, loaded) {
		t.Errorf("got %v, want %v", loaded, docs)
	}
}

func TestGOBStorageRoundtrip(t *testing.T) {
	docs := []crawler.Document{
		{ID: 10, URL: "https://pkg.go.dev", Title: "Pkg Go Dev"},
		{ID: 20, URL: "https://go.dev/doc", Title: "Go Documentation"},
	}

	buf := new(bytes.Buffer)

	if err := StoreGOB(buf, docs); err != nil {
		t.Fatalf("StoreGOB() failed: %v", err)
	}

	loaded, err := LoadGOB(buf)
	if err != nil {
		t.Fatalf("LoadGOB() failed: %v", err)
	}

	if !reflect.DeepEqual(docs, loaded) {
		t.Errorf("got %v, want %v", loaded, docs)
	}
}

func TestStorage_Errors(t *testing.T) {
	t.Run("nil writer and reader for JSON", func(t *testing.T) {
		if err := Store(nil, []crawler.Document{}); err == nil {
			t.Errorf("Store(nil, ...) expected error, got nil")
		}
		if _, err := Load(nil); err == nil {
			t.Errorf("Load(nil) expected error, got nil")
		}
	})

	t.Run("nil writer and reader for GOB", func(t *testing.T) {
		if err := StoreGOB(nil, []crawler.Document{}); err == nil {
			t.Errorf("StoreGOB(nil, ...) expected error, got nil")
		}
		if _, err := LoadGOB(nil); err == nil {
			t.Errorf("LoadGOB(nil) expected error, got nil")
		}
	})

	t.Run("invalid JSON input", func(t *testing.T) {
		r := strings.NewReader("not a valid json")
		if _, err := Load(r); err == nil {
			t.Errorf("Load(invalid) expected error, got nil")
		}
	})

	t.Run("invalid GOB input", func(t *testing.T) {
		r := strings.NewReader("not a valid gob stream")
		if _, err := LoadGOB(r); err == nil {
			t.Errorf("LoadGOB(invalid) expected error, got nil")
		}
	})
}
