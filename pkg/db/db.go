package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/karabasBaRaBaS-1276/go-basic-final/pkg/models"
	_ "modernc.org/sqlite"
)

const (
	schema = `
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			date CHAR(8) NOT NULL DEFAULT "",
			title VARCHAR(128) NOT NULL,
			comment TEXT,
			repeat VARCHAR(128)
		);
		
		CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);
		`
	addTaskQuery = `
		INSERT INTO scheduler (
			date, title, comment, repeat) 
			VALUES (:date, :title, :comment, :repeat
		);`
)

type Repository struct {
	DBase *sql.DB
}

// Инициализация экземпляра структуры Handle
func New(dbFile string, dbDriver string) (*Repository, error) {
	db, err := Init(dbFile, dbDriver)
	return &Repository{DBase: db}, err
}

// Инициализация базы данных
// Принимает на вход указатель на логгер, имя БД и имя драйвера БД
// Возвращает:
//   - указатель на БД (*sql.DB)
//   - ошибку (error)
func Init(dbFile, dbDriver string) (*sql.DB, error) {

	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open(dbDriver, dbFile)
	if err != nil {
		return nil, err
	}

	if install {
		if _, err = db.Exec(schema); err != nil {
			return nil, err
		}
	}

	return db, nil
}

// Добавить новую задачу в планировщик
func (repository *Repository) AddTask(task *models.Task) (string, error) {

	result, err := repository.DBase.Exec(
		addTaskQuery,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return "", fmt.Errorf("failed to add record: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("failed to add record: %w", err)
	}

	return strconv.Itoa(int(id)), nil
}
