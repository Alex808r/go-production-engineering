package dip_test

import (
	"context"
	"errors"
	"testing"

	"go-core-4/lessons/17-system-design/homework-17/solid/dip"
)

// Тестовые реализации (моки), подтверждающие независимость бизнес-логики от деталей реализации
type mockPayment struct {
	shouldFail bool
	charged    int64
}

func (m *mockPayment) Charge(_ context.Context, _ string, amount int64) error {
	if m.shouldFail {
		return errors.New("insufficient funds")
	}
	m.charged += amount
	return nil
}

type mockRepo struct {
	saved []dip.Ticket
}

func (m *mockRepo) Save(_ context.Context, t dip.Ticket) error {
	m.saved = append(m.saved, t)
	return nil
}

func (m *mockRepo) GetByID(_ context.Context, _ int64) (dip.Ticket, error) {
	if len(m.saved) == 0 {
		return dip.Ticket{}, errors.New("not found")
	}
	return m.saved[0], nil
}

func TestTicketBookingService_DIP(t *testing.T) {
	ctx := context.Background()

	// 1. Успешный сценарий
	p := &mockPayment{}
	r := &mockRepo{}
	svc := dip.NewTicketBookingService(p, r)

	ticket := dip.Ticket{ID: 1, MovieID: 10, UserEmail: "buyer@example.com", Price: 500}
	if err := svc.BookTicket(ctx, ticket); err != nil {
		t.Fatalf("expected successful booking, got: %v", err)
	}

	if p.charged != 500 {
		t.Errorf("expected charged 500, got %d", p.charged)
	}
	if len(r.saved) != 1 || !r.saved[0].Paid {
		t.Errorf("expected ticket to be saved as paid")
	}

	// 2. Сценарий ошибки оплаты
	pFail := &mockPayment{shouldFail: true}
	rFail := &mockRepo{}
	svcFail := dip.NewTicketBookingService(pFail, rFail)

	if err := svcFail.BookTicket(ctx, ticket); err == nil {
		t.Error("expected booking to fail when payment fails")
	}
	if len(rFail.saved) != 0 {
		t.Error("ticket should NOT be saved if payment failed")
	}
}
