package middleware

import (
	"3-validation-api/internal/auth"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "НЕТ. Забыл заголовок или 'Bearer '", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		email, err := auth.ParseToken(tokenString)
		if err != nil {
			http.Error(w, "НЕТ. Токен кривой: "+err.Error(), http.StatusUnauthorized)
			return
		}
		r.Header.Set("X-User-Email", email)

		next.ServeHTTP(w, r)
	}
}
