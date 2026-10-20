// Package spider реализует сканер содержимого веб-сайтов.
package spider

import (
	"net/http"
	"strings"

	"go-core-4/lessons/06-oop/homework-06/pkg/crawler"

	"golang.org/x/net/html"
)

// Service - служба поискового робота spider.
type Service struct{}

// New - конструктор службы поискового робота.
func New() *Service {
	return &Service{}
}

// Scan осуществляет рекурсивный обход ссылок сайта, указанного в URL.
func (s *Service) Scan(url string, depth int) (data []crawler.Document, err error) {
	pages := make(map[string]string)

	parse(url, url, depth, pages)

	for pageURL, title := range pages {
		item := crawler.Document{
			URL:   pageURL,
			Title: title,
		}
		data = append(data, item)
	}

	return data, nil
}

// BatchScan заглушка для интерфейса.
func (s *Service) BatchScan(urls []string, depth int, workers int) (<-chan crawler.Document, <-chan error) {
	ch := make(chan crawler.Document)
	errCh := make(chan error)
	close(ch)
	close(errCh)
	return ch, errCh
}

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
		if strings.HasPrefix(link, "/") && len(link) > 1 {
			link = baseurl + link
		}
		if strings.HasPrefix(link, baseurl) && pages[link] == "" {
			parse(link, baseurl, depth-1, pages)
		}
	}

	return nil
}

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
