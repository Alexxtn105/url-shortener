// internal/storage/sqlite/sqlite.go
// Для установки sqlite:
// go get github.com/mattn/go-sqlite3
package sqlite

import (
	"database/sql"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/mattn/go-sqlite3"
)

// Структура объекта Storage
type Storage struct {
	db *sql.DB //из пакета "database/sql"
}

// Конструктор объекта Storage
func NewStorage(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage" // Имя текущей функции для логов и ошибок

	// Подключаемся к БД
	db, err := sql.Open("sqlite3", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// создаем таблицу, если ее еще нет
	// TODO: можно прикрутить миграции, для тренировки
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
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

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	// Подготавливаем запрос (проверка корректности синтаксиса)
	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	//выполняем запрос
	res, err := s.db.Exec(urlToSave, alias)
	if err != nil {
		// Здесь мы приводим полученную ошибку ко внутреннему типу библиотеки sqlite3,
		// чтобы посмотреть, не является ли эта ошибка sqlite3.ErrConstraintUnique.
		// Если это так, значит, мы попытались добавить дубликат имеющейся записи. Об этом мы сообщим в вызывающую функцию, вернув уже свою ошибку для данной ситуации: storage.ErrURLExists. Получив ее, сервер сможет сообщить клиенту о том, что такой alias у нас уже есть.
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	//Возвращаем ID
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"
}
