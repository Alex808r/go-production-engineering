// Package slogdemo демонстрирует возможности структурированного логирования log/slog (Go 1.21+).
package slogdemo

import (
	"bytes"
	"context"
	"io"
	"log/slog"
)

// Password реализует интерфейс slog.LogValuer для автоматического маскирования
// конфиденциальных данных в логах.
type Password string

// LogValue маскирует пароль звездочками при выводе в лог.
func (p Password) LogValue() slog.Value {
	return slog.StringValue("******")
}

// User представляет пользователя системы.
type User struct {
	ID       int
	Email    string
	Password Password
}

// LogValue группирует атрибуты пользователя и маскирует пароль.
func (u User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("id", u.ID),
		slog.String("email", u.Email),
		slog.Any("password", u.Password),
	)
}

// NewJSONLogger создает логгер с выводом в формате JSON в переданный io.Writer.
func NewJSONLogger(w io.Writer, level slog.Level) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: false,
	}
	handler := slog.NewJSONHandler(w, opts)
	return slog.New(handler)
}

// NewTextLogger создает логгер с человекочитаемым текстовым форматом.
func NewTextLogger(w io.Writer, level slog.Level) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewTextHandler(w, opts)
	return slog.New(handler)
}

// LogUserLogin логирует событие авторизации пользователя с контекстом и атрибутами.
func LogUserLogin(ctx context.Context, logger *slog.Logger, u User, ip string) {
	logger.InfoContext(ctx, "user authenticated",
		slog.Any("user", u),
		slog.String("remote_ip", ip),
		slog.Group("session",
			slog.String("device", "desktop"),
			slog.Bool("remember_me", true),
		),
	)
}

// LogToBuffer выполняет логирование в буфер и возвращает накопленную строку (для тестов).
func LogToBuffer(u User, ip string) string {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, slog.LevelInfo)
	LogUserLogin(context.Background(), logger, u, ip)
	return buf.String()
}
