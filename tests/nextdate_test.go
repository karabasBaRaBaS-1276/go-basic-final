package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/karabasBaRaBaS-1276/go-basic-final/pkg/api"
	"github.com/stretchr/testify/assert"
)

type nextDateFunc struct {
	date      string
	repeat    string
	want      string
	wantError string
}

var errorFormatDate = errors.New("incorrect date format in 'dstart' variable")
var errorRepeatIsEmpty = errors.New("variable 'repeat' cannot be empty")
var errorRepeatFormat = errors.New("incorrect format in 'repeat' variable")

func TestNextDateFunc(t *testing.T) {
	tbl := []nextDateFunc{
		{"20240126", "", "", errorRepeatIsEmpty.Error()},
		{"20240126", "k 34", "", errorRepeatFormat.Error()},
		{"20240126", "ooops", "", errorRepeatFormat.Error()},
		{"15000156", "y", "", errorFormatDate.Error()},
		{"ooops", "y", "", errorFormatDate.Error()},
		{"16890220", "y", `20240220`, ""},
		{"16890220", "y 1", `20240220`, errorRepeatFormat.Error()},
		{"20250701", "y", `20260701`, ""},
		{"20240101", "y", `20250101`, ""},
		{"20231231", "y", `20241231`, ""},
		{"20240229", "y", `20250301`, ""},
		{"20240301", "y", `20250301`, ""},
		{"20240113", "d", "", errorRepeatFormat.Error()},
		{"20240320", "d tst", "", errorRepeatFormat.Error()},
		{"20240113", "d 7", `20240127`, ""},
		{"20240120", "d 20", `20240209`, ""},
		{"20240202", "d 30", `20240303`, ""},
		{"20240320", "d 401", "", errorRepeatFormat.Error()},
		{"20231225", "d 12", `20240130`, ""},
		{"20240228", "d 1", "20240229", ""},
	}
	check := func() {
		for _, v := range tbl {
			now, _ := time.Parse("20060102", "20240126") // Если сегодня 26 января 2024 года
			result, err := api.NextDate(now, v.date, v.repeat)

			if err != nil {
				assert.ErrorContains(t, err, v.wantError, "Данные для проверки: {Начальное время: %q, Правило для повтора: %q, Ждем в ошибке: %q}",
					v.date, v.repeat, v.wantError)
			} else {
				assert.Equal(t, v.want, result, "Данные для проверки: {Начальное время: %q, Правило для повтора: %q, Ожидаемое значение: %q}",
					v.date, v.repeat, v.want)
			}
		}
	}
	check()
	tbl = []nextDateFunc{
		{"20231106", "m", "", errorRepeatFormat.Error()},
		{"20231106", "m -1,u 5,t", "", errorRepeatFormat.Error()},
		{"20231106", "m 0", "", errorRepeatFormat.Error()},
		{"20231106", "m 13", "20240213", ""},
		{"20240120", "m 40,11,19", "", errorRepeatFormat.Error()},
		{"20240116", "m 16,5", "20240205", ""},
		{"20240126", "m 25,26,7", "20240207", ""},
		{"20240409", "m 31", "20240531", ""},
		{"20240329", "m 10,17 12,8,1", "20240810", ""},
		{"20230311", "m 07,19 05,6", "20240507", ""},
		{"20230311", "m 1 1,2", "20240201", ""},
		{"20240127", "m -1", "20240131", ""},
		{"20240222", "m -2", "20240228", ""},
		{"20240222", "m -2,-3", "", errorRepeatFormat.Error()},
		{"20240326", "m -1,-2", "20240330", ""},
		{"20240201", "m -1,18", "20240218", ""},
		{"20240125", "w", "", errorRepeatFormat.Error()},
		{"20240125", "w week", "", errorRepeatFormat.Error()},
		{"20240125", "w 1,3,6,i", "", errorRepeatFormat.Error()},
		{"20240125", "w 1,2,3", "20240129", ""},
		{"20240126", "w 7", "20240128", ""},
		{"20230126", "w 4,5", "20240201", ""},
		{"20230226", "w 8,4,5", "", errorRepeatFormat.Error()},
	}
	check()
}
