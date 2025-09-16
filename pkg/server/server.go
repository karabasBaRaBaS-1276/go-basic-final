package server

import (
	"log"
	"net/http"
	"time"

	api "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/api"
	dbase "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/db"
)

// Запускает веб сервер
// Принимает на вход указатель на логгер, адрес для запуска сервера и указатель на хранилище
// Возвращает:
//   - указатель на настроенный сервер для старта (*http.Server)
func Get(log *log.Logger, address string, repository *dbase.Repository) *http.Server {

	webDir := "./web" // Путь относительно рабочей дирректории
	// http-роутер
	mux := http.NewServeMux()

	api.Init(log, mux, repository)                     // Обработчики API
	mux.Handle("/", http.FileServer(http.Dir(webDir))) // Обработчик статики

	server := &http.Server{
		Addr:         address,          // адрес для запуска
		Handler:      mux,              // http-роутер
		ErrorLog:     log,              // указатель на логер
		ReadTimeout:  5 * time.Second,  // таймаут для чтения
		WriteTimeout: 10 * time.Second, // таймаут для записи
		IdleTimeout:  15 * time.Second, // таймаут ожидания следующего запроса
	}

	return server
}
