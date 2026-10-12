// Package storage реализует сохранение и загрузку документов через контракты io.Writer и io.Reader.
package storage

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"

	"go-core-4/06-oop/homework-06/pkg/crawler"
)

// Store выполняет сериализацию документов в формате JSON и записывает их в io.Writer.
func Store(w io.Writer, docs []crawler.Document) error {
	if w == nil {
		return fmt.Errorf("storage.Store: writer не может быть nil")
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(docs); err != nil {
		return fmt.Errorf("storage.Store: ошибка кодирования JSON: %w", err)
	}

	return nil
}

// Load считывает и десериализует срез документов из io.Reader в формате JSON.
func Load(r io.Reader) ([]crawler.Document, error) {
	if r == nil {
		return nil, fmt.Errorf("storage.Load: reader не может быть nil")
	}

	var docs []crawler.Document
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&docs); err != nil {
		return nil, fmt.Errorf("storage.Load: ошибка декодирования JSON: %w", err)
	}

	return docs, nil
}

// StoreGOB сохраняет документы в бинарном формате GOB.
func StoreGOB(w io.Writer, docs []crawler.Document) error {
	if w == nil {
		return fmt.Errorf("storage.StoreGOB: writer не может быть nil")
	}

	encoder := gob.NewEncoder(w)
	if err := encoder.Encode(docs); err != nil {
		return fmt.Errorf("storage.StoreGOB: ошибка кодирования GOB: %w", err)
	}

	return nil
}

// LoadGOB загружает документы из бинарного формата GOB.
func LoadGOB(r io.Reader) ([]crawler.Document, error) {
	if r == nil {
		return nil, fmt.Errorf("storage.LoadGOB: reader не может быть nil")
	}

	var docs []crawler.Document
	decoder := gob.NewDecoder(r)
	if err := decoder.Decode(&docs); err != nil {
		return nil, fmt.Errorf("storage.LoadGOB: ошибка декодирования GOB: %w", err)
	}

	return docs, nil
}
