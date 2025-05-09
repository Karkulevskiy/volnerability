package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrUserExists        = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserLevelNotFound = fmt.Errorf("user level not found")
)

func IsSyntaxError(err error) bool {
	var sqliteErr sqlite3.Error
	if ok := errors.As(err, &sqliteErr); ok {
		return strings.Contains(sqliteErr.Error(), "syntax error")
	}
	return false
}
