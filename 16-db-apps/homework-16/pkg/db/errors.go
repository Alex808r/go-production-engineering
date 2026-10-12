// Package db определяет ошибки, общие для всех реализаций хранилища.
package db

import "errors"

// Sentinel errors — доменные ошибки хранилища.
// Слой бизнес-логики и API должен проверять ошибки через errors.Is(err, db.ErrNotFound),
// а не привязываться к конкретной реализации (pgx.ErrNoRows, sql.ErrNoRows).
var (
	// ErrNotFound возвращается при запросе несуществующей записи.
	ErrNotFound = errors.New("запись не найдена")
	// ErrDuplicateKey возвращается при попытке создать дубликат уникального ключа.
	ErrDuplicateKey = errors.New("запись с таким ключом уже существует")
)
