package service_task_done

import (
	"log"
	"time"

	dbase "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/db"
	service_date_next "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/service/date_next"
	service_task_del "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/service/task_del"
	service_task_edit "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/service/task_edit"
	service_task_info "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/service/task_info"
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

// Отметить задачу выполненной
// Принимает на вход:
//   - идентификатор задачи
//
// Возвращает:
//   - Ошибка
func (service *Service) Done(idTask string) error {

	log := service.log
	log.Printf("   Service 'Done' Begin with idTask = %s\n", idTask)
	// Получить информацию о задаче
	task, err := service_task_info.New(log, service.repository).Info(idTask)
	if err != nil {
		log.Printf("   Service return: err = %s\n", err.Error())
		return err
	}
	// Анализ полученной информации и принятия решения, что делаем
	if task.Repeat == "" {
		// Удаляем запись, которая не требует повторов
		err = service_task_del.New(log, service.repository).Delete(task.Id)
	} else {
		// Рассчитаем следующую дату повтора и сохраним изменения
		task.Date, err = service_date_next.New(log).NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			log.Printf("   Service return: err = %s\n", err.Error())
			return err
		}
		err = service_task_edit.New(log, service.repository).Edit(&task)
	}

	return err
}
