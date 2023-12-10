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

	stmt, err := db.Prepare(
		`CREATE TABLE IF NOT EXIST url(
		id SERIAL PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL
	); 
	CREATE INDEX IF NOT EXIST idx_alias ON url(alias)`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	if _, err := stmt.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &Storage{db: db}, nil
}
