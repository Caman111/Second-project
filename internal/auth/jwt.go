package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("JWT_SECRET")

func ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Join(errors.New("unexpected signing method"))
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		phone, _ := claims["phone"].(string)
		return phone, nil
	}

	return "", errors.New("invalid token")
}
