package service_tasks_get

import (
	"time"

	dbase "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/db"
	models "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/models"
)

const limit = 50

// Структура сервиса
type Service struct {
	repository *dbase.Repository // указатель на хранилище
}

// Инициализация экземпляра структуры Service
func New(repository *dbase.Repository) *Service {
	return &Service{repository: repository}
}

// Получить список ближайших задач планировщика
// Принимает на вход:
//   - search - строка для поиска
//
// Возвращает:
//   - Массив найденных записей
//   - Ошибка
func (service *Service) Get(search string) (models.TaskList, error) {
	// Проверяем на корректность указанных в задаче данных
	var (
		likeExp string
		dateExp string
	)
	if search != "" {
		// Может передана дата?
		date, err := time.Parse("02.01.2006", search)
		if err != nil {
			likeExp = search
		} else {
			dateExp = date.Format("20060102")
		}
	}

	// Получаем список задач
	taskList, err := service.repository.GetTasks(likeExp, dateExp, limit)

	return taskList, err
}
