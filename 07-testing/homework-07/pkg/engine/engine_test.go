package engine

import (
	"bytes"
	"testing"

	"go-core-4/07-testing/homework-07/pkg/crawler"
	"go-core-4/07-testing/homework-07/pkg/crawler/membot"
	"go-core-4/07-testing/homework-07/pkg/index"
)

func TestEngine_ScanAndSearch(t *testing.T) {
	bot := membot.New()
	idx := index.New()
	eng := New(bot, idx)

	err := eng.Scan([]string{"https://go.dev"}, 1)
	if err != nil {
		t.Fatalf("eng.Scan() failed: %v", err)
	}

	if eng.Count() == 0 {
		t.Fatalf("eng.Count() = 0, want > 0")
	}

	// Поиск слова "golang"
	results := eng.Search("golang")
	if len(results) == 0 {
		t.Errorf("eng.Search(\"golang\") returned 0 results, want > 0")
	}

	// Поиск отсутствующего слова
	nonExistent := eng.Search("kubernetes")
	if len(nonExistent) != 0 {
		t.Errorf("eng.Search(\"kubernetes\") = %v, want empty", nonExistent)
	}
}

func TestEngine_SaveAndLoad(t *testing.T) {
	bot := membot.New()
	eng1 := New(bot, nil)
	_ = eng1.Scan([]string{""}, 1)

	buf := new(bytes.Buffer)
	if err := eng1.Save(buf); err != nil {
		t.Fatalf("eng1.Save() failed: %v", err)
	}

	// Создаем второй движок без краулера (чистый поисковик из кэша)
	eng2 := New(nil, nil)
	if err := eng2.Load(buf); err != nil {
		t.Fatalf("eng2.Load() failed: %v", err)
	}

	if eng2.Count() != eng1.Count() {
		t.Errorf("eng2.Count() = %d, want %d", eng2.Count(), eng1.Count())
	}

	results := eng2.Search("yandex")
	if len(results) == 0 {
		t.Errorf("eng2.Search(\"yandex\") returned 0 results, want > 0")
	}
}

func TestEngine_Encapsulation(t *testing.T) {
	eng := New(nil, nil)
	docs := []crawler.Document{
		{ID: 0, Title: "Test Document", URL: "https://example.com"},
	}
	eng.SetDocuments(docs)

	retrieved := eng.Documents()
	retrieved[0].Title = "Modified Title"

	internalDocs := eng.Documents()
	if internalDocs[0].Title != "Test Document" {
		t.Errorf("Внутреннее состояние движка было мутировано извне! Title = %q", internalDocs[0].Title)
	}
}

func TestEngine_NilCrawlerError(t *testing.T) {
	eng := New(nil, nil)
	err := eng.Scan([]string{"https://example.com"}, 1)
	if err == nil {
		t.Errorf("eng.Scan() with nil crawler expected error, got nil")
	}
}

func TestEngine_SaveLoadErrors(t *testing.T) {
	eng := New(nil, nil)
	if err := eng.Save(nil); err == nil {
		t.Errorf("eng.Save(nil) expected error, got nil")
	}
	if err := eng.Load(nil); err == nil {
		t.Errorf("eng.Load(nil) expected error, got nil")
	}
}
