package service_task_add

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

// Добавить задачу в планировщик
// Принимает на вход:
//   - task - информация о новой задаче
//
// Возвращает:
//   - Идентификатор новой записи планировщика
//   - Ошибка
func (service *Service) Add(task *models.Task) (string, error) {
	log := service.log
	log.Printf("   Service 'Add' Begin with task = %#v\n", task)
	// Проверяем на корректность указанных в задаче данных
	task, err := task.CheckAndEnrichNewTask(log)
	if err != nil {
		log.Printf("   Service return: err = %s\n", err.Error())
		return "", err
	}
	// Добавляем новую таску в БД
	result, err := service.repository.AddTask(task)
	return result, err
}
