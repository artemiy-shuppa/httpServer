package postgres

import (
	"database/sql"
	"myHttpServer/internal/domain"

	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

func Open(dataSourceName string) (*Database, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{db}, nil
}
func (db *Database) Close() {
	db.DB.Close()
}

func (db *Database) MessageStorage() domain.MessageService {
	return &MessageStorage{db: db}
}
