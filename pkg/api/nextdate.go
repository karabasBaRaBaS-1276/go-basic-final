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
//   - now — время, от которого ищется ближайшая дата;
//   - dstart — исходное время в формате 20060102, от которого начинается отсчёт повторений;
//   - repeat — правило повторения в формате:
//     d <число> - задача переносится на указанное число дней.
//     Максимально допустимое число равно 400
//     y - задача выполняется ежегодно.
//     При выполнении задачи дата перенесётся на год вперёд.
//     w <через запятую от 1 до 7> - задача назначается в указанные дни недели,
//     где 1 — понедельник, 7 — воскресенье
//     m <через запятую от 1 до 31, -1, -2> [через запятую от 1 до 12] задача назначается в указанные дни месяца.
//     При этом вторая последовательность чисел опциональна и указывает на определённые месяцы
//
// Возвращает:
//   - Строка в формате 20060102
//   - Ошибка
//
// todo. Переписать внутри. Видится, что стоит разделить внутри на подфункции для каждого символа из 'dmyw'.
// Это сделает код более читаемым
func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	dbegin, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("%w: %w", errorFormatDate, err)
	}

	typeRule, dayNumber, weekDays, monthDays, monthNumbers, err := checkAndSpliteRepeat(repeat)

	if err != nil {
		fmt.Printf("dstart = %s; err = %s\n", dbegin.Format("2006-01-02"), err.Error())
		return "", fmt.Errorf("%w: %w", errorRepeatFormat, err)
	}

	date := dbegin
	switch typeRule {
	case "d": // Дни
		for {
			date = date.AddDate(0, 0, dayNumber)
			if date.After(now) {
				break
			}
		}
	case "y": // Годы
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}
	case "w": // Указанные дни недели
		for {
			// Важно, чтобы новая дата была в день недели, указанный в массиве weekDays (например, повторения каждый понедельник)
			if date.After(now) {
				weekDay := int(date.Weekday())
				if weekDay == 0 { // вск у нас это 7
					weekDay = 7
				}
				if contains(weekDays, weekDay) {
					break
				}
			} else {
				date = now
			}

			date = date.AddDate(0, 0, 1)
		}
	case "m": // указанные дни месяца
		// <через запятую от 1 до 31, -1, -2> [через запятую от 1 до 12] задача назначается в указанные дни месяца.
		//     При этом вторая последовательность чисел опциональна и указывает на определённые месяцы
		for {
			// Важно, чтобы новая дата была в день месяца, указанный в массиве monthDays и происходила в указанный месяц monthNumbers
			// (например, повторения каждое первое число второго месяца)
			if date.After(now) {
				month := int(date.Month()) // Номер месяца
				day := date.Day()          // Номер дня
				// todo. Что делать, если в monthDays указаны -2 или -1
				//
				if contains(monthDays, day) || checkOnLastDays(monthDays, date) { // Дни совпали
					if len(monthNumbers) > 0 { // Есть указание на конкретные месяцы
						if contains(monthNumbers, month) { // Месяцы совпали
							break
						}
					} else { // Нет указания на конкретный месяц
						break
					}
				}

			} else {
				date = now
			}

			date = date.AddDate(0, 0, 1)
		}
	}

	fmt.Printf("dstart = %s; typeRule = %s; dayNumber = %d; weekDays = %d; monthDays = %d; monthNumbers = %d\n ndate = %s\n",
		dbegin.Format("2006-01-02"), typeRule, dayNumber, weekDays, monthDays, monthNumbers, date.Format("2006-01-02"))

	return date.Format("20060102"), err
}

// Проверить и разделить правила для повтора
//
// Принимает на вход:
//   - repeat - строка с описанием правила повтора. См. @NextDate
//
// Возвращает:
//
//   - typeRule - тип повторяемого правила:
//     d - повторы основаны на днях;
//     w - повторы основаны на неделях;
//     m - повторы основаны на месяцах;
//     y - повторы основаны на годах.
//   - dayNumber - смещение в днях для повтора
//   - weekDays - смещения на дни недели
//   - monthDays - смещения на дни месяца
//   - monthNumbers - смещения на номера месяцев
//   - err - ошибка, если не удалось разделить правила для повтора
func checkAndSpliteRepeat(repeat string) (typeRule string, dayNumber int, weekDays []int, monthDays []int, monthNumbers []int, err error) {

	if repeat == "" {
		err = errorRepeatIsEmpty
		return
	}

	rules := strings.Split(repeat, " ")
	if !strings.Contains("dywm", rules[0]) {
		err = fmt.Errorf("%w: %s", errorRepeatFormat, "it should start with 'd', or 'w', or 'm', or 'y'")
		return
	}
	typeRule = rules[0]

	var intNumber int

	switch typeRule {
	case "d": // дни
		if len(rules) != 2 {
			err = fmt.Errorf("after 'd' must be single integer")
			break
		}
		dayNumber, err = strconv.Atoi(rules[1])
		if err != nil {
			err = fmt.Errorf("after 'd' value '%s' is not integer", rules[1])
			break
		}
		if (dayNumber > 400) || (dayNumber < 1) {
			err = fmt.Errorf("after 'd' value '%s' is not between 1 and 400", rules[1])
			break
		}
	case "y": // год
		if len(rules) != 1 {
			err = fmt.Errorf("instead '%s' must be 'y'", repeat)
			break
		}
	case "w": // недели
		if len(rules) != 2 {
			err = fmt.Errorf("after 'w' must be integer number")
			break
		}
		weeks := strings.Split(rules[1], ",")
		for _, v := range weeks {
			intNumber, err = strconv.Atoi(v)
			if err != nil {
				err = fmt.Errorf("after 'w' value '%s' is not integer", v)
				break
			}
			if (intNumber < 1) || (intNumber > 7) {
				err = fmt.Errorf("after 'w' value '%s' is not between 1 and 7", v)
				break
			}
			weekDays = append(weekDays, intNumber) // День недели
		}
	case "m": // Дни месяца
		if (len(rules) < 2) || (len(rules) > 3) {
			err = fmt.Errorf("after 'm' must be one or two integer ranges")
			break
		}
		days := strings.Split(rules[1], ",")
		for _, v := range days {
			intNumber, err = strconv.Atoi(v)
			if err != nil {
				err = fmt.Errorf("after 'm' value '%s' is not integer", v)
				break
			}
			if intNumber < -2 {
				err = fmt.Errorf("after 'm' value '%s' must be greater than -3", v)
				break
			}
			if intNumber == 0 {
				err = fmt.Errorf("after 'm' value of '%s' must not be equal to 0", v)
				break
			}
			if intNumber > 31 {
				err = fmt.Errorf("after 'm' the value of '%s' must be less than 31", v)
				break
			}
			monthDays = append(monthDays, intNumber) // День месяца
		}
		if err == nil && len(rules) == 3 {
			months := strings.Split(rules[2], ",")
			for _, v := range months {
				intNumber, err = strconv.Atoi(v)
				if err != nil {
					err = fmt.Errorf("after 'm' value '%s' is not integer", v)
					break
				}
				if (intNumber < 1) || (intNumber > 12) {
					err = fmt.Errorf("after 'm' value '%s' must be greater than -3", v)
					break
				}
				monthNumbers = append(monthNumbers, intNumber) // Номер месяца
			}
		}
	}

	return
}

// Присутствует ли target значение в массиве arr
func contains(arr []int, target int) bool {
	for _, num := range arr {
		if num == target {
			return true
		}
	}
	return false
}

// Вернет true, если в массиве arr присутсвуют значения -1 или -2 и
// дата date представляет собой последний или предпоследний день месяца
func checkOnLastDays(arr []int, date time.Time) (result bool) {

	result = false
	for _, num := range arr {
		if num == -1 { // Последний день месяца
			result = date.AddDate(0, 0, 1).Month() != date.Month()
		} else {
			if num == -2 { // Предпоследний день месяца
				result = date.AddDate(0, 0, 2).Month() != date.Month()
			}
		}
	}

	return result
}
