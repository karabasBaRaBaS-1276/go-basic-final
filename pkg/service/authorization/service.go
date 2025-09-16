package service_authorization

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"slices"

	"github.com/golang-jwt/jwt/v5"
	models "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/models"
)

const (
	issuer   = "github.com/karabasBaRaBaS-1276/go-basic-final"
	audience = "web-app"
	subject  = "admin"
)

// Структура сервиса
type Service struct {
	log    *log.Logger // Логгер
	secret []byte      // Секретный ключ для подписи
}

// Claims - структура для хранения данных в токене
type Claims struct {
	jwt.RegisteredClaims // Базовый тип
}

// Инициализация экземпляра структуры Service
func New(log *log.Logger) *Service {
	return &Service{log: log, secret: []byte(os.Getenv("TODO_JWT_SECRET_KEY"))}
}

// Авторизация пользователя
// Принимает на вход:
//   - auth - информация для авторизации
//
// Возвращает:
//   - Токен доступа
//   - Ошибка
func (service *Service) Signin(auth *models.Auth) (string, error) {
	log := service.log
	log.Printf("   Service 'Signin' Begin\n")
	// Проверяем на корректность указанных данных
	if auth.Password != os.Getenv("TODO_PASSWORD") {
		return "", errors.New("incorrect password")
	}
	now := time.Now()

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,                                    // Субъект, которому токен выдается
			Audience:  jwt.ClaimStrings{audience},                 // Приложение, которое запросило токен
			Issuer:    issuer,                                     // Кто выпустил токен
			IssuedAt:  jwt.NewNumericDate(now),                    // Время выпуска
			ExpiresAt: jwt.NewNumericDate(now.Add(8 * time.Hour)), // Время истечения
			NotBefore: jwt.NewNumericDate(now),                    // Время, с которого токен действителен
			ID:        fmt.Sprintf("%s-%d", subject, now.Unix()),  // Уникальный ID токена
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// получаем подписанный токен
	signedToken, err := jwtToken.SignedString(service.secret)
	if err != nil {
		return "", err
	}
	service.log.Printf("token with jti = '%s' generated", claims.ID)
	// Выпускаем токен
	return signedToken, nil
}

// Проверить токен
//
// Принимает на вход:
//   - токен для проверки
//
// Возвращает:
//   - true если токен успешно проверен
//   - Ошибка
func (service *Service) ValidateToken(tokenString string) (bool, error) {

	claims := Claims{}
	if tokenString == "" {
		return false, errors.New("token is invalid")
	}

	//
	jwtToken, err := jwt.ParseWithClaims(tokenString, &claims,
		func(token *jwt.Token) (interface{}, error) {
			// желательно проверять используемый метод
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return service.secret, nil
		})
	if err != nil {
		return false, err
	}
	if !jwtToken.Valid {
		return false, errors.New("token is invalid")
	}

	now := time.Now()
	if claims.ExpiresAt.Time.Before(now) {
		return false, errors.New("token expired")
	}

	if claims.NotBefore.Time.After(now) {
		return false, errors.New("token is not active yet")
	}

	if claims.Issuer != issuer {
		return false, errors.New("invalid token issuer")
	}

	if len(claims.Audience) > 0 {
		audienceValid := slices.Contains(claims.Audience, audience)
		if !audienceValid {
			return false, errors.New("invalid token audience")
		}
	}

	if claims.Subject != subject {
		return false, errors.New("invalid token subject")
	}
	service.log.Printf("token with jti = '%s' is valid", claims.ID)

	return true, nil
}
