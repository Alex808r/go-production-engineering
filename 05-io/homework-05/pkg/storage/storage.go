// Package storage реализует сохранение и загрузку документов с использованием интерфейсов io.Writer и io.Reader.
package storage

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"

	"go-core-4/05-io/homework-05/pkg/crawler"
)

// Store выполняет сериализацию документов в формате JSON и записывает их в io.Writer.
// Задача №1: Механизм записи данных сканера с использованием io.Writer.
func Store(w io.Writer, docs []crawler.Document) error {
	if w == nil {
		return fmt.Errorf("storage.Store: writer не может быть nil")
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ") // Красивое форматирование для удобства чтения файла
	if err := encoder.Encode(docs); err != nil {
		return fmt.Errorf("storage.Store: ошибка кодирования JSON: %w", err)
	}

	return nil
}

// Load считывает и десериализует срез документов из io.Reader в формате JSON.
// Задача №2: Механизм загрузки данных из файла с использованием io.Reader.
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

// StoreGOB сохраняет документы в бинарном формате GOB с использованием io.Writer.
// GOB — это нативный бинарный формат сериализации Go-to-Go, разобранный в лекции.
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

// LoadGOB загружает документы из бинарного формата GOB с использованием io.Reader.
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
