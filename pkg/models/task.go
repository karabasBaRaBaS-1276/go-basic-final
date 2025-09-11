package models

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	service_date_next "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/service/date_next"
)

const (
	titleRegex   = "^[A-Za-zА-Яа-яёЁ0-9 -]{0,128}$"
	commentRegex = "^[A-Za-zА-Яа-яёЁ0-9 -,.!?:]{0,512}$"
)

var errorTitleIsEmpty = errors.New("variable 'title' cannot be empty")
var errorTitleFormat = errors.New("incorrect format in 'title' variable")
var errorCommentFormat = errors.New("incorrect format in 'comment' variable")
var errorInvalidRegex = errors.New("value does not match regular expression")

var validTitleRegex = regexp.MustCompile(titleRegex)
var validCommentRegex = regexp.MustCompile(commentRegex)

// Структура, полностью описывающая задачу планировщика
type Task struct {
	Id      string `json:"id"`      // Идентификатор задачи
	Title   string `json:"title"`   // заголовок задачи. Обязательное поле
	Date    string `json:"date"`    // дата задачи в формате 20060102
	Comment string `json:"comment"` // комментарий к задаче
	Repeat  string `json:"repeat"`  // правило повторения
}

// Структура, описывающая Id задачи планировщика
type TaskId struct {
	Id string `json:"id"` // Идентификатор задачи
}

// Структура, описывающая список задачи планировщика
type TaskList struct {
	Tasks []Task `json:"tasks"`
}

func (task *Task) CheckAndEnrichNewTask() (*Task, error) {

	if task.Title == "" {
		return task, errorTitleIsEmpty
	} else {
		if !validTitleRegex.MatchString(task.Title) {
			return task, fmt.Errorf("%w: %w: %s %s", errorTitleFormat, errorInvalidRegex, task.Title, titleRegex)
		}
	}
	if task.Comment != "" {
		if !validCommentRegex.MatchString(task.Comment) {
			return task, fmt.Errorf("%w: %w: %s %s", errorCommentFormat, errorInvalidRegex, task.Comment, commentRegex)
		}
	}
	if task.Repeat != "" {
		_, _, _, _, _, err := service_date_next.CheckAndSpliteRepeat(task.Repeat)
		if err != nil {
			return task, err
		}
	}

	if task.Date != "" {
		dateTime, err := time.Parse("20060102", task.Date)
		now := time.Now()
		if err != nil {
			return task, err
		}
		if dateTime.Before(now) {
			if task.Repeat != "" {
				date, err := service_date_next.New().NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					return task, err
				}
				task.Date = date
			} else { // Дату нельзя оставить меньше сегодняшней если нет повторов
				task.Date = now.Format("20060102")
			}
		}
	} else {
		task.Date = time.Now().Format("20060102")
	}
	return task, nil
}
