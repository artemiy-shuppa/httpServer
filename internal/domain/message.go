package domain

import (
	"time"
)

type Message struct {
	ID        int       `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageSlice []Message

type MessageService interface {
	Message(id int) (*Message, error)
	Messages() (*MessageSlice, error)
	CreateMessage(m *Message) error
	UpdateMessage(id int, m *Message) error
	DeleteMessage(id int) error
}
