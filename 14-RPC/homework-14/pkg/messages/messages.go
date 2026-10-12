package messages

import "time"

// Message представляет сообщение в RPC-службе.
// Содержит обязательные поля по условию домашнего задания №14:
// идентификатор, время отправки и текст.
type Message struct {
	ID   int       `json:"id"`
	Time time.Time `json:"time"`
	Text string    `json:"text"`
}
