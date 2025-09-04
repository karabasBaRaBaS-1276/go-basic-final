package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var errorFormatDate = errors.New("incorrect date format in 'dstart' variable")
var errorRepeatIsEmpty = errors.New("variable 'repeat' cannot be empty")
var errorRepeatFormat = errors.New("incorrect format in 'repeat' variable")

// Вернуть следующую дату
// Принимает на вход:
//  - now — время, от которого ищется ближайшая дата;
//  - dstart — исходное время в формате 20060102, от которого начинается отсчёт повторений;
//  - repeat — правило повторения в формате:
//		d <число> - задача переносится на указанное число дней.
// 			Максимально допустимое число равно 400
//		y - задача выполняется ежегодно.
// 			При выполнении задачи дата перенесётся на год вперёд.
//	- w <через запятую от 1 до 7> - задача назначается в указанные дни недели,
// 			где 1 — понедельник, 7 — воскресенье
//	- m <через запятую от 1 до 31, -1, -2> [через запятую от 1 до 12] задача назначается в указанные дни месяца.
// 			При этом вторая последовательность чисел опциональна и указывает на определённые месяцы
// Возвращает:
//   - Строка в формате 20060102
//   - Ошибка

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", errorRepeatIsEmpty
	}

	dbegin, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("%w: %w", errorFormatDate, err)
	}

	rules := strings.Split(repeat, " ")
	if !strings.Contains("dywm", rules[0]) {
		return "", fmt.Errorf("%w: %s", errorRepeatFormat, "it should start with 'd', or 'w', or 'm', or 'y'")
	}

	var dayNumber int
	var weekNumbers []int
	switch rules[0] {
	case "d":
		if len(rules) != 2 {
			err = fmt.Errorf("after 'd' must be single integer")
			break
		}
		dayNumber, err = strconv.Atoi(rules[1])
		if err != nil {
			err = fmt.Errorf("after 'd' value '%s' is not integer", rules[1])
			break
		}
		if dayNumber > 400 {
			err = fmt.Errorf("after 'd' value '%s' is not integer", rules[1])
			break
		}
	case "y":
		if len(rules) != 1 {
			err = fmt.Errorf("instead '%s' must be 'y'", repeat)
			break
		}
	case "w":
		if len(rules) != 2 {
			err = fmt.Errorf("after 'w' must be integer number")
			break
		}
		weeks := strings.Split(rules[1], ",")
		for _, v := range weeks {
			dayNumber, err = strconv.Atoi(v)
			if err != nil {
				err = fmt.Errorf("after 'w' value '%s' is not integer", v)
				break
			}
			if (dayNumber < 1) || (dayNumber > 7) {
				err = fmt.Errorf("after 'w' value '%s' is not between 1 and 7", v)
				break
			}
			weekNumbers = append(weekNumbers, dayNumber)

		}
	}

	if err != nil {
		return "", fmt.Errorf("%w: %w", errorRepeatFormat, err)
	}

	fmt.Printf("dstart = %s; dayNumber = %d; weekNumbers = %d\n", dbegin, dayNumber, weekNumbers)

	return "", err
}
