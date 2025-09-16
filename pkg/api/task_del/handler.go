package api_task_del

import (
	"encoding/json"
	"log"
	"net/http"

	dbase "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/db"
	models "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/models"
	service_task_del "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/service/task_del"
)

// Структура обработчика
type Handle struct {
	log     *log.Logger               // логгер
	service *service_task_del.Service // указатель на сервис, реализующий бизнес логику
}

// Инициализация экземпляра структуры Handle
func New(log *log.Logger, repository *dbase.Repository) *Handle {
	return &Handle{log: log, service: service_task_del.New(log, repository)}
}

// Обработка http запроса удаления задачи
func (handler *Handle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	log := handler.log
	log.Println("=== Delete Task Begin ===")
	log.Printf("Receive %s: '%s'\n", request.Method, request.URL.String())

	var (
		responseError models.Error // Модель, которую нужно вернуть в случае ошибки
		err           error
	)

	// Пробуем получить задачу из запроса
	idFind := request.FormValue("id") // идентификатор задачи
	// Отдаем данные в бизнес логику
	log.Printf("Data for business logic: id = %s\n", idFind)
	err = handler.service.Delete(idFind)

	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		log.Printf("Error: '%s'\n", err.Error())
		responseError.Error = err.Error()
		resp, err := json.Marshal(responseError)
		if err != nil {
			log.Printf("json.Marshal return error: '%s'\n", err.Error())
		}
		http.Error(writer, string(resp), http.StatusBadRequest)
		return
	}
	log.Printf("Success: task with id = %s deleted", idFind)

	//writer.WriteHeader(http.StatusNoContent)
	_, err = writer.Write([]byte("{}")) // ради тестов добавим {} в ответ. Хотя достаточно вернуть код ответа 204...
	if err != nil {
		log.Printf("Error write response: '%s'\n", err.Error())
	}
	log.Println("===  << End >> ===")
}
