package service_task_info

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

// Получть информацию о задаче
// Принимает на вход:
//   - идентификатор задачи
//
// Возвращает:
//   - task - информация о задаче
//   - Ошибка
func (service *Service) Info(idTask string) (models.Task, error) {
	// Получаем задачу из БД
	task, err := service.repository.InfoTask(idTask)
	return task, err
}
