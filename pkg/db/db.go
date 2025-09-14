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
		INSERT INTO scheduler 
			(date, title, comment, repeat) 
			VALUES (:date, :title, :comment, :repeat);
		`
	getTasksByTitleOrCommentQuery = `
		SELECT id, date, title, comment, repeat FROM scheduler
		WHERE title LIKE :likeExp OR comment LIKE :likeExp
		ORDER BY date LIMIT :limit;
	`
	getTasksByDateQuery = `
		SELECT id, date, title, comment, repeat FROM scheduler
		WHERE date = :dateExp
		LIMIT :limit;
	`
	getTasksQuery = `
		SELECT id, date, title, comment, repeat FROM scheduler
		ORDER BY date LIMIT :limit;
	`
	getTaskById = `
		SELECT id, date, title, comment, repeat FROM scheduler
		WHERE id = :id;
	`
	editTaskQueryById = `
		UPDATE scheduler 
		SET 
			date = :date, 
			title = :title, 
			comment = :comment, 
			repeat = :repeat
		WHERE id = :id;
	`
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

// Получить список ближайших задач планировщика
func (repository Repository) GetTasks(likeExp string, dateExp string, limit int) (models.TaskList, error) {

	var (
		rows     *sql.Rows
		err      error
		tasks    []models.Task
		taskList models.TaskList
	)

	switch {
	case likeExp != "":
		pattern := "%" + likeExp + "%"
		rows, err = repository.DBase.Query(
			getTasksByTitleOrCommentQuery,
			sql.Named("likeExp", pattern),
			sql.Named("limit", limit),
		)
	case dateExp != "":
		rows, err = repository.DBase.Query(
			getTasksByDateQuery,
			sql.Named("dateExp", dateExp),
			sql.Named("limit", limit),
		)
	default:
		rows, err = repository.DBase.Query(
			getTasksQuery,
			sql.Named("limit", limit),
		)
	}

	if err != nil {
		taskList.Tasks = tasks
		return taskList, fmt.Errorf("failed to find tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id      string
			title   string
			date    string
			comment string
			repeat  string
		)

		err := rows.Scan(&id, &date, &title, &comment, &repeat)
		if err != nil {
			taskList.Tasks = tasks
			return taskList, fmt.Errorf("failed to read record: %w", err)
		}

		tasks = append(
			tasks,
			models.Task{
				Id:      id,
				Title:   title,
				Date:    date,
				Comment: comment,
				Repeat:  repeat,
			},
		)

	}
	taskList.Tasks = tasks
	if taskList.Tasks == nil {
		taskList.Tasks = make([]models.Task, 0)
	}
	return taskList, nil
}

// Получить информацию о задаче
func (repository Repository) InfoTask(idTask string) (models.Task, error) {
	if idTask == "" {
		return models.Task{}, fmt.Errorf("failed to read record: 'id' cannot be empty")
	}

	row := repository.DBase.QueryRow(
		getTaskById,
		sql.Named("id", idTask),
	)

	var (
		id      string
		title   string
		date    string
		comment string
		repeat  string
	)

	err := row.Scan(&id, &date, &title, &comment, &repeat)
	if err != nil {
		return models.Task{}, fmt.Errorf("failed to read record (id = %s): %w", idTask, err)
	}

	task := models.Task{
		Id:      id,
		Title:   title,
		Date:    date,
		Comment: comment,
		Repeat:  repeat,
	}

	return task, nil

}

// Редактировать задачу в планировщике
func (repository *Repository) EditTask(task *models.Task) error {

	result, err := repository.DBase.Exec(
		editTaskQueryById,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.Id),
	)
	if err != nil {
		return fmt.Errorf("failed to edit record: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to edit record: %w", err)
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for edit task`)
	}

	return nil
}
