package main

import (
	"log"

	"github.com/joho/godotenv"
	server "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/server"
)

func main() {
	log := log.Default()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	webServer := server.Get(log)

	log.Println("Запускаем веб сервер")
	if err := webServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
