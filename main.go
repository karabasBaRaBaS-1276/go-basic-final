package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	database "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/db"
	server "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/server"
)

var db *sql.DB

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
	setDefaultEnv("TODO_PORT", "7540")           // Порт для запуска веб сервера по умолчанию
	setDefaultEnv("TODO_DBDRIVER", "sqlite")     // Драйвер для работы с БД
	setDefaultEnv("TODO_DBFILE", "scheduler.db") // Имя базы данных
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

	// Переменные окружения, которые нам нужны
	loadEnv(log)

	// База данных, с которой нам предстоит работать
	log.Println("Инициализируем базу данных")
	db, err := database.Init(os.Getenv("TODO_DBFILE"), os.Getenv("TODO_DBDRIVER"))
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()
	defer func() {
		if err = db.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Веб сервер, с которым нам предстоит работать
	webServer := server.Get(log, os.Getenv("TODO_PORT"))

	log.Println("Запускаем веб сервер")
	if err := webServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
