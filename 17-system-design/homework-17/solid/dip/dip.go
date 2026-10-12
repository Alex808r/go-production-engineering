// Package dip демонстрирует принцип инверсии зависимостей (Dependency Inversion Principle) в Go.
// Модули высокого уровня (бизнес-логика / use cases) не должны зависеть от модулей низкого уровня (БД, платёжные шлюзы, сеть).
// Оба уровня должны зависеть от абстракций (интерфейсов).
// Детали (конкретные адаптеры) зависят от абстракций, заданных ядром.
package dip

import (
	"context"
	"errors"
	"fmt"
)

// Ticket — доменная сущность билета.
type Ticket struct {
	ID        int64
	MovieID   int64
	UserEmail string
	Price     int64
	Paid      bool
}

// PaymentGateway — абстракция платежного шлюза (принадлежит бизнес-логике).
type PaymentGateway interface {
	Charge(ctx context.Context, email string, amount int64) error
}

// TicketRepository — абстракция хранилища билетов (принадлежит бизнес-логике).
type TicketRepository interface {
	Save(ctx context.Context, ticket Ticket) error
	GetByID(ctx context.Context, id int64) (Ticket, error)
}

// TicketBookingService — модуль высокого уровня (бизнес-логика).
// Не зависит ни от Stripe, ни от PostgreSQL — зависит только от интерфейсов!
type TicketBookingService struct {
	payment PaymentGateway
	repo    TicketRepository
}

func NewTicketBookingService(p PaymentGateway, r TicketRepository) *TicketBookingService {
	return &TicketBookingService{payment: p, repo: r}
}

// BookTicket выполняет бизнес-сценарий покупки билета:
// 1. Списание средств через платежный интерфейс
// 2. Сохранение билета в репозиторий со статусом Paid = true
func (s *TicketBookingService) BookTicket(ctx context.Context, ticket Ticket) error {
	if ticket.Price <= 0 {
		return errors.New("ticket price must be positive")
	}

	// 1. Оплата
	if err := s.payment.Charge(ctx, ticket.UserEmail, ticket.Price); err != nil {
		return fmt.Errorf("payment failed: %w", err)
	}

	// 2. Сохранение
	ticket.Paid = true
	if err := s.repo.Save(ctx, ticket); err != nil {
		return fmt.Errorf("saving ticket failed: %w", err)
	}

	return nil
}
