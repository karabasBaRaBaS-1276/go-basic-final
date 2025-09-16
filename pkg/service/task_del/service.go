package service_task_del

import (
	"log"

	dbase "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/db"
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

// Удалить задачу
// Принимает на вход:
//   - идентификатор задачи
//
// Возвращает:
//   - Ошибка
func (service *Service) Delete(idTask string) error {
	log := service.log
	log.Printf("   Service 'Delete' Begin with id = %s\n", idTask)
	// Удаляем таску в БД
	err := service.repository.DelTask(idTask)
	return err
}
