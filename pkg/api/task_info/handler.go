package api_task_info

import (
	"encoding/json"
	"log"
	"net/http"

	dbase "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/db"
	models "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/models"
	service_task_info "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/service/task_info"
)

// Структура обработчика
type Handle struct {
	log     *log.Logger                // логгер
	service *service_task_info.Service // указатель на сервис, реализующий бизнес логику
}

// Инициализация экземпляра структуры Handle
func New(log *log.Logger, repository *dbase.Repository) *Handle {
	return &Handle{log: log, service: service_task_info.New(log, repository)}
}

// Обработка http запроса получения информации о задаче
func (handler *Handle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	log := handler.log
	log.Println("=== Info Task Begin ===")
	log.Printf("Receive %s: '%s'\n", request.Method, request.URL.String())

	var (
		task          models.Task  // Модель, которую нужно вернуть в случае успеха
		responseError models.Error // Модель, которую нужно вернуть в случае ошибки
		err           error
	)

	// Пробуем получить id задачи из запроса
	idFind := request.FormValue("id") // идентификатор задачи
	// Отдаем данные в бизнес логику
	log.Printf("Data for business logic: id = %s\n", idFind)
	task, err = handler.service.Info(idFind)

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
	log.Printf("Success: %#v\n", task)
	resp, err := json.Marshal(task)
	if err != nil {
		log.Printf("json.Marshal return error: '%s'\n", err.Error())
	}
	_, err = writer.Write(resp)
	if err != nil {
		log.Printf("Error write response: '%s'\n", err.Error())
	}

	log.Println("===  << End >> ===")
}
