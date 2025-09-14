package api_authorization

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	models "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/models"
	service_authorization "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/service/authorization"
)

// Структура обработчика
type Handle struct {
	log     *log.Logger                    // логгер
	service *service_authorization.Service // указатель на сервис, релизующий бизнес логику
}

// Инициализация экземпляра структуры Handle
func New(log *log.Logger) *Handle {
	return &Handle{log: log, service: service_authorization.New(log)}
}

// Обработка http запроса авторизации
func (handler *Handle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	log := handler.log
	log.Println("=== Give Authorization Begin ===")
	log.Printf("Receive %s: '%s'\n", request.Method, request.URL.String())

	var (
		auth          models.Auth  // Модель, которую ожидаем в запросе
		responseError models.Error // Модель, которую нужно вернуть в случае ошибки
		JWT           models.JWT   // Модель, которую нужно вернуть в случае успеха
		err           error
		buf           bytes.Buffer
	)

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
	if err = json.Unmarshal(buf.Bytes(), &auth); err != nil {
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
	log.Printf("Data for business logic: ***auth*** struct\n")
	result, err := handler.service.Signin(&auth)

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
	log.Printf("Success: '%.8s...'\n", result)
	JWT.Token = result
	resp, err := json.Marshal(JWT)
	if err != nil {
		log.Printf("json.Marshal return error: '%s'\n", err.Error())
	}
	_, err = writer.Write(resp)
	if err != nil {
		log.Printf("Error write response: '%s'\n", err.Error())
	}

	log.Println("===  << End >> ===")
}
