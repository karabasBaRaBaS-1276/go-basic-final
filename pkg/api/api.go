package api

import (
	"log"
	"net/http"

	authorization "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/api/authorization"
	nextdate "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/api/date_next"
	addtask "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/api/task_add"
	deltask "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/api/task_del"
	donetask "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/api/task_done"
	edittask "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/api/task_edit"
	infotask "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/api/task_info"
	gettasks "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/api/tasks_get"
	dbase "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/db"
)

// Инициализация API обработчиков
func Init(log *log.Logger, mux *http.ServeMux, repository *dbase.Repository) {

	// Авторизация
	handlerAuthorization := authorization.New(log)
	mux.Handle("POST /api/signin", PanicMiddleware(http.HandlerFunc(handlerAuthorization.ServeHTTP)))

	// Получить следующую дату
	handlerNextDate := nextdate.New(log)
	mux.Handle("GET /api/nextdate", PanicMiddleware(http.HandlerFunc(handlerNextDate.ServeHTTP)))

	// Добавить новую задачу
	handlerAddTask := addtask.New(log, repository)
	mux.Handle("POST /api/task", PanicMiddleware(JWTAuthMiddleware(log, http.HandlerFunc(handlerAddTask.ServeHTTP))))

	// Получить данные о задаче
	handlerInfoTask := infotask.New(log, repository)
	mux.Handle("GET /api/task", PanicMiddleware(JWTAuthMiddleware(log, http.HandlerFunc(handlerInfoTask.ServeHTTP))))

	// Изменить данные о задаче
	handlerEditTask := edittask.New(log, repository)
	mux.Handle("PUT /api/task", PanicMiddleware(JWTAuthMiddleware(log, http.HandlerFunc(handlerEditTask.ServeHTTP))))

	// Получить список ближайших задач
	handlerGetTasks := gettasks.New(log, repository)
	mux.Handle("GET /api/tasks", PanicMiddleware(JWTAuthMiddleware(log, http.HandlerFunc(handlerGetTasks.ServeHTTP))))

	// Удалить задачу
	handlerDelTask := deltask.New(log, repository)
	mux.Handle("DELETE /api/task", PanicMiddleware(JWTAuthMiddleware(log, http.HandlerFunc(handlerDelTask.ServeHTTP))))

	// Отместить задачу выполненной
	handlerDoneTask := donetask.New(log, repository)
	mux.Handle("POST /api/task/done", PanicMiddleware(JWTAuthMiddleware(log, http.HandlerFunc(handlerDoneTask.ServeHTTP))))
}
