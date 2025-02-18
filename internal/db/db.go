package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	const op = "storage.slqite.New"
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		fmt.Println("db file already exists")
		return &Storage{}, nil
	}

	dbFile, err := os.Create(dbPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	dbFile.Close()

	fmt.Println("db file was created")
	return &Storage{}, nil
}

func (s *Storage) Init(migrationPath, storagePath string) error {
	const op = "storage.db.init"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer db.Close()

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://internal/db", "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Printf("db already created: %s\n", err)
			return nil
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	fmt.Println("db was initialized")
	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Login() error {
	// TODO implement me
	return nil
}

func (s *Storage) Register() error {
	// TODO implement me
	return nil
}
