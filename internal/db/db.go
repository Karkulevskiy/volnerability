package db

import (
	"database/sql"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	// const op = "storage.slqite.New"

	// db, err := sql.Open("sqlite3", storagePath)
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	// stmt, err := db.Prepare(``) //TODO добавить инициализацию таблиц
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	// if _, err = stmt.Exec(); err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	// return &Storage{db: db}, nil
	return &Storage{}, nil
}

func (s *Storage) Login() error {
	// TODO implement me
	return nil
}

func (s *Storage) Register() error {
	// TODO implement me
	return nil
}
