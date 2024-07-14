package sqlite

import (
	"api_server/internal/storage"
	"database/sql"
	// "errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
    db *sql.DB
}

func New(storagePath string) (*Storage, error) {
    const op = "storage.sqlite.NewStorage" // Имя текущей функции для логов и ошибок

    db, err := sql.Open("sqlite3", storagePath) // Подключаемся к БД
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    // Создаем таблицу, если ее еще нет
    stmt, err := db.Prepare(`
    CREATE TABLE IF NOT EXISTS connection(
        id INTEGER PRIMARY KEY,
        device TEXT NOT NULL,
        datetime DATETIME CURRENT_TIMESTAMP);
    CREATE INDEX IF NOT EXISTS idx_device ON connection(device);
    `)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    _, err = stmt.Exec()
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return &Storage{db: db}, nil
}

func (s *Storage) SaveGameLaunch(deviceToSave string) (int64, error) {
    const op = "storage.sqlite.SaveURL"

    // Подготавливаем запрос
    stmt, err := s.db.Prepare("INSERT INTO connection(device, datetime) values(?, CURRENT_TIMESTAMP)")
    if err != nil {
        return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
    }

    // Выполняем запрос
    res, err := stmt.Exec(deviceToSave)
    if err != nil {
        if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
            return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
        }

        return 0, fmt.Errorf("%s: execute statement: %w", op, err)
    }

    // Получаем ID созданной записи
    id, err := res.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
    }

    // Возвращаем ID
    return id, nil
}

// func (s *Storage) GetURL(alias string) (string, error) {
//     const op = "storage.sqlite.GetURL"

//     stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
//     if err != nil {
//         return "", fmt.Errorf("%s: prepare statement: %w", op, err)
//     }

//     var resURL string
    
//     err = stmt.QueryRow(alias).Scan(&resURL)
//     if errors.Is(err, sql.ErrNoRows) {
//         return "", storage.ErrURLNotFound
//     }
//     if err != nil {
//         return "", fmt.Errorf("%s: execute statement: %w", op, err)
//     }

//     return resURL, nil
// }