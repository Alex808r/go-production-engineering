// Package traps демонстрирует нетривиальные подводные камни Go, выравнивание памяти в структурах,
// паттерн Functional Options и особенности работы мап/таймеров (из Uber Go Style Guide и 100 Go Mistakes).
package traps

import (
	"context"
	"slices"
	"time"
	"unsafe"
)

// --- 1. Выравнивание структур (Data Alignment & Struct Padding) ---

// BadStruct имеет неоптимальный порядок полей.
// На 64-битной архитектуре bool занимает 1 байт, но из-за 8-байтового выравнивания int64
// компилятор вставляет 7 байт padding между a и b, и еще 7 байт после c.
// Итог: 1 + (7 pad) + 8 + 1 + (7 pad) = 24 байта!
type BadStruct struct {
	A bool  // 1 байт (+7 pad)
	B int64 // 8 байт
	C bool  // 1 байт (+7 pad)
}

// GoodStruct имеет поля, отсортированные по убыванию размера.
// Итог: 8 + 1 + 1 + (6 pad) = 16 байт (экономия 33% памяти!).
type GoodStruct struct {
	B int64 // 8 байт
	A bool  // 1 байт
	C bool  // 1 байт (+6 pad)
}

// GetStructSizes возвращает размеры структур в байтах.
func GetStructSizes() (badSize, goodSize uintptr) {
	return unsafe.Sizeof(BadStruct{}), unsafe.Sizeof(GoodStruct{})
}

// --- 2. Утечка памяти при срезах (Slice Reslicing Leak) ---

// LeakSubslice демонстрирует удержание большого массива через срез.
func LeakSubslice(huge []int) []int {
	return huge[:2] // Срез удерживает весь базовый массив в памяти от сборщика мусора!
}

// SafeSubslice копирует нужные элементы в новый срез, позволяя GC освободить исходный массив.
func SafeSubslice(huge []int) []int {
	return slices.Clone(huge[:2])
}

// --- 3. Паттерн Functional Options ---

// ServerConfig хранит конфигурацию сервера.
type ServerConfig struct {
	Host    string
	Port    int
	Timeout time.Duration
	MaxConn int
}

// Option определяет функцию модификации конфигурации.
type Option func(*ServerConfig)

// WithPort задает порт сервера.
func WithPort(port int) Option {
	return func(c *ServerConfig) {
		c.Port = port
	}
}

// WithTimeout задает таймаут запросов.
func WithTimeout(timeout time.Duration) Option {
	return func(c *ServerConfig) {
		c.Timeout = timeout
	}
}

// WithMaxConn задает максимальное число соединений.
func WithMaxConn(maxConn int) Option {
	return func(c *ServerConfig) {
		c.MaxConn = maxConn
	}
}

// NewServerConfig создает конфигурацию с дефолтными значениями и применяет опции.
func NewServerConfig(host string, opts ...Option) ServerConfig {
	cfg := ServerConfig{
		Host:    host,
		Port:    8080,
		Timeout: 30 * time.Second,
		MaxConn: 1000,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

// --- 4. Безопасный таймер в циклах вместо time.After (Timer Leak Prevention) ---

// ProcessWithTimeout демонстрирует безопасное использование переиспользуемого time.Timer
// вместо time.After, который создает новый канал на каждой итерации и удерживает память до истечения TTL.
func ProcessWithTimeout(ctx context.Context, ch <-chan int, timeout time.Duration) (int, bool) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case val := <-ch:
		return val, true
	case <-timer.C:
		return 0, false
	case <-ctx.Done():
		return 0, false
	}
}

// --- 5. Проверка соответствия интерфейсу во время компиляции (Compile-time interface verification) ---

// Runner определяет интерфейс сервиса.
type Runner interface {
	Run() string
}

// Service реализует интерфейс Runner.
type Service struct{}

func (s *Service) Run() string {
	return "service running"
}

// Compile-time assertion: если сигнатура Run изменится, компиляция упадет прямо здесь.
// Это идиома из Uber Go Style Guide. В рантайме не создает аллокаций.
var _ Runner = (*Service)(nil)

// --- 6. Очистка мапы и освобождение памяти бакетов ---

// RecreateMap демонстрирует освобождение памяти после очистки большой мапы.
// delete() очищает ключи, но число бакетов в структуре hmap не уменьшается!
// Для возврата памяти сборщику мусора мапу нужно пересоздать (или занулить ссылку).
func RecreateMap[K comparable, V any](m map[K]V) map[K]V {
	// Очищаем ссылку на старую мапу, чтобы GC мог забрать бакеты
	return make(map[K]V)
}

