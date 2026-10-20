// Package spider реализует сканер содержимого веб-сайтов.
// Пакет позволяет получить список ссылок и заголовков страниц внутри веб-сайта по его URL.
package spider

import (
	"net/http"
	"strings"

	"go-core-4/lessons/03-algorithms/homework-03/pkg/crawler"

	"golang.org/x/net/html"
)

// Service - служба поискового робота.
type Service struct{}

// New - конструктор службы поискового робота.
func New() *Service {
	return &Service{}
}

// Scan осуществляет рекурсивный обход ссылок сайта, указанного в URL,
// с учётом глубины перехода по ссылкам, переданной в depth.
func (s *Service) Scan(url string, depth int) (data []crawler.Document, err error) {
	pages := make(map[string]string)

	parse(url, url, depth, pages)

	for url, title := range pages {
		item := crawler.Document{
			URL:   url,
			Title: title,
		}
		data = append(data, item)
	}

	return data, nil
}

// parse рекурсивно обходит ссылки на страницах, начиная с url,
// замешивая найденные результаты в карту pages.
func parse(url, baseurl string, depth int, pages map[string]string) error {
	if depth == 0 {
		return nil
	}

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	page, err := html.Parse(response.Body)
	if err != nil {
		return err
	}

	title := pageTitle(page)
	if title != "" {
		pages[url] = title
	}

	links := pageLinks(nil, page)
	for _, link := range links {
		link = strings.TrimSuffix(link, "/")
		// относительная ссылка
		if strings.HasPrefix(link, "/") && len(link) > 1 {
			link = baseurl + link
		}
		// ссылка на тот же сервер
		if strings.HasPrefix(link, baseurl) && pages[link] == "" {
			parse(link, baseurl, depth-1, pages)
		}
	}

	return nil
}

// pageTitle осуществляет рекурсивный обход HTML-дерева и возвращает значение первого тега <title>.
func pageTitle(n *html.Node) string {
	var title string
	if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
		return n.FirstChild.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title = pageTitle(c)
		if title != "" {
			break
		}
	}
	return title
}

// pageLinks осуществляет рекурсивный обход HTML-дерева и возвращает срез ссылок в тегах <a>.
func pageLinks(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				if !sliceContains(links, a.Val) {
					links = append(links, a.Val)
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = pageLinks(links, c)
	}
	return links
}

func sliceContains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
