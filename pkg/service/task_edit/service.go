package service_task_edit

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

// Изменить задачу
// Принимает на вход:
//   - task - новая информация о задаче
//
// Возвращает:
//   - Ошибка
func (service *Service) Edit(task *models.Task) error {
	log := service.log
	log.Printf("   Service 'Edit' Begin with task = %#v\n", task)
	// Проверяем на корректность указанных в задаче данных
	task, err := task.CheckAndEnrichNewTask(log)
	if err != nil {
		return err
	}
	// Редактируем таску в БД
	err = service.repository.EditTask(task)
	return err
}
