package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"volnerability-game/internal/common"
	"volnerability-game/internal/domain"
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
    pass_hash BLOB,
	oauth_id TEXT,
	is_oauth BOOLEAN DEFAULT FALSE,
	total_attempts INTEGER DEFAULT 0,
	pass_levels INTEGER DEFAULT 0
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
    description TEXT,
	expected_input TEXT
);`,
	`CREATE TABLE IF NOT EXISTS hints
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    level_id INTEGER NOT NULL REFERENCES levels(id) ON DELETE CASCADE,
    hint_text TEXT NOT NULL
);
	`,
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

func (s *Storage) IsQueryValid(query string) error {
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return nil
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.slqite.New"

	_, err := os.Stat(storagePath)
	if err == nil {
		return OpenDb(storagePath)
	}

	if !os.IsNotExist(err) {
		return nil, err
	}

	if err := CreateFileDb(storagePath); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return Init(storagePath)
}

func OpenDb(storagePath string) (*Storage, error) {
	fmt.Println("db file already initialized")
	const op = "storage.db.init"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func CreateFileDb(storagePath string) error {
	const op = "storage.slqite.New"
	fmt.Println("creating db file")
	dbFile, err := os.Create(storagePath)
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

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte, isOauth bool) (uid int64, err error) {
	const op = "storage.sqlite.SaveUser"
	query := "INSERT INTO users(email, pass_hash, is_oauth) VALUES(?, ?, ?)"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return -1, fmt.Errorf("%s: %s", op, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash, isOauth)
	if err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) && sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, fmt.Errorf("%s: %s", op, ErrUserExists)
		}
		return -1, fmt.Errorf("%s: %s", op, err)
	}

	uid, err = res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return uid, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	fmt.Println("ALLOOOOO")
	const op = "storage.sqlite.User"
	query := "SELECT * FROM users WHERE email = ?"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %s", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)
	var user models.User
	if err = row.Scan(&user.ID, &user.Email, &user.PassHash, &user.OauthID, &user.IsOauth, &user.TotalAttempts, &user.PassLevels); err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %s", op, ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %s", op, err)
	}
	return user, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user models.User) error {
	needUpdate := make([]string, 0, 3)
	if len(user.PassHash) > 0 {
		needUpdate = append(needUpdate, "pass_hash = "+string(user.PassHash))
	}
	if user.TotalAttempts > 0 {
		needUpdate = append(needUpdate, "total_attempts = "+strconv.Itoa(user.TotalAttempts))
	}
	if user.PassLevels > 0 {
		needUpdate = append(needUpdate, "pass_levels = "+strconv.Itoa(user.PassLevels))
	}

	query := fmt.Sprintf(`
	UPDATE users
	SET %s
	WHERE email = ?
	`, strings.Join(needUpdate, ", "))

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec()
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (s *Storage) UserLevels(ctx context.Context, email string) ([]models.UserLevel, error) {
	const op = "storage.sqlite.UserLevels"
	query := `
	SELECT l.level_id, l.user_id l.is_completed, l.last_input, l.attempt_response, l.attempts
	FROM user_levels l
	LEFT JOIN users u on u.id = l.user_id
	WHERE u.email = ?
	`

	rows, err := s.db.Query(query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user levels: %s %w", op, err)
	}
	defer rows.Close()

	userLevels := []domain.UserLevel{}

	for rows.Next() {
		var (
			levelId, userId, attempts  int
			isCompleted                bool
			lastInput, attemptResponse sql.NullString
		)

		if err := rows.Scan(&levelId, &userId, &isCompleted, &lastInput, &attemptResponse, &attempts); err != nil {
			return nil, fmt.Errorf("scan error: %s %w", op, err)
		}

		userLevel := domain.UserLevel{
			LevelId:         levelId,
			UserId:          userId,
			IsCompleted:     isCompleted,
			LastInput:       lastInput.String,
			AttemptResponse: attemptResponse.String,
			Attempts:        attempts,
		}

		userLevels = append(userLevels, userLevel)
	}

	return userLevels, nil
}

func (s *Storage) UpdateUserAndLevel(ctx context.Context, user models.User, ul models.UserLevel) error {
	const op = "storage.sqlite.UpdateUserAndLevel"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: begin transaction error: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Обновляем пользователя
	_, err = tx.ExecContext(ctx, `
		UPDATE users
		SET total_attempts = ?, pass_levels = ?
		WHERE id = ?
	`, user.TotalAttempts, user.PassLevels, user.ID)
	if err != nil {
		return fmt.Errorf("%s: update user error: %w", op, err)
	}

	// Проверяем, существует ли уже запись user_level
	var exists bool
	err = tx.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM user_levels
			WHERE user_id = ? AND level_id = ?
		)
	`, ul.UserId, ul.LevelId).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: check user_level existence error: %w", op, err)
	}

	if exists {
		// Обновляем существующую запись
		_, err = tx.ExecContext(ctx, `
			UPDATE user_levels
			SET is_completed = ?, last_input = ?, attempt_response = ?, attempts = ?
			WHERE user_id = ? AND level_id = ?
		`, ul.IsCompleted, ul.LastInput, ul.AttemptResponse, ul.Attempts, ul.UserId, ul.LevelId)
		if err != nil {
			return fmt.Errorf("%s: update user_level error: %w", op, err)
		}
	} else {
		// Вставляем новую запись
		_, err = tx.ExecContext(ctx, `
			INSERT INTO user_levels (user_id, level_id, is_completed, last_input, attempt_response, attempts)
			VALUES (?, ?, ?, ?, ?, ?)
		`, ul.UserId, ul.LevelId, ul.IsCompleted, ul.LastInput, ul.AttemptResponse, ul.Attempts)
		if err != nil {
			return fmt.Errorf("%s: insert user_level error: %w", op, err)
		}
	}

	// Подтверждаем транзакцию
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%s: commit error: %w", op, err)
	}

	return nil
}

// UserStartLevel creates new row in user_levels
func (s *Storage) UserStartLevel(ctx context.Context, userId, levelId int) error {
	const op = "storage.sqlite.UserStartLevel"
	query := `INSERT INTO user_levels(user_id, level_id) VALUES(?, ?)`

	res, err := s.db.Exec(query, userId, levelId)
	if err != nil {
		return fmt.Errorf("failed to add new user_level. op: %s, err: %w", op, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows, op: %s, err: %w", op, err)
	}

	if count == 0 {
		return fmt.Errorf("no any rows inserted into user_levels. op: %s", op)
	}

	return nil
}

func (s *Storage) Level(ctx context.Context, id int) (models.Level, error) {
	const op = "storage.sqlite.LevelByID"
	query := `
	SELECT 
		l.id, l.name, l.description, l.expected_input, l.start_input
	FROM levels l
	WHERE l.id = ?
	`

	row := s.db.QueryRowContext(ctx, query, id)

	var (
		levelId                     int
		name, desc, expected, start sql.NullString
	)

	if err := row.Scan(&levelId, &name, &desc, &expected, &start); err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return models.Level{}, fmt.Errorf("%s: level not found: %w", op, err)
		}
		return models.Level{}, fmt.Errorf("%s: scan error: %w", op, err)
	}

	hints, err := s.HintsByLevelId(ctx, id)
	if err != nil {
		fmt.Println(err.Error())
		return models.Level{}, fmt.Errorf("%s: failed to query hints for level id: %d", op, id)
	}

	return models.Level{
		Id:            levelId,
		Name:          name.String,
		Description:   desc.String,
		ExpectedInput: expected.String,
		StartInput:    start.String,
		Hints:         common.Map(func(hint domain.Hint) string { return hint.Text }, hints...),
	}, nil
}

// Если пригодится, не тестировано еще
// func (s *Storage) Levels(ctx context.Context, ids ...int) ([]models.Level, error) {
// 	const op = "storage.sqlite.Levels"
//
// 	if len(ids) == 0 {
// 		return nil, nil
// 	}
//
// 	// Создаём плейсхолдеры (?, ?, ?) и аргументы для запроса
// 	placeholders := make([]string, len(ids))
// 	args := make([]interface{}, len(ids))
// 	for i, id := range ids {
// 		placeholders[i] = "?"
// 		args[i] = id
// 	}
//
// 	query := fmt.Sprintf(`
// 	SELECT
// 		l.id, l.name, l.description, l.expected_input, h.hint_text
// 	FROM levels l
// 	LEFT JOIN hints h ON l.id = h.level_id
// 	WHERE l.id IN (%s)
// 	`, strings.Join(placeholders, ","))
//
// 	rows, err := s.db.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: query error: %w", op, err)
// 	}
// 	defer rows.Close()
//
// 	// Используем map для группировки уровней по id
// 	levelMap := make(map[int]*models.Level)
//
// 	for rows.Next() {
// 		var (
// 			id              int
// 			name            string
// 			description     sql.NullString
// 			expectedInput   sql.NullString
// 			hintText        sql.NullString
// 		)
//
// 		if err := rows.Scan(&id, &name, &description, &expectedInput, &hintText); err != nil {
// 			return nil, fmt.Errorf("%s: scan error: %w", op, err)
// 		}
//
// 		// Если уровень уже существует в карте, просто добавляем подсказку
// 		level, exists := levelMap[id]
// 		if !exists {
// 			level = &models.Level{
// 				Id:            id,
// 				Name:          name,
// 				Description:   description.String,
// 				ExpectedInput: expectedInput.String,
// 				Hints:         []string{},
// 			}
// 			levelMap[id] = level
// 		}
//
// 		if hintText.Valid && hintText.String != "" {
// 			level.Hints = append(level.Hints, hintText.String)
// 		}
// 	}
//
// 	// Конвертируем map в срез
// 	levels := make([]models.Level, 0, len(levelMap))
// 	for _, lvl := range levelMap {
// 		levels = append(levels, *lvl)
// 	}
//
// 	return levels, nil
// }

func (s *Storage) Levels(ctx context.Context, ids ...int) ([]models.Level, error) {
	const op = "storage.sqlite.Levels"
	query := `
	SELECT 
		l.id, l.name, l.description, l.expected_input, h.hint_text
	FROM levels l
	LEFT JOIN hints h on l.id = h.level_id
	WHERE l.id in (%s)
	`

	// map to id1,id2,id3 for db selection
	strIds := common.Map(strconv.Itoa, ids...)
	selectBy := strings.Join(strIds, ",")

	rows, err := s.db.Query(query, selectBy)
	if err != nil {
		return nil, fmt.Errorf("failed to get levels: %s %w", op, err)
	}
	defer rows.Close()

	levels := make([]domain.Level, 0, len(ids))

	for rows.Next() {
		var (
			lvlId                            int
			lvlName                          string
			lvlDesc, hintText, expectedInput sql.NullString
		)

		if err := rows.Scan(&lvlId, &lvlName, &lvlDesc, &expectedInput, &hintText); err != nil {
			return nil, fmt.Errorf("scan error: %s %w", op, err)
		}

		level := domain.Level{
			Id:            lvlId,
			Name:          lvlName,
			Description:   lvlDesc.String,
			Hints:         []string{hintText.String},
			ExpectedInput: expectedInput.String,
		}

		levels = append(levels, level)
	}
	return levels, nil
}

func (s *Storage) Hint(ctx context.Context, id int) (models.Hint, error) {
	const op = "storage.sqlite.Hint"
	query := `
	SELECT 
		h.id, h.level_id, h.hint_text
	FROM hints h
	LEFT JOIN levels l on l.id = h.level_id
	WHERE h.id = ?
	`

	rows, err := s.db.Query(query, id)
	if err != nil {
		return models.Hint{}, fmt.Errorf("failed to get hint: %s %w", op, err)
	}

	defer rows.Close()

	hint := domain.Hint{}
	for rows.Next() {
		var (
			hintId, levelId int
			text            sql.NullString
		)
		if err := rows.Scan(&hintId, &levelId, &text); err != nil {
			return domain.Hint{}, fmt.Errorf("scan error: %s %w", op, err)
		}
		hint.Id = hintId
		hint.LevelId = levelId
		hint.Text = text.String
	}

	return hint, nil
}

func (s *Storage) HintsByLevelId(ctx context.Context, id int) ([]models.Hint, error) {
	const op = "storage.sqlite.HintsByLevelID"
	query := `
	SELECT 
		h.id, h.level_id, h.hint_text
	FROM hints h
	WHERE h.level_id = ?
	`

	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: query error: %w", op, err)
	}
	defer rows.Close()

	var hints []models.Hint
	for rows.Next() {
		var (
			hintId, lvlId int
			text          sql.NullString
		)
		if err := rows.Scan(&hintId, &lvlId, &text); err != nil {
			return nil, fmt.Errorf("%s: scan error: %w", op, err)
		}
		hints = append(hints, models.Hint{
			Id:      hintId,
			LevelId: lvlId,
			Text:    text.String,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return hints, nil
}
