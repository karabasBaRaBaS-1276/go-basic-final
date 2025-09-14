package service_task_info

import (
	"log"

	dbase "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/db"
	models "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/models"
)

// Структура сервиса
type Service struct {
	log        *log.Logger       // логгер
	repository *dbase.Repository // указатель на хранилище
}

// Инициализация экземпляра структуры Service
func New(log *log.Logger, repository *dbase.Repository) *Service {
	return &Service{log: log, repository: repository}
}

// Получить информацию о задаче
// Принимает на вход:
//   - идентификатор задачи
//
// Возвращает:
//   - task - информация о задаче
//   - Ошибка
func (service *Service) Info(idTask string) (models.Task, error) {
	log := service.log
	log.Printf("   Service 'Info' Begin with idTask = %s\n", idTask)
	// Получаем задачу из БД
	task, err := service.repository.InfoTask(idTask)
	return task, err
}
