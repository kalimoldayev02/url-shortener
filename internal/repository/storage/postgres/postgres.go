package postgres

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const operation = "storage.postgres.New"

	fmt.Println(storagePath)
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS url (
			id SERIAL PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);`,
	}
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			// обработка ошибки, если что-то пошло не так
			return nil, fmt.Errorf("%s: %w", operation, err)

		}
	}

	return &Storage{db: db}, nil
}
