package api_tasks_get

import (
	"encoding/json"
	"log"
	"net/http"

	dbase "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/db"
	"github.com/karabasBaRaBaS-1276/go-basic-final/pkg/models"
	service_task_get "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/service/tasks_get"
)

// Структура обработчика
type Handle struct {
	log     *log.Logger               // логгер
	service *service_task_get.Service // указатель на сервис, реализующий бизнес логику
}

// Инициализация экземпляра структуры Handle
func New(log *log.Logger, repository *dbase.Repository) *Handle {
	return &Handle{log: log, service: service_task_get.New(repository)}
}

// Обработка http запроса получения списка ближайших задач
func (handler *Handle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	log := handler.log
	log.Println("=== Get Tasks Begin ===")
	log.Printf("Receive %s: '%s'\n", request.Method, request.URL.String())

	var (
		responseError models.Error // Модель, которую нужно вернуть в случае ошибки
	)
	search := request.FormValue("search") // строка поиска

	// Отдаем данные в бизнес логику
	log.Printf("Data for business logic: search = '%s'\n", search)
	taskList, err := handler.service.Get(search)
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
	log.Printf("Success: %d record found\n", len(taskList.Tasks))

	resp, err := json.Marshal(taskList)
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
