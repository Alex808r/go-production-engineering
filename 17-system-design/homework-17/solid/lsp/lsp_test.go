package lsp_test

import (
	"context"
	"testing"

	"go-core-4/17-system-design/homework-17/solid/lsp"
)

// Тест проверяет взаимозаменяемость реализаций Notifier (LSP)
func TestNotificationService_LSP(t *testing.T) {
	ctx := context.Background()

	email := &lsp.EmailNotifier{}
	sms := &lsp.SMSNotifier{}

	notifiers := []struct {
		name      string
		notifier  lsp.Notifier
		recipient string
		msg       string
	}{
		{"Email", email, "user@example.com", "Hello via Email"},
		{"SMS", sms, "+79991234567", "Hello via SMS"},
	}

	for _, tc := range notifiers {
		t.Run(tc.name, func(t *testing.T) {
			svc := lsp.NewNotificationService(tc.notifier)
			if err := svc.Notify(ctx, tc.recipient, tc.msg); err != nil {
				t.Fatalf("unexpected error for %s: %v", tc.name, err)
			}
		})
	}

	if len(email.SentMessages) != 1 {
		t.Errorf("expected 1 email, got %d", len(email.SentMessages))
	}
	if len(sms.SentMessages) != 1 {
		t.Errorf("expected 1 sms, got %d", len(sms.SentMessages))
	}
}
