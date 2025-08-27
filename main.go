package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	server "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/server"
)

// Загрузка переменных окружения
// Принимает на вход указатель на логгер
func loadEnv(log *log.Logger) {
	if err := godotenv.Load(); err != nil {
		// Если файла нет - это нормально для production/CI
		if !os.IsNotExist(err) {
			// Прекращаем работу если иные проблемы с файлом
			log.Fatal("Ошибка загрузки .env файла:", err)
		}
	}
	// Установка переменных по умолчанию
	setDefaultEnv("TODO_PORT", "7540") // Порт для запуска веб сервера по умолчанию
}

// Установка окружений по умолчанию, если их нет
func setDefaultEnv(key, defaultValue string) {
	if os.Getenv(key) == "" {
		log.Printf("Умолчание для %s = %s", key, defaultValue)
		os.Setenv(key, defaultValue)
	}
}

func main() {
	log := log.Default()

	loadEnv(log)

	webServer := server.Get(log)

	log.Println("Запускаем веб сервер")
	if err := webServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
