CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,o
    pass_hash BLOB NOT NULL,
);

CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS apps (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE,
);

CREATE TABLE IF NOT EXISTS levels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
);

CREATE TABLE IF NOT EXISTS user_levels (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    level_id INTEGER REFERECES levels(id) ON DELETE CASCADE,
    is_completed BOOLEAN DEFAULT FALSE, --пройден ли уровень
    last_input TEXT, --последний ввод пользователя (например: последний веденный код)
    attempt_response TEXT, --последний ответ сервера на попытку пройти уровень 
    attempts INTEGER DEFAULT 0, --количество попыток пройти уровень (будет увеличиваться при неудаче)
);