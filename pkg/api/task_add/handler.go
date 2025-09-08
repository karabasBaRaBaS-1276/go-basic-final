package api_task_add

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	dbase "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/db"
	models "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/models"
	service_task_add "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/service/task_add"
)

// Структура обработчика
type Handle struct {
	log     *log.Logger               // логгер
	service *service_task_add.Service // указатель на сервис, реализующий бизнес логику
}

// Инициализация экземпляра структуры Handle
func New(log *log.Logger, repository *dbase.Repository) *Handle {
	return &Handle{log: log, service: service_task_add.New(repository)}
}

// Обработка http запроса добавления новой задачи
func (handler *Handle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	log := handler.log
	log.Println("=== Add Task Begin ===")
	log.Printf("Receive %s: '%s'\n", request.Method, request.URL.String())

	// Пробуем получить задачу из запроса
	var task models.Task
	var responseError models.Error
	var responseId models.TaskId
	var err error
	var buf bytes.Buffer
	//
	_, err = buf.ReadFrom(request.Body)
	if err != nil {
		log.Printf("Read from body return error: '%s'\n", err.Error())
		responseError.Error = err.Error()
		resp, err := json.Marshal(responseError)
		if err != nil {
			log.Printf("json.Marshal return error: '%s'\n", err.Error())
		}
		http.Error(writer, string(resp), http.StatusBadRequest)
		return
	}
	// десериализуем JSON в Task
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		log.Printf("json.Unmarshal return error: '%s'\n", err.Error())
		responseError.Error = err.Error()
		resp, err := json.Marshal(responseError)
		if err != nil {
			log.Printf("json.Marshal return error: '%s'\n", err.Error())
		}
		http.Error(writer, string(resp), http.StatusBadRequest)
		return
	}
	// Отдаем данные в бизнес логику
	log.Printf("Data for business logic: %#v\n", task)
	result, err := handler.service.Add(&task)
	if err != nil {
		log.Printf("Error: %s'\n", err.Error())
		responseError.Error = err.Error()
		resp, err := json.Marshal(responseError)
		if err != nil {
			log.Printf("json.Marshal return error: '%s'\n", err.Error())
		}
		http.Error(writer, string(resp), http.StatusBadRequest)
		return
	}
	log.Printf("Success: %s'\n", result)
	responseId.Id = result
	resp, err := json.Marshal(responseId)
	if err != nil {
		log.Printf("json.Marshal return error: '%s'\n", err.Error())
	}
	_, err = writer.Write(resp)
	if err != nil {
		log.Printf("Error write response: '%s'\n", err.Error())
	}

	log.Println("===  << End >> ===")

	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
}
