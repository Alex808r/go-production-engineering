// Package webapp предоставляет веб-интерфейс для просмотра содержимого
// инвертированного поискового индекса и хранилища проиндексированных документов.
package webapp

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"

	"github.com/gorilla/mux"
)

// Document представляет проиндексированный документ для веб-отображения.
type Document struct {
	ID    int    `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
}

// IndexEntry представляет запись поискового индекса для шаблонизатора HTML.
type IndexEntry struct {
	Word    string `json:"word"`
	DocIDs  []int  `json:"doc_ids"`
	Count   int    `json:"count"`
}

// DataProvider определяет контракт источника данных для веб-приложения.
type DataProvider interface {
	Documents() []Document
	Index() map[string][]int
}

// Service представляет HTTP-службу веб-приложения.
type Service struct {
	router   *mux.Router
	provider DataProvider
	tmplDocs *template.Template
	tmplIdx  *template.Template
}

// New создает и настраивает новый экземпляр веб-приложения webapp.
func New(provider DataProvider) *Service {
	s := &Service{
		router:   mux.NewRouter(),
		provider: provider,
	}

	s.initTemplates()
	s.routes()
	return s
}

// Router возвращает настроенный HTTP-маршрутизатор (реализует http.Handler).
func (s *Service) Router() *mux.Router {
	return s.router
}

// routes регистрирует эндпоинты в маршрутизаторе gorilla/mux.
func (s *Service) routes() {
	s.router.HandleFunc("/docs", s.handleDocs).Methods(http.MethodGet)
	s.router.HandleFunc("/index", s.handleIndex).Methods(http.MethodGet)
}

// handleDocs обрабатывает запросы на получение списка проиндексированных документов.
// Доступен по адресу GET /docs. Поддерживает HTML (по умолчанию в браузере) и JSON.
func (s *Service) handleDocs(w http.ResponseWriter, r *http.Request) {
	if s.provider == nil {
		http.Error(w, "провайдер данных не инициализирован", http.StatusInternalServerError)
		return
	}

	docs := s.provider.Documents()

	// Сортируем документы по возрастанию ID для стабильного и предсказуемого отображения
	sortedDocs := make([]Document, len(docs))
	copy(sortedDocs, docs)
	sort.Slice(sortedDocs, func(i, j int) bool {
		return sortedDocs[i].ID < sortedDocs[j].ID
	})

	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(sortedDocs)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := struct {
		Title string
		Count int
		Docs  []Document
	}{
		Title: "Хранилище документов GoSearch",
		Count: len(sortedDocs),
		Docs:  sortedDocs,
	}

	if err := s.tmplDocs.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("ошибка рендеринга шаблона: %v", err), http.StatusInternalServerError)
	}
}

// handleIndex обрабатывает запросы на получение содержимого обратного поискового индекса.
// Доступен по адресу GET /index. Поддерживает HTML (по умолчанию) и JSON.
func (s *Service) handleIndex(w http.ResponseWriter, r *http.Request) {
	if s.provider == nil {
		http.Error(w, "провайдер данных не инициализирован", http.StatusInternalServerError)
		return
	}

	idxMap := s.provider.Index()

	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(idxMap)
		return
	}

	// Преобразуем словарь в отсортированный срез для красивой таблицы
	words := make([]string, 0, len(idxMap))
	for w := range idxMap {
		words = append(words, w)
	}
	sort.Strings(words)

	entries := make([]IndexEntry, 0, len(words))
	for _, word := range words {
		ids := idxMap[word]
		entries = append(entries, IndexEntry{
			Word:   word,
			DocIDs: ids,
			Count:  len(ids),
		})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := struct {
		Title   string
		Count   int
		Entries []IndexEntry
	}{
		Title:   "Инвертированный индекс GoSearch",
		Count:   len(entries),
		Entries: entries,
	}

	if err := s.tmplIdx.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("ошибка рендеринга шаблона: %v", err), http.StatusInternalServerError)
	}
}

// wantsJSON определяет, запросил ли клиент формат JSON через Accept-заголовок или параметр ?format=json.
func wantsJSON(r *http.Request) bool {
	if r.URL.Query().Get("format") == "json" {
		return true
	}
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "application/json")
}

// initTemplates компилирует безопасные HTML-шаблоны с помощью html/template (защита от XSS).
func (s *Service) initTemplates() {
	const docsHTML = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }}</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; margin: 30px; background: #f8fafc; color: #1e293b; }
        h1 { color: #0284c7; }
        .nav { margin-bottom: 20px; }
        .nav a { margin-right: 15px; color: #2563eb; text-decoration: none; font-weight: 600; }
        .nav a:hover { text-decoration: underline; }
        table { border-collapse: collapse; width: 100%; background: #ffffff; box-shadow: 0 1px 3px rgba(0,0,0,0.1); border-radius: 8px; overflow: hidden; }
        th, td { text-align: left; padding: 12px 16px; border-bottom: 1px solid #e2e8f0; }
        th { background: #f1f5f9; font-weight: 600; color: #475569; }
        tr:hover { background: #f8fafc; }
        .badge { background: #e0f2fe; color: #0369a1; padding: 2px 8px; border-radius: 9999px; font-size: 0.85em; font-weight: 500; }
        .url a { color: #2563eb; text-decoration: none; word-break: break-all; }
        .url a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <div class="nav">
        <a href="/docs">📁 Документы</a>
        <a href="/index">🔍 Инвертированный индекс</a>
    </div>
    <h1>{{ .Title }}</h1>
    <p>Всего проиндексировано документов: <span class="badge">{{ .Count }}</span></p>
    <table>
        <thead>
            <tr>
                <th style="width: 80px;">ID</th>
                <th>Заголовок</th>
                <th>URL адрес</th>
            </tr>
        </thead>
        <tbody>
            {{ range .Docs }}
            <tr>
                <td><span class="badge">#{{ .ID }}</span></td>
                <td><strong>{{ .Title }}</strong></td>
                <td class="url"><a href="{{ .URL }}" target="_blank" rel="noopener">{{ .URL }}</a></td>
            </tr>
            {{ else }}
            <tr>
                <td colspan="3" style="text-align: center; color: #64748b;">Хранилище документов пусто</td>
            </tr>
            {{ end }}
        </tbody>
    </table>
</body>
</html>`

	const indexHTML = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }}</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; margin: 30px; background: #f8fafc; color: #1e293b; }
        h1 { color: #0284c7; }
        .nav { margin-bottom: 20px; }
        .nav a { margin-right: 15px; color: #2563eb; text-decoration: none; font-weight: 600; }
        .nav a:hover { text-decoration: underline; }
        table { border-collapse: collapse; width: 100%; background: #ffffff; box-shadow: 0 1px 3px rgba(0,0,0,0.1); border-radius: 8px; overflow: hidden; }
        th, td { text-align: left; padding: 12px 16px; border-bottom: 1px solid #e2e8f0; }
        th { background: #f1f5f9; font-weight: 600; color: #475569; }
        tr:hover { background: #f8fafc; }
        .badge { background: #e0f2fe; color: #0369a1; padding: 2px 8px; border-radius: 9999px; font-size: 0.85em; font-weight: 500; }
        .doc-tag { display: inline-block; background: #f1f5f9; color: #334155; padding: 2px 6px; border-radius: 4px; margin-right: 4px; font-family: monospace; }
    </style>
</head>
<body>
    <div class="nav">
        <a href="/docs">📁 Документы</a>
        <a href="/index">🔍 Инвертированный индекс</a>
    </div>
    <h1>{{ .Title }}</h1>
    <p>Уникальных терминов в индексе: <span class="badge">{{ .Count }}</span></p>
    <table>
        <thead>
            <tr>
                <th style="width: 250px;">Ключевое слово (Токен)</th>
                <th style="width: 120px;">Вхождений</th>
                <th>Идентификаторы документов (Document IDs)</th>
            </tr>
        </thead>
        <tbody>
            {{ range .Entries }}
            <tr>
                <td><strong>{{ .Word }}</strong></td>
                <td><span class="badge">{{ .Count }}</span></td>
                <td>
                    {{ range .DocIDs }}
                    <span class="doc-tag">#{{ . }}</span>
                    {{ end }}
                </td>
            </tr>
            {{ else }}
            <tr>
                <td colspan="3" style="text-align: center; color: #64748b;">Поисковый индекс пуст</td>
            </tr>
            {{ end }}
        </tbody>
    </table>
</body>
</html>`

	s.tmplDocs = template.Must(template.New("docs").Parse(docsHTML))
	s.tmplIdx = template.Must(template.New("index").Parse(indexHTML))
}
