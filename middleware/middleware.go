package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"3-validation-api/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte("JWT_SECRET")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Отсутствует заголовок или неверный формат", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
			}
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Невалидный токен", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Ошибка чтения claims", http.StatusUnauthorized)
			return
		}
		phone, ok := claims["phone"].(string)
		if !ok || phone == "" {
			http.Error(w, "В токене отсутствует телефон", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), models.UserPhoneKey, phone)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
