package service_authorization

import (
	"errors"
	"log"
	"os"

	models "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/models"
)

// Структура сервиса
type Service struct {
	log *log.Logger // логгер
}

// Инициализация экземпляра структуры Service
func New(log *log.Logger) *Service {
	return &Service{log: log}
}

// Авторизация пользователя
// Принимает на вход:
//   - task - информация о новой задаче
//
// Возвращает:
//   - Токен доступа
//   - Ошибка
func (service *Service) Signin(auth *models.Auth) (string, error) {
	log := service.log
	log.Printf("   Service 'Signin' Begin\n")
	// Проверяем на корректность указанных данных
	if auth.Password != os.Getenv("TODO_PASSWORD") {
		return "", errors.New("incorrect password")
	}

	// Выпускаем токен
	return "todo_result_token", nil
}
