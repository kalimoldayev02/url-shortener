package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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

func (s *Storage) SaveUrl(urlToSave, alias string) (int64, error) {
	const operation = "storage.postgres.SaveUrl"
	query := "INSERT INTO url(url, alias) VALUES ($1, $2) RETURNING id"

	stmp, err := s.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	var urlId int64

	err = stmp.QueryRow(urlToSave, alias).Scan(&urlId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	defer stmp.Close()

	return urlId, nil
}

func (s *Storage) GetUrlByAlias(alias string) (string, error) {
	const operation = "storage.postgres.GetUrlByAlias"
	query := "SELECT url FROM url WHERE alias = $1"

	stmp, err := s.db.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	var resUrl string

	err = stmp.QueryRow(alias).Scan(&resUrl)
	if err != nil {
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	defer stmp.Close()

	return resUrl, nil
}

func (s *Storage) DeleteUrlByAlias(alias string) error {
	const operation = "storage.postgres.DeleteUrlByAlias"
	query := "DELETE FROM url WHERE alias = $1"

	stmp, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	err = stmp.QueryRow(alias).Scan()
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}
