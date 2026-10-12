// Package lsp демонстрирует принцип подстановки Барбары Лисков (Liskov Substitution Principle) в Go.
// В Go нет наследования типов, поэтому LSP выражается через реализацию интерфейсов:
// если функция принимает интерфейс, любая его конкретная реализация должна сохранять
// контракт и семантические инварианты, не вызывая неожиданного поведения у вызывающего кода.
package lsp

import (
	"context"
	"fmt"
	"strings"
)

// Notifier определяет контракт отправки уведомлений.
// Инвариант: метод Send не должен паниковать и должен возвращать понятную ошибку
// при невозможности отправки, сохраняя контекст отмены.
type Notifier interface {
	Send(ctx context.Context, recipient string, message string) error
}

// EmailNotifier отправляет уведомления по электронной почте.
type EmailNotifier struct {
	SentMessages []string
}

func (e *EmailNotifier) Send(_ context.Context, recipient string, message string) error {
	if !strings.Contains(recipient, "@") {
		return fmt.Errorf("invalid email address: %s", recipient)
	}
	e.SentMessages = append(e.SentMessages, fmt.Sprintf("EMAIL to %s: %s", recipient, message))
	return nil
}

// SMSNotifier отправляет SMS-уведомления.
type SMSNotifier struct {
	SentMessages []string
}

func (s *SMSNotifier) Send(_ context.Context, recipient string, message string) error {
	if len(recipient) < 5 {
		return fmt.Errorf("invalid phone number: %s", recipient)
	}
	s.SentMessages = append(s.SentMessages, fmt.Sprintf("SMS to %s: %s", recipient, message))
	return nil
}

// NotificationService — сервис оповещения пользователей.
// Следует LSP: работает с любым Notifier одинаково корректно и взаимозаменяемо.
type NotificationService struct {
	notifier Notifier
}

func NewNotificationService(n Notifier) *NotificationService {
	return &NotificationService{notifier: n}
}

func (s *NotificationService) Notify(ctx context.Context, recipient, msg string) error {
	return s.notifier.Send(ctx, recipient, msg)
}
