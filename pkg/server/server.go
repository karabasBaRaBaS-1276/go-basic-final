package server

import (
	"log"
	"net/http"
	"time"
)

// Запускает веб сервер
// Принимает на вход указатель на логгер и порт для запуска сервера
// Возвращает:
//   - указатель на настроенный сервер для старта (*http.Server)
func Get(log *log.Logger, port string) *http.Server {

	webDir := "./web" // Путь относительно рабочей дирректории
	// http-роутер
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	server := &http.Server{
		Addr:         ":" + port,       // порт был ранее определен
		Handler:      mux,              // http-роутер
		ErrorLog:     log,              // указатель на логер
		ReadTimeout:  5 * time.Second,  // таймаут для чтения
		WriteTimeout: 10 * time.Second, // таймаут для записи
		IdleTimeout:  15 * time.Second, // таймаут ожидания следующего запроса
	}

	return server
}
