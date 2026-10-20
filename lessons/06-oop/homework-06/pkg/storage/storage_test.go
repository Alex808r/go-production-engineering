package storage

import (
	"bytes"
	"reflect"
	"testing"

	"go-core-4/lessons/06-oop/homework-06/pkg/crawler"
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
