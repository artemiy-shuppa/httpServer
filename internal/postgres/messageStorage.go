package postgres

import (
	"myHttpServer/internal/domain"
)

// MessageStorage implements domain.MessageService using postgress database
type MessageStorage struct {
	db *Database
}

var _ domain.MessageService = (*MessageStorage)(nil)

func (s *MessageStorage) Message(id int) (*domain.Message, error) {
	row := s.db.QueryRow(
		"SELECT * FROM Message WHERE id = $1",
		id,
	)
	if err := row.Err(); err != nil {
		return nil, err
	}
	var m domain.Message
	if err := row.Scan(&m.ID, &m.Body, &m.CreatedAt); err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *MessageStorage) Messages() (*domain.MessageSlice, error) {
	rows, err := s.db.Query("SELECT * FROM Message")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var mm domain.MessageSlice
	for rows.Next() {
		var m domain.Message
		if err := rows.Scan(&m.ID, &m.Body, &m.CreatedAt); err != nil {
			return &mm, err
		}
		mm = append(mm, m)
	}
	if err := rows.Err(); err != nil {
		return &mm, err
	}
	return &mm, nil
}

func (s *MessageStorage) CreateMessage(m *domain.Message) error {
	_, err := s.db.Exec(
		"INSERT INTO Message (body, created_at) VALUES ($1, $2)",
		m.Body,
		m.CreatedAt,
	)

	return err
}

func (s *MessageStorage) UpdateMessage(id int, m *domain.Message) error {
	_, err := s.db.Exec(
		"UPDATE Message SET body = $1 WHERE id = $2",
		m.Body,
		m.ID,
	)
	return err
}

func (s *MessageStorage) DeleteMessage(id int) error {
	_, err := s.db.Exec("DELETE FROM Message WHERE id = $1", id)
	return err
}
