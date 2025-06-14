#!/bin/bash

DB_FILE="storage.db"

# Проверка: существует ли база данных
# Смотри задания тута:   https://docs.google.com/document/d/1HTJQ8QaDV1WNj_JcSdu62ZJJvmGZMzyCKzePJ7Nm7Tw/edit?usp=sharing
if [[ ! -f "$DB_FILE" ]]; then
    echo "📂 Файл базы данных '$DB_FILE' не найден. Создаю..."
    sqlite3 "$DB_FILE" "" || { echo "❌ Ошибка при создании базы данных."; exit 1; }
else
    echo "📁 Используем существующую базу данных '$DB_FILE'"
fi

# Выполнение SQL скрипта (создание таблиц + вставка данных)
sqlite3 "$DB_FILE" <<'EOF'
-- ========== СОЗДАНИЕ ТАБЛИЦ ==========

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    pass_hash BLOB NOT NULL,
    oauth_id TEXT,
	is_oauth BOOLEAN DEFAULT FALSE,
    total_attempts INTEGER DEFAULT 0,
    pass_levels INTEGER DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS apps (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS levels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    expected_input TEXT,
    start_input TEXT
);

CREATE TABLE IF NOT EXISTS hints (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    level_id INTEGER NOT NULL REFERENCES levels(id) ON DELETE CASCADE,
    hint_text TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_levels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    level_id INTEGER REFERENCES levels(id) ON DELETE CASCADE,
    is_completed BOOLEAN DEFAULT FALSE,
    last_input TEXT,
    attempt_response TEXT,
    attempts INTEGER DEFAULT 0
);

-- ========== ЗАПОЛНЕНИЕ УРОВНЕЙ ==========

INSERT INTO levels (id, name, description, expected_input) VALUES
(0, 'Начало', 'Вы — хакер-новичок, нанятый таинственным заказчиком.
Ваша цель — как можно сильнее навредить компании-конкуренту. 
Сначала вы должны пробраться в систему через уязвимые веб-интерфейсы, а
затем подорвать устойчивость серверного ПО через самописные программы на питоне.
Каждое успешно выполненное задание приближает вас к настоящему хаосу.', ''),

(1, 'Кто они такие?', 'Мы знаем, что наш конкурент хостит свой сервис на сайте с доменом /hosting.
Давай узнаем об этой компании больше, "дернув" эндпоинт /about с помощью curl.', 'curl -X GET localhost:9086/about'),

(2, 'У них не защищен сервер', 'Удивительно, но они забыли использовать простейшие схемы защиты. Мы узнали следующее об их системе:
admin@yandex.ru -> почта администратора их внутренней дашборды
localhost:9086/db/about -> эндпоинт для получения информации о схеме базы данных 
localhost:9086/files -> эндпоинт для получения файлов на системе
localhost:9086/login -> эндпоинт для авторизации
Давай попробуем узнать их БД, наверное стоит заходить под админом:', 'curl -X GET localhost:9086/db/about?user=admin'),

(3, 'Брутфорс', 'Хм, похоже нужно обладать правами суперпользователя для просмотра БД. Давай попробуем тогда
залогиниться под админом, но нужен пароль. Раз у них не настроена защита, то и пароль должен быть простым,
например 5-ти значный.', 'curl -X POST http://localhost:9086/login -d "user=admin&password=12345"'),

(4, 'Файловая система', 'А что если посмотреть их файлы на системе. Давай сделаем это, отправив query запрос с командой просмотра файлов в UNIX системах.', 'curl -X POST http://localhost:9086/files -d "cmd=ls"'),

(5, 'Дырявая система', 'Сработало. Мы получили следующие файлы: main.go, db.go, db.sql
Давай посмотрим схему БД, просто "катнув" ее.', 'curl -X POST http://localhost:9086/files -d "cmd=cat db.sql"'),

(6, 'АЗЫ SQL', 'Что ж, мы получили следующую информацию:
- У конкурентов есть таблицы users
Давай попробуем найти рут пользователя, использовав SQL инъекцию.
Мы узнали из файла db.go, что там используется такой код:
username = input("Username: ")
query = "SELECT * FROM users WHERE username = ''{username}''"', ''' OR ''a'' = ''a'''),

(7, 'В поисках пароля', 'Что ж. Мы узнали, что в системе есть рут юзер. Теперь нужно получить его пароль. Но давай сразу узнаем пароли всех пользователей. Из того же файла мы видим, что в коде используется такая строчка:
SELECT name, email FROM clients WHERE name LIKE ''%$search%''.
Давай используем таблицу password, чтобы узнать пароли.', ''' UNION SELECT password FROM users --'),

(8, 'Делай грязь', 'А теперь нам нужно устроить настоящий хаос в системе конкурентов. Давай дропнем таблицу их пользователей. Большая вероятность, что они не делают снапшоты их БД. Берем тот же код из файла для работы с БД и видим: INSERT INTO feedback (text) VALUES (''$comment'')', '''; DROP TABLE users;--'),

(9, 'Stack Overflow', 'Ну давай напишем код. Давай сломаем их основной сервер, сделав переполнение стека. Кажется, что это можно сделать простой рекурсией в питоне.', 'def f(): f()\nf()'),

(10, 'OOM', 'А теперь заполним всю доступную память и крашнем всю систему. Давай также напишем скрипт, который вызовет Memory Error.', 'a = []\nwhile True: a.append(''X''*10**6)');

-- ========== ПОДСКАЗКИ ДЛЯ УРОВНЕЙ ==========

INSERT INTO hints (level_id, hint_text) VALUES
(1, 'Попробуй использовать curl'),
(1, '-X GET для запроса'),

(2, 'Попробуй использовать curl'),
(2, '-X GET для запроса'),
(2, 'Используй ? для query'),

(3, '-X POST для запроса'),
(3, 'Используй -d для передачи form'),
(3, '& для передачи нескольких параметров'),
(3, 'Параметры admin, password'),

(4, '-X POST для запроса'),
(4, 'Используй ? для query'),
(4, 'Используй ls'),

(5, '-X POST для запроса'),
(5, 'Используй ? для query'),
(5, 'Используй cat'),

(6, 'Используй конструкцию OR'),
(6, 'Попробуй сравнить два одинаковых слова'),

(7, 'Используй конструкцию UNION SELECT'),
(7, 'Используй -- для комментирования'),

(8, 'Используй конструкцию DROP TABLE'),
(8, 'Используй -- для комментирования'),

(9, 'Используй рекурсию'),
(9, 'Пусть функция вызывает сама себя'),

(10, 'Давай создавать много массивов'),
(10, 'Давай создавать много строк');

EOF

echo "✅ База данных успешно создана и заполнена."

