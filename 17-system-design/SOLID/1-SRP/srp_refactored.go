// Package srp демонстрирует рефакторинг с соблюдением принципа единственной ответственности (SRP).
// Вместо одного «всезнающего» монолитного клиента каждый провайдер облачного хранилища
// инкапсулирует только свою логику авторизации и передачи данных.
package srp

import (
	"context"
	"fmt"
)

// CloudStorage — абстрактный контракт взаимодействия с облачными хранилищами.
type CloudStorage interface {
	GetFile(ctx context.Context, name string) ([]byte, error)
	UploadFile(ctx context.Context, name string, data []byte) error
}

// YandexDiskStorage отвечает исключительно за взаимодействие с Яндекс.Диском.
type YandexDiskStorage struct {
	Token string
}

// NewYandexDiskStorage создает экземпляр клиента Яндекс.Диска.
func NewYandexDiskStorage(token string) *YandexDiskStorage {
	return &YandexDiskStorage{Token: token}
}

// GetFile скачивает файл из Яндекс.Диска.
func (y *YandexDiskStorage) GetFile(ctx context.Context, name string) ([]byte, error) {
	if name == "" {
		return nil, fmt.Errorf("empty file name")
	}
	return []byte("yandex-data:" + name), nil
}

// UploadFile загружает файл на Яндекс.Диск.
func (y *YandexDiskStorage) UploadFile(ctx context.Context, name string, data []byte) error {
	if name == "" {
		return fmt.Errorf("empty file name")
	}
	return nil
}

// GoogleDriveStorage отвечает исключительно за взаимодействие с Google Drive.
type GoogleDriveStorage struct {
	APIKey string
}

// NewGoogleDriveStorage создает экземпляр клиента Google Drive.
func NewGoogleDriveStorage(apiKey string) *GoogleDriveStorage {
	return &GoogleDriveStorage{APIKey: apiKey}
}

// GetFile скачивает файл из Google Drive.
func (g *GoogleDriveStorage) GetFile(ctx context.Context, name string) ([]byte, error) {
	if name == "" {
		return nil, fmt.Errorf("empty file name")
	}
	return []byte("google-data:" + name), nil
}

// UploadFile загружает файл на Google Drive.
func (g *GoogleDriveStorage) UploadFile(ctx context.Context, name string, data []byte) error {
	if name == "" {
		return fmt.Errorf("empty file name")
	}
	return nil
}
