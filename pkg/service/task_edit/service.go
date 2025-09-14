package service_task_edit

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

// Изменить задачу
// Принимает на вход:
//   - task - новая информация о задаче
//
// Возвращает:
//   - Ошибка
func (service *Service) Edit(task *models.Task) error {
	// Проверяем на корректность указанных в задаче данных
	task, err := task.CheckAndEnrichNewTask()
	if err != nil {
		return err
	}
	// Редактируем таску в БД
	err = service.repository.EditTask(task)
	return err
}
