package models

// Структура, описывающая ошибку
type Auth struct {
	Password string `json:"password"` // пароль
}

type JWT struct {
	Token string `json:"token"`
}
