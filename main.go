package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	dbase "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/db"
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
	setDefaultEnv("TODO_HOST", "localhost")                            // Хост для запуска веб сервера по умолчанию
	setDefaultEnv("TODO_PORT", "7540")                                 // Порт для запуска веб сервера по умолчанию
	setDefaultEnv("TODO_DBDRIVER", "sqlite")                           // Драйвер для работы с БД
	setDefaultEnv("TODO_DBFILE", "scheduler.db")                       // Имя базы данных
	setDefaultEnv("TODO_PASSWORD", "12345")                            // Пароль пользователя
	setDefaultEnv("TODO_JWT_SECRET_KEY", "Секретный ключ для подписи") // Секретный ключ для подписи jwt
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
	repository, err := dbase.New(os.Getenv("TODO_DBFILE"), os.Getenv("TODO_DBDRIVER"))
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()
	defer func() {
		if err = repository.DBase.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Веб сервер, с которым нам предстоит работать
	webServer := server.Get(log, fmt.Sprintf("%s:%s", os.Getenv("TODO_HOST"), os.Getenv("TODO_PORT")), repository)

	log.Printf("Запускаем веб сервер на http://%s\n", webServer.Addr)
	if err := webServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
