package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	models "volnerability-game/internal/domain"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mattn/go-sqlite3"
)

var queries = []string{
	`CREATE TABLE IF NOT EXISTS users
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    pass_hash BLOB NOT NULL
);`,
	`CREATE INDEX IF NOT EXISTS idx_email ON users (email);`,
	`CREATE TABLE IF NOT EXISTS apps
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,	
    name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);`,
	`CREATE TABLE IF NOT EXISTS levels 
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT
);`,
	`CREATE TABLE IF NOT EXISTS user_levels 
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    level_id INTEGER REFERENCES levels(id) ON DELETE CASCADE,
    is_completed BOOLEAN DEFAULT FALSE, --пройден ли уровень
    last_input TEXT, --последний ввод пользователя (например: последний веденный код)
    attempt_response TEXT, --последний ответ сервера на попытку пройти уровень 
    attempts INTEGER DEFAULT 0 --количество попыток пройти уровень (будет увеличиваться при неудаче)
);`,
}

type Storage struct {
	db *sql.DB
}

func New(dbPath string) error {
	const op = "storage.slqite.New"
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		fmt.Println("db file already exists")
		return nil
	}

	dbFile, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	dbFile.Close()

	fmt.Println("db file was created")
	return nil
}

func Init(storagePath string) (*Storage, error) {
	const op = "storage.db.init"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for _, q := range queries {
		fmt.Println(q)
		if _, err := db.Exec(q); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	fmt.Println("tables were created")
	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error) {
	const op = "storage.sqlite.SaveUser"
	query := "INSERT INTO users(email, pass_hash) VALUES(?, ?)"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return -1, fmt.Errorf("%s: %s", op, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) && sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, fmt.Errorf("%s: %s", op, ErrUserExists)
		}
	}

	uid, err = res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return uid, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.sqlite.User"
	query := "SELECT * FROM users WHERE email = ?"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %s", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)
	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %s", op, ErrUserNotFound)
		}
	}

	return user, nil
}
