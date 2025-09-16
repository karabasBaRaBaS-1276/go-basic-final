package api

import (
	"log"
	"net/http"
	"runtime/debug"

	service_authorization "github.com/karabasBaRaBaS-1276/go-basic-final/pkg/service/authorization"
)

// Middleware обработчик с целью проверки валидности токена и вызова следующего обработчика
func JWTAuthMiddleware(log *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var tokenString string
		// получаем куку
		cookie, err := request.Cookie("token")
		if err == nil {
			tokenString = cookie.Value
		}
		valid, err := service_authorization.New(log).ValidateToken(tokenString)

		if !valid {
			// возвращаем ошибку авторизации 401
			log.Printf("JWTAuthMiddleware: %s\n", err.Error())
			http.Error(writer, err.Error(), http.StatusUnauthorized)
			return
		}
		log.Println("JWTAuthMiddleware: token is valid")
		next.ServeHTTP(writer, request)
	})
}

// Middleware обработчик panic ошибки, которая может быть получена в нижестоящих обработчиках
func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if p := recover(); p != nil {
				log.Printf("panic recovered: %v, stack: %s\n", p, debug.Stack())

				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte(`{"error": "internal server error"}`))
			}
		}()

		next.ServeHTTP(writer, request)
	})
}
