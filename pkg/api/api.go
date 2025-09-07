package api

import (
	"log"
	"net/http"

	nextdate "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/api/next_date"
)

// Инициализация API обработчиков
func Init(log *log.Logger, mux *http.ServeMux) {

	// Получить следующую дату
	handlerNextDate := nextdate.New(log)
	mux.HandleFunc("GET /api/nextdate", handlerNextDate.ServeHTTP)
}
