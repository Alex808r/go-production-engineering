// Package fuzzing демонстрирует возможности встроенного фаззинг-тестирования (testing.F).
package fuzzing

import (
	"errors"
	"strings"
	"unicode/utf8"
)

var (
	ErrEmptyInput    = errors.New("input is empty")
	ErrInvalidFormat = errors.New("invalid format: expected key=value")
	ErrInvalidUTF8   = errors.New("input is not valid UTF-8")
)

// QueryParams хранит разобранные параметры строки запроса вида k1=v1&k2=v2.
type QueryParams map[string]string

// ParseQuery безопасно парсит строку запроса.
// Функция гарантирует отсутствие паник даже при передаче абсолютно некорректных байт.
func ParseQuery(input string) (QueryParams, error) {
	if !utf8.ValidString(input) {
		return nil, ErrInvalidUTF8
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return nil, ErrEmptyInput
	}

	params := make(QueryParams)
	pairs := strings.Split(input, "&")

	for _, pair := range pairs {
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, ErrInvalidFormat
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if key == "" || val == "" {
			return nil, ErrInvalidFormat
		}
		params[key] = val
	}

	if len(params) == 0 {
		return nil, ErrEmptyInput
	}

	return params, nil
}
