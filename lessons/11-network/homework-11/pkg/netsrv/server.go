package netsrv

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"go-core-4/lessons/07-testing/homework-07/pkg/crawler"
)

// Searcher определяет контракт поискового движка.
type Searcher interface {
	Search(query string) []crawler.Document
}

// Server представляет сетевую службу поисковика GoSearch.
type Server struct {
	addr     string
	searcher Searcher
	listener net.Listener
	timeout  time.Duration
	mu       sync.Mutex
	closed   bool
}

// New создает экземпляр сетевой службы.
// Если addr пустой, по умолчанию используется ":8000" (прослушивание всех интерфейсов 0.0.0.0:8000).
func New(addr string, searcher Searcher) *Server {
	if addr == "" {
		addr = ":8000"
	}
	return &Server{
		addr:     addr,
		searcher: searcher,
		timeout:  5 * time.Minute, // 5 минут таймаут на бездействие
	}
}

// SetTimeout устанавливает таймаут неактивности клиентского соединения.
func (s *Server) SetTimeout(d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.timeout = d
}

// Start регистрирует сетевую службу на указанном адресе и запускает цикл обработки подключений.
func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("netsrv: не удалось запустить сервер на %s: %w", s.addr, err)
	}

	s.mu.Lock()
	s.listener = l
	s.mu.Unlock()

	return s.Serve(l)
}

// Serve принимает подключения из существующего listener.
func (s *Server) Serve(l net.Listener) error {
	s.mu.Lock()
	s.listener = l
	s.mu.Unlock()

	for {
		conn, err := l.Accept()
		if err != nil {
			s.mu.Lock()
			closed := s.closed
			s.mu.Unlock()
			if closed {
				return nil
			}
			continue
		}

		// Каждое клиентское подключение обрабатывается в отдельной горутине
		go s.HandleConn(conn)
	}
}

// Stop корректно завершает работу сетевой службы.
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}
	s.closed = true

	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// Addr возвращает сетевой адрес запущенной службы.
func (s *Server) Addr() net.Addr {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listener != nil {
		return s.listener.Addr()
	}
	return nil
}

// HandleConn обрабатывает сессию интерактивного сетевого подключения.
func (s *Server) HandleConn(conn net.Conn) {
	defer conn.Close()

	// Приветственное сообщение
	fmt.Fprintf(conn, "=== Сетевая служба GoSearch (netsrv) ===\n")
	fmt.Fprintf(conn, "Введите поисковый запрос (или 'exit' для завершения):\n\n> ")

	scanner := bufio.NewScanner(conn)
	for {
		s.mu.Lock()
		timeout := s.timeout
		s.mu.Unlock()

		if timeout > 0 {
			_ = conn.SetDeadline(time.Now().Add(timeout))
		}

		if !scanner.Scan() {
			break
		}

		query := strings.TrimSpace(scanner.Text())
		if query == "" {
			fmt.Fprintf(conn, "> ")
			continue
		}

		// Команда выхода
		if strings.EqualFold(query, "exit") || strings.EqualFold(query, "quit") {
			fmt.Fprintf(conn, "Сеанс завершен. До свидания!\n")
			return
		}

		// Поиск документов через внедренный движок
		docs := s.searcher.Search(query)
		if len(docs) == 0 {
			fmt.Fprintf(conn, "По запросу %q ничего не найдено.\n\n> ", query)
			continue
		}

		fmt.Fprintf(conn, "Найдено документов: %d\n", len(docs))
		for i, doc := range docs {
			fmt.Fprintf(conn, "  %d. [%d] %s — %s\n", i+1, doc.ID, doc.Title, doc.URL)
		}
		fmt.Fprintf(conn, "\n> ")
	}
}
