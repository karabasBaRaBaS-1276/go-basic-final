package service_task_add

import (
	dbase "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/db"
	models "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/models"
)

// Структура сервиса
type Service struct {
	repository *dbase.Repository // указатель на хранилище
}

// Инициализация экземпляра структуры Service
func New(repository *dbase.Repository) *Service {
	return &Service{repository: repository}
}

// Добавить задачу в планировщик
// Принимает на вход:
//   - task - информация о новой задаче
//
// Возвращает:
//   - Идентификатор новой записи планировщика
//   - Ошибка
func (service *Service) Add(task *models.Task) (string, error) {
	// Проверяем на корректность указанных в задаче данных
	task, err := task.CheckAndEnrichNewTask()
	if err != nil {
		return "", err
	}
	// Добавляем новую таску в БД
	result, err := service.repository.AddTask(task)
	return result, err
}
