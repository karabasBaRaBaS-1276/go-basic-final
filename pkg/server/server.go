package server

import (
	"log"
	"net/http"
	"os"
	"time"
)

// Запускает веб сервер
// Принимает на вход указатель на логгер
// Возвращает:
//   - указатель на настроенный сервер для старта (*http.Server)
func Get(log *log.Logger) *http.Server {

	webDir := "./web"              // Путь относительно рабочей дирректории
	port := os.Getenv("TODO_PORT") // Порт из переменной окружения TODO_PORT
	if port == "" {
		port = "7540" // Порт по умолчанию
	}
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
