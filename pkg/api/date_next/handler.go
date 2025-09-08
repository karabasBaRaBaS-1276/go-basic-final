package api_date_next

import (
	"fmt"
	"log"
	"net/http"
	"time"

	service_date_next "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/service/date_next"
)

// Структура обработчика
type Handle struct {
	log     *log.Logger                // логгер
	service *service_date_next.Service // указатель на сервис, релизующий бизнес логику
}

// Инициализация экземпляра структуры Handle
func New(log *log.Logger) *Handle {
	return &Handle{log: log, service: service_date_next.New()}
}

// Обработка http запроса получения новой даты
func (handler *Handle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	log := handler.log
	log.Println("=== Give Next Date Begin ===")
	log.Printf("Receive %s: '%s'\n", request.Method, request.URL.String())

	nowSt := request.FormValue("now")     // Текущая дата
	date := request.FormValue("date")     // Дата в задании планировщика
	repeat := request.FormValue("repeat") // Правило для повтора задания

	var now time.Time
	var err error

	if nowSt == "" {
		now = time.Now()
	} else {
		now, err = time.Parse("20060102", nowSt)
	}
	if err != nil {
		log.Printf("Error time parse 'now': '%s'\n", err.Error())
		http.Error(writer, err.Error(), http.StatusBadRequest)
	} else {
		result, err := handler.service.NextDate(now, date, repeat)
		if err != nil {
			log.Printf("Error: %s'\n", err.Error())
			http.Error(writer, err.Error(), http.StatusBadRequest)
		} else {
			log.Printf("Success: %s'\n", result)
			fmt.Fprint(writer, result)
		}
	}
	log.Println("===  << End >> ===")

	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
}
