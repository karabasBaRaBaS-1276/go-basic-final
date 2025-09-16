package api_task_add

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/karabasBaRaBaS-1276/go-basic-final/pkg/api/api_common"
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
	return &Handle{log: log, service: service_task_add.New(log, repository)}
}

// Обработка http запроса добавления новой задачи
func (handler *Handle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	log := handler.log
	log.Println("=== Add Task Begin ===")
	log.Printf("Receive %s: '%s'\n", request.Method, request.URL.String())

	var (
		task          models.Task   // Модель, которую ожидаем в запросе
		responseError models.Error  // Модель, которую нужно вернуть в случае ошибки
		responseId    models.TaskId // Модель, которую нужно вернуть в случае успеха
		err           error
	)

	apiCommon := api_common.New(log)
	// Пробуем получить задачу из запроса
	task, err = apiCommon.GetTaskFromJson(request)

	if err != nil {
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
	log.Printf("Success: '%s'\n", result)
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
}
