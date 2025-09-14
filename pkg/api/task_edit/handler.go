package api_task_info

import (
	"encoding/json"
	"log"
	"net/http"

	api_common "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/api/api_common"
	dbase "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/db"
	models "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/models"
	service_task_edit "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/service/task_edit"
)

// Структура обработчика
type Handle struct {
	log     *log.Logger                // логгер
	service *service_task_edit.Service // указатель на сервис, реализующий бизнес логику
}

// Инициализация экземпляра структуры Handle
func New(log *log.Logger, repository *dbase.Repository) *Handle {
	return &Handle{log: log, service: service_task_edit.New(repository)}
}

// Обработка http запроса редактирования задачи
func (handler *Handle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	log := handler.log
	log.Println("=== Edit Task Begin ===")
	log.Printf("Receive %s: '%s'\n", request.Method, request.URL.String())

	var (
		task          models.Task  // Модель, которую ожидаем в запросе
		responseError models.Error // Модель, которую нужно вернуть в случае ошибки
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
	err = handler.service.Edit(&task)

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
	log.Printf("Success: task with id = %s changed", task.Id)

	//writer.WriteHeader(http.StatusNoContent)
	_, err = writer.Write([]byte("{}")) // ради тестов добавим {} в ответ. Хотя достаточно вернуть код ответа 204...
	if err != nil {
		log.Printf("Error write response: '%s'\n", err.Error())
	}
	log.Println("===  << End >> ===")
}
