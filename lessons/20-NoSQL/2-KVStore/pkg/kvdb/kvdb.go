// Package kvdb реализует потокобезопасное хранилище «ключ-значение» в оперативной памяти.
package kvdb

import "sync"

// DB представляет потокобезопасную базу данных «ключ-значение» в оперативной памяти.
type DB struct {
	mu   sync.RWMutex
	data map[string]string
}

// New инициализирует и возвращает новое KV-хранилище.
func New() *DB {
	return &DB{
		data: make(map[string]string),
	}
}

// GET возвращает значение по ключу (или пустую строку, если ключ отсутствует).
func (db *DB) GET(key string) string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.data[key]
}

// Get возвращает значение и булев флаг, указывающий на наличие ключа.
func (db *DB) Get(key string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	val, ok := db.data[key]
	return val, ok
}

// SET сохраняет значение по ключу под исключительной блокировкой записи.
func (db *DB) SET(key string, val string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = val
}

// DELETE удаляет ключ из хранилища.
func (db *DB) DELETE(key string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.data, key)
}

// Len возвращает текущее количество записей в хранилище.
func (db *DB) Len() int {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return len(db.data)
}
