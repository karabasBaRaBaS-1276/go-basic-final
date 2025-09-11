package api

import (
	"log"
	"net/http"

	nextdate "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/api/date_next"
	addtask "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/api/task_add"
	gettasks "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/api/tasks_get"
	dbase "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/db"
)

// Инициализация API обработчиков
func Init(log *log.Logger, mux *http.ServeMux, repository *dbase.Repository) {

	// Получить следующую дату
	handlerNextDate := nextdate.New(log)
	mux.HandleFunc("GET /api/nextdate", handlerNextDate.ServeHTTP)

	// Добавить новую задачу
	handlerAddTask := addtask.New(log, repository)
	mux.HandleFunc("POST /api/task", handlerAddTask.ServeHTTP)

	// Получить список ближайших задач
	handlerGetTasks := gettasks.New(log, repository)
	mux.HandleFunc("GET /api/tasks", handlerGetTasks.ServeHTTP)
}
