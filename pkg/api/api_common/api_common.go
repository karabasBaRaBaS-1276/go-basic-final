package api_common

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/karabasBaRaBaS-1276/go-basic-final/pkg/models"
)

// Структура обработчика с методами, которые переиспользуются в иных обработчиках
type ApiCommon struct {
	log *log.Logger // логгер
}

func New(log *log.Logger) *ApiCommon {
	return &ApiCommon{log}
}

// Получить информацию о задаче из запроса
func (ApiCommon) GetTaskFromJson(request *http.Request) (models.Task, error) {
	var (
		task models.Task // Модель, которую ожидаем в запросе
		err  error
		buf  bytes.Buffer
	)
	_, err = buf.ReadFrom(request.Body)
	if err != nil {
		log.Printf("Read from body return error: '%s'\n", err.Error())
		return models.Task{}, err
	}
	// десериализуем JSON в Task
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		log.Printf("json.Unmarshal return error: '%s'\n", err.Error())
		return models.Task{}, err
	}
	return task, nil
}
